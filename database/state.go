package database

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Transaction

	dbFile *os.File
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gen, err := loadGenesis(filepath.Join(cwd, "database", "genesis.json"))
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}
	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	state := &State{Balances: balances, txMempool: make([]Transaction, 0), dbFile: f}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Transaction
		json.Unmarshal(scanner.Bytes(), &tx)
		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}
	return state, nil
}

func (s *State) apply(tx Transaction) error {

}
