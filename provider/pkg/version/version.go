// File: provider/pkg/version/version.go
package version

// Version is initialized by the Go linker to contain the semver of this build.
var Version string = "0.0.0"

// GetVersion returns the current version of the package.
func GetVersion() string {
	return Version
}
