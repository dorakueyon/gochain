package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
)

func (b *Blockchain) ProofOfWork(lastProof Proof) Proof {
	var p Proof = 0
	for !validProof(lastProof, p) {
		p += 1

	}
	return p
}

func validProof(lastProof, proof Proof) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	h := hex.EncodeToString(guessHash[:])
	if h[:4] == "0000" {
		return true
	}
	return false
}

func hash(b Block) string {
	s, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	guessHash := sha256.Sum256([]byte(s))
	hexedHash := hex.EncodeToString(guessHash[:])
	return hexedHash
}
