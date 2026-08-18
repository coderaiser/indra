package tape

import "testing"

// TestRules guards tape.go's Rules() entry point, which the registry calls at
// package init. Without an in-package test it reports 0% coverage.
func TestRules(t *testing.T) {
	rules := Rules()
	if len(rules) == 0 {
		t.Fatal("Rules() returned no rules")
	}
	for _, r := range rules {
		if r.Name == "" {
			t.Fatal("rule has empty name")
		}
		if r.Plugin == nil {
			t.Fatal("rule has nil plugin")
		}
	}
}
