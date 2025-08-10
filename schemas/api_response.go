package schemas

import (
	"context"
)

type APIResponse interface {
	SetError(err error)
}

func MakeErrorResponse[T APIResponse](ctx context.Context, resp T, err error) T {
	_, span := tracer.Start(ctx, "build error response")
	defer span.End()
	resp.SetError(err)
	return resp
}
