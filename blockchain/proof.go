package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	h := hash(guess)
	if h[:4] == "0000" {
		return true
	}
	return false
}

func hash(s string) string {
	guessHash := sha256.Sum256([]byte(s))
	hexedHash := hex.EncodeToString(guessHash[:])
	return hexedHash
}
