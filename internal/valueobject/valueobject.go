package valueobject

import "fmt"

type Wrapped[T comparable] struct {
	value T
}

func (v *Wrapped[T]) Value() T {
	return v.value
}

func (v *Wrapped[T]) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *Wrapped[T]) Validate() error {
	return nil
}

func Wrap[T comparable](value T) Wrapped[T] {
	return Wrapped[T]{value}
}
