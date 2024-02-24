package common

func Pointer[T any](element T) *T {
	return &element
}
