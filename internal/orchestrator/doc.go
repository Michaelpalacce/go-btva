/*
The `orchestrator` package determines which components to fetch based on the OS.

# Purpose

The `orchestrator` package is the main driving force of the application as it's responsible for combinig most if not all
components and executing different tasks with them

# Usage

```go

	o = orchestrator.NewOrchestrator(osPtr, statePtr)
	err = o.Tasks(orchestrator.WithOptions())

	if err != nil {
		slog.Error("Error while adding tasks.", "err", err)
		os.Exit(1)
	}
	err = o.Run()

	if err != nil {
		slog.Error("Error while running tasks.", "err", err)
		os.Exit(1)
	}

```
*/
package orchestrator
