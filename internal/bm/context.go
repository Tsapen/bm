package bm

import (
	"context"
)

type cxtKey int

const (
	reqIDKey cxtKey = iota
)

// WithReqID adds request id into context.
func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

// WithReqID gets request id from context.
func ReqIDFromCtx(ctx context.Context) string {
	return ctx.Value(reqIDKey).(string)
}
