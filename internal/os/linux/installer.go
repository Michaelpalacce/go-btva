package linux

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type Installer struct {
	OS      *os.OS
	Options *args.Options
}

// GetAllSoftware will return all the software in order that need to be installed
func (i *Installer) GetAllSoftware() []software.Software {
	return []software.Software{i.Java(), i.Mvn(), i.Node(), i.VsCode()}
}
