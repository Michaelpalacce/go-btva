package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/os/darwin"
	"github.com/Michaelpalacce/go-btva/internal/os/linux"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
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
func NewHandler(os *os.OS, state *state.State, options *args.Options) (*Handler, error) {
	handler := &Handler{os: os, state: state, options: options}

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
// Each software has it's own state.
func (h *Handler) SetupSoftware() error {
	for _, software := range h.installer.GetAllSoftware() {
		if err := h.installSoftware(software); err != nil {
			return err
		}
	}

	return nil
}

// END Software Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Local Env Block

func (h *Handler) SetupLocalEnv() error {
	if state.Get(h.state, envDone()) {
		slog.Info("Environment setup already done, skipping...")
		return nil
	}

	if h.options.Local.SetupM2 {
		if err := h.prepareSettingsXml(h.os, h.options, h.state); err != nil {
			return err
		}
	}

	h.state.Set(
		state.WithDone(ENV_STATE, true),
		state.WithMsg(ENV_STATE, "Finished environment setup."),
		state.WithErr(ENV_STATE, nil),
	)

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

	// if infraDone(h.state) {
	// 	slog.Info("Minimal infrastructure already done, skipping...")
	// 	return nil
	// }

	slog.Info("Setting up minimal infrastructure on vm", "vmIp", h.options.Infra.SSHVMIP)

	slog.Info("Trying to connect to VM via ssh", "vmIp", h.options.Infra.SSHVMIP)
	client, err := h.getClient()
	if err != nil {
		return err
	}

	defer client.Close()

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_CONNECTION),
		state.WithMsg(INFRA_STATE, fmt.Sprintf("Connected to VM (%s) via ssh", h.options.Infra.SSHVMIP)),
	)

	if err := h.runMinimalInfra(client); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	if err := h.fetchGitlabPassword(client); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	if err := h.createGitlabPat(client); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	if err := h.getRunnerAuthToken(); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	if err := h.registerGitlabRunner(client); err != nil {
		h.state.Set(state.WithErr(INFRA_STATE, err))
		return err
	}

	if err := h.fetchNexusPassword(client); err != nil {
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

// Final Block

// Final will print out some instructions to the user
// If it was done already, it won't log anything
func (h *Handler) Final() error {
	// if state.Get(h.state, finalDone()) == true {
	// 	return nil
	// }

	if err := h.NexusInstructions(); err != nil {
		h.state.Set(state.WithErr(FINAL_STATE, err))
		return err
	}

	if err := h.GitlabInstructions(); err != nil {
		h.state.Set(state.WithErr(FINAL_STATE, err))
		return err
	}

	h.state.Set(
		state.WithDone(FINAL_STATE, true),
		state.WithMsg(FINAL_STATE, "Finished entire setup"),
		state.WithErr(FINAL_STATE, nil),
	)

	return nil
}
