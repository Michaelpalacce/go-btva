package args

import (
	"flag"
	"fmt"
	"strings"
)

type opt struct {
	defaultValue  any
	description   string
	shorthandFlag string
	longhandFlag  string
}

var usage = `go-btva is a tool that allows you to:
- [x] Setup local environment
- [x] Install needed software
- [x] Setup minimal infrastructure on a linux machine
- [x] Connect local env with existing Artifact Manager

Usage:
`

var opts = make([]opt, 0)

func AddNewVar(varPtr any, longhand string, shorthand string, defaultValue any, message string) {
	opts = append(opts, opt{
		defaultValue:  defaultValue,
		shorthandFlag: shorthand,
		longhandFlag:  longhand,
		description:   message,
	})

	switch v := varPtr.(type) {
	case *string:
		if dv, ok := defaultValue.(string); ok {
			if shorthand != "" {
				flag.StringVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				flag.StringVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been a string", defaultValue))
		}
	case *bool:
		if dv, ok := defaultValue.(bool); ok {
			if shorthand != "" {
				flag.BoolVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				flag.BoolVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been a bool", defaultValue))
		}
	default:
		panic("Var must be a pointer of string or bool")
	}
}

func Usage() {
	maxShort := 0
	maxLong := 0

	for _, o := range opts {
		if o.shorthandFlag != "" && len(o.shorthandFlag) > maxShort {
			maxShort = len(o.shorthandFlag)
		}
		if o.longhandFlag != "" && len(o.longhandFlag) > maxLong {
			if len(o.longhandFlag) > maxLong {
				maxLong = len(o.longhandFlag)
			}
		}
	}

	fmt.Println(usage)
	fmt.Println("Options:")
	for _, o := range opts {
		short := ""
		if o.shorthandFlag != "" {
			short = fmt.Sprintf("-%-*s", maxShort, o.shorthandFlag)
		} else {
			short = strings.Repeat(" ", maxShort+1)
		}

		long := ""
		if o.longhandFlag != "" {
			long = fmt.Sprintf("--%-*s", maxLong, o.longhandFlag)
		} else {
			long = strings.Repeat(" ", maxLong+2)
		}

		desc := o.description

		if o.defaultValue != nil && o.defaultValue != "" {
			desc += fmt.Sprintf(" (default: %v)", o.defaultValue)
		}

		fmt.Printf("    %s    %s    %s\n", short, long, desc)
	}
}
