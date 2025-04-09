package reset

import (
	"flag"
	"os"

	"github.com/Michaelpalacce/go-btva/cmd/reset/reset_options"
	"github.com/Michaelpalacce/go-btva/internal/run/state"
	"github.com/Michaelpalacce/go-btva/pkg/args"
)

var options = &reset_options.ResetOptions{}

var usage = `Resets the state. WARN: You could lose data!`

func (c *ResetCommand) Args() *reset_options.ResetOptions {
	if options.Parsed {
		return options
	}

	args, err := args.NewArgs(
		os.Args[2:],
		args.WithUsage(usage),
		args.WithFs(flag.NewFlagSet("run", flag.ExitOnError)),
	)
	if err != nil {
		panic(err)
	}

	args.AddVar(&options.AssumeYes, "", "y", false, "Assume yes, don't promp.")
	args.AddVar(&options.StateFile, "", "f", state.JSON_STORAGE_FILE, "State file to delete")

	args.Parse()

	options.Parsed = true

	return options
}
