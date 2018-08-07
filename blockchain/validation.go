package blockchain

import (
	"fmt"
)

func (b *Blockchain) ValidChain(blocks []Block) bool {
	LastBlock := blocks[0]
	fmt.Printf("len is %d\n", len(blocks))
	var c int = 1
	for c < len(blocks) {
		bc := blocks[c]
		fmt.Println(LastBlock)
		fmt.Println(bc)
		fmt.Println("--------------")

		if bc.PreviousHash != hash(LastBlock) {
			return false
		}

		if ok := validProof(LastBlock.Proof, bc.Proof); !ok {
			return false
		}
		LastBlock = bc

		//count up
		c = c + 1
	}
	return true
}
