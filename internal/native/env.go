package native

import (
	"fmt"
	"log/slog"
	osz "os"
	"text/template"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/file"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

const (
	ENV_STEP_SETTINGS_XML = iota + 1
)

const (
	ENV_STATE = "Env"
)

type artifactory struct {
	ReleaseRepo  string
	SnapshotRepo string
	GroupRepo    string
}

type infra struct {
	Artifactory artifactory
}

type nexus struct {
	Password string
}

type gitlab struct {
	Password string
}

type settingsTemplate struct {
	Nexus  nexus
	Gitlab gitlab

	Infra infra
}

// prepareSettingsXml will replace the `settings.xml` in your `~/.m2` dir
func (h *Handler) prepareSettingsXml(os *os.OS, options *args.Options, s *state.State) error {
	if state.Get(h.state, envStep()) >= ENV_STEP_SETTINGS_XML {
		slog.Info("Skipping settings.xml configuration. Already done.")
		return nil
	}

	slog.Info("Configuring `settings.xml`.")

	templateSettingsPath := fmt.Sprintf("%s/configs/settings.xml", os.Cwd)
	tmpSettingsPath := fmt.Sprintf("%s/settings.xml", os.TempDir)
	if err := file.DeleteIfExists(tmpSettingsPath); err != nil {
		h.state.Set(state.WithErr(ENV_STATE, err))
		return fmt.Errorf("there was an existing settings file at: %s. Could not remove it. Err was %w", tmpSettingsPath, err)
	}

	if _, err := file.Copy(templateSettingsPath, fmt.Sprintf("%s/settings.xml", os.TempDir)); err != nil {
		h.state.Set(state.WithErr(ENV_STATE, err))
		return fmt.Errorf("could not copy settings.xml file to a temp dir. Err was %w", err)
	}

	baseURL := fmt.Sprintf("http://%s/nexus/repository/", options.Infra.SSHVMIP)

	templateVars := settingsTemplate{
		Nexus: nexus{
			Password: state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_NEXUS_PASSWORD_KEY)),
		},
		Gitlab: gitlab{
			Password: state.Get(s, state.GetContextProp(INFRA_STATE, INFRA_GITLAB_PASSWORD_KEY)),
		},
		Infra: infra{
			Artifactory: artifactory{
				ReleaseRepo:  baseURL + "maven-releases",
				SnapshotRepo: baseURL + "maven-snapshots",
				GroupRepo:    baseURL + "maven-public",
			},
		},
	}

	settingsBytes, err := osz.ReadFile(templateSettingsPath)
	if err != nil {
		h.state.Set(state.WithErr(ENV_STATE, err))
		return fmt.Errorf("could not read settings.xml file. Err was %w", err)
	}

	template, err := template.New("settings.xml").Parse(string(settingsBytes))
	if err != nil {
		h.state.Set(state.WithErr(ENV_STATE, err))
		return fmt.Errorf("could not parse settings.xml file. Err was %w", err)
	}

	m2SettingsPath := fmt.Sprintf("%s/.m2/settings.xml", os.HomeDir)

	fo, err := osz.Create(m2SettingsPath)
	if err != nil {
		h.state.Set(state.WithErr(ENV_STATE, err))
		return fmt.Errorf("could not open file %s for writing. Err was %w", m2SettingsPath, err)
	}

	err = template.Execute(fo, templateVars)
	if err != nil {
		h.state.Set(state.WithErr(ENV_STATE, err))
		return fmt.Errorf("could replace template vars. Err was %w", err)
	}

	h.state.Set(
		state.WithMsg(ENV_STATE, "Finished configuring settings.xml"),
		state.WithStep(ENV_STATE, ENV_STEP_SETTINGS_XML),
		state.WithErr(ENV_STATE, nil),
	)

	return nil
}

func envDone() state.GetSuccessStateOption {
	return state.GetDone(ENV_STATE)
}

func envStep() state.GetStepStateOption {
	return state.GetStep(ENV_STATE)
}
