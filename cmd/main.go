package main

import (
	"log"

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

	softwareChan := make(chan error)
	localEnvChan := make(chan error)
	infraChan := make(chan error)

	go handler.SetupSoftware(softwareChan)
	go handler.SetupLocalEnv(localEnvChan)
	go handler.SetupInfra(infraChan)

	// Result Block. Handles the final result of the program

	// Run 3 times to allow for all 3 chans to return sth
	// @TODO: Set state
	for i := 0; i < 3; i++ {
		select {
		case err := <-softwareChan:
			if err != nil {
				log.Printf("Software setup error: %v", err)
			} else {
				log.Printf("Software setup done")
			}
		case err := <-localEnvChan:
			if err != nil {
				log.Printf("Local environment setup error: %v", err)
			} else {
				log.Printf("Local Environemnt setup done")
			}
		case err := <-infraChan:
			if err != nil {
				log.Printf("Infrastructure setup error: %v", err)
			} else {
				log.Printf("Infra setup done")
			}
		}
	}
}
