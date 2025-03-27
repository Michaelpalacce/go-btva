package main

import (
	"log"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/native"
	"github.com/Michaelpalacce/go-btva/pkg/logger"
	"github.com/Michaelpalacce/go-btva/pkg/os"
	"github.com/Michaelpalacce/go-btva/pkg/state"
)

func main() {
	// Logger Block. Will configure the `slog` logger
	logger.ConfigureLogging()

	// Variables block. Init vars

	var (
		handler *native.Handler
		err     error
		osPtr   *os.OS
		s       *state.State
	)

	// Init Block. Used for fetching and creating needed structs

	if s, err = state.NewState(state.WithDefaultJsonStorage(true)); err != nil {
		slog.Error("Error while loading state.", "err", err)
		return
	}

	if s.Options == nil {
		slog.Info("State file missing or options are not present. Reading arguments.")
		s.Options = args.Args()
	} else {
		slog.Info("State file detected and options loaded. Ignoring arguments passed.")
	}

	osPtr = os.GetOS()

	if handler, err = native.NewHandler(osPtr, s, s.Options); err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	// Execution Block. Handles the actual execution of the program

	if err := handler.SetupSoftware(); err != nil {
		slog.Error("Software setup error", "err", err)
		return
	}

	if err := handler.SetupInfra(); err != nil {
		slog.Error("Infrastructure setup error", "err", err)
		return
	}

	if err := handler.SetupLocalEnv(); err != nil {
		slog.Error("Local environment setup error", "err", err)
		return
	}

	if err := handler.Final(); err != nil {
		slog.Error("Error while displaying final instructions", "err", err)
		return
	}
}
