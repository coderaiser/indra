package convert_equal_to_not_ok_test

import (
	"testing"

	"coderaiser/indra/internal/plugin_tape/convert_equal_to_not_ok"
	. "coderaiser/indra/internal/test"
)

var Test = CreateTest("convert-equal-to-not-ok", convert_equal_to_not_ok.Plugin{})

func TestConvertEqualToNotOk(t *testing.T) {
	Test(t, "convert-equal-to-not-ok: report: convert-equal-to-not-ok", func(t *T) {
		t.Report("convert-equal-to-not-ok", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: convert-equal-to-not-ok", func(t *T) {
		t.Transform("convert-equal-to-not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: not-ok", func(t *T) {
		t.NoReport("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: not-ok", func(t *T) {
		t.NoTransform("not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report: convert-equal-false-to-not-ok", func(t *T) {
		t.Report("convert-equal-false-to-not-ok", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: convert-equal-false-to-not-ok", func(t *T) {
		t.Transform("convert-equal-false-to-not-ok")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: empty-string", func(t *T) {
		t.Transform("empty-string")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report: zero-int", func(t *T) {
		t.Report("zero-int", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: zero-int", func(t *T) {
		t.Transform("zero-int")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: zero-float", func(t *T) {
		t.NoReport("zero-float")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: zero-float", func(t *T) {
		t.NoTransform("zero-float")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report: named-false", func(t *T) {
		t.Report("named-false", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: named-false", func(t *T) {
		t.Transform("named-false")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: report: named-zero", func(t *T) {
		t.Report("named-zero", "convert Equal(x, nil/false) to NotOk(x)")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: transform: named-zero", func(t *T) {
		t.Transform("named-zero")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: non-empty-string", func(t *T) {
		t.NoReport("non-empty-string")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: non-empty-string", func(t *T) {
		t.NoTransform("non-empty-string")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: non-falsy-args", func(t *T) {
		t.NoReport("non-falsy-args")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: non-falsy-args", func(t *T) {
		t.NoTransform("non-falsy-args")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no report: imaginary-literal", func(t *T) {
		t.NoReport("imaginary-literal")
		t.End()
	})

	Test(t, "convert-equal-to-not-ok: no transform: imaginary-literal", func(t *T) {
		t.NoTransform("imaginary-literal")
		t.End()
	})
}
