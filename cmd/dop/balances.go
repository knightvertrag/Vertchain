package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vertchain.com/dop/database"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageError()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	balancesCmd.AddCommand(balancesListCmd())
	return balancesCmd
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances",
		Run: func(cmd *cobra.Command, args []string) {
			dataDir, _ := cmd.Flags().GetString(flagDataDir)
			state, err := database.NewStateFromDisk(dataDir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Printf("Account balances at: %x\n", state.LatestBlockHash())
			fmt.Println("_______________")
			fmt.Println("")
			for account, balance := range state.Balances {
				fmt.Printf("%s: %d\n", account, balance)
			}
		},
	}
	return balancesListCmd
}
