/*
The `native` package determines which components to fetch based on the OS.

# Purpose

The `native` package is the main driving force of the application as it's responsible for combinig most if not all
components and executing different tasks with them

# Usage

```go
native.NewHandler(os, options)
```

After fetching the native handler, all the os ops must be handled through it
*/
package native
