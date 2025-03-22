package native

import (
	"fmt"

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
	installer software.Installer
}

// NewHandler will return a new native Handler that will be used to manage and execute os operations
func NewHandler(os *os.OS, options *args.Options) (*Handler, error) {
	handler := &Handler{os: os, state: state.NewState(), options: options}

	switch os.Distro {
	case "linux":
		handler.installer = &linux.LinuxInstaller{OS: os, Options: options}
	case "windows":
	case "darwin":
	default:
		return nil, fmt.Errorf("OS %s is not supported", os.Distro)
	}

	return handler, nil
}

func (h *Handler) SetupSoftware(c chan error) {
	if h.state.GetDone(software.IsJavaInstalled(false)) && h.options.Software.InstallJava {
		err := h.installer.InstallJava()
		if err != nil {
			h.state.Set(software.JavaInstalled(err))
			c <- err

			return
		}
	}
}
func (h *Handler) SetupLocalEnv(c chan error) {}
func (h *Handler) SetupInfra(c chan error)    {}
