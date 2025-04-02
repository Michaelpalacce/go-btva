package orchestrator

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////////////////////////// Minimal INFRA

// WithFullMinimalInfrastructure will setup the entire minimal infrastructure stack and configure your env accordingly
func WithFullMinimalInfrastructure() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		if err := WithPartialMinimalInfrastructureSetup()(o); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureGitlab()(o); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureNexus()(o); err != nil {
			return err
		}

		if err := WithPartialMinimalInfrastructureSettingsXml()(o); err != nil {
			return err
		}

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will setup Gitlab using the minimal infra installer. Make sure to have added
// WithPartialMinimalInfrastructureSetup first
// You don't need to pass infraComponent, it will be created
func WithPartialMinimalInfrastructureGitlab() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

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
func WithPartialMinimalInfrastructureNexus() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

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
// Since this does valdation, it will flush the state to storage. It's done in a go routine, you don't need to wait for it
func WithPartialMinimalInfrastructureSetup() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

		if err := o.Options.ValidateMinimalInfra(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		go o.State.Flush()

		o.InfraTasks = append(o.InfraTasks, []TaskFunc{
			infraComponent.RunMinimalInfra,
		}...)

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will run replace your settings.xml file with the minimal infra settings
// You don't need to pass infraComponent, it will be created
func WithPartialMinimalInfrastructureSettingsXml() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

		if err := o.Options.ValidateAriaAutomation(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		o.EnvTasks = append(o.EnvTasks, []TaskFunc{
			infraComponent.MinimalInfraSettingsXml,
			infraComponent.InfraSettingsXml,
		}...)

		return nil
	}
}
