package env

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/components/infra"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/file"
	osl "github.com/Michaelpalacce/go-btva/pkg/os"
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
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

type ariaInventory struct {
	FQDN        string
	Port        string
	Username    string
	Password    string
	OrgName     string
	ProjectName string
}

type settingsInventory struct {
	Nexus  nexusInventory
	Gitlab gitlabInventory

	Infra infraInventory
	Aria  ariaInventory
}

//go:embed templates/*
var templates embed.FS

// SettingsXml will replace the `settings.xml` in your `~/.m2` dir
func (e *Env) SettingsXml() error {
	if state.Get(e.state, envStep()) >= ENV_STEP_SETTINGS_XML {
		return nil
	}
	slog.Info("Configuring `settings.xml`.")

	m2SettingsPath := fmt.Sprintf("%s/.m2/settings.xml", e.os.HomeDir)
	baseURL := fmt.Sprintf("http://%s/nexus/repository/", e.options.Infra.SSHVMIP)

	if file.Exists(m2SettingsPath) {
		var a bool
		var err error

		if a, err = prompt.AskYesNo(fmt.Sprintf("settings.xml file found in %s. Are you sure you want to replace it?", m2SettingsPath)); err != nil {
			return fmt.Errorf("could not get an answer. Err was: %w", err)
		}

		if !a {
			e.state.Set(
				state.WithMsg(ENV_STATE, "User has settings.xml already and terminated the replacement"),
				state.WithStep(ENV_STATE, ENV_STEP_SETTINGS_XML),
				state.WithErr(ENV_STATE, nil),
			)
			return nil
		}
	}

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
		Aria: getAriaInventory(),
	}

	template, err := template.New("settings.xml").ParseFS(templates, "templates/settings.xml")
	if err != nil {
		return fmt.Errorf("could not parse settings.xml file. Err was %w", err)
	}

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

// getAriaInventory will prompt the user a series of question needed to build the aria inventory
func getAriaInventory() ariaInventory {
	inv := ariaInventory{FQDN: "vra-l-01a.corp.local", Port: "443", Username: "configurationadmin", Password: "", OrgName: "vidm-l-01a", ProjectName: "dev"}

	if ans, err := prompt.AskText(fmt.Sprintf("What is Aria Automation's FQDN without `https://`. Default (%s)", inv.FQDN)); err == nil {
		inv.FQDN = ans
	}

	if ans, err := prompt.AskText(fmt.Sprintf("What is Aria Automation's port? Default (%s)", inv.Port)); err == nil {
		inv.Port = ans
	}

	if ans, err := prompt.AskText(fmt.Sprintf("What is the username of the account for Aria Automation? Default (%s)", inv.Username)); err == nil {
		inv.Username = ans
	}

	if ans, err := prompt.AskPass("What is the password of the account for Aria Automation?"); err == nil {
		inv.Password = ans
	}

	if ans, err := prompt.AskText(fmt.Sprintf("What is the org name used in Aria Automation? Default (%s)", inv.OrgName)); err == nil {
		inv.OrgName = ans
	}

	if ans, err := prompt.AskText(fmt.Sprintf("What is the default project name in Aria Automation you want to push automation code to? Default (%s)", inv.ProjectName)); err == nil {
		inv.ProjectName = ans
	}

	return inv
}

func envStep() state.GetStepStateOption {
	return state.GetStep(ENV_STATE)
}
