package sql

import (
	"context"
	"fmt"
	"github.com/iyarkov/foundation/telemetry"
	"github.com/jackc/pgx/v5"
)

type OpenTelemetryTracer struct {
}

func (t *OpenTelemetryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	var traceName string
	if data.SQL == ";" {
		traceName = "sql:ping"
	} else {
		traceName = fmt.Sprintf("sql:%s", data.SQL)
	}
	ctx, _ = telemetry.StartSpan(ctx, traceName)
	return ctx
}

func (t *OpenTelemetryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	telemetry.SpanFromContext(ctx).End()
}
