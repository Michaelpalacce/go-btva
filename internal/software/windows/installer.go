package windows

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type WindowsInstaller struct {
	OS      *os.OS
	Options *args.Options
}

func (i *WindowsInstaller) GetOS() *os.OS {
	return i.OS
}
