package blockchain

import (
	"encoding/json"
	"testing"
)

func TestFirstBlock(t *testing.T) {
	n := NewBlockchain()
	if n.index != 1 {
		t.Fatalf("index num is wrong.got=%d", n.index)
	}
	tests := []struct {
		proof        Proof
		expectedLeng int
	}{
		{1, 2},
	}
	for _, tt := range tests {
		lengthOfCurrentTransactions := len(n.currentTransactions)
		nb := n.NewBlock(
			tt.proof,
		)
		length := len(n.chain)
		if length != tt.expectedLeng {
			t.Fatalf("length of chain is wrong.got=%d, want=%d. %+v", length, tt.expectedLeng, n.chain)
		}
		if len(nb.Transactions) != lengthOfCurrentTransactions {
			t.Fatalf("length of transactions is wrong.got=%d", len(nb.Transactions))
		}
	}
}

func TestNewTransaction(t *testing.T) {
	b := NewBlockchain()
	index := b.NewTransaction("sender", "recipient", 100)
	expectedVal := 2
	if index != expectedVal {
		t.Fatalf("num is wrong. got=%d, want=%d", index, expectedVal)
	}

	ct := b.currentTransactions[len(b.currentTransactions)-1]

	if ct.sender != "sender" {
		t.Fatalf("expected val is wrong")
	}
	if ct.recipient != "recipient" {
		t.Fatalf("expected val is wrong")
	}

}

func TestValidProof(t *testing.T) {
	tests := []struct {
		pp           Proof
		p            Proof
		expectedBool bool
	}{
		{100, 35293, true},
		{102, 58312, true},
		{102, 58112, false},
	}
	for _, tt := range tests {
		result := validProof(tt.pp, tt.p)
		if result != tt.expectedBool {
			t.Errorf("expected result is wrong.got=%t, proof=%d", result, tt.p)
		}
	}
}

func TestProofOfWork(t *testing.T) {
	bc := NewBlockchain()
	p := bc.ProofOfWork(102)
	if p != 58312 {
		t.Fatalf("ProofOfWork doesn't work properly.got=%d", p)
	}
	p = bc.ProofOfWork(100)
	if p != 35293 {
		t.Fatalf("ProofOfWork doesn't work properly.got=%d", p)
	}

}

func TestHash(t *testing.T) {
	input := Block{
		Index:        1,
		TimeStamp:    1,
		Proof:        1,
		PreviousHash: "100",
	}
	expectedVal := "7f8b201e0065c8a96ff7684c413681235ec286cdb3b33d230afbf5fe2c103540"
	b, _ := json.Marshal(input)
	v := hash(string(b))
	if v != expectedVal {
		t.Fatalf("hash is wrong.got=%s", v)
	}
}
