package outage

import (
	"fmt"
	"os"
)

var (
	// the project name
	PROJ_NAME = "outage"
)

var (
	// the version meta
	MAJOR = 0
	MINOR = 1
	MACRO = 0
)

type Version bool

func (Version) String() (ver string) {
	ver = fmt.Sprintf("%v (v%d.%d.%d)", PROJ_NAME, MAJOR, MINOR, MACRO)
	return
}

func (ver Version) BeforeApply() (err error) {
	fmt.Println(ver)
	os.Exit(0)
	return
}
