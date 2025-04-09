package options

import (
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

type ResetOptions struct {
	AssumeYes bool
	State     bool
	StateFile string

	Parsed bool
}

func (options *ResetOptions) ValidateState() error {
	if options.State && options.StateFile == "" {
		var err error
		if options.StateFile, err = prompt.AskText("`StateFile` set to empty, but `State` is true. Please provide a value"); err != nil {
			return err
		}
	}

	return nil
}
