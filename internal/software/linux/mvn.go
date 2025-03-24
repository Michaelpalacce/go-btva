package linux

import (
	"fmt"
	"os/exec"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/software"
	"github.com/Michaelpalacce/go-btva/pkg/file"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// MvnSoftware is responsible for installing, removing and checking if mvn is installed
type MvnSoftware struct {
	OS      *os.OS
	Options *args.Options

	initialized bool
}

var mvnSoftware *MvnSoftware = &MvnSoftware{}

// Install will install mvn with apt
// @NOTE: Make sure you run the go process as sudo
func (s *MvnSoftware) Install() error {
	if err := s.removeTempFiles(); err != nil {
		return err
	}

	if err := s.ensureInstallZipExists(); err != nil {
		return err
	}

	if err := untar(fmt.Sprintf("/tmp/apache-maven-%s-bin.tar.gz", s.Options.Software.LinuxMvnVersion), "/opt"); err != nil {
		return err
	}

	if err := s.symlinkMvn(); err != nil {
		return err
	}

	if err := s.removeTempFiles(); err != nil {
		return err
	}

	return nil
}

// Remove will remove mvn
func (s *MvnSoftware) Remove() error {
	return nil
}

// Exists verifies if mvn is already installed.
// Relies on `which`, which returns exit code 0 if the program is found and 1 if not
func (s *MvnSoftware) Exists() bool {
	cmd := exec.Command("which", "mvn")
	_, err := cmd.CombinedOutput()

	return err == nil
}

// removeTempFiles is a helper that will remove the downloaded tar.gz files pre and post install
func (s *MvnSoftware) removeTempFiles() error {
	mvnRmCmd := exec.Command("rm", "-rf", fmt.Sprintf("/tmp/apache-maven-%s-bin.tar.gz", s.Options.Software.LinuxMvnVersion))

	if output, err := mvnRmCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error while running command. Error was %w, output was %s", err, output)
	}

	return nil
}

// ensureInstallZipExists will download the mvn tar.gz file if it does not exist
func (s *MvnSoftware) ensureInstallZipExists() error {
	if file.Exists(s.getInstallZipPath()) {
		return nil
	}

	mvnWgetCmd := exec.Command(
		"wget",
		fmt.Sprintf(
			"https://downloads.apache.org/maven/maven-3/%s/binaries/apache-maven-%s-bin.tar.gz",
			s.Options.Software.LinuxMvnVersion,
			s.Options.Software.LinuxMvnVersion,
		),
		"-P",
		"/tmp",
	)

	if output, err := mvnWgetCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error while running command. Error was %w, output was %s", err, output)
	}

	return nil
}

// getInstallZipPath is an internal function that will give us the download installer zip location
func (s *MvnSoftware) getInstallZipPath() string {
	return fmt.Sprintf("/tmp/apache-maven-%s-bin.tar.gz", s.Options.Software.LinuxMvnVersion)
}

// symlinkMvn will symlink the mvn binary to /usr/bin/mvn
func (s *MvnSoftware) symlinkMvn() error {
	mvnLnCmd := exec.Command("ln", "-sf", fmt.Sprintf("/opt/apache-maven-%s/bin/mvn", s.Options.Software.LinuxMvnVersion), "/usr/bin/mvn")

	if output, err := mvnLnCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error while running command. Error was %w, output was %s", err, output)
	}

	return nil
}

func (s *MvnSoftware) GetName() string    { return software.MvnSoftwareKey }
func (s *MvnSoftware) GetVersion() string { return s.Options.Software.LinuxMvnVersion }

// Java will return the MvnSoftware object that can be used to install, remove or check if mvn exists
// Only a single instance of the MvnSoftware will be returned
func (i *LinuxInstaller) Mvn() software.Software {
	if !mvnSoftware.initialized {
		mvnSoftware.OS = i.OS
		mvnSoftware.Options = i.Options
		mvnSoftware.initialized = true
	}

	return mvnSoftware
}

// Helper funcs

// untar works on a tar.gz file and untars it
func untar(filename string, destination string) error {
	mvnExtractCmd := exec.Command("tar", "xf", filename, "-C", destination)

	if output, err := mvnExtractCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error untarring file. Error was %w, output was %s", err, output)
	}

	return nil
}
