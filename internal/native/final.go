package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/native/components/infra"
	"github.com/Michaelpalacce/go-btva/internal/state"
)

const (
	FINAL_INSTRUCTIONS_NEXUS_STEP = iota + 1
	FINAL_INSTRUCTIONS_GITLAB_STEP
)

func (h *Handler) nexusInstructions() error {
	if state.Get(h.state, finalStep()) >= FINAL_INSTRUCTIONS_NEXUS_STEP {
		return nil
	}

	nexusPassword := infra.NexusAdminPassword(h.state)
	if nexusPassword == "" {
		return fmt.Errorf("nexus password is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	slog.Info("===============================================================")
	slog.Info("============================ NEXUS ============================")
	slog.Info("===============================================================")
	slog.Warn("Nexus has an initial setup wizard that needs to be followed through the UI.")
	slog.Info(fmt.Sprintf("Please visit: http://%s:8081/nexus", h.options.Infra.SSHVMIP))
	slog.Info("Username: admin")
	slog.Info(fmt.Sprintf("Password: %s", nexusPassword))

	h.state.Set(
		state.WithStep(FINAL_STATE, FINAL_INSTRUCTIONS_NEXUS_STEP),
		state.WithQuietMsg(FINAL_STATE, "Printed Nexus instructions"),
	)

	return nil
}

func (h *Handler) gitlabInstructions() error {
	if state.Get(h.state, finalStep()) >= FINAL_INSTRUCTIONS_GITLAB_STEP {
		return nil
	}

	gitlabPassword := infra.GitlabAdminPassword(h.state)
	if gitlabPassword == "" {
		return fmt.Errorf("gitlab password is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	slog.Info("===============================================================")
	slog.Info("============================ GITLAB ===========================")
	slog.Info("===============================================================")
	slog.Info("Gitlab setup with a CI/CD runner")
	slog.Info(fmt.Sprintf("Gitlab: http://%s:8081/gitlab", h.options.Infra.SSHVMIP))
	slog.Info("Username: root")
	slog.Info(fmt.Sprintf("Password: %s", gitlabPassword))

	h.state.Set(
		state.WithStep(FINAL_STATE, FINAL_INSTRUCTIONS_GITLAB_STEP),
		state.WithQuietMsg(FINAL_STATE, "Printed Gitlab instructions"),
	)

	return nil
}

func finalStep() state.GetStepStateOption {
	return state.GetStep(FINAL_STATE)
}
