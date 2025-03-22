/*
This package is used to parse the app configuration/options. Uses the native `flag` package of golang for CLI parsing

# Supported Sources

- [x] CLI
- [ ] File
- [ ] Web

# Usage

You can parse them like:

```go
args.Args()
```
This will return you with a `args.Options` that will contain all the arguments that were either defaulted or specified by the user.

This package will also display usage if called with `--help`
*/
package args
