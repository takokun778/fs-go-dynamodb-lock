package gateway

import (
	"context"

	"github.com/google/uuid"
)

type key string

const ctxRID = key("rid")

func SetRIDCtx(parent context.Context) context.Context {
	id := uuid.New().String()

	return context.WithValue(parent, ctxRID, id)
}

func GetRIDCtx(ctx context.Context) string {
	v := ctx.Value(ctxRID)

	id, ok := v.(string)

	if !ok {
		return ""
	}

	return id
}
