package support

import "context"

type ctxKey struct{}

func WithContextId(ctx context.Context, ctxId string) context.Context {
	return context.WithValue(ctx, &ctxKey{}, ctxId)
}

func ContextId(ctx context.Context) string {
	if str, ok := ctx.Value(&ctxKey{}).(string); ok {
		return str
	}
	return ""
}
