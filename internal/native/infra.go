package native

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/pkg/state"
	"github.com/melbahja/goph"
)

const (
	INFRA_STEP_CONNECTION = iota + 1
	INFRA_STEP_SETUP
	INFRA_STEP_INFO_FETCHED
)

const (
	INFRA_STATE               = "MinimalInfra"
	INFRA_GITLAB_PASSWORD_KEY = "gitlabPassword"
)

// getClient will retrieve a client, using the Infra options
func (h *Handler) getClient() (*goph.Client, error) {
	infraOptions := h.options.Infra
	return ssh.GetClient(infraOptions.SSHVMIP, infraOptions.SSHUsername, infraOptions.SSHPassword, infraOptions.SSHPrivateKey, infraOptions.SSHPrivateKeyPassphrase)
}

// runMinimalInfra will fetch the BTVA minimal infra installer and run it
// @TODO: Fix the branch
func (h *Handler) runMinimalInfra(client *goph.Client) error {
	if h.state.GetStep(h.getMinimalInfraStep()) >= INFRA_STEP_SETUP {
		slog.Info("Skipping minimal infrastructure installer, step already done.")
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Running the minimal infrastructure installer. This may take a few minutes as it waits for services to be healthy."))

	out, err := client.Run("curl -o- https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/refactor/minimal-infra-simplified-setup/infrastructure/install.sh | bash")
	if err != nil {
		return fmt.Errorf("minimal infrastructure installer exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_SETUP),
		state.WithMsg(INFRA_STATE, "Minimal infrastructure installer successfully set up."),
	)

	return nil
}

// fetchGitlabPassword will fetch the password for Gitlab and store it in the context store
func (h *Handler) fetchGitlabPassword(client *goph.Client) error {
	if h.state.GetStep(h.getMinimalInfraStep()) >= INFRA_STEP_INFO_FETCHED {
		slog.Info("Skipping password fetching, step already done.")
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Fetching gitlab admin password"))

	out, err := client.Run("docker exec gitlab-ce grep 'Password:' /etc/gitlab/initial_root_password | awk '{print $2}'")
	if err != nil {
		return fmt.Errorf("gitlab admin password fetching exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_INFO_FETCHED),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Gitlab admin password fetched successfully."),
	)

	return nil
}

// isMinimalInfraDone will give us a state.GetSuccessStateOption that will check if the minimal infra is done
func (h *Handler) isMinimalInfraDone() state.GetSuccessStateOption {
	return state.GetDone(INFRA_STATE)
}

func (h *Handler) getMinimalInfraStep() state.GetStepStateOption {
	return state.GetStep(INFRA_STATE)
}

func (h *Handler) getGitlabAdminPassword() state.GetContextPropStateOption {
	return state.GetContextProp(INFRA_STATE, INFRA_GITLAB_PASSWORD_KEY)
}
