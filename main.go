package main

import (
	"log/slog"

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

	if s, err = state.NewState(state.WithDefaultJsonStorage(true), state.WithCliArgs()); err != nil {
		slog.Error("Error while loading state.", "err", err)
		return
	}

	osPtr = osl.GetOS()

	h = handler.NewHandler(osPtr, s, s.Options)

	// Execution Block. Handles the actual execution of the program

	if err := h.SetupSoftware(); err != nil {
		slog.Error("Software setup error", "err", err)
		return
	}

	if err := h.SetupInfra(); err != nil {
		slog.Error("Infrastructure setup error", "err", err)
		return
	}

	if err := h.SetupLocalEnv(); err != nil {
		slog.Error("Local environment setup error", "err", err)
		return
	}

	if err := h.Final(); err != nil {
		slog.Error("Error while displaying final instructions", "err", err)
		return
	}

	h.Setup(
		func(h *handler.Handler) error {
			return nil
		},
	)
}
