package observability

import "context"

const traceParentKey contextKey = "trace_parent"

func WithTraceParent(ctx context.Context, traceParent string) context.Context {
	return context.WithValue(ctx, traceParentKey, traceParent)
}

func TraceParentFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, ok := ctx.Value(traceParentKey).(string)
	if !ok {
		return ""
	}
	return value
}
