package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/state"
)

const (
	FINAL_STATE = "Final"

	FINAL_INSTRUCTIONS_NEXUS_STEP = iota + 1
)

func (h *Handler) Instructions() error {
	if state.Get(h.state, finalStep()) >= FINAL_INSTRUCTIONS_NEXUS_STEP {
		return nil
	}

	nexusPassword := state.Get(h.state, state.GetContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY))
	if nexusPassword == "" {
		return fmt.Errorf("nexus password is an empty string. Was it deleted? Rerunning the infra may help.")
	}

	slog.Info("==========================================================================")
	slog.Info("==========================================================================")
	slog.Info("==========================================================================")
	slog.Info("Everything is setup.")
	slog.Info("Nexus has an initial setup wizard that needs to be followed through the UI.")
	slog.Info(fmt.Sprintf("Please visit: http://%s:8081/nexus", h.options.Infra.SSHVMIP))
	slog.Info("Username: admin")
	slog.Info(fmt.Sprintf("Password: %s", nexusPassword))

	h.state.Set(
		state.WithStep(FINAL_STATE, FINAL_INSTRUCTIONS_NEXUS_STEP),
		state.WithQuietMsg(FINAL_STATE, "Printed Nexus instructions"),
	)

	return nil
}

// finalDone will give us a state.GetSuccessStateOption that will check if the final part was ran before
func finalDone() state.GetSuccessStateOption {
	return state.GetDone(FINAL_STATE)
}

func finalStep() state.GetStepStateOption {
	return state.GetStep(FINAL_STATE)
}
