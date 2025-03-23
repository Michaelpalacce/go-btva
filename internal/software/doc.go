/*
Software package contains:
- Software managers based on the OS and the software to install

# Example Usage

## To create a os agnostic installer

```go

	switch os.Distro {
	case "linux":
		handler.installer = &linux.LinuxInstaller{OS: os, Options: options}
	case "windows":
		fallthrough
	case "darwin":
		fallthrough
	default:
		return nil, fmt.Errorf("OS %s is not supported", os.Distro)
	}

```

## Using the created installer to install java

```go

	err := h.installer.Java().Install()

```

## Using the created installer to remove java

```go

	err := h.installer.Java().Remove()

```

## Using the created installer to check if Java is installed

```go

	exists := h.installer.Java().Exists()

```

@TODO: In the future, this package will be able to work with local installation files
*/
package software
