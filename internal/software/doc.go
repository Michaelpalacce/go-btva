/*
Software package contains:
- Installer that accepts an *os.OS and is used to call the correct software installers
- Installers for software based on OS

# How to use?

# Getting an installer based on your OS

```go

	os := os.GetOS()
	installer := sofware.NewInstaller(os)
	installer.installJava()

```
*/
package software
