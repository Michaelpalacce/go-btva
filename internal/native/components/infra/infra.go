package infra

import (
	"fmt"
	"strings"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/gitlab"
	"github.com/Michaelpalacce/go-btva/pkg/os"
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

const (
	BTVA_INSTALL_DIR_INFRA         = "/opt/build-tools-for-vmware-aria/infrastructure"
	BTVA_MINIMAL_INFRA_INSTALL_URL = "https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/refactor/minimal-infra-simplified-setup/infrastructure/install.sh"
)

type Infra struct {
	os      *os.OS
	state   *state.State
	options *args.Options
	client  *goph.Client
}

func NewInfra(os *os.OS, state *state.State, options *args.Options, client *goph.Client) *Infra {
	return &Infra{os: os, state: state, options: options, client: client}
}

// GetClient will ssh into the machine and give you a goph.Client pointer you can use to run commands.
// @WARN: Make sure to defer client.Close()
// @NOTE: There is probably a beter place for this
func GetClient(options *args.Options, s *state.State) (*goph.Client, error) {
	infraOptions := options.Infra

	client, err := ssh.GetClient(infraOptions.SSHVMIP, infraOptions.SSHUsername, infraOptions.SSHPassword, infraOptions.SSHPrivateKey, infraOptions.SSHPrivateKeyPassphrase)
	if err != nil {
		return nil, fmt.Errorf("could not create a new ssh client and/or ssh into machine. Err was: %w", err)
	}

	s.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_CONNECTION),
		state.WithMsg(INFRA_STATE, fmt.Sprintf("Connected to VM (%s) via ssh", options.Infra.SSHVMIP)),
	)

	return client, nil
}

// runMinimalInfra will fetch the BTVA minimal infra installer and run it
// @TODO: Fix the branch
func (i *Infra) RunMinimalInfra() error {
	if infraStep(i.state) >= INFRA_STEP_SETUP {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Running the minimal infrastructure installer. This may take a few minutes as it waits for services to be healthy."))

	out, err := i.client.Run(fmt.Sprintf("curl -o- %s | bash -s -- %s %q", BTVA_MINIMAL_INFRA_INSTALL_URL, i.options.Infra.DockerUsername, i.options.Infra.DockerPAT))
	if err != nil {
		return fmt.Errorf("minimal infrastructure installer exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	if out, err := i.client.Run(fmt.Sprintf("sed -i \"s|external_url 'http://infra.corp.local/gitlab'|external_url '%q'|\" %s/docker-compose.yml", gitlabUrl(*i.options), BTVA_INSTALL_DIR_INFRA)); err != nil {
		return fmt.Errorf("failed to modify compose file. err was %w, output was:\n%s", err, out)
	}

	if out, err := i.client.Run(fmt.Sprintf("docker compose -f %s/docker-compose.yml up -d --wait", BTVA_INSTALL_DIR_INFRA)); err != nil {
		return fmt.Errorf("failed to start containers. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_SETUP),
		state.WithMsg(INFRA_STATE, "Minimal infrastructure installer successfully set up."),
	)

	return nil
}

// fetchGitlabPassword will fetch the password for Gitlab and store it in the context store
// Command looks a bit big, but it's all so we can fail in case the file doesn't exists or the container is not started
func (i *Infra) FetchGitlabPassword() error {
	if infraStep(i.state) >= INFRA_STEP_INFO_FETCHED_GITLAB && GitlabAdminPassword(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Fetching gitlab admin password"))

	out, err := i.client.Run("docker exec gitlab-ce test -f /etc/gitlab/initial_root_password && docker exec gitlab-ce grep 'Password:' /etc/gitlab/initial_root_password | awk '{print $2}'")
	if err != nil {
		return fmt.Errorf("gitlab admin password fetching exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_INFO_FETCHED_GITLAB),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Gitlab admin password fetched successfully."),
	)

	return nil
}

// createGitlabPat with the help of ruby on the gitlab container will generate a new Public Access Token
func (i *Infra) CreateGitlabPat() error {
	if infraStep(i.state) >= INFRA_STEP_GITLAB_PAT_CREATED && GitlabPat(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Creating a new Gitlab Public access token."))

	gitlabPassword := state.Get(i.state, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY))
	if gitlabPassword == "" {
		return fmt.Errorf("gitlab password is an empty string. Was it deleted? Rerunning the infra may help.")
	}
	gitlabPat := gitlabPassword[:20]

	out, err := i.client.Run(fmt.Sprintf("docker exec gitlab-ce gitlab-rails runner 'token = User.find_by_username(\"root\").personal_access_tokens.create(scopes: [:read_user, :read_repository, :api, :create_runner, :manage_runner, :sudo, :admin_mode], name: \"Automation token\", expires_at: 356.days.from_now); token.set_token(\"%s\"); token.save! '", gitlabPat))
	if err != nil {
		if !isDuplicateKeyGitlab(string(out)) {
			return fmt.Errorf("gitlab admin public access token creation exited unsuccessfully. err was %w, output was:\n%s", err, out)
		}
	}

	i.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_GITLAB_PAT_CREATED),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PAT_KEY, gitlabPat),
		state.WithMsg(INFRA_STATE, "Gitlab admin public access token created successfully."),
	)

	return nil
}

