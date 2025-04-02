package infra

import "github.com/Michaelpalacce/go-btva/internal/templates"

// InfraSettingsXml will replace the ~/.m2/settings.xml file on the os
func (i *InfraComponent) InfraSettingsXml() error {
	return templates.SettingsXml(i.os.HomeDir, i.options.Artifactory, i.options.Aria.Automation)
}
