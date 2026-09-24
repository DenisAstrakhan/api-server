package domain

// generic type
type Nullable[T any] struct {
	Value *T
	Set   bool
}
