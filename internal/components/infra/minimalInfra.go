package infra

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/internal/templates"
	"github.com/Michaelpalacce/go-btva/pkg/gitlab"
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
	"github.com/melbahja/goph"
)

// @NOTE: This file contains the definition of all the tasks for minimal infrastructure

const (
	// Private
	_INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY = "gitlabRunnerAuthToken"
	_INFRA_GITLAB_RUNNER_REGISTERED_KEY = "gitlabRunnerRegistered"
	_INFRA_SETUP_DONE_KEY               = "infraSetup"
)

const (
	_BTVA_INSTALL_DIR_INFRA         = "/opt/build-tools-for-vmware-aria/infrastructure"
	_BTVA_MINIMAL_INFRA_INSTALL_URL = "https://raw.githubusercontent.com/vmware/build-tools-for-vmware-aria/refs/heads/refactor/minimal-infra-simplified-setup/infrastructure/install.sh"
)

// getClient will ssh into the machine and give you a goph.Client pointer you can use to run commands.
// @WARN: Make sure to defer client.Close()
func getClient(options *args.Options) (*goph.Client, error) {
	infraOptions := options.Infra

	client, err := ssh.GetClient(infraOptions.SSHVMIP, infraOptions.SSHUsername, infraOptions.SSHPassword, infraOptions.SSHPrivateKey, infraOptions.SSHPrivateKeyPassphrase)
	if err != nil {
		return nil, fmt.Errorf("could not create a new ssh client and/or ssh into machine. Err was: %w", err)
	}

	return client, nil
}

