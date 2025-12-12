package version

import (
	"fmt"
	"runtime"
)

var (
	// These will be set by goreleaser at build time
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get returns version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("LazyService %s (%s) built on %s with %s for %s",
		i.Version, i.Commit[:8], i.Date, i.GoVersion, i.Platform)
}

// Short returns a short version string
func (i Info) Short() string {
	return fmt.Sprintf("LazyService %s", i.Version)
}
