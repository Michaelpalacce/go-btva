/*
The `native` package determines which components to fetch based on the OS.

# Purpose

The `native` package is the main driving force of the application as it's responsible for combinig most if not all
components and executing different tasks with them

# Usage

```go
handler, err := native.NewHandler(os, options)

if err != nil { panic(err) }

softwareChan := make(chan error)
localEnvChan := make(chan error)
infraChan := make(chan error)

go handler.SetupSoftware(softwareChan)
go handler.SetupLocalEnv(localEnvChan)
go handler.SetupInfra(infraChan)
```

After fetching the native handler, all the os ops must be handled through it
*/
package native
