package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/iyarkov/foundation/auth"
	"github.com/iyarkov/foundation/logger"
	"github.com/iyarkov/foundation/support"
	"github.com/iyarkov/foundation/telemetry"
	"github.com/iyarkov/foundation/tls"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var contextIdMeta = "contextId"
var authTokenMeta = "authToken"

type connectionInfoCtxKey struct{}

// ServerContextId Server Side Interceptor
func ServerContextId(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// Get request ID from the request
	meta, ok := metadata.FromIncomingContext(ctx)
	var contextId string
	if ok {
		contextIdSlice := meta.Get(contextIdMeta)
		if len(contextIdSlice) > 0 {
			contextId = contextIdSlice[0]
		}
	}
	if contextId == "" {
		contextId = uuid.New().String()
	}

	ctx = logger.WithContextIdAndLogger(ctx, contextId)
	return handler(ctx, req)
}

// ClientContextId Client Side Interceptor
func ClientContextId(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	zerolog.Ctx(ctx).Debug().Msg("ClientContextId")
	contextId := support.ContextId(ctx)
	if contextId != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, contextIdMeta, contextId)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

// ServerConnectionInfo Client Side Interceptor
func ServerConnectionInfo(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := zerolog.Ctx(ctx)
	client, ok := peer.FromContext(ctx)
	tlsInfo, ok := client.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("unable to authentificate using %s", client.AuthInfo.AuthType()))
	}
	var connectionInfo tls.ConnectionInfo
	connectionInfo.FromPeerCertificate(tlsInfo.State.PeerCertificates[0])
	log.Debug().Msgf("Connection Info %v", connectionInfo)
	ctx = context.WithValue(ctx, &connectionInfoCtxKey{}, connectionInfo)
	return handler(ctx, req)
}

func ConnectionInfo(ctx context.Context) tls.ConnectionInfo {
	if ci, ok := ctx.Value(&connectionInfoCtxKey{}).(tls.ConnectionInfo); ok {
		return ci
	}
	return tls.ConnectionInfo{}
}

func NewServerAuthInterceptor(initCtx context.Context, conf *auth.Configuration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		connectionInfo := ConnectionInfo(ctx)
		trusted := false
		for _, p := range conf.TrustedPeers {
			if p == connectionInfo.Peer {
				trusted = true
			}
		}

		var encodedToken string
		meta, ok := metadata.FromIncomingContext(ctx)
		if ok {
			tokenSlice := meta.Get(authTokenMeta)
			if len(tokenSlice) > 0 {
				encodedToken = tokenSlice[0]
			}
		}
		log := zerolog.Ctx(ctx)
		if encodedToken == "" {
			log.Debug().Msgf("auth: client trusted: %t, not authenticated", trusted)
			return handler(ctx, req)
		}
		if trusted {
			ctx, err = auth.WithStringToken(ctx, encodedToken)
			if err != nil {
				log.Debug().Err(err).Msgf("invalid token")
				return nil, status.Error(codes.InvalidArgument, "unable to parse auth token")
			}
		} else {
			log.Debug().Msgf("auth: client not trusted, not  implemented yet")
			return nil, status.Error(codes.Unimplemented, "non-trusted peers are not implemented yet")
		}
		logEvent := log.Debug()
		if logEvent.Enabled() {
			token := auth.AuthToken(ctx)
			log.Debug().Err(err).Msgf("auth: client trusted: %t, token %+v", trusted, token)
		}
		return handler(ctx, req)
	}
}

func ClientAuth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	zerolog.Ctx(ctx).Debug().Msg("ClientAuth")
	token := auth.AuthToken(ctx)
	if token.IsAuthenticated() {
		ctx = metadata.AppendToOutgoingContext(ctx, authTokenMeta, token.WriteToString())
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

func ServerTrace(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx, span := telemetry.StartSpan(ctx, fmt.Sprintf("grpc%s", info.FullMethod))
	defer span.End()
	contextId := support.ContextId(ctx)
	if contextId != "" {
		span.SetAttributes(attribute.String("contextId", contextId))
	}
	return handler(ctx, req)
}

func ClientTrace(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	zerolog.Ctx(ctx).Debug().Msg("ClientTrace")
	ctx, span := telemetry.StartSpan(ctx, fmt.Sprintf("grpc%s", method))
	defer span.End()
	return invoker(ctx, method, req, reply, cc, opts...)
}
