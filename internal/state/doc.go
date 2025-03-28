/*
The `state` package is used to contain the execution state of the program. It will save whatever is given to it and may or may not redirect it.

# How to use?

## Creating a new State object with JsonStorage

```go

	 state := state.NewState()

	// This will also load the state file
	state.Modify(state.WithJsonStorage("file.json", true))

```
*/
package state
