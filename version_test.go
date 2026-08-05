package indra_test

import (
	"regexp"
	"testing"

	indra "coderaiser/indra"

	. "github.com/coderaiser/go-tape"
)

func TestVersionFromJSON(t *testing.T) {
	Test(t, "version: returns version string", func(t *T) {
		result := indra.VersionFromJSON([]byte("{\"version\":\"1.2.3\"}"))
		t.Equal(result, "1.2.3")
		t.End()
	})

	Test(t, "version: invalid JSON returns unknown", func(t *T) {
		result := indra.VersionFromJSON([]byte("{invalid"))
		t.Equal(result, "unknown")
		t.End()
	})

	Test(t, "version: empty version returns unknown", func(t *T) {
		result := indra.VersionFromJSON([]byte("{\"version\":\"\"}"))
		t.Equal(result, "unknown")
		t.End()
	})
}

func TestVersionLine(t *testing.T) {
	Test(t, "version: VersionLine contains binary name", func(t *T) {
		t.Match(indra.VersionLine(), "indra ")
		t.End()
	})
}

func TestVersion(t *testing.T) {
	Test(t, "version: reads package.json at runtime", func(t *T) {
		semver := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
		t.Ok(semver.MatchString(indra.Version()))
		t.End()
	})

	Test(t, "version: missing package.json returns unknown", func(t *T) {
		t.TB().Chdir(t.TB().TempDir())
		result := indra.Version()
		t.Equal(result, "unknown")
		t.End()
	})
}
