package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"vertchain.com/tbb/database"
)

const httpPort = 8080

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash             `json:"block_hash"`
	Balances map[database.Account]uint `json:"balaces"`
}

type TxAddReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint   `json:"amount"`
	Data   string `json:"data"`
}

type TxAddRes struct {
	Hash database.Hash `json:"block_hash"`
}

// Start the server with datadir as the database directory
func Run(dataDir string) error {
	fmt.Printf("Listening on HTTP port: %d\n", httpPort)

	state, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	http.HandleFunc("/balances/list", func(rw http.ResponseWriter, r *http.Request) {
		listBalancesHandler(rw, r, state)
	})

	http.HandleFunc("/tx/add", func(rw http.ResponseWriter, r *http.Request) {
		txAddHandler(rw, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

// Handler to write all user balances to response
func listBalancesHandler(rw http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(rw, BalancesRes{state.LatestBlockHash(), state.Balances})
}

// Handler to add transaction and write block hash to response
func txAddHandler(rw http.ResponseWriter, r *http.Request, state *database.State) {
	req := TxAddReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	tx := database.NewTransaction(database.NewAccount(req.From), database.NewAccount(req.To), req.Amount, req.Data)

	err = state.AddTx(tx)
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	hash, err := state.Persist()
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	writeRes(rw, TxAddRes{hash})
}

// Write error to response
func writeErrRes(rw http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(jsonErrRes)
}

// Marshal content to json and write it to response
func writeRes(rw http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(contentJson)
}

// Unmarshal request body and put in reqBody
func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body. %s", err.Error())
	}

	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body. %s", err.Error())
	}

	return nil
}
