package native

import (
	"fmt"
	"strings"

	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/gitlab"
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
	"github.com/melbahja/goph"
)

const (
	INFRA_STEP_CONNECTION = iota + 1
	INFRA_STEP_SETUP
	INFRA_STEP_INFO_FETCHED_GITLAB
	INFRA_STEP_GITLAB_PAT_CREATED
	INFRA_STEP_GITLAB_RUNNER_AUTH_TOKEN
	INFRA_STEP_GITLAB_RUNNER_REGISTERED
	INFRA_STEP_INFO_FETCHED_NEXUS
)

const (
	INFRA_STATE = "MinimalInfra"

	INFRA_GITLAB_ADMIN_PASSWORD_KEY    = "gitlabPassword"
	INFRA_GITLAB_ADMIN_PAT_KEY         = "gitlabPat"
	INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY = "gitlabRunnerAuthToken"
	INFRA_NEXUS_PASSWORD_KEY           = "nexusPassword"
)

const BTVA_INSTALL_DIR_INFRA = "/opt/build-tools-for-vmware-aria/infrastructure"

// getClient will retrieve a client, using the Infra options
func (h *Handler) getClient() (*goph.Client, error) {
	infraOptions := h.options.Infra
	return ssh.GetClient(infraOptions.SSHVMIP, infraOptions.SSHUsername, infraOptions.SSHPassword, infraOptions.SSHPrivateKey, infraOptions.SSHPrivateKeyPassphrase)
}

