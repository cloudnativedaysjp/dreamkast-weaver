package valueobject_test

import (
	"dreamkast-weaver/internal/valueobject"
	"fmt"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Foo struct {
	valueobject.Wrapped[string]
}

func NewFoo(v string) (*Foo, error) {
	o := &Foo{valueobject.Wrap(v)}
	if err := o.Validate(); err != nil {
		return nil, err
	}
	return o, nil
}

func (v *Foo) Validate() error {
	return validation.Validate(v.Value(),
		validation.Length(2, 3),
	)
}

func TestWrap(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		given := "foo"

		got, err := NewFoo(given)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Value() != given {
			t.Errorf("not equal: want=%#v, got=%#v", given, got.Value())
		}
		if fmt.Sprintf("%s", got) != given {
			t.Errorf("not equal: want=%#v, got=%#v", given, got.String())
		}
	})

	errtests := []struct {
		name  string
		given string
	}{
		{
			name:  "too short",
			given: "a",
		},
		{
			name:  "too long",
			given: "long",
		},
	}
	for _, tt := range errtests {
		t.Run("err: "+tt.name, func(t *testing.T) {
			got, err := NewFoo(tt.given)
			if err == nil {
				t.Errorf("error not raised: %#v", got)
			}
		})
	}
}
