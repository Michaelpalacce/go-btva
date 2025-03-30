package darwin

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/os/software"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type DarwinInstaller struct {
	OS      *os.OS
	Options *args.Options
}

// GetAllSoftware will return all the software in order that need to be installed
func (i *DarwinInstaller) GetAllSoftware() []software.Software {
	return []software.Software{i.Java(), i.Mvn(), i.Node()}
}
