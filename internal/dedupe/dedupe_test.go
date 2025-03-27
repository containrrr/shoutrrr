package dedupe_test

import (
	"reflect"
	"testing"

	"github.com/nicholas-fedor/shoutrrr/internal/dedupe"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := map[string]struct {
		input []string
		want  []string
	}{
		"no duplicates":                             {input: []string{"a", "b", "c"}, want: []string{"a", "b", "c"}},
		"duplicate inside slice":                    {input: []string{"a", "b", "a", "c"}, want: []string{"a", "b", "c"}},
		"duplicate at end of slice":                 {input: []string{"a", "b", "c", "a"}, want: []string{"a", "b", "c"}},
		"duplicate next to each other inside slice": {input: []string{"a", "b", "b", "c"}, want: []string{"a", "b", "c"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := dedupe.RemoveDuplicates(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
		})
	}
}
