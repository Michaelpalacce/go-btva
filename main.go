package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Michaelpalacce/go-btva/cmd/reset"
	"github.com/Michaelpalacce/go-btva/cmd/run"
	"github.com/Michaelpalacce/go-btva/pkg/logger"
)

type Runner interface {
	Run() error
	Name() string
}

func root(args []string) error {
	cmds := []Runner{
		&run.RunCommand{},
		&reset.ResetCommand{},
	}

	availableSubcommandsArr := make([]string, 0)

	for _, runner := range cmds {
		availableSubcommandsArr = append(availableSubcommandsArr, runner.Name())
	}

	availableSubcommands := strings.Join(availableSubcommandsArr, " ")

	if len(args) < 1 {
		runCommand := &run.RunCommand{}
		os.Args = append(os.Args, runCommand.Name())
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown sub-command. Available Commands: %s", availableSubcommands)
}

func main() {
	// Logger Block. Will configure the `slog` logger
	logger.ConfigureLogging()

	if err := root(os.Args[1:]); err != nil {
		slog.Error("Error executing tool", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}