// runMinimalInfra will fetch the BTVA minimal infra installer and run it
// @TODO: Fix the branch
func (h *Handler) runMinimalInfra(client *goph.Client) error {
	if infraStep(h.state) >= INFRA_STEP_SETUP {
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Running the minimal infrastructure installer. This may take a few minutes as it waits for services to be healthy."))

	out, err := client.Run(fmt.Sprintf("curl -o- https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/refactor/minimal-infra-simplified-setup/infrastructure/install.sh | bash -s -- %s %q", h.options.Infra.DockerUsername, h.options.Infra.DockerPAT))
	if err != nil {
		return fmt.Errorf("minimal infrastructure installer exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	url := fmt.Sprintf("http://%s:8082/gitlab", h.options.Infra.SSHVMIP)

	if out, err := client.Run(fmt.Sprintf("sed -i \"s|external_url 'http://infra.corp.local/gitlab'|external_url '%q'|\" %s/docker-compose.yml", url, BTVA_INSTALL_DIR_INFRA)); err != nil {
		return fmt.Errorf("failed to modify compose file. err was %w, output was:\n%s", err, out)
	}

	if out, err := client.Run(fmt.Sprintf("docker compose -f %s/docker-compose.yml up -d --wait", BTVA_INSTALL_DIR_INFRA)); err != nil {
		return fmt.Errorf("failed to start containers. err was %w, output was:\n%s", err, out)
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_SETUP),
		state.WithMsg(INFRA_STATE, "Minimal infrastructure installer successfully set up."),
	)

	return nil
}

// fetchGitlabPassword will fetch the password for Gitlab and store it in the context store
// Command looks a bit big, but it's all so we can fail in case the file doesn't exists or the container is not started
func (h *Handler) fetchGitlabPassword(client *goph.Client) error {
	if infraStep(h.state) >= INFRA_STEP_INFO_FETCHED_GITLAB {
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Fetching gitlab admin password"))

	out, err := client.Run("docker exec gitlab-ce test -f /etc/gitlab/initial_root_password && docker exec gitlab-ce grep 'Password:' /etc/gitlab/initial_root_password | awk '{print $2}'")
	if err != nil {
		return fmt.Errorf("gitlab admin password fetching exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_INFO_FETCHED_GITLAB),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Gitlab admin password fetched successfully."),
	)

	return nil
}

// createGitlabPat with the help of ruby on the gitlab container will generate a new Public Access Token
func (h *Handler) createGitlabPat(client *goph.Client) error {
	if infraStep(h.state) >= INFRA_STEP_GITLAB_PAT_CREATED {
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Creating a new Gitlab Public access token."))

	gitlabPassword := state.Get(h.state, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY))
	if gitlabPassword == "" {
		return fmt.Errorf("gitlab password is an empty string. Was it deleted? Rerunning the infra may help.")
	}
	gitlabPat := gitlabPassword[:20]

	out, err := client.Run(fmt.Sprintf("docker exec gitlab-ce gitlab-rails runner 'token = User.find_by_username(\"root\").personal_access_tokens.create(scopes: [:read_user, :read_repository, :api, :create_runner, :manage_runner, :sudo, :admin_mode], name: \"Automation token\", expires_at: 356.days.from_now); token.set_token(\"%s\"); token.save! '", gitlabPat))
	if err != nil {
		return fmt.Errorf("gitlab admin public access token creation exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_GITLAB_PAT_CREATED),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PAT_KEY, gitlabPat),
		state.WithMsg(INFRA_STATE, "Gitlab admin public access token created successfully."),
	)

	return nil
}

// getRunnerAuthToken will fetch an auth token that can be used to register a new gitlab runner
func (h *Handler) getRunnerAuthToken() error {
	if infraStep(h.state) >= INFRA_STEP_GITLAB_RUNNER_AUTH_TOKEN {
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Creating a auth token for the gitlab runner."))

	gitlabPat := state.Get(h.state, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PAT_KEY))
	if gitlabPat == "" {
		return fmt.Errorf("gitlab pat is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	client := gitlab.NewGitlabClient(fmt.Sprintf("http://%s:8082/gitlab", h.options.Infra.SSHVMIP), gitlabPat)

	token, err := client.GetRunnerAuthToken("instance_type")
	if err != nil {
		return fmt.Errorf("error while trying to fetch runner auth token. Err was: %s")
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_GITLAB_RUNNER_AUTH_TOKEN),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY, token),
		state.WithMsg(INFRA_STATE, "Gitlab auth token successfully created."),
	)

	return nil
}

// getRunnerAuthToken will fetch an auth token that can be used to register a new gitlab runner
func (h *Handler) registerGitlabRunner(client *goph.Client) error {
	if infraStep(h.state) >= INFRA_STEP_GITLAB_RUNNER_REGISTERED {
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Registering gitlab runner with generated auth token"))

	runnerAuthToken := state.Get(h.state, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY))
	if runnerAuthToken == "" {
		return fmt.Errorf("runner auth token is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	url := fmt.Sprintf("http://%s:8082/gitlab", h.options.Infra.SSHVMIP)

	out, err := client.Run(fmt.Sprintf("docker exec gitlab-runner gitlab-runner register --non-interactive --url \"%s\" --token \"%s\" --executor \"docker\" --docker-image alpine:latest --description \"docker-runner\"", url, runnerAuthToken))
	if err != nil {
		return fmt.Errorf("registering a gitlab runner exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_GITLAB_RUNNER_REGISTERED),
		state.WithMsg(INFRA_STATE, "Gitlab runner registered successfully."),
	)

	return nil
}

// fetchNexusPassword will fetch the password for Nexus and store it in the context store
func (h *Handler) fetchNexusPassword(client *goph.Client) error {
	if infraStep(h.state) >= INFRA_STEP_INFO_FETCHED_NEXUS {
		return nil
	}

	h.state.Set(state.WithMsg(INFRA_STATE, "Fetching nexus admin password"))

	out, err := client.Run("docker exec nexus cat /nexus-data/admin.password")
	if err != nil {
		if isNoSuchFileOrDirectoryErr(string(out)) {
			pass, err := prompt.AskPass("In order to continue execution, please provide nexus password manually:")
			if err != nil {
				return fmt.Errorf("error while providing nexus password. err was %w", err, out)
			}

			out = []byte(pass)
		} else {
			return fmt.Errorf("nexus admin password fetching exited unsuccessfully. err was %w, output was:\n%s", err, out)
		}
	}

	h.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_INFO_FETCHED_NEXUS),
		state.WithContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Nexus admin password fetched successfully."),
	)

	return nil
}

// infraDone will give us a state.GetSuccessStateOption that will check if the minimal infra is done
func infraDone(s *state.State) bool {
	return state.Get(s, state.GetDone(INFRA_STATE))
}

// infraStep gets the current step for the infra setup that we are on
func infraStep(s *state.State) int {
	return state.Get(s, state.GetStep(INFRA_STATE))
}

func gitlabAdminPassword(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY))
}

func nexusAdminPassword(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY))
}

func isNoSuchFileOrDirectoryErr(msg string) bool {
	return strings.Contains(msg, "No such file or directory")
}
