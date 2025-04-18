package orchestrator

import (
	"fmt"
)

///////////////////////////////////////////////////////////////////////////////////////////////////// Infra

// WithPartialMinimalInfrastructureGitlab will run replace your settings.xml file with the minimal infra settings
// You don't need to pass infraComponent, it will be created
func WithSettingsXml() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

		if err := o.State.Options.ValidateAriaAutomation(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		if err := o.State.Options.ValidateArtifactManagerArguments(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		o.EnvTasks = append(o.EnvTasks, []TaskFunc{
			infraComponent.InfraSettingsXml,
		}...)

		return nil
	}
}

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
// Since this does valdation, it will flush the state to storage.
func WithPartialMinimalInfrastructureSetup() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

		if err := o.State.Options.ValidateMinimalInfra(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		o.State.Flush()

		o.InfraTasks = append(o.InfraTasks, []TaskFunc{
			infraComponent.RunMinimalInfra,
		}...)

		return nil
	}
}

// WithPartialMinimalInfrastructureGitlab will run replace your settings.xml file with the minimal infra settings
func WithPartialMinimalInfrastructureSettingsXml() func(*Orchestrator) error {
	return func(o *Orchestrator) error {
		infraComponent := o.components.infraComponent

		if err := o.State.Options.ValidateAriaAutomation(); err != nil {
			return fmt.Errorf("error trying to validate passed options. Err was: %w", err)
		}

		o.EnvTasks = append(o.EnvTasks, []TaskFunc{
			infraComponent.MinimalInfraSettingsXml,
		}...)

		return nil
	}
}
