package stacktrace_test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"dreamkast-weaver/internal/pkg/stacktrace"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestShowStackTrace(t *testing.T) {
	t.Run("nested error", func(t *testing.T) {
		err1 := func() error {
			return fmt.Errorf("bar: %w", stacktrace.With(fmt.Errorf("foo")))
		}
		err2 := func() error {
			return fmt.Errorf("baz: %w", stacktrace.With(err1()))
		}
		err := stacktrace.With(err2())

		buf := &bytes.Buffer{}
		stacktrace.Show(err, buf)

		found := regexp.MustCompile("testing.tRunner").FindAllIndex(buf.Bytes(), -1)
		assert.Equal(t, 1, len(found))
	})

	t.Run("single error", func(t *testing.T) {
		// pkg/errors.New has stacktrace originally.
		err := errors.New("foo")

		buf := &bytes.Buffer{}
		stacktrace.Show(err, buf)

		found := regexp.MustCompile("testing.tRunner").FindAllIndex(buf.Bytes(), -1)
		assert.Equal(t, 1, len(found))
	})

	t.Run("nil", func(t *testing.T) {
		buf := &bytes.Buffer{}
		stacktrace.Show(nil, buf)

		assert.Equal(t, 0, len(buf.Bytes()))
	})
}
