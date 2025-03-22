package linux

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type LinuxInstaller struct {
	OS      *os.OS
	Options *args.Options
}

func (i *LinuxInstaller) GetOS() *os.OS {
	return i.OS
}
