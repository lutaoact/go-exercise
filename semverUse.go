package main

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/golang/glog"
)

func main() {
	fmt.Println(isOldVersion("8.2.0"))
}

// 如果版本解析出问题，则都认为是旧版本
func isOldVersion(version string) bool {
	c, _ := semver.NewConstraint("< 8.3.0")
	v, err := semver.NewVersion(version)
	if err != nil {
		glog.Errorf("semver.NewVersion: %+v", err)
		return true
	}

	ok, msgs := c.Validate(v)
	for _, m := range msgs {
		fmt.Println(m)
	}
	return ok
}
