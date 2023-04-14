package value_test

import (
	"dreamkast-weaver/internal/dkui/value"
	"testing"
)

func TestConfName(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		given := "cicd2023"

		got, err := value.NewConfName(value.ConferenceKind(given))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != value.CICD2023 {
			t.Errorf("not equal: want=%#v, got=%#v", value.CICD2023, got)
		}
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
			got, err := value.NewConfName(value.ConferenceKind(tt.given))
			if err == nil {
				t.Errorf("error not raised: %#v", got)
			}
		})
	}

}
