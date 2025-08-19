package utility

// Ptr returns a pointer to the provided value
func Ptr[T any](s T) *T {
	return &s
}
