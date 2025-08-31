package schemas

type APIResponse interface {
	// Given a response, set all the error attributes, and clear the success attributes
	SetErrorAttributes(err error)
}
