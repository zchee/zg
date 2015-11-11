package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"
)

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'zg help'.
var commands = []*Command{
	cmdAdd,
	cmdVersion,
}

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// Usage is the one-line usage message.
	// The first word in the line is taken to be the command name.
	Usage string

	// Short is the short description shown in the 'godep help' output.
	Short string

	// Long is the long message shown in the
	// 'godep help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Usage = usageExit
	flag.Parse()
	log.SetFlags(0)
	log.SetPrefix("%zg: ")

	args := flag.Args()
	if len(args) < 1 {
		usageExit()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.UsageExit() }
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s: unknown command %q\n", os.Args[0], args[0])
	fmt.Fprintf(os.Stderr, "Run '%s help' for usage.\n", os.Args[0])
	os.Exit(2)
}

// Name returns the name of a command.
func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: zg help command\n\n")
		fmt.Fprintf(os.Stderr, "Too many arguments given.\n")
		os.Exit(2)
	}
	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}
}

// UsageExit prints usage information and exits.
func (c *Command) UsageExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s %s\n\n", os.Args[0], c.Usage)
	fmt.Fprintf(os.Stderr, "Run '%s help %s' for help.\n", os.Args[0], c.Name())
	os.Exit(2)
}

func usageExit() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
	})
	template.Must(t.Parse(strings.TrimSpace(text) + "\n\n"))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var usageTemplate = `
zg is the new z, yo

Usage:
    zg command [arguments]

The commands are:{{range .}}
    {{.Name | printf "%-8s"}} {{.Short}}{{end}}

Use "zg help [command]" for more information about a command.
`

var helpTemplate = `
Usage: zg {{.Usage}}

{{.Long | trim}}
`
