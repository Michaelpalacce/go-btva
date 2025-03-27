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
	if h.state.GetDone(software.SoftwareDone(soft)) {
		slog.Info("Software already installed, skipping...", "name", soft.GetName(), "version", soft.GetVersion())
		return nil
	}

	if soft.Exists() {
		h.state.Set(software.WithSoftwareInstalled(soft, nil))
		return nil
	}

	slog.Info("Software is not installed, installing", "name", soft.GetName(), "version", soft.GetVersion())

	err := soft.Install()
	if err != nil {
		h.state.Set(software.WithSoftwareInstalled(soft, err))
		return err
	}

	h.state.Set(software.WithSoftwareInstalled(soft, nil))

	return nil
}
