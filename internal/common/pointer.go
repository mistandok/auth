package common

// Pointer get pointer for any type.
func Pointer[T any](element T) *T {
	return &element
}
