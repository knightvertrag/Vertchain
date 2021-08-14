package database

type Account string

func NewAccount(value string) Account {
	return Account(value)
}

type Transaction struct {
	From   Account `json:"from"`
	To     Account `json:"to"`
	Amount uint    `json:"Amount"`
	Data   string  `json:"data"`
}

func NewTransaction(from Account, to Account, amount uint, data string) Transaction {
	return Transaction{from, to, amount, data}
}

func (t Transaction) IsReward() bool {
	return t.Data == "reward"
}
