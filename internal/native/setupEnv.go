package native

import (
	"embed"
	"fmt"
	"log/slog"
	osz "os"
	"text/template"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/state"
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

type settingsInventory struct {
	Nexus  nexus
	Gitlab gitlab

	Infra infra
}

//go:embed templates/*
var templates embed.FS

// prepareSettingsXml will replace the `settings.xml` in your `~/.m2` dir
func (h *Handler) prepareSettingsXml(os *os.OS, options *args.Options, s *state.State) error {
	if state.Get(h.state, envStep()) >= ENV_STEP_SETTINGS_XML {
		slog.Info("Skipping settings.xml configuration. Already done.")
		return nil
	}

	slog.Info("Configuring `settings.xml`.")
	baseURL := fmt.Sprintf("http://%s/nexus/repository/", options.Infra.SSHVMIP)

	templateVars := settingsInventory{
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

	template, err := template.New("settings.xml").ParseFS(templates, "templates/settings.xml")
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
	defer fo.Close()

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
