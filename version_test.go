package indra_test

import (
	"regexp"
	"testing"

	indra "coderaiser/indra"

	tape "github.com/coderaiser/go-tape"
)

func TestVersionFromJSON(t *testing.T) {
	tape.Test(t, "version: returns version string", func(t *tape.T) {
		t.Equal(indra.VersionFromJSON([]byte("{\"version\":\"1.2.3\"}")), "1.2.3")
		t.End()
	})

	tape.Test(t, "version: invalid JSON returns unknown", func(t *tape.T) {
		t.Equal(indra.VersionFromJSON([]byte("{invalid")), "unknown")
		t.End()
	})

	tape.Test(t, "version: empty version returns unknown", func(t *tape.T) {
		t.Equal(indra.VersionFromJSON([]byte("{\"version\":\"\"}")), "unknown")
		t.End()
	})
}

func TestVersionLine(t *testing.T) {
	tape.Test(t, "version: VersionLine contains binary name", func(t *tape.T) {
		t.Match(indra.VersionLine(), "indra ")
		t.End()
	})
}

func TestVersion(t *testing.T) {
	tape.Test(t, "version: reads package.json at runtime", func(t *tape.T) {
		semver := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
		t.Ok(semver.MatchString(indra.Version()))
		t.End()
	})

	tape.Test(t, "version: missing package.json returns unknown", func(t *tape.T) {
		t.TB().Chdir(t.TB().TempDir())
		t.Equal(indra.Version(), "unknown")
		t.End()
	})
}
