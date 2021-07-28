package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Transaction
	snapshot  Snapshot
	dbFile    *os.File
}

type Snapshot [32]byte

func (s *State) takeSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}
	s.snapshot = sha256.Sum256(data)
	return nil
}

func (s *State) LatestSnapshot() Snapshot {
	return s.snapshot
}

func InitializeState() error {
	current_path, _ := os.Getwd()
	os.Create(filepath.Join(current_path, "database", "state.json"))
	data, _ := loadGenesis(filepath.Join(current_path, "database", "genesis.json"))
	to_write, _ := json.MarshalIndent(data, "", "    ")
	err := ioutil.WriteFile(filepath.Join(current_path, "database", "state.json"), to_write, 0644)
	return err
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
	state := &State{Balances: balances, txMempool: make([]Transaction, 0), dbFile: f, snapshot: Snapshot{}}

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

	err = state.takeSnapshot()
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (s *State) Add(tx Transaction) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

func (s *State) apply(tx Transaction) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) Persist() (Snapshot, error) {
	mempool := make([]Transaction, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			return Snapshot{}, err
		}
		fmt.Println("Persisting new State to disk:")
		fmt.Printf("\t%s\n", txJson)
		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return Snapshot{}, err
		}

		err = s.takeSnapshot()
		if err != nil {
			return Snapshot{}, nil
		}
		fmt.Printf("New DB Snapshot: %x\n", s.snapshot)
		s.txMempool = s.txMempool[1:]
	}
	return s.snapshot, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}
