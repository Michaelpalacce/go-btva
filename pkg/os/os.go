package os

import (
	"fmt"
	"runtime"

	osz "os"
)

// OS contains what the current OS is as well as different configuration that may be usefull
// WARN: You should not use this directly, but instead use `os.GetOS()`. This struct is provided for mocking or to be used as an
// argument
type OS struct {
	Distro      string
	initialized bool

	Shell        string
	ShellProfile string
}

// os has a single instance.
var os *OS = &OS{
	Distro:      runtime.GOOS, // `linux`, `windows` or `darwin` ref: https://pkg.go.dev/runtime#pkg-constants
	initialized: false,
}

// GetOs will return information about the OS, like distro, arch and configuration if needed.
// This will panic if any errors are hit
func GetOS() *OS {
	if !os.initialized {
		if err := initializeOS(); err != nil {
			panic(err)
		}
	}

	return os
}

// initializeOS will initialize the OS with data from the current OS
func initializeOS() error {
	if err := determineShell(); err != nil {
		return fmt.Errorf("error while trying to determine shell: %w", err)
	}

	os.initialized = true

	return nil
}

// determineShell works for `linux` and `darwin` and fetches the active shell and it's profile file
// Profile file is the file that is read on shell start
func determineShell() error {
	switch os.Distro {
	case "darwin":
		fallthrough
	case "linux":
		shell := osz.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}

		var profile string
		switch shell {
		case "/bin/zsh":
			profile = "$HOME/.zshrc"
		case "/bin/fish":
			profile = "$HOME/.config/fish/config.fish"
		case "/bin/bash":
			profile = "$HOME/.bashrc"
		default:
			return fmt.Errorf("shell %s is not supported", shell)
		}

		os.Shell = shell
		os.ShellProfile = profile
	}

	return nil
}
