package final

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/components/infra"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type Final struct {
	os      *os.OS
	state   *state.State
	options *args.Options
}

func NewFinal(os *os.OS, state *state.State, options *args.Options) *Final {
	return &Final{os: os, state: state, options: options}
}

// MinimalInfraNexusInstructions will print out details for nexus
func (f *Final) MinimalInfraNexusInstructions() error {
	nexusPassword := infra.NexusAdminPassword(f.state)
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

func (f *Final) MinimalInfraGitlabInstructions() error {
	gitlabPassword := infra.GitlabAdminPassword(f.state)
	if gitlabPassword == "" {
		return fmt.Errorf("gitlab password is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	gitlabPat := infra.GitlabPat(f.state)
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
