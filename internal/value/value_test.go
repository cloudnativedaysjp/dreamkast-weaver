package value_test

import (
	"testing"

	"dreamkast-weaver/internal/value"

	"github.com/stretchr/testify/assert"
)

func TestConfName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		given := "cicd2023"

		got, err := value.NewConfName(value.ConferenceKind(given))
		assert.NoError(t, err)
		assert.Equal(t, value.CICD2023, got)
	})

	errTests := []struct {
		name  string
		given string
	}{
		{
			name:  "not exist",
			given: "cicd2022",
		},
	}

	for _, tt := range errTests {
		t.Run("err: "+tt.name, func(t *testing.T) {
			_, err := value.NewConfName(value.ConferenceKind(tt.given))
			assert.Error(t, err)
		})
	}
}
