package args

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

// ValidateMinimalInfra runs validation if all the inputs for running minimal infrastructure are set
func (options *Options) ValidateMinimalInfra() error {
	var err error
	if options.Infra.MinimalInfrastructure {
		if options.Infra.SSHVMIP == "" {
			if options.Infra.SSHVMIP, err = prompt.AskText("MinimalInfrastructure selected, but you did not provide sshVmIp, please type in the IP: "); err != nil {
				return fmt.Errorf("sshVmIp must be provided. Err: %w", err)
			}
		}

		if options.Infra.SSHPrivateKey == "" && options.Infra.SSHPassword == "" {
			if options.Infra.SSHPassword, err = prompt.AskPass("MinimalInfrastructure selected, but you did not provide sshPassword or sshPrivateKey, please type in password: "); err != nil {
				return fmt.Errorf("sshPassword must be provided. Err: %w", err)
			}
		}

		if options.Infra.SSHUsername == "" {
			if options.Infra.SSHUsername, err = prompt.AskText("MinimalInfrastructure selected, but you did not provide sshUsername, please type in the username of root or a passwordless sudo user"); err != nil {
				return fmt.Errorf("sshUsername must be provided. Err: %w", err)
			}
		}
	}

	if options.Infra.DockerUsername != "" && options.Infra.DockerPAT == "" {
		if options.Infra.DockerPAT, err = prompt.AskPass("dockerUsername passed, but you did not provide dockerPat, please type in public access token: "); err != nil {
			return fmt.Errorf("dockerPat must be provided with dockerUsername. Err: %w", err)
		}
	}

	return nil
}

// ValidateAriaAutomation will prompt the user a series of question needed to build the aria inventory if the settings are missing
func (options *Options) ValidateAriaAutomation() error {
	var err error
	if options.Aria.Automation.FQDN == "" {
		if options.Aria.Automation.FQDN, err = prompt.AskText(fmt.Sprintf("What is Aria Automation's FQDN without `https://`. Current (%s)", options.Aria.Automation.FQDN)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.Port == "" {
		if options.Aria.Automation.Port, err = prompt.AskText(fmt.Sprintf("What is Aria Automation's port? Current (%s)", options.Aria.Automation.Port)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.Username == "" {
		if options.Aria.Automation.Username, err = prompt.AskText(fmt.Sprintf("What is the username of the account for Aria Automation? Current (%s)", options.Aria.Automation.Username)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.Password == "" {
		if options.Aria.Automation.Password, err = prompt.AskPass("What is the password of the account for Aria Automation?"); err != nil {
			return err
		}
	}

	if options.Aria.Automation.OrgName == "" {
		if options.Aria.Automation.OrgName, err = prompt.AskText(fmt.Sprintf("What is the org name used in Aria Automation? Current (%s)", options.Aria.Automation.OrgName)); err != nil {
			return err
		}
	}

	if options.Aria.Automation.ProjectName == "" {
		if options.Aria.Automation.ProjectName, err = prompt.AskText(fmt.Sprintf("What is the default project name in Aria Automation you want to push automation code to? Current (%s)", options.Aria.Automation.ProjectName)); err != nil {
			return err
		}
	}

	return nil
}
