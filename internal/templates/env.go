package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/pkg/file"
)

type settingsInventory struct {
	ArtifactManager args.ArtifactManager
	Aria            args.AriaAutomation
}

//go:embed templates/*
var templates embed.FS

func SettingsXml(homeDir string, artifactoryInventory args.ArtifactManager, ariaInventory args.AriaAutomation) error {
	m2SettingsPath := fmt.Sprintf("%s/.m2/settings.xml", homeDir)

	if file.Exists(m2SettingsPath) {
		return nil
		// answ, err := prompt.AskYesNo(fmt.Sprintf("%s already exists, do you want to replace it?", m2SettingsPath))
		// if err != nil {
		// 	return fmt.Errorf("could not get an answer if we should replace settings.xml")
		// }
		//
		// if !answ {
		// 	slog.Info("Selected to skip setting", "settingsXml", m2SettingsPath)
		// 	return nil
		// }
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
