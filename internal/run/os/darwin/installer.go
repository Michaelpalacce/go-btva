package darwin

import (
	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/internal/run/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type Installer struct {
	OS      *os.OS
	Options *options.RunOptions
}

// GetAllSoftware will return all the software in order that need to be installed
func (i *Installer) GetAllSoftware() []software.Software {
	return []software.Software{i.Java(), i.Mvn(), i.Node(), i.VsCode()}
}
