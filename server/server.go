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

const PORT = ":8082"
const MINER = "0"
const MINER_AWARD = 1

type JsonNewTransactionRequest struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

type NodeRegisterResponse struct {
	Message    string   `json:"message"`
	TotalNodes []string `json:"total_nodes"`
}

type JsonNodeRequest struct {
	Node []string `json:"nodes"`
}

type TransactionResponse struct {
	Message string `json: "message"`
}

type ChainResponse struct {
	Chain  []blockchain.Block `json:"chain"`
	Length int                `json:length`
}

type MiningResponse struct {
	Message      string                   `json: "message"`
	Index        int                      `json: "index"`
	Transactions []blockchain.Transaction `json:"transactions"`
	Proof        blockchain.Proof         `json:"proof"`
	PreviousHash string                   `json:"previoushash"`
}

type ConsensusResponse struct {
	Message string             `json: "message"`
	Chain   []blockchain.Block `json:"chain"`
}

func StartHandler() {
	fmt.Println("Starting Server...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HelloWorld)
	router.HandleFunc("/transactions/new", newTransaction)
	router.HandleFunc("/mine", mining)
	router.HandleFunc("/chain", fullChain)
	router.HandleFunc("/nodes/register", registerNodes)
	router.HandleFunc("/nodes/resolve", consensus)

	fmt.Printf("http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
	bcn := &blockchain.Blockchain{}
	err := blockchain.LoadBlockChain(bcn)
	if err != nil {
		panic(err)
	}

	body := JsonNewTransactionRequest{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if body.Sender == "" || body.Recipient == "" || body.Amount == 0 {
		http.Error(w, "値がありません", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	index := bcn.NewTransaction(
		blockchain.Sender(body.Sender),
		blockchain.Recipient(body.Recipient),
		body.Amount,
	)
	err = blockchain.StoreBlockChain(bcn)
	if err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("トランザクションはブロック %d に追加されました", index)
	res := TransactionResponse{Message: msg}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)
}

func mining(w http.ResponseWriter, r *http.Request) {
	bcn := &blockchain.Blockchain{}
	err := blockchain.LoadBlockChain(bcn)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-type", "application/json")
	lb := bcn.LastBlock()
	lp := lb.Proof
	p := bcn.ProofOfWork(lp)

	bcn.NewTransaction(MINER, newNodeIdentifier(), MINER_AWARD)
	b := bcn.NewBlock(p)
	err = blockchain.StoreBlockChain(bcn)
	if err != nil {
		panic(err)
	}

	res := MiningResponse{
		Message:      "新しいブロックを採掘しました",
		Index:        b.Index,
		Transactions: b.Transactions,
		Proof:        b.Proof,
		PreviousHash: b.PreviousHash,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)

}

func newNodeIdentifier() blockchain.Recipient {
	guid := xid.New()
	return blockchain.Recipient(guid.String())
}

func fullChain(w http.ResponseWriter, r *http.Request) {
	bcn := &blockchain.Blockchain{}
	err := blockchain.LoadBlockChain(bcn)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-type", "application/json")
	len := len(bcn.Chain)
	res := ChainResponse{
		Chain:  bcn.Chain,
		Length: len,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)
}

func registerNodes(w http.ResponseWriter, r *http.Request) {
	bcn := &blockchain.Blockchain{}
	err := blockchain.LoadBlockChain(bcn)
	if err != nil {
		panic(err)
	}
	body := JsonNodeRequest{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body.Node) == 0 {
		http.Error(w, "有効でないノードのリストです", http.StatusBadRequest)
		return
	}
	for _, node := range body.Node {
		bcn.RegisterNode(node)
	}
	err = blockchain.StoreBlockChain(bcn)
	if err != nil {
		panic(err)
	}

	nl := []string{}
	for _, node := range bcn.Nodes {
		nl = append(nl, string(node))
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	res := NodeRegisterResponse{
		Message:    "新しいノードが追加されました",
		TotalNodes: nl,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)

}

func consensus(w http.ResponseWriter, r *http.Request) {
	bcn := &blockchain.Blockchain{}
	err := blockchain.LoadBlockChain(bcn)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-type", "application/json")
	res := ConsensusResponse{}
	ok := bcn.ResolveConflict()
	if ok {
		res.Message = "チェーンが置き換えられました"
		res.Chain = []blockchain.Block{}

	} else {
		res.Message = "チェーンが確認されました"
		res.Chain = []blockchain.Block{}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&res)

}
