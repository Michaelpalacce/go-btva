package handler

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/components/env"
	"github.com/Michaelpalacce/go-btva/internal/components/final"
	"github.com/Michaelpalacce/go-btva/internal/components/infra"
	"github.com/Michaelpalacce/go-btva/internal/os/darwin"
	"github.com/Michaelpalacce/go-btva/internal/os/linux"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type taskFunc func() error

// Handler is a struct that orchestrates the setup process based on OS
type Handler struct {
	os      *os.OS
	state   *state.State
	options *args.Options

	// WIP
	softwareTasks []func() error
	infraTasks    []func() error
	envTasks      []func() error
	finalTasks    []func() error
	// WIP
}

// WIP
type AddTaskOption func(h *Handler) error

func (h *Handler) AddTask(options ...AddTaskOption) error {
	for _, option := range options {
		if err := option(h); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) RunTasks() {
	for _, task := range h.softwareTasks {
		task()
	}

	for _, task := range h.infraTasks {
		task()
	}

	for _, task := range h.envTasks {
		task()
	}

	for _, task := range h.finalTasks {
		task()
	}
}

// WIP

// NewHandler will return a new Handler that will be used to manage and execute os operations
func NewHandler(os *os.OS, state *state.State, options *args.Options) *Handler {
	return &Handler{os: os, state: state, options: options}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Software Block

// SetupSoftware will install all the needed software based on the os and options
// Each software has it's own state.
func (h *Handler) SetupSoftware() error {
	var installer installer
	switch h.os.Distro {
	case "linux":
		installer = &linux.Installer{OS: h.os, Options: h.options}
	case "darwin":
		installer = &darwin.Installer{OS: h.os, Options: h.options}
	case "windows":
		fallthrough
	default:
		return fmt.Errorf("OS %s is not supported", h.os.Distro)
	}

	for _, software := range installer.GetAllSoftware() {
		if err := h.installSoftware(software); err != nil {
			h.state.Set(withSoftwareInstalled(software, err))
			return err
		}
	}

	return nil
}

// END Software Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Infra Block
func (h *Handler) SetupInfra() error {
	if h.options.Infra.MinimalInfrastructure == false {
		return nil
	}

	infraComponent := infra.NewInfra(h.os, h.state, h.options)

	h.tasks([]taskFunc{
		infraComponent.RunMinimalInfra,
		infraComponent.FetchGitlabPassword,
		infraComponent.CreateGitlabPat,
		infraComponent.GetRunnerAuthToken,
		infraComponent.RegisterGitlabRunner,
		infraComponent.FetchNexusPassword,
	})

	h.state.Set(
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
	finalComponent := final.NewFinal(h.os, h.state, h.options)

	h.tasks([]taskFunc{
		finalComponent.MinimalInfraNexusInstructions,
		finalComponent.MinimalInfraGitlabInstructions,
	})

	return nil
}

// Final Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Local Env Block

func (h *Handler) SetupLocalEnv() error {
	envComponent := env.NewNev(h.os, h.state, h.options)

	h.tasks([]taskFunc{
		envComponent.MinimalInfraSettingsXml,
	})

	return nil
}

// END Setup Local Env Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// tasks will run the given tasks one by one
func (h *Handler) tasks(tasks []taskFunc) error {
	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}
	}

	return nil
}
