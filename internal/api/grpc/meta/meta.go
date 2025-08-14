package meta

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	RequestIDCtxKey
	UserCtxKey
)

func RequestID(ctx context.Context) string {
	if ctx == nil {
		return "unknow"
	}

	requestID, ok := ctx.Value(RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return "unknow"
	}

	return requestID.String()
}
