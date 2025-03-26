package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/pkg/state"
	"github.com/melbahja/goph"
)

const (
	MINIMAL_INFRA_STEP_CONNECTION = iota + 1
	MINIMAL_INFRA_STEP_SETUP
	MINIMAL_INFRA_STEP_UP
	MINIMAL_INFRA_STEP_INFO_FETCHED
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
	if h.state.GetStep(h.getMinimalInfraStep()) >= MINIMAL_INFRA_STEP_SETUP {
		slog.Info("Skipping minimal infrastructure installer, step already done.")
		return nil
	}

	h.state.Set(state.WithMsg(MINIMAL_INFRA_STATE, "Running the minimal infrastructure installer."))
	slog.Info("Running the minimal infrastructure installer.")

	out, err := client.Run("curl -o- https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/refactor/minimal-infra-simplified-setup/infrastructure/install.sh | bash")
	if err != nil {
		return fmt.Errorf("minimal infrastructure installer exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	slog.Info("Minimal infrastructure installer successfully set up.")

	h.state.Set(state.WithStep(MINIMAL_INFRA_STATE, MINIMAL_INFRA_STEP_SETUP))

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

func (h *Handler) getMinimalInfraStep() state.GetStepStateOption {
	return func(s *state.State) int {
		value := s.GetValue(MINIMAL_INFRA_STATE)
		if value == nil {
			return 0
		}

		return value.Step
	}
}
