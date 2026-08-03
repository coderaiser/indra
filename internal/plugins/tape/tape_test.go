package tape_test

import (
	"testing"

)

func TestRulesKeys(t *testing.T) {
	for _, key := range []string{"remove-skip", "add-t-end"} {
		if _, ok := Rules[key]; !ok {
			t.Fatalf("Rules missing key %q", key)
		}
	}
}

func TestRulesValues(t *testing.T) {
	for k, v := range Rules {
		if v == nil {
			t.Fatalf("Rules[%q] is nil", k)
		}
	}
}
