package darwin

import (
	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

type DarwinInstaller struct {
	OS      *os.OS
	Options *args.Options
}

func (i *DarwinInstaller) GetOS() *os.OS {
	return i.OS
}
