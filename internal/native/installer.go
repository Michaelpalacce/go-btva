package native

import (
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/software"
)

// Installer is a common interface implemented by the installers of all the major systems
type Installer interface {
	Java() software.Software
	Mvn() software.Software
	Node() software.Software
}

// installSoftware is an internal function that can be used to install any software. It will run through a set of commands
func (h *Handler) installSoftware(soft software.Software) error {
	if h.state.GetDone(software.IsSoftwareNotInstalled(soft)) && !soft.Exists() {
		slog.Info("Software is not installed, installing", "name", soft.GetName(), "version", soft.GetVersion())

		err := soft.Install()
		if err != nil {
			if err := h.state.Set(software.SoftwareInstalled(soft, err)); err != nil {
				slog.Error("Error setting state", err)
			}
			return err
		}

		if err := h.state.Set(software.SoftwareInstalled(soft, nil)); err != nil {
			slog.Error("Error setting state", err)
		}

		slog.Info("Software successfully installed", "name", soft.GetName(), "version", soft.GetVersion())
	} else {
		slog.Info("Software already installed, skipping...", "name", soft.GetName(), "version", soft.GetVersion())
	}

	return nil
}
