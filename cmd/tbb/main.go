package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"vertchain.com/tbb/database"
)

func main() {

	path, _ := os.Getwd()
	if _, err := os.Stat(filepath.Join(path, "database", "tx.db")); os.IsNotExist(err) {
		os.Create(filepath.Join(path, "database", "tx.db"))
	}

	if _, err := os.Stat(filepath.Join(path, "database", "genesis.json")); os.IsNotExist(err) {
		data := map[database.Account]uint{
			"Anurag": 100000,
			"Doge":   100000,
			"Cheems": 100000,
		}
		database.InitializeGenesis(data)
	}

	if _, err := os.Stat(filepath.Join(path, "database", "state.json")); os.IsNotExist(err) {
		os.Create(filepath.Join(path, "database", "state.json"))
	}

	var tbbCmd = &cobra.Command{
		Use:   "tbb",
		Short: "The Blockchain Bar CLI",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tbbCmd.AddCommand(versionCmd)
	tbbCmd.AddCommand(balancesCmd())
	tbbCmd.AddCommand(txCmd())

	err := tbbCmd.Execute()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func incorrectUsageError() error {
	return fmt.Errorf("invalid usage")
}
