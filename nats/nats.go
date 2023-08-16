package nats

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type Configuration struct {
	Host     string
	Port     uint16
	Stream   string
	Consumer struct {
		Name    string
		Subject string
		Queue   string
	}
	Workers int
}

type MessageHandler func(ctx context.Context, msg *nats.Msg) error

type Messenger struct {
	conn         *nats.Conn
	jet          nats.JetStreamContext
	subscription *nats.Subscription

	handlers map[string]MessageHandler

	closeChan chan bool
	workers   []chan context.Context
}

func NewMessenger(cfg *Configuration) *Messenger {
	messenger := Messenger{
		handlers:  make(map[string]MessageHandler, 0),
		closeChan: make(chan bool, 1),
		workers:   make([]chan context.Context, 0, cfg.Workers),
	}
	return &messenger
}

func (m *Messenger) AddHandler(subject string, handler MessageHandler) error {
	if _, ok := m.handlers[subject]; ok {
		return fmt.Errorf("handler for subject [%s] already registered", subject)
	}
	m.handlers[subject] = handler
	return nil
}

func (m *Messenger) Start(ctx context.Context, cfg *Configuration) error {
	url := fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port)
	conn, err := nats.Connect(url, nats.ClosedHandler(func(conn *nats.Conn) {
		close(m.closeChan)
	}))
	if err != nil {
		return fmt.Errorf("nats connectioni failed %w", err)
	}
	m.conn = conn
	jet, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf("jet stream init failed: %w", err)
	}
	m.jet = jet

	ch := make(chan *nats.Msg, cfg.Workers)
	_, err = jet.ChanQueueSubscribe(cfg.Consumer.Subject, cfg.Consumer.Queue, ch, nats.Bind(cfg.Stream, cfg.Consumer.Name))
	if err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	log := zerolog.Ctx(ctx)
	for i := 0; i < cfg.Workers; i++ {
		workerId := i
		go func() {
			m.worker(workerId, ch)
		}()
	}
	log.Info().Msg("messenger started")
	return nil
}

func (m *Messenger) Stop(ctx context.Context) {
	log := zerolog.Ctx(ctx)

	if err := m.conn.Drain(); err != nil {
		log.Error().Err(err).Msg("failed to drain nats connection")
	} else {
		log.Info().Msg("nats connection drain started")
	}

	for _, worker := range m.workers {
		close(worker)
	}

	select {
	case <-m.closeChan:
		log.Info().Msg("nats connection closed")
	case <-ctx.Done():
		log.Warn().Msg("nats connection drain timeout")
	}
}

func (m *Messenger) worker(id int, inputChan chan *nats.Msg) {
	// Close worker when done
	closeChan := make(chan context.Context, 1)
	m.workers = append(m.workers, closeChan)

	handleSafe := func(ctx context.Context, msg *nats.Msg, handler MessageHandler) error {
		defer func() {
			log := zerolog.Ctx(ctx)
			if r := recover(); r != nil {
				log.Error().Any("recovered", r).Msg("message handler panicked")
			}
		}()
		return handler(ctx, msg)
	}
	for {
		select {
		case msg, ok := <-inputChan:
			{
				if !ok {
					// channel closed - do nothing
					continue
				}
				contextId := msg.Header.Get("contextId")
				if contextId == "" {
					contextId = uuid.New().String()
				}
				ctx := support.WithContextId(context.Background(), contextId)
				log := zerolog.DefaultContextLogger.With().
					Str("contextId", contextId).
					Int("msgWorker", id).
					Logger()
				ctx = log.WithContext(ctx)
				handler, found := m.handlers[msg.Subject]
				if !found {
					log.Error().Str("subject", msg.Subject).Msg("handler not found")
					continue
				}

				if err := handleSafe(ctx, msg, handler); err != nil {
					log.Error().Err(err).Msg("handler returned an error")
					return
				}

				if err := msg.Ack(nats.ContextOpt{Context: ctx}); err != nil {
					log.Error().Err(err).Msg("Ack timeout")
					return
				}
			}
		case ctx := <-closeChan:
			log := zerolog.Ctx(ctx)
			log.Debug().Int("worker", id).Msg("msg worker stopped")
			return
		}
	}

}
