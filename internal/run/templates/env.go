package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/pkg/file"
)

type settingsInventory struct {
	ArtifactManager options.ArtifactManager
	Aria            options.AriaAutomation
}

//go:embed templates/*
var templates embed.FS

// SettingsXml will create a settings.xml file at teh given location
// This will skip the creation if it exists and warn the user
func SettingsXml(m2SettingsPath string, artifactoryInventory options.ArtifactManager, ariaInventory options.AriaAutomation) error {
	if file.Exists(m2SettingsPath) {
		slog.Warn(fmt.Sprintf("%s already exists. Skipping replacement...", m2SettingsPath))
		return nil
	}

	slog.Info("Configuring `settings.xml`.")

	templateVars := settingsInventory{
		ArtifactManager: artifactoryInventory,
		Aria:            ariaInventory,
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
