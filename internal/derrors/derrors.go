package derrors

import (
	"errors"
	"fmt"
)

type UserError struct {
	msg string
	err error
}

func (e *UserError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	} else {
		return e.msg
	}
}

func (e *UserError) Wrap(err error) error {
	//nolint
	ee := NewUserError(e.msg).(*UserError)
	ee.err = err
	return ee
}

func (e *UserError) Unwrap() error {
	return e.err
}

func NewUserError(msg string) error {
	return &UserError{
		msg: msg,
	}
}

func IsUserError(err error) bool {
	ue := (*UserError)(nil)
	return errors.As(err, &ue)
}
