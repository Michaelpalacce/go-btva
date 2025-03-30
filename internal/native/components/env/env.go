package env

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/native/components/infra"
	"github.com/Michaelpalacce/go-btva/internal/state"
	osl "github.com/Michaelpalacce/go-btva/pkg/os"
)

type Env struct {
	os      *osl.OS
	state   *state.State
	options *args.Options
}

func NewNev(os *osl.OS, state *state.State, options *args.Options) *Env {
	return &Env{os: os, state: state, options: options}
}

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

type infraInventory struct {
	Artifactory artifactory
}

type nexusInventory struct {
	Password string
}

type gitlabInventory struct {
	Password string
}

type settingsInventory struct {
	Nexus  nexusInventory
	Gitlab gitlabInventory

	Infra infraInventory
}

//go:embed templates/*
var templates embed.FS

// SettingsXml will replace the `settings.xml` in your `~/.m2` dir
func (e *Env) SettingsXml() error {
	if state.Get(e.state, envStep()) >= ENV_STEP_SETTINGS_XML {
		return nil
	}

	slog.Info("Configuring `settings.xml`.")
	baseURL := fmt.Sprintf("http://%s/nexus/repository/", e.options.Infra.SSHVMIP)

	templateVars := settingsInventory{
		Nexus: nexusInventory{
			Password: infra.NexusAdminPassword(e.state),
		},
		Gitlab: gitlabInventory{
			Password: infra.GitlabAdminPassword(e.state),
		},
		Infra: infraInventory{
			Artifactory: artifactory{
				ReleaseRepo:  baseURL + "maven-releases",
				SnapshotRepo: baseURL + "maven-snapshots",
				GroupRepo:    baseURL + "maven-public",
			},
		},
	}

	template, err := template.New("settings.xml").ParseFS(templates, "templates/settings.xml")
	if err != nil {
		return fmt.Errorf("could not parse settings.xml file. Err was %w", err)
	}

	m2SettingsPath := fmt.Sprintf("%s/.m2/settings.xml", e.os.HomeDir)

	fo, err := os.Create(m2SettingsPath)
	if err != nil {
		return fmt.Errorf("could not open file %s for writing. Err was %w", m2SettingsPath, err)
	}
	defer fo.Close()

	err = template.Execute(fo, templateVars)
	if err != nil {
		return fmt.Errorf("could replace template vars. Err was %w", err)
	}

	e.state.Set(
		state.WithMsg(ENV_STATE, "Finished configuring settings.xml"),
		state.WithStep(ENV_STATE, ENV_STEP_SETTINGS_XML),
		state.WithErr(ENV_STATE, nil),
	)

	return nil
}

func envStep() state.GetStepStateOption {
	return state.GetStep(ENV_STATE)
}
