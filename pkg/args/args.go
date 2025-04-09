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

type Args struct {
	Usage     string
	Examples  string
	FS        *flag.FlagSet
	Arguments []string

	opts []opt
}

type ArgsOptions func(*Args) error

func NewArgs(arguments []string, options ...ArgsOptions) (*Args, error) {
	args := &Args{
		Arguments: arguments,
	}

	if err := args.Modify(options...); err != nil {
		return nil, err
	}

	return args, nil
}

// Modify is used to Modify the Args with different stuff like examples or usage
func (s *Args) Modify(options ...ArgsOptions) error {
	for _, option := range options {
		if err := option(s); err != nil {
			return err
		}
	}

	return nil
}

func (a *Args) getCorrectFs() *flag.FlagSet {
	if a.FS != nil {
		return a.FS
	} else {
		return flag.CommandLine
	}
}

func (a *Args) AddVar(varPtr any, longhand string, shorthand string, defaultValue any, message string) {
	fs := a.getCorrectFs()

	a.opts = append(a.opts, opt{
		defaultValue:  defaultValue,
		shorthandFlag: shorthand,
		longhandFlag:  longhand,
		description:   message,
	})

	switch v := varPtr.(type) {
	case *string:
		if dv, ok := defaultValue.(string); ok {
			if shorthand != "" {
				fs.StringVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				fs.StringVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been a string", defaultValue))
		}
	case *bool:
		if dv, ok := defaultValue.(bool); ok {
			if shorthand != "" {
				fs.BoolVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				fs.BoolVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been a bool", defaultValue))
		}
	default:
		panic("Var must be a pointer of string or bool")
	}
}

func (a *Args) GetUsage() {
	maxShort := 0
	maxLong := 0

	for _, o := range a.opts {
		if o.shorthandFlag != "" && len(o.shorthandFlag) > maxShort {
			maxShort = len(o.shorthandFlag)
		}
		if o.longhandFlag != "" && len(o.longhandFlag) > maxLong {
			if len(o.longhandFlag) > maxLong {
				maxLong = len(o.longhandFlag)
			}
		}
	}

	fmt.Println(a.Usage)

	if len(a.opts) > 0 {
		fmt.Println("Options:")
	}

	for _, o := range a.opts {
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
	if a.Examples != "" {
		fmt.Println(a.Examples)
	}
}

func (a *Args) Parse() {
	fs := a.getCorrectFs()

	fs.Parse(a.Arguments)
}

func WithUsage(usage string) ArgsOptions {
	return func(a *Args) error {
		a.Usage = usage
		return nil
	}
}

func WithExamples(examples string) ArgsOptions {
	return func(a *Args) error {
		a.Examples = examples
		return nil
	}
}

func WithFs(fs *flag.FlagSet) ArgsOptions {
	return func(a *Args) error {
		a.FS = fs
		fs.Usage = a.GetUsage
		return nil
	}
}
