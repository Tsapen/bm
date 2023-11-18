package bm

import (
	"context"
)

type cxtKey int

const (
	reqIDKey cxtKey = iota
)

func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

func ReqIDFromCtx(ctx context.Context) string {
	return ctx.Value(reqIDKey).(string)
}
