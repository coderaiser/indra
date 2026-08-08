package convert_switch_to_if_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_indra/convert_switch_to_if"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-switch-to-if", convert_switch_to_if.Plugin{})

func TestConvertSwitchToIf(t *testing.T) {
	Test(t, "convert-switch-to-if: report convert-switch-to-if", func(t *T) {
		t.Report("convert-switch-to-if", "use 'if' instead of 'switch'")
		t.End()
	})

	Test(t, "convert-switch-to-if: transform convert-switch-to-if", func(t *T) {
		t.Transform("convert-switch-to-if")
		t.End()
	})

	Test(t, "convert-switch-to-if: no report has-default", func(t *T) {
		t.NoReport("has-default")
		t.End()
	})

	Test(t, "convert-switch-to-if: no report no-return", func(t *T) {
		t.NoReport("no-return")
		t.End()
	})

	Test(t, "convert-switch-to-if: no report multi-value-case", func(t *T) {
		t.NoReport("multi-value-case")
		t.End()
	})

	Test(t, "convert-switch-to-if: no report has-fallthrough", func(t *T) {
		t.NoReport("has-fallthrough")
		t.End()
	})

	Test(t, "convert-switch-to-if: no report has-init", func(t *T) {
		t.NoReport("has-init")
		t.End()
	})

	Test(t, "convert-switch-to-if: report has-break", func(t *T) {
		t.Report("has-break", "use 'if' instead of 'switch'")
		t.End()
	})

	Test(t, "convert-switch-to-if: transform has-break", func(t *T) {
		t.Transform("has-break")
		t.End()
	})

	Test(t, "convert-switch-to-if: transform nested-switch", func(t *T) {
		t.Transform("nested-switch")
		t.End()
	})
}
