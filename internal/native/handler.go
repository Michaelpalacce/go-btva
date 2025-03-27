package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software/darwin"
	"github.com/Michaelpalacce/go-btva/internal/software/linux"
	"github.com/Michaelpalacce/go-btva/pkg/os"
	"github.com/Michaelpalacce/go-btva/pkg/state"
)

// Handler is a struct that orchestrates the setup process based on OS
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
	case "darwin":
		handler.installer = &darwin.DarwinInstaller{OS: os, Options: options}
	case "windows":
		fallthrough
	default:
		return nil, fmt.Errorf("OS %s is not supported", os.Distro)
	}

	return handler, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Software Block

// SetupSoftware will install all the needed software based on the os and options
func (h *Handler) SetupSoftware() error {
	if h.options.Software.InstallJava {
		if err := h.installSoftware(h.installer.Java()); err != nil {
			return err
		}
	}

	if h.options.Software.InstallMvn {
		if err := h.installSoftware(h.installer.Mvn()); err != nil {
			return err
		}
	}

	if h.options.Software.InstallNode {
		if err := h.installSoftware(h.installer.Node()); err != nil {
			return err
		}
	}

	return nil
}

// END Software Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Local Env Block

// @TODO: Finish
func (h *Handler) SetupLocalEnv() error {
	return nil
}

// END Setup Local Env Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Infra Block
// @TODO: Minimal infra needs to be moved to a strategy of sorts.
func (h *Handler) SetupInfra() error {
	if h.options.Infra.MinimalInfrastructure == false {
		return nil
	}

	if h.state.GetDone(h.infraDone()) {
		slog.Info("Minimal infrastructure already done, skipping...")
		return nil
	}

	slog.Info("Setting up minimal infrastructure on vm", "vmIp", h.options.Infra.SSHVMIP)

	slog.Info("Trying to connect to VM via ssh", "vmIp", h.options.Infra.SSHVMIP)
	client, err := h.getClient()
	if err != nil {
		return err
	}

	defer client.Close()
	slog.Info("Connected to VM via ssh", "vmIp", h.options.Infra.SSHVMIP)

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_CONNECTION),
		state.WithMsg(INFRA_STATE, "Connected to VM"),
	)

	if err := h.runMinimalInfra(client); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	if err := h.fetchGitlabPassword(client); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	h.state.Set(
		state.WithDone(INFRA_STATE, true),
		state.WithMsg(INFRA_STATE, "Finished infra setup"),
		state.WithErr(INFRA_STATE, nil),
	)

	return nil
}

// END Setup Infra Block
