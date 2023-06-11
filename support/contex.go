package support

import "context"

type ctxKey struct{}

func SetContextId(ctx context.Context, ctxId string) context.Context {
	return context.WithValue(ctx, &ctxKey{}, ctxId)
}

func GetContextId(ctx context.Context) string {
	if str, ok := ctx.Value(&ctxKey{}).(string); ok {
		return str
	}
	return ""
}
