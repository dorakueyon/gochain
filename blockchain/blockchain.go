package blockchain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Blockchain struct {
	Index               int
	Chain               []Block
	Hash                string
	Nodes               []Node
	CurrentTransactions []Transaction
}

type Block struct {
	Index        int           `json:"index"`
	TimeStamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transaction"`
	Proof        Proof         `json:proof`
	PreviousHash string        `json:previoushash`
}

type Transaction struct {
	Sender    Sender    `json: "sender"`
	Recipient Recipient `json: "recipient"`
	Amount    int       `json:"amount"`
}

type Node string
type Recipient string
type Sender string
type Proof int

type ChainResponse struct {
	Chain  []Block `json:"chain"`
	Length int     `json:length`
}

func NewBlockchain() *Blockchain {
	b := &Blockchain{
		Index:               1,
		CurrentTransactions: []Transaction{},
	}

	b.createGenesiBlock()

	return b
}

func (b *Blockchain) createGenesiBlock() {
	b.NewBlock(1)
}

func (b *Blockchain) NewTransaction(s Sender, r Recipient, i int) int {
	t := Transaction{
		Sender:    s,
		Recipient: r,
		Amount:    i,
	}
	b.CurrentTransactions = append(b.CurrentTransactions, t)
	lb := b.LastBlock()
	return lb.Index + 1
}

func (b *Blockchain) LastBlock() Block {
	return b.Chain[len(b.Chain)-1]
}

func (b *Blockchain) NewBlock(pf Proof) Block {
	ph := b.PreviousHash()
	block := Block{
		Index:        len(b.Chain) + 1,
		TimeStamp:    time.Now().UnixNano(),
		Transactions: b.CurrentTransactions,
		Proof:        pf,
		PreviousHash: ph,
	}
	b.CurrentTransactions = []Transaction{}
	b.Chain = append(b.Chain, block)
	return block
}

func (b *Blockchain) PreviousHash() string {
	var ph string
	length := len(b.Chain)
	if length == 0 {
		ph = "100"
	} else {
		block := b.Chain[length-1]
		//		b, _ := json.Marshal(block)
		//		ph = hash(string(b))
		ph = hash(block)
	}
	return ph
}

func (b *Blockchain) RegisterNode(address string) {
	u, err := url.Parse(address)
	if err != nil {
		log.Fatal(err)
	}
	//hostname := Node(u.Hostname())
	host := Node(u.Host)
	//b.Nodes = append(b.Nodes, hostname)
	b.Nodes = append(b.Nodes, host)
}

func (b *Blockchain) ResolveConflict() bool {
	ns := b.Nodes
	maxLength := len(b.Chain)
	for _, node := range ns {
		url := fmt.Sprintf("http://%s/chain", node)
		fmt.Println(url)
		res, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("invalid code: %d", res.StatusCode)
		}
		resp := ChainResponse{}
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			log.Printf("errpr: %#v, res: %#v", err, res)
		}

		if resp.Length > maxLength && b.ValidChain(resp.Chain) {
			b.Chain = resp.Chain
			return true
		}
	}
	return false
}
