package main

import (
	"github.com/dorakueyon/gochain/blockchain"
	"github.com/dorakueyon/gochain/server"
)

func main() {
	blockchain.InitBlockChain()
	server.StartHandler()
}
