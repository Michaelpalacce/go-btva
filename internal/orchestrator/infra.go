package orchestrator

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/internal/components/infra"
)

///////////////////////////////////////////////////////////////////////////////////////////////////// Minimal INFRA

// WithFullMinimalInfrastructure will setup the entire minimal infrastructure stack and configure your env accordingly
func WithFullMinimalInfrastructure() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := infra.NewInfraComponent(o.OS, o.State)

		if err := WithPartialMinimalInfrastructureSetup(infraComponent)(o); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureGitlab(infraComponent)(o); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureNexus(infraComponent)(o); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureSettingsXml(infraComponent)(o); err != nil {
			return err
		}

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will setup Gitlab using the minimal infra installer. Make sure to have added
// WithPartialMinimalInfrastructureSetup first
// You don't need to pass infraComponent, it will be created
func WithPartialMinimalInfrastructureGitlab(infraComponent *infra.InfraComponent) func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		if infraComponent == nil {
			infraComponent = infra.NewInfraComponent(o.OS, o.State)
		}

		o.InfraTasks = append(o.InfraTasks, []TaskFunc{
			infraComponent.FetchGitlabPassword,
			infraComponent.CreateGitlabPat,
			infraComponent.GetRunnerAuthToken,
			infraComponent.RegisterGitlabRunner,
		}...)

		o.FinalTasks = append(o.FinalTasks, []TaskFunc{
			infraComponent.MinimalInfraGitlabInstructions,
		}...)

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will setup Nexus using the minimal infra installer. Make sure to have added
// WithPartialMinimalInfrastructureSetup first
// You don't need to pass infraComponent, it will be created
func WithPartialMinimalInfrastructureNexus(infraComponent *infra.InfraComponent) func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		if infraComponent == nil {
			infraComponent = infra.NewInfraComponent(o.OS, o.State)
		}

		o.InfraTasks = append(o.InfraTasks, []TaskFunc{
			infraComponent.FetchNexusPassword,
		}...)

		o.FinalTasks = append(o.FinalTasks, []TaskFunc{
			infraComponent.MinimalInfraNexusInstructions,
		}...)

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will run the minimal infra installer
// You don't need to pass infraComponent, it will be created
func WithPartialMinimalInfrastructureSetup(infraComponent *infra.InfraComponent) func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		if infraComponent == nil {
			infraComponent = infra.NewInfraComponent(o.OS, o.State)
		}

		if err := o.Options.ValidateMinimalInfra(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		o.InfraTasks = append(o.InfraTasks, []TaskFunc{
			infraComponent.RunMinimalInfra,
		}...)

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will run replace your settings.xml file with the minimal infra settings
// You don't need to pass infraComponent, it will be created
func WithPartialMinimalInfrastructureSettingsXml(infraComponent *infra.InfraComponent) func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		if infraComponent == nil {
			infraComponent = infra.NewInfraComponent(o.OS, o.State)
		}

		o.EnvTasks = append(o.EnvTasks, []TaskFunc{
			infraComponent.MinimalInfraSettingsXml,
		}...)

		return nil
	}
}

