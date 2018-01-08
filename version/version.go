package version

import "github.com/coreos/go-semver/semver"

var (
	// verMajor is for an API incompatible changes.
	verMajor int64
	// verMinor is for functionality in a backwards-compatible manner.
	verMinor int64 = 1
	// verPatch is for backwards-compatible bug fixes.
	verPatch int64 = 0
	// verPre indicates pre-release.
	verPre string
	// verDev indicates development branch. Releases will be empty string.
	verDev string
)

// Version is the specification version that the package types support.
var Version = semver.Version{
	Major:      verMajor,
	Minor:      verMinor,
	Patch:      verPatch,
	PreRelease: semver.PreRelease(verPre),
	Metadata:   verDev,
}
