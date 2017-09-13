package version

import (
	"bytes"
	"fmt"
)

var (
	// GitCommit - The git commit that was compiled. This will be filled in by the compiler.
	GitCommit string
	// GitDescribe - ?
	GitDescribe string

	// Version is populated in the init function
	Version = "unknown"
	// VersionPrerelease is populated in the init function
	VersionPrerelease = "unknown"
)

// Info holds version information
type Info struct {
	Revision          string
	Version           string
	VersionPrerelease string
}

// GetVersion returns a VersionInfo object
func GetVersion() *Info {
	ver := Version
	rel := VersionPrerelease
	if GitDescribe != "" {
		ver = GitDescribe
	}
	if GitDescribe == "" && rel == "" && VersionPrerelease != "" {
		rel = "dev"
	}

	return &Info{
		Revision:          GitCommit,
		Version:           ver,
		VersionPrerelease: rel,
	}
}

func (c *Info) String() string {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "v%s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)

		if c.Revision != "" {
			fmt.Fprintf(&versionString, " (%s)", c.Revision)
		}
	}

	return versionString.String()
}
