package convert_equal_to_deep_equal_test

import (
	"path/filepath"
	"runtime"
	"testing"

	indratest "coderaiser/indra/internal/test"
)

var (
	_, _file, _, _ = runtime.Caller(0)
	_dir           = filepath.Join(filepath.Dir(_file), "fixture")
	Test           = indratest.CreateTest("coderaiser/indra/internal/plugins/convert-equal-to-deep-equal", _dir)
)

func TestConvertEqualToDeepEqual(t *testing.T) {
	Test(t, "convert-equal-to-deep-equal: report array second arg", func(t *indratest.T) {
		t.Report("equal-slice", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: transform array second arg", func(t *indratest.T) {
		t.Transform("equal-slice")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: report array first arg", func(t *indratest.T) {
		t.Report("equal-slice-first", "Equal: use DeepEqual for slices")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no report for DeepEqual", func(t *indratest.T) {
		t.NoReport("deep-equal")
		t.End()
	})

	Test(t, "convert-equal-to-deep-equal: no transform for DeepEqual", func(t *indratest.T) {
		t.NoTransform("deep-equal")
		t.End()
	})
}
