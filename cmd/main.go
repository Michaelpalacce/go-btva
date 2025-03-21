package main

import (
	"fmt"

	"github.com/Michaelpalacce/go-btva/internal/args"
)

func main() {
	cliOptions := args.ParseArguments()

	fmt.Println(cliOptions)
}
