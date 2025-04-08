/*
The `state` package is used to contain the execution state of the program. It will save whatever is given to it and may or may not redirect it.

# How to use?

## Creating a new State object with JsonStorage

```go

	if s, err = state.NewState(state.WithDefaultJsonStorage(true)); err != nil {
		slog.Error("Error while loading state.", "err", err)
		return
	}

```

## Creating a state object with no storage, but loading options from CLI arguments

```go

	if s, err = state.NewState(state.WithCliArgs()); err != nil {
		slog.Error("Error while loading state.", "err", err)
		os.Exit(1)
	}

```
*/
package state
