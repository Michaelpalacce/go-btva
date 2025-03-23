package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/internal/software/linux"
	"github.com/Michaelpalacce/go-btva/pkg/os"
	"github.com/Michaelpalacce/go-btva/pkg/state"
)

// Handler is a struct that can
type Handler struct {
	os      *os.OS
	state   *state.State
	options *args.Options

	// installer is a pointer
	installer Installer
}

// NewHandler will return a new native Handler that will be used to manage and execute os operations
func NewHandler(os *os.OS, options *args.Options) (*Handler, error) {
	handler := &Handler{os: os, state: state.NewState(), options: options}

	if options.Local.SaveState {
		handler.state.Set(state.WithJsonStorage(options.Local.StateJson, true))
	}

	switch os.Distro {
	case "linux":
		handler.installer = &linux.LinuxInstaller{OS: os, Options: options}
	case "windows":
		fallthrough
	case "darwin":
		fallthrough
	default:
		return nil, fmt.Errorf("OS %s is not supported", os.Distro)
	}

	return handler, nil
}

// SetupSoftware will install all the needed software based on the os and options
// @NOTE: This is meant to be ran async
func (h *Handler) SetupSoftware(c chan error) {
	if h.state.GetDone(software.IsJavaInstalled(false)) && h.options.Software.InstallJava {
		slog.Info("Java is not installed, installing")
		err := h.installer.InstallJava()
		if err != nil {
			if err := h.state.Set(software.JavaInstalled(err)); err != nil {
				slog.Error("Error setting state", err)
			}
			c <- err

			return
		}

		if err := h.state.Set(software.JavaInstalled(nil)); err != nil {
			slog.Error("Error setting state", err)
		}

		slog.Info("Successfully installed Java")
	} else {
		slog.Info("Java is already installed, skipping...")
	}

	c <- nil
}

// @TODO: Finish
func (h *Handler) SetupLocalEnv(c chan error) {
	c <- nil
}

// @TODO: Finish
func (h *Handler) SetupInfra(c chan error) {
	c <- nil
}
