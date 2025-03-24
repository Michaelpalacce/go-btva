package os

import "runtime"

// OS contains what the current OS is as well as different configuration that may be usefull
// WARN: You should not use this directly, but instead use `os.GetOS()`. This struct is provided for mocking or to be used as an
// argument
type OS struct {
	Distro      string
	initialized bool
}

// os has a single instance.
var os *OS = &OS{
	Distro:      runtime.GOOS, // `linux`, `windows` or `darwin` ref: https://pkg.go.dev/runtime#pkg-constants
	initialized: false,
}

// GetOs will return information about the OS, like distro, arch and configuration if needed.
// This will
func GetOS() *OS {
	if !os.initialized {
		initializeOS(os)
	}

	return os
}

func initializeOS(os *OS) {
	// Do more config fetching here if needed
	os.initialized = true
}
