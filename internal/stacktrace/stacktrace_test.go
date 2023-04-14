package stacktrace_test

import (
	"bytes"
	"dreamkast-weaver/internal/stacktrace"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestShowStackTrace(t *testing.T) {
	t.Run("nested error", func(t *testing.T) {
		err1 := errors.New("foo")
		err2 := errors.Wrap(err1, "bar")
		err3 := errors.Wrap(err2, "baz")

		buf := &bytes.Buffer{}
		stacktrace.Show(err3, buf)

		found := regexp.MustCompile("testing.tRunner").FindAllIndex(buf.Bytes(), -1)
		assert.Equal(t, 1, len(found))
	})

	t.Run("single error", func(t *testing.T) {
		err := errors.New("foo")

		buf := &bytes.Buffer{}
		stacktrace.Show(err, buf)

		found := regexp.MustCompile("testing.tRunner").FindAllIndex(buf.Bytes(), -1)
		assert.Equal(t, 1, len(found))
	})

	t.Run("nil", func(t *testing.T) {
		buf := &bytes.Buffer{}
		stacktrace.Show(nil, buf)
		t.Log(buf)

		assert.Equal(t, 0, len(buf.Bytes()))
	})
}
