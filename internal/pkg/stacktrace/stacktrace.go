package stacktrace

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// Show shows the stacktrace of the original error only.
func Show(err error, writer ...io.Writer) {
	var w io.Writer
	if len(writer) > 0 {
		w = writer[0]
	} else {
		w = os.Stdout
	}
	if st := Get(err); st != "" {
		fmt.Fprintf(w, "StackTrace: %s\n\n", st)
	}
}

// Get returns the stacktrace of the original error only.
func Get(err error) string {
	st := bottomStackTrace(err)
	if st != nil {
		return fmt.Sprintf("%+v", st.StackTrace())
	}
	return ""
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func bottomStackTrace(err error) stackTracer {
	nestedErr := errors.Unwrap(err)
	if nestedErr != nil {
		if st := bottomStackTrace(nestedErr); st != nil {
			return st
		}
	}
	// type check after checking all nested errors are not stackTracer
	if e, ok := err.(stackTracer); ok { // nolint
		return e
	}
	return nil
}

var (
	With = errors.WithStack
)
