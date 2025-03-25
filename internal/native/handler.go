package native

import (
	"fmt"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/internal/software/linux"
	"github.com/Michaelpalacce/go-btva/internal/ssh"
	"github.com/Michaelpalacce/go-btva/pkg/os"
	"github.com/Michaelpalacce/go-btva/pkg/state"
)

// Handler is a struct that orchestrates the setup process based on OS
type Handler struct {
	os      *os.OS
	state   *state.State
	options *args.Options

	// installer is a pointer
	installer Installer
}

// NewHandler will return a new native Handler that will be used to manage and execute os operations
func NewHandler(os *os.OS, options *args.Options) (*Handler, error) {
	handler := &Handler{os: os, state: state.NewState(), options: options}

	if options.Local.SaveState {
		if err := handler.state.Modify(state.WithJsonStorage(options.Local.StateJson, true)); err != nil {
			return nil, err
		}
	}

	switch os.Distro {
	case "linux":
		handler.installer = &linux.LinuxInstaller{OS: os, Options: options}
	case "windows":
		fallthrough
	case "darwin":
		fallthrough
	default:
		return nil, fmt.Errorf("OS %s is not supported", os.Distro)
	}

	return handler, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Software Block

// SetupSoftware will install all the needed software based on the os and options
// @NOTE: This is meant to be ran async
func (h *Handler) SetupSoftware(c chan error) {
	if h.options.Software.InstallJava {
		if err := h.installSoftware(h.installer.Java()); err != nil {
			c <- err
			return
		}
	}

	if h.options.Software.InstallMvn {
		if err := h.installSoftware(h.installer.Mvn()); err != nil {
			c <- err
			return
		}
	}

	if h.options.Software.InstallNode {
		if err := h.installSoftware(h.installer.Node()); err != nil {
			c <- err
			return
		}
	}

	c <- nil
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

// END Software Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Local Env Block

// @TODO: Finish
func (h *Handler) SetupLocalEnv(c chan error) {
	c <- nil
}

// END Setup Local Env Block

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Setup Infra Block

// @TODO: Finish
func (h *Handler) SetupInfra(c chan error) {
	infraOptions := h.options.Infra
	client, err := ssh.GetClient(infraOptions.SSHVMIP, infraOptions.SSHUsername, infraOptions.SSHPassword, infraOptions.SSHPrivateKey, infraOptions.SSHPrivateKeyPassphrase)
	if err != nil {
		c <- fmt.Errorf("could not create client. err was %w", err)
		return
	}

	defer client.Close()

	out, err := client.Run("ls -lah /tmp")
	if err != nil {
		c <- fmt.Errorf("process exited with error. err was %w, output was:\n%s", err, out)
		return
	}

	fmt.Println(string(out))

	c <- nil
}

// END Setup Infra Block
