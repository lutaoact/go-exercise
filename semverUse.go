package main

import (
	"fmt"

	"github.com/coreos/go-semver/semver"
)

const (
	VersionMajor int64 = 1 << iota
	VersionMinor
	VersionPatch
	VersionPre string = "xxx"
	VersionDev string = "yyy"
)

// Version is the specification version that the package types support.
var Version = semver.Version{
	Major:      VersionMajor,
	Minor:      VersionMinor,
	Patch:      VersionPatch,
	PreRelease: semver.PreRelease(VersionPre),
	Metadata:   VersionDev,
}

func main() {
	fmt.Println(Version.String())
}
