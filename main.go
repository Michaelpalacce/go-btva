package main

import (
	"log/slog"
	"os"

	"github.com/Michaelpalacce/go-btva/internal/orchestrator"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/logger"
	osl "github.com/Michaelpalacce/go-btva/pkg/os"
)

func main() {
	// Logger Block. Will configure the `slog` logger
	logger.ConfigureLogging()

	// Variables block. Init vars

	var (
		o     *orchestrator.Orchestrator
		err   error
		osPtr *osl.OS
		s     *state.State
	)

	// Init Block. Used for fetching and creating needed structs

	if s, err = state.NewState(state.WithDefaultJsonStorage(), state.WithCliArgs()); err != nil {
		slog.Error("Error while loading state.", "err", err)
		os.Exit(1)
	}

	osPtr = osl.GetOS()

	o = orchestrator.NewOrchestrator(osPtr, s, s.Options)

	err = o.Tasks(
		orchestrator.WithAllSoftware(),
		orchestrator.WithFullMinimalInfrastructure(),
	)
	if err != nil {
		slog.Error("Error while adding tasks.", "err", err)
		os.Exit(1)
	}

	err = o.Run()
	if err != nil {
		slog.Error("Error while running tasks.", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}
