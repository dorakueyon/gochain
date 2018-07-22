package blockchain

import (
	"encoding/json"
	"time"
)

type Blockchain struct {
	index               int
	chain               []Block
	hash                string
	currentTransactions []Transaction
}

type Block struct {
	Index        int           `json:"index"`
	TimeStamp    int64         `json:"timestamp"`
	Transactions []Transaction `json:"transaction"`
	Proof        Proof         `json:proof`
	PreviousHash string        `json:previoushash`
}

type Transaction struct {
	sender    Sender
	recipient Recipient
	amount    int
}

type Recipient string
type Sender string
type Proof int

func NewBlockchain() *Blockchain {
	b := &Blockchain{
		index:               1,
		currentTransactions: []Transaction{},
	}

	b.createGenesiBlock()

	return b
}

func (b *Blockchain) createGenesiBlock() {
	b.NewBlock(1)
}

func (b *Blockchain) NewTransaction(s Sender, r Recipient, i int) int {
	t := Transaction{
		sender:    s,
		recipient: r,
		amount:    i,
	}
	b.currentTransactions = append(b.currentTransactions, t)
	lb := b.LastBlock()
	return lb.Index + 1
}

func (b *Blockchain) LastBlock() Block {
	return b.chain[len(b.chain)-1]
}

func (b *Blockchain) NewBlock(pf Proof) Block {
	ph := b.PreviousHash()
	block := Block{
		Index:        len(b.chain) + 1,
		TimeStamp:    time.Now().UnixNano(),
		Transactions: b.currentTransactions,
		Proof:        pf,
		PreviousHash: ph,
	}
	b.currentTransactions = []Transaction{}
	b.chain = append(b.chain, block)
	return block
}

func (b *Blockchain) PreviousHash() string {
	var ph string
	length := len(b.chain)
	if length == 0 {
		ph = "100"
	} else {
		block := b.chain[length-1]
		b, _ := json.Marshal(block)
		ph = hash(string(b))
	}
	return ph
}
