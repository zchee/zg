package main

import "fmt"

var cmdVersion = &Command{
	Usage: "version",
	Short: "Print the version",
	Long: `
Print the version
`,
	Run: runVersion,
}

var (
	Version   = "0.0.1"
	GitCommit = "HEAD"
)

func runVersion(cmd *Command, args []string) {
	fmt.Println("version: " + Version + " (" + GitCommit + ")")
}
