package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func writeGenesisToDisk(path string) error {
	genesis := map[Account]uint{
		"Anurag": 100000,
		"Doge":   100000,
		"Cheems": 100000,
	}
	current_path, _ := os.Getwd()
	genesisTime := time.Now()
	os.Create(filepath.Join(current_path, "database", "genesis.json"))
	balances := genesis
	data := struct {
		GenesisTime time.Time        `json:"genesis_time"`
		Chain_Id    string           `json:"chain_id"`
		Balances    map[Account]uint `json:"balances"`
	}{
		GenesisTime: genesisTime,
		Chain_Id:    "the-blockchain-bar-ledger",
		Balances:    balances,
	}
	genesisJson, err := json.MarshalIndent(data, "", "    ")
	ioutil.WriteFile(path, genesisJson, 0644)
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
