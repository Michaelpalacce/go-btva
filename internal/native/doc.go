/*
The `native` package determines which components to fetch based on the OS.

# Purpose

The `native` package is the main driving force of the application as it's responsible for combinig most if not all
components and executing different tasks with them

# Example Usage

## To create a os agnostic installer

```go

	switch os.Distro {
	    case "linux":
	        handler.installer = &linux.LinuxInstaller{OS: os, Options: options}
	    case "darwin":
	        handler.installer = &darwin.DarwinInstaller{OS:os, Options:options}
	    case "windows":
	        fallthrough
	    default:
	        return nil, fmt.Errorf("OS %s is not supported", os.Distro)
	}

```

## Using the created installer to install java

```go

	err := h.installer.Java().Install()

```

## Using the created installer to check if Java is installed

```go

	exists := h.installer.Java().Exists()

```
*/
package native
