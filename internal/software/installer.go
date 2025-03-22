package software

import (
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

// Installer is a common interface implemented by the installers of all the major systems
type Installer interface {
	InstallNode() error
	InstallJava() error
	InstallMvn() error
	GetOS() *os.OS
}
