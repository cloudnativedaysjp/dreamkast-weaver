package value

import "fmt"

type valueObject[T comparable] struct {
	value T
}

func (v *valueObject[T]) Value() T {
	return v.value
}

func (v *valueObject[T]) String() string {
	return fmt.Sprintf("%v", v.value)
}

func (v *valueObject[T]) Validate() error {
	return nil
}

func wrap[T comparable](value T) valueObject[T] {
	return valueObject[T]{value}
}
