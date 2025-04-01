package infra

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/handler"
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

///////////////////////////////////////////////////////////////////////////////////////////////////// Minimal INFRA

func WithFullMinimalInfrastructure() func(*handler.Handler) error {
	return func(h *handler.Handler) error {
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

func WithPartialMinimalInfrastructureGitlab() func(*handler.Handler) error {
	return func(h *handler.Handler) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.InfraTasks = append(h.InfraTasks, []handler.TaskFunc{
			infraComponent.FetchGitlabPassword,
			infraComponent.CreateGitlabPat,
			infraComponent.GetRunnerAuthToken,
			infraComponent.RegisterGitlabRunner,
		}...)

		h.FinalTasks = append(h.FinalTasks, []handler.TaskFunc{
			infraComponent.MinimalInfraGitlabInstructions,
		}...)

		return nil
	}
}

func WithPartialMinimalInfrastructureNexus() func(*handler.Handler) error {
	return func(h *handler.Handler) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.InfraTasks = append(h.InfraTasks, []handler.TaskFunc{
			infraComponent.FetchNexusPassword,
		}...)

		h.FinalTasks = append(h.FinalTasks, []handler.TaskFunc{
			infraComponent.MinimalInfraNexusInstructions,
		}...)

		return nil
	}
}

func WithPartialMinimalInfrastructureSetup() func(*handler.Handler) error {
	return func(h *handler.Handler) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.InfraTasks = append(h.InfraTasks, []handler.TaskFunc{
			infraComponent.RunMinimalInfra,
		}...)

		return nil
	}
}

func WithPartialMinimalInfrastructureSettingsXml() func(*handler.Handler) error {
	return func(h *handler.Handler) error {
		infraComponent := NewInfraComponent(h.OS, h.State)

		h.EnvTasks = append(h.EnvTasks, []handler.TaskFunc{
			infraComponent.MinimalInfraSettingsXml,
		}...)

		return nil
	}
}
