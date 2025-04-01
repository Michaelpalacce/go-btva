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
