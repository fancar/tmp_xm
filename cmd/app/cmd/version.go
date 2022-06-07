package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Prints version of %s", appName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
