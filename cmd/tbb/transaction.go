package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vertchain.com/tbb/database"
)

const flagFrom = "from"
const flagTo = "to"
const flagAmount = "amount"
const flagData = "data"

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with transactions",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageError()
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	txsCmd.AddCommand(txAddCmd())
	return txsCmd
}

func txAddCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "add",
		Short: "Adds new Transactions to database",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			amount, _ := cmd.Flags().GetUint(flagAmount)
			data, _ := cmd.Flags().GetString(flagData)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTransaction(fromAcc, toAcc, amount, data)

			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			defer state.Close()

			err = state.AddTx(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			_, err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println("Transaction successfully added to ledger.")
		},
	}

	cmd.Flags().String(flagFrom, "", "Account to send tokens from")
	cmd.MarkFlagRequired(flagFrom)

	cmd.Flags().String(flagTo, "", "Account to send tokens to")
	cmd.MarkFlagRequired(flagTo)

	cmd.Flags().Uint(flagAmount, 0, "Amount of tokens to send")
	cmd.MarkFlagRequired(flagAmount)

	cmd.Flags().String(flagData, "", "Possible values: 'reward'")

	return cmd
}