// getRunnerAuthToken will fetch an auth token that can be used to register a new gitlab runner
func (i *Infra) GetRunnerAuthToken() error {
	if infraStep(i.state) >= INFRA_STEP_GITLAB_RUNNER_AUTH_TOKEN && GitlabRunnerAuthToken(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Creating a auth token for the gitlab runner."))

	gitlabPat := GitlabPat(i.state)
	if gitlabPat == "" {
		return fmt.Errorf("gitlab pat is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	client := gitlab.NewGitlabClient(gitlabUrl(*i.options), gitlabPat)

	token, err := client.GetRunnerAuthToken("instance_type")
	if err != nil {
		return fmt.Errorf("error while trying to fetch runner auth token. Err was: %s")
	}

	i.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_GITLAB_RUNNER_AUTH_TOKEN),
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY, token),
		state.WithMsg(INFRA_STATE, "Gitlab auth token successfully created."),
	)

	return nil
}

// getRunnerAuthToken will fetch an auth token that can be used to register a new gitlab runner
func (i *Infra) RegisterGitlabRunner() error {
	if infraStep(i.state) >= INFRA_STEP_GITLAB_RUNNER_REGISTERED {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Registering gitlab runner with generated auth token"))

	runnerAuthToken := state.Get(i.state, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY))
	if runnerAuthToken == "" {
		return fmt.Errorf("runner auth token is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	out, err := i.client.Run(fmt.Sprintf("docker exec gitlab-runner gitlab-runner register --non-interactive --url \"%s\" --token \"%s\" --executor \"docker\" --docker-image alpine:latest --description \"docker-runner\"", gitlabUrl(*i.options), runnerAuthToken))
	if err != nil {
		return fmt.Errorf("registering a gitlab runner exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_GITLAB_RUNNER_REGISTERED),
		state.WithMsg(INFRA_STATE, "Gitlab runner registered successfully."),
	)

	return nil
}

// fetchNexusPassword will fetch the password for Nexus and store it in the context store
func (i *Infra) FetchNexusPassword() error {
	if infraStep(i.state) >= INFRA_STEP_INFO_FETCHED_NEXUS && NexusAdminPassword(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Fetching nexus admin password"))

	out, err := i.client.Run("docker exec nexus cat /nexus-data/admin.password")
	if err != nil {
		if isNoSuchFileOrDirectoryErr(string(out)) {
			pass, err := prompt.AskPass("You've already went through the nexus initial wizard.", "In order to continue execution, please provide nexus password manually:")
			if err != nil {
				return fmt.Errorf("error while providing nexus password. err was %w", err, out)
			}

			out = []byte(pass)
		} else {
			return fmt.Errorf("nexus admin password fetching exited unsuccessfully. err was %w, output was:\n%s", err, out)
		}
	}

	i.state.Set(
		state.WithStep(INFRA_STATE, INFRA_STEP_INFO_FETCHED_NEXUS),
		state.WithContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Nexus admin password fetched successfully."),
	)

	return nil
}

// infraStep gets the current step for the infra setup that we are on
func infraStep(s *state.State) int {
	return state.Get(s, state.GetStep(INFRA_STATE))
}

func GitlabAdminPassword(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY))
}

func GitlabPat(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PAT_KEY))
}

func GitlabRunnerAuthToken(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY))
}

func NexusAdminPassword(s *state.State) string {
	return state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY))
}

func isNoSuchFileOrDirectoryErr(msg string) bool {
	return strings.Contains(msg, "No such file or directory")
}

// isDuplicateKeyGitlab is a case where the gitlabPat was deleted for some reason from the state file and when we are creating it we are
// getting an error that it's a duplicate. If that is the case, we can assume that it is the one we are trying to pass anyway
func isDuplicateKeyGitlab(msg string) bool {
	return strings.Contains(msg, "duplicate key value violates unique constraint")
}

func gitlabUrl(opts args.Options) string {
	return fmt.Sprintf("http://%s:8082/gitlab", opts.Infra.SSHVMIP)
}
