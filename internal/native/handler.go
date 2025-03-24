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
		if err := handler.state.Modify(state.WithJsonStorage(options.Local.StateJson, true)); err != nil {
			return nil, err
		}
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
	//                  in state                           is wanted in arguments            already installed
	if h.state.GetDone(software.IsJavaInstalled(false)) && h.options.Software.InstallJava && !h.installer.Java().Exists() {
		slog.Info("Java is not installed, installing")

		javaSoftware := h.installer.Java()
		err := javaSoftware.Install()
		if err != nil {
			if err := h.state.Set(software.SoftwareInstalled(javaSoftware, err)); err != nil {
				slog.Error("Error setting state", err)
			}
			c <- err

			return
		}

		if err := h.state.Set(software.SoftwareInstalled(javaSoftware, nil)); err != nil {
			slog.Error("Error setting state", err)
		}

		slog.Info("Successfully installed Java")
	} else {
		slog.Info("Java is already installed, skipping...")
	}

	//                  in state                           is wanted in arguments            already installed
	if h.state.GetDone(software.IsMvnInstalled(false)) && h.options.Software.InstallMvn && !h.installer.Mvn().Exists() {
		slog.Info("Maven is not installed, installing")

		mvnSoftware := h.installer.Mvn()
		err := mvnSoftware.Install()
		if err != nil {
			if err := h.state.Set(software.SoftwareInstalled(mvnSoftware, err)); err != nil {
				slog.Error("Error setting state", err)
			}
			c <- err

			return
		}

		if err := h.state.Set(software.SoftwareInstalled(mvnSoftware, nil)); err != nil {
			slog.Error("Error setting state", err)
		}

		slog.Info("Successfully installed Maven1")
	} else {
		slog.Info("Maven is already installed, skipping...")
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
