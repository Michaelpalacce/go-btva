package reset

import (
	"flag"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/options"
	"github.com/Michaelpalacce/go-btva/pkg/args"
)

var resetOptions = &options.ResetOptions{}

var usage = `Resets different parts of the state. WIP, currently just deletes the state file`

func (c *ResetCommand) Args() *options.ResetOptions {
	if resetOptions.Parsed {
		return resetOptions
	}

	args, err := args.NewArgs(
		os.Args[2:],
		args.WithUsage(usage),
		args.WithFs(flag.NewFlagSet("run", flag.ExitOnError)),
	)
	if err != nil {
		panic(err)
	}

	args.AddVar(&resetOptions.AssumeYes, "", "y", false, "Assume yes, don't promp.")
	args.AddVar(&resetOptions.State, "state", "s", true, "Reset the entire state.")
	args.AddVar(&resetOptions.StateFile, "stateFile", "sf", options.JSON_STORAGE_FILE, "State file modify or delete.")

	args.Parse()

	resetOptions.Parsed = true

	return resetOptions
}
