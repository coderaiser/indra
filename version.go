package indra

import (
	"encoding/json"
	"fmt"
	"os"
)

// VersionFromJSON returns the version string from package.json bytes.
func VersionFromJSON(packageJSONBytes []byte) string {
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(packageJSONBytes, &pkg); err != nil {
		return "unknown"
	}
	if pkg.Version == "" {
		return "unknown"
	}
	return pkg.Version
}

// Version reads the version from package.json at runtime.
func Version() string {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return "unknown"
	}
	return VersionFromJSON(data)
}

// VersionLine returns "indra x.y.z" for -v output.
func VersionLine() string {
	return fmt.Sprintf("v%s", Version())
}
