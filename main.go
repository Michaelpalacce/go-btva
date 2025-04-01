package main

import (
	"log/slog"
	"os"

	infra_component "github.com/Michaelpalacce/go-btva/internal/components/infra"
	software_component "github.com/Michaelpalacce/go-btva/internal/components/software"
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
		h     *orchestrator.Orchestrator
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

	h = orchestrator.NewOrchestrator(osPtr, s, s.Options)

	// TODO: Add correct tasks based on what is decided by the user
	err = h.RunTasks(
		software_component.WithAllSoftware(),
		infra_component.WithFullMinimalInfrastructure(),
	)
	if err != nil {
		slog.Error("Error while running tasks.", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}
