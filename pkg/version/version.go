// Package version contains package global variables
// goreleaser usage: https://goreleaser.com/cookbooks/using-main.version/
package version

import (
	"fmt"
)

var (
	// Version is the build version
	Version = "dev"
	// Commit is the build commit
	Commit = "none"
	// Date is the build date
	Date = "unknown"
	// Build carries the build information of the current build. It is initialized in the package init
	Build BuildInfo
)

// BuildInfo bundles the build information
// swagger:response buildInfo
type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// String implements the Stringer interface
func (info BuildInfo) String() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nCommit Date: %s\n", info.Version, info.Commit, info.Date)
}

func init() {
	Build.Version = Version
	Build.Commit = Commit
	Build.Date = Date
}
