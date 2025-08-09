package schemas

type ErrorResponse interface {
	SetError(err error)
}

func BuildErrorResponse[T ErrorResponse](resp T, err error) T {
	resp.SetError(err)
	return resp
}
