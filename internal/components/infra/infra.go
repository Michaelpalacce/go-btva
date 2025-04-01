package infra

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/orchestrator"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type InfraComponent struct {
	os      *os.OS
	state   *state.State
	options *args.Options
}

func NewInfraComponent(os *os.OS, state *state.State) *InfraComponent {
	return &InfraComponent{os: os, state: state, options: state.Options}
}

const (
	INFRA_STATE = "Infra"

	// Public
	INFRA_GITLAB_ADMIN_PASSWORD_KEY = "gitlabPassword"
	INFRA_GITLAB_ADMIN_PAT_KEY      = "gitlabPat"
	INFRA_NEXUS_PASSWORD_KEY        = "nexusPassword"
)

///////////////////////////////////////////////////////////////////////////////////////////////////// Minimal INFRA

func WithFullMinimalInfrastructure() func(*orchestrator.Orchestrator) error {
	return func(h *orchestrator.Orchestrator) error {
		//@TODO: MOVE
		if h.Options.Infra.MinimalInfrastructure == false {
			return nil
		}

		if err := WithPartialMinimalInfrastructureSetup()(h); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureGitlab()(h); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureNexus()(h); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureSettingsXml()(h); err != nil {
			return err
		}

		return nil
	}
}

func WithPartialMinimalInfrastructureGitlab() func(*orchestrator.Orchestrator) error {
	return func(h *orchestrator.Orchestrator) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.InfraTasks = append(h.InfraTasks, []orchestrator.TaskFunc{
			infraComponent.FetchGitlabPassword,
			infraComponent.CreateGitlabPat,
			infraComponent.GetRunnerAuthToken,
			infraComponent.RegisterGitlabRunner,
		}...)

		h.FinalTasks = append(h.FinalTasks, []orchestrator.TaskFunc{
			infraComponent.MinimalInfraGitlabInstructions,
		}...)

		return nil
	}
}

func WithPartialMinimalInfrastructureNexus() func(*orchestrator.Orchestrator) error {
	return func(h *orchestrator.Orchestrator) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.InfraTasks = append(h.InfraTasks, []orchestrator.TaskFunc{
			infraComponent.FetchNexusPassword,
		}...)

		h.FinalTasks = append(h.FinalTasks, []orchestrator.TaskFunc{
			infraComponent.MinimalInfraNexusInstructions,
		}...)

		return nil
	}
}

func WithPartialMinimalInfrastructureSetup() func(*orchestrator.Orchestrator) error {
	return func(h *orchestrator.Orchestrator) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.InfraTasks = append(h.InfraTasks, []orchestrator.TaskFunc{
			infraComponent.RunMinimalInfra,
		}...)

		return nil
	}
}

func WithPartialMinimalInfrastructureSettingsXml() func(*orchestrator.Orchestrator) error {
	return func(h *orchestrator.Orchestrator) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.EnvTasks = append(h.EnvTasks, []orchestrator.TaskFunc{
			infraComponent.MinimalInfraSettingsXml,
		}...)

		return nil
	}
}
