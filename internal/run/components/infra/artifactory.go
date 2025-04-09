package infra

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/internal/run/templates"
)

// InfraSettingsXml will replace the ~/.m2/settings.xml file on the os
func (i *InfraComponent) InfraSettingsXml() error {
	m2SettingsPath := fmt.Sprintf("%s/.m2/settings.xml", i.os.HomeDir)

	return templates.SettingsXml(m2SettingsPath, i.options.ArtifactManager, i.options.Aria.Automation)
}
