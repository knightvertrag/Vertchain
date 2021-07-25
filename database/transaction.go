package database

type Account string

func NewAccount(value string) Account {
	return Account(value)
}

type Transaction struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func NewTransaction(from Account, to Account, value uint, data string) Transaction {
	return Transaction{from, to, value, data}
}

func (t Transaction) IsReward() bool {
	return t.Data == "reward"
}
