/*
The `state` package is used to contain the execution state of the program. It will save whatever is given to it and may or may not redirect it.

# How to use?

## Creating a new State object with JsonStorage

```go

	if s, err = state.NewState(state.WithDefaultJsonStorage(true), state.WithCliArgs()); err != nil {
		slog.Error("Error while loading state.", "err", err)
		return
	}

```
*/
package state
