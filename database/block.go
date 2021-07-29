package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte

func (h Hash) Encode() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) Decode(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Block struct {
	Header BlockHeader   `json:"header"`
	TXs    []Transaction `json:"payload"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(parent Hash, time uint64, txs []Transaction) Block {
	return Block{BlockHeader{parent, time}, txs}
}

func (block Block) Hash() (Hash, error) {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return Hash{}, nil
	}

	return sha256.Sum256(blockJSON), nil
}
