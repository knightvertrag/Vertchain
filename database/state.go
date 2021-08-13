package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Transaction

	dbFile          *os.File
	latestBlockHash Hash
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

	f, err := os.OpenFile(filepath.Join(cwd, "database", "block.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	state := &State{Balances: balances, txMempool: make([]Transaction, 0), dbFile: f, latestBlockHash: Hash{}}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err := json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}
		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) AddTx(tx Transaction) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

func (s *State) apply(tx Transaction) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Amount
		return nil
	}

	if tx.Amount > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Amount
	s.Balances[tx.To] += tx.Amount

	return nil
}

func (s *State) Persist() (Hash, error) {
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{blockHash, block}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	if _, err = s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		return Hash{}, err
	}
	s.latestBlockHash = blockHash

	s.txMempool = []Transaction{}

	return s.latestBlockHash, nil

}

func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}
