package main

import (
	"github.com/blackironj/easycommit/cmd"
)

func main() {
	cmd.SetGlobalFlag()

	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
