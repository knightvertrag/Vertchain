package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const flagDataDir = "datadir"

func main() {

	var tbbCmd = &cobra.Command{
		Use:   "dop",
		Short: "The Doge's Pub CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(runCmd())

	err := tbbCmd.Execute()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//Add datadir flag and mark it as required to cmd
func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to the node data dir where the DB will be/is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func incorrectUsageError() error {
	return fmt.Errorf("invalid usage")
}
