package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/pkg/state"
	"github.com/melbahja/goph"
)

const MINIMAL_INFRA_STATE = "MinimalInfra"

// getClient will retrieve a client, using the Infra options
func (h *Handler) getClient() (*goph.Client, error) {
	infraOptions := h.options.Infra
	return ssh.GetClient(infraOptions.SSHVMIP, infraOptions.SSHUsername, infraOptions.SSHPassword, infraOptions.SSHPrivateKey, infraOptions.SSHPrivateKeyPassphrase)
}

// runMinimalInfra will fetch the BTVA minimal infra installer and run it
// @TODO: Fix the branch
func (h *Handler) runMinimalInfra(client *goph.Client) error {
	h.state.Set(
		state.WithStep(MINIMAL_INFRA_STATE, 2),
		state.WithMsg(MINIMAL_INFRA_STATE, "Running the minimal infrastructure installer."),
	)
	slog.Info("Running the minimal infrastructure installer.")

	out, err := client.Run("curl -o- https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/refactor/minimal-infra-simplified-setup/infrastructure/install.sh | bash")
	if err != nil {
		return fmt.Errorf("minimal infrastructure installer exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	slog.Info("Minimal infrastructure installer successfully set up.")

	return nil
}

// isMinimalInfraDone will give us a state.GetSuccessStateOption that will check if the minimal infra is done
func (h *Handler) isMinimalInfraDone() state.GetSuccessStateOption {
	return func(s *state.State) bool {
		value := s.GetValue(MINIMAL_INFRA_STATE)
		if value == nil {
			return false
		}

		return value.Done
	}
}
