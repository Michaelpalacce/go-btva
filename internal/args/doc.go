/*
This package is used to parse the command line arguments. Uses the native `flag` package of golang

You can parse them like:

```go
args.ParseArguments()
```
This will return you with a `args.Options` that will contain all the arguments that were either defaulted or specified by the user.

This package will also display usage if called with `--help`
*/
package args
