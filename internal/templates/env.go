package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/go-btva/pkg/file"
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

type ArtifactoryInventory struct {
	ReleaseRepo  string
	SnapshotRepo string
	GroupRepo    string
	Password     string
}

type infraInventory struct {
	Artifactory ArtifactoryInventory
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
	Infra infraInventory
	Aria  ariaInventory
}

//go:embed templates/*
var templates embed.FS

func SettingsXml(homeDir string, artifactoryInventory ArtifactoryInventory) error {
	m2SettingsPath := fmt.Sprintf("%s/.m2/settings.xml", homeDir)

	if file.Exists(m2SettingsPath) {
		return nil
	}

	slog.Info("Configuring `settings.xml`.")

	templateVars := settingsInventory{
		Infra: infraInventory{Artifactory: artifactoryInventory},
		Aria:  getAriaInventory(),
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

	slog.Info("Finished configuring settings.xml")

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
