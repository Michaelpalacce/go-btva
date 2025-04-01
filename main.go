package main

import (
	"log/slog"

	infra_component "github.com/Michaelpalacce/go-btva/internal/components/infra"
	software_component "github.com/Michaelpalacce/go-btva/internal/components/software"
	"github.com/Michaelpalacce/go-btva/internal/handler"
	"github.com/Michaelpalacce/go-btva/internal/state"
	"github.com/Michaelpalacce/go-btva/pkg/logger"
	osl "github.com/Michaelpalacce/go-btva/pkg/os"
)

func main() {
	// Logger Block. Will configure the `slog` logger
	logger.ConfigureLogging()

	// Variables block. Init vars

	var (
		h     *handler.Handler
		err   error
		osPtr *osl.OS
		s     *state.State
	)

	// Init Block. Used for fetching and creating needed structs

	if s, err = state.NewState(state.WithDefaultJsonStorage(), state.WithCliArgs()); err != nil {
		slog.Error("Error while loading state.", "err", err)
		return
	}

	osPtr = osl.GetOS()

	h = handler.NewHandler(osPtr, s, s.Options)

	// TODO: Add correct tasks based on what is decided by the user
	h.AddTasks(
		software_component.WithAllSoftware(),
		infra_component.WithFullMinimalInfrastructure(),
	)

	h.RunTasks()
}
