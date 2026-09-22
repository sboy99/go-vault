// Package version holds the go-vault release version.
// Release tags stamp this via ldflags; scripts/set-version.sh updates the default.
package version

// Version is the semver string reported by the CLI and used as app.version default.
// Overridden at link time: -X github.com/sboy99/go-vault/internal/version.Version=X.Y.Z
var Version = "0.3.0"
