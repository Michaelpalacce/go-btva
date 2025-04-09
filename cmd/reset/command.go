package reset

import (
	"log/slog"

	"github.com/Michaelpalacce/go-btva/pkg/file"
	"github.com/Michaelpalacce/go-btva/pkg/prompt"
)

// ResetCommand will reset the state (for now)
type ResetCommand struct{}

func (c *ResetCommand) Name() string {
	return "reset"
}

func (c *ResetCommand) Run() error {
	options := c.Args()

	if !options.AssumeYes {
		a, err := prompt.AskYesNo("Are you sure? You may lose data!")

		if err != nil || !a {
			slog.Warn("Skipping state deletion.")
			return nil
		} else {
			slog.Info("Deleting state.")
		}
	}

	return file.DeleteIfExists(options.StateFile)
}
