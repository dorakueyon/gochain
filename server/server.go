package server

import (
	"encoding/json"
	"fmt"
	"github.com/dorakueyon/gochain/blockchain"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	_ "html"
	"log"
	"net/http"
)

const PORT = ":8080"
const MINER = "0"
const MINER_AWARD = 0

type JsonNewTransactionRequest struct {
	sender    string `json:"sender"`
	recipient string `json:"recipient"`
	amount    int    `json:"amount"`
}

type TransactionResponse struct {
	Message string `json: "message"`
}

type MiningResponse struct {
	Message      string                   `json: "message"`
	Index        int                      `json: "index"`
	Transactions []blockchain.Transaction `json:"transactions"`
	Proof        blockchain.Proof         `json:"proof"`
	PreviousHash string                   `json:"previoushash"`
}

func StartHandler() {
	env := Env{}
	env.bc = blockchain.NewBlockchain()
	fmt.Println("Starting Server...")
	router := mux.NewRouter().StrictSlash(true)
	//	router.HandleFunc("/", HelloWorld)
	//router.HandleFunc("/", Handler{env, handler.HelloWorld})
	router.HandleFunc("/transactions/new", newTransaction)
	router.HandleFunc("/mine", mining)
	//http.Handle("/", Handler{env, HelloWorld})
	//http.ListenAndServe("localhost:8080", nil)

	fmt.Printf("http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))

}

func HelloWorld(env *Env, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
	body := JsonNewTransactionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-type", "application/json")

	msg := fmt.Sprintf("トランザクションはブロック %d に追加されました", 1)
	res := TransactionResponse{Message: msg}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)
}

func mining(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	bc := blockchain.NewBlockchain()
	lb := bc.LastBlock()
	lp := lb.Proof
	p := bc.ProofOfWork(lp)

	bc.NewTransaction(MINER, newNodeIdentifier(), MINER_AWARD)
	b := bc.NewBlock(p)
	res := MiningResponse{
		Message:      "新しいブロックを採掘しました",
		Index:        b.Index,
		Transactions: []blockchain.Transaction{},
		Proof:        b.Proof,
		PreviousHash: b.PreviousHash,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)

}

func newNodeIdentifier() blockchain.Recipient {
	id := xid.New()
	return blockchain.Recipient(id.String())
}
