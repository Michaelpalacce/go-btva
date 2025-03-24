package native

import (
	"github.com/Michaelpalacce/go-btva/internal/software"
)

// Installer is a common interface implemented by the installers of all the major systems
type Installer interface {
	Java() software.Software
	Mvn() software.Software
	Node() software.Software
}
