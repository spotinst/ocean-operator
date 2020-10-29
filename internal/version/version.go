package version

import (
	"strings"

	"github.com/hashicorp/go-version"
)

var (
	// Version represents the main version number.
	//
	// Populated at compile time.
	// Read-only.
	Version string

	// Prerelease represents an optional pre-release label for the version.
	// If this is "" (empty string) then it means that it is a final release.
	// Otherwise, this is a pre-release such as "beta", "rc1", etc.
	//
	// Populated at compile time.
	// Read-only.
	Prerelease string

	// Metadata represents an optional build metadata.
	//
	// Populated at compile time.
	// Read-only.
	Metadata string

	// SemVer is an instance of SemVer version (https://semver.org) used to
	// verify that the full version is a proper semantic version.
	//
	// Populated at runtime.
	// Read-only.
	SemVer *version.Version
)

func init() {
	v := Version

	if Prerelease != "" {
		v += "-" + strings.TrimPrefix(Prerelease, "-")
	}

	if Metadata != "" {
		v += "+" + strings.TrimPrefix(Metadata, "+")
	}

	// Parse and verify the given version.
	SemVer = version.Must(version.NewSemver(v))
}

// String returns the complete version string, including prerelease and metadata.
func String() string {
	return SemVer.String()
}