// RunMinimalInfra will fetch the BTVA minimal infra installer and run it
// @TODO: Fix the branch
func (i *InfraComponent) RunMinimalInfra() error {
	if state.Get(i.state, state.GetContextProp(INFRA_STATE, _INFRA_SETUP_DONE_KEY)) == "true" {
		i.state.Set(
			state.WithWarn(INFRA_STATE, "Minimal Infrastructure already setup. Skipping"),
		)
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Running the minimal infrastructure installer. This may take a few minutes as it waits for services to be healthy."))

	client, err := getClient(i.options)
	if err != nil {
		return err
	}
	defer client.Close()

	out, err := client.Run(fmt.Sprintf("curl -o- %s | bash -s -- %s %q", _BTVA_MINIMAL_INFRA_INSTALL_URL, i.options.Infra.DockerUsername, i.options.Infra.DockerPAT))
	if err != nil {
		return fmt.Errorf("minimal infrastructure installer exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	// Fixes the compose install
	if out, err := client.Run(fmt.Sprintf("sed -i \"s|external_url 'http://infra.corp.local/gitlab'|external_url '%q'|\" %s/docker-compose.yml", gitlabUrl(*i.options), _BTVA_INSTALL_DIR_INFRA)); err != nil {
		return fmt.Errorf("failed to modify compose file. err was %w, output was:\n%s", err, out)
	}

	if out, err := client.Run(fmt.Sprintf("docker compose -f %s/docker-compose.yml up -d --wait", _BTVA_INSTALL_DIR_INFRA)); err != nil {
		return fmt.Errorf("failed to start containers. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithContextProp(INFRA_STATE, _INFRA_SETUP_DONE_KEY, "true"),
		state.WithMsg(INFRA_STATE, "Minimal infrastructure installer successfully set up."),
	)

	return nil
}

// FetchGitlabPassword will fetch the password for Gitlab and store it in the context store
// Command looks a bit big, but it's all so we can fail in case the file doesn't exists or the container is not started
func (i *InfraComponent) FetchGitlabPassword() error {
	if GitlabAdminPassword(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Fetching gitlab admin password"))

	client, err := getClient(i.options)
	if err != nil {
		return err
	}
	defer client.Close()

	out, err := client.Run("docker exec gitlab-ce test -f /etc/gitlab/initial_root_password && docker exec gitlab-ce grep 'Password:' /etc/gitlab/initial_root_password | awk '{print $2}'")
	if err != nil {
		return fmt.Errorf("gitlab admin password fetching exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Gitlab admin password fetched successfully."),
	)

	return nil
}

// CreateGitlabPat with the help of ruby on the gitlab container will generate a new Public Access Token
func (i *InfraComponent) CreateGitlabPat() error {
	if GitlabPat(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Creating a new Gitlab Public access token."))

	gitlabPassword := state.Get(i.state, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PASSWORD_KEY))
	if gitlabPassword == "" {
		return fmt.Errorf("gitlab password is an empty string. Was it deleted? Rerunning the infra may help.")
	}
	gitlabPat := gitlabPassword[:20]

	client, err := getClient(i.options)
	if err != nil {
		return err
	}
	defer client.Close()

	out, err := client.Run(fmt.Sprintf("docker exec gitlab-ce gitlab-rails runner 'token = User.find_by_username(\"root\").personal_access_tokens.create(scopes: [:read_user, :read_repository, :api, :create_runner, :manage_runner, :sudo, :admin_mode], name: \"Automation token\", expires_at: 356.days.from_now); token.set_token(\"%s\"); token.save! '", gitlabPat))
	if err != nil && !isDuplicateKeyGitlab(string(out)) {
		return fmt.Errorf("gitlab admin public access token creation exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithContextProp(INFRA_STATE, INFRA_GITLAB_ADMIN_PAT_KEY, gitlabPat),
		state.WithMsg(INFRA_STATE, "Gitlab admin public access token created successfully."),
	)

	return nil
}

// GetRunnerAuthToken will fetch an auth token that can be used to register a new gitlab runner
func (i *InfraComponent) GetRunnerAuthToken() error {
	if GitlabRunnerAuthToken(i.state) != "" {
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
		state.WithContextProp(INFRA_STATE, _INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY, token),
		state.WithMsg(INFRA_STATE, "Gitlab auth token successfully created."),
	)

	return nil
}

// getRunnerAuthToken will fetch an auth token that can be used to register a new gitlab runner
// Re-running this is ok, it will say it's already valid
func (i *InfraComponent) RegisterGitlabRunner() error {
	if state.Get(i.state, state.GetContextProp(INFRA_STATE, _INFRA_GITLAB_RUNNER_REGISTERED_KEY)) == "true" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Registering gitlab runner with generated auth token"))

	runnerAuthToken := state.Get(i.state, state.GetContextProp(INFRA_STATE, _INFRA_GITLAB_RUNNER_AUTH_TOKEN_KEY))
	if runnerAuthToken == "" {
		return fmt.Errorf("runner auth token is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	client, err := getClient(i.options)
	if err != nil {
		return err
	}
	defer client.Close()

	out, err := client.Run(fmt.Sprintf("docker exec gitlab-runner gitlab-runner register --non-interactive --url \"%s\" --token \"%s\" --executor \"docker\" --docker-image alpine:latest --description \"docker-runner\"", gitlabUrl(*i.options), runnerAuthToken))
	if err != nil {
		return fmt.Errorf("registering a gitlab runner exited unsuccessfully. err was %w, output was:\n%s", err, out)
	}

	i.state.Set(
		state.WithContextProp(INFRA_STATE, _INFRA_GITLAB_RUNNER_REGISTERED_KEY, "true"),
		state.WithMsg(INFRA_STATE, "Gitlab runner registered successfully."),
	)

	return nil
}

// FetchNexusPassword will fetch the password for Nexus and store it in the context store
func (i *InfraComponent) FetchNexusPassword() error {
	if NexusAdminPassword(i.state) != "" {
		return nil
	}

	i.state.Set(state.WithMsg(INFRA_STATE, "Fetching nexus admin password"))

	client, err := getClient(i.options)
	if err != nil {
		return err
	}
	defer client.Close()

	out, err := client.Run("docker exec nexus cat /nexus-data/admin.password")
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
		state.WithContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY, strings.TrimSpace(string(out))),
		state.WithMsg(INFRA_STATE, "Nexus admin password fetched successfully."),
	)

	return nil
}

// MinimalInfraNexusInstructions will print out details for nexus
func (f *InfraComponent) MinimalInfraNexusInstructions() error {
	nexusPassword := NexusAdminPassword(f.state)
	if nexusPassword == "" {
		return fmt.Errorf("nexus password is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	slog.Info("===============================================================")
	slog.Info("============================ NEXUS ============================")
	slog.Info("===============================================================")
	slog.Warn("Nexus has an initial setup wizard that needs to be followed through the UI.")
	slog.Info(fmt.Sprintf("Please visit: http://%s:8081/nexus", f.options.Infra.SSHVMIP))
	slog.Info("Username: admin")
	slog.Info(fmt.Sprintf("Password: %s", nexusPassword))

	return nil
}

// MinimalInfraGitlabInstructions will print out instructions for gitlab
func (f *InfraComponent) MinimalInfraGitlabInstructions() error {
	gitlabPassword := GitlabAdminPassword(f.state)
	if gitlabPassword == "" {
		return fmt.Errorf("gitlab password is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	gitlabPat := GitlabPat(f.state)
	if gitlabPat == "" {
		return fmt.Errorf("gitlab public access token is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	slog.Info("===============================================================")
	slog.Info("============================ GITLAB ===========================")
	slog.Info("===============================================================")
	slog.Info("Gitlab setup with a CI/CD runner")
	slog.Info(fmt.Sprintf("Gitlab: http://%s:8081/gitlab", f.options.Infra.SSHVMIP))
	slog.Info("Username: root")
	slog.Info(fmt.Sprintf("Password: %s", gitlabPassword))
	slog.Info(fmt.Sprintf("Public Access Token: %s", gitlabPat))

	return nil
}

// MinimalInfraSettingsXml will replace the `settings.xml` in your `~/.m2` dir
func (i *InfraComponent) MinimalInfraSettingsXml() error {
	baseURL := fmt.Sprintf("http://%s/nexus/repository/", i.options.Infra.SSHVMIP)

	return templates.SettingsXml(i.os.HomeDir, templates.ArtifactoryInventory{
		ReleaseRepo:  baseURL + "maven-releases",
		SnapshotRepo: baseURL + "maven-snapshots",
		GroupRepo:    baseURL + "maven-public",
		Password:     NexusAdminPassword(i.state),
	})
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
