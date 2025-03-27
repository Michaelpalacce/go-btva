package main

import (
	"log"
	"log/slog"

	"github.com/Michaelpalacce/go-btva/internal/args"
	"github.com/Michaelpalacce/go-btva/internal/native"
	"github.com/Michaelpalacce/go-btva/pkg/logger"
	"github.com/Michaelpalacce/go-btva/pkg/os"
)

func main() {
	// Logger Block. Will configure the `slog` logger
	logger.ConfigureLogging()

	// Variables block. Init vars

	var (
		handler *native.Handler
		err     error
	)

	// Init Block. Used for fetching and creating needed structs

	opts := args.Args()
	os := os.GetOS()

	if handler, err = native.NewHandler(os, opts); err != nil {
		log.Fatalf("Error creating handler: %v", err)
	}

	// Execution Block. Handles the actual execution of the program

	if err := handler.SetupSoftware(); err != nil {
		slog.Error("Software setup error", "err", err)
	} else {
		slog.Info("Software setup done")
	}

	if err := handler.SetupInfra(); err != nil {
		slog.Error("Infrastructure setup error", "err", err)
	} else {
		slog.Info("Infra setup done")
	}

	if err := handler.SetupLocalEnv(); err != nil {
		slog.Error("Local environment setup error", "err", err)
	} else {
		slog.Info("Local Environemnt setup done")
	}
}
