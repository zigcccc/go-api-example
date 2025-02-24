package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "go-api-example",
	Short: "A simple example application",
}

func Execute() error {
	return RootCmd.Execute()
}