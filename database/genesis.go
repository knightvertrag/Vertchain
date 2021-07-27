package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type genesisFile struct {
	GenesisTime time.Time        `json:"genesis_time"`
	Chain_Id    string           `json:"chain_id"`
	Balances    map[Account]uint `json:"balances"`
}

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func InitializeGenesis(holders map[Account]uint) error {
	path, _ := os.Getwd()
	genesisTime := time.Now()
	os.Create(filepath.Join(path, "database", "genesis.json"))
	balances := holders
	data := genesisFile{
		GenesisTime: genesisTime,
		Chain_Id:    "the-blockchain-bar-ledger",
		Balances:    balances,
	}
	file, err := json.MarshalIndent(data, "", "    ")
	ioutil.WriteFile(filepath.Join(path, "database", "genesis.json"), file, 0644)
	return err
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}
