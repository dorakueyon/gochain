package blockchain

import (
	"fmt"
	"testing"
)

func TestFirstBlock(t *testing.T) {
	n := NewBlockchain()
	if n.Index != 1 {
		t.Fatalf("index num is wrong.got=%d", n.Index)
	}
	tests := []struct {
		proof        Proof
		expectedLeng int
	}{{1, 2}}
	for _, tt := range tests {
		lengthOfCurrentTransactions := len(n.CurrentTransactions)
		nb := n.NewBlock(
			tt.proof,
		)
		length := len(n.Chain)
		if length != tt.expectedLeng {
			t.Fatalf("length of chain is wrong.got=%d, want=%d. %+v", length, tt.expectedLeng, n.Chain)
		}
		if len(nb.Transactions) != lengthOfCurrentTransactions {
			t.Fatalf("length of transactions is wrong.got=%d", len(nb.Transactions))
		}
	}
}

func TestNewTransactions(t *testing.T) {
	b := NewBlockchain()
	index := b.NewTransaction("sender1", "recipient1", 100)
	expectedVal := 2
	if index != expectedVal {
		t.Fatalf("num is wrong. got=%d, want=%d", index, expectedVal)
	}

	ct := b.CurrentTransactions[len(b.CurrentTransactions)-1]

	if ct.Sender != "sender1" {
		t.Fatalf("expected val is wrong")
	}
	if ct.Recipient != "recipient1" {
		t.Fatalf("expected val is wrong")
	}

	index = b.NewTransaction("sender2", "recipient2", 200)
	expectedVal = 2
	if index != expectedVal {
		t.Fatalf("num is wrong. got=%d, want=%d", index, expectedVal)
	}
	lb := b.LastBlock()
	p := b.ProofOfWork(lb.Proof)
	b.NewBlock(p)
	index = b.NewTransaction("sender3", "recipient3", 500)
	expectedVal = 3
	if index != expectedVal {
		t.Fatalf("num is wrong. got=%d, want=%d", index, expectedVal)
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
	v := hash(input)
	if v != expectedVal {
		t.Fatalf("hash is wrong.got=%s", v)
	}
}

func TestStoreAndLoadBlockChain(t *testing.T) {
	data := NewBlockchain()
	err := StoreBlockChain(data)
	if err != nil {
		t.Fatalf("Store got error. err: %s", err)
	}

	bc := &Blockchain{}
	err = LoadBlockChain(bc)
	if err != nil {
		t.Fatalf("Load got error. err: %s", err)
	}
	if bc.Index != data.Index {
		t.Fatalf("Store and Load got error.")
	}
}

func TestAddNodeLists(t *testing.T) {
	bc := NewBlockchain()
	tests := []struct {
		input       string
		expectedVal Node
	}{
		{"http://192.168.0.5:5000", "192.168.0.5:5000"},
		{"http://192.168.0.3:5000/chain", "192.168.0.3:5000"},
		{"http://www.google.com/serch", "www.google.com"},
	}
	for _, tt := range tests {
		bc.RegisterNode(tt.input)
		addedNode := bc.Nodes[len(bc.Nodes)-1]
		if addedNode != tt.expectedVal {
			t.Fatalf("register node failed.got=%v, want=%v", addedNode, tt.expectedVal)
		}

	}
}

func TestValidChain(t *testing.T) {
	bc := NewBlockchain()
	lb := bc.LastBlock()
	lastProof := lb.Proof
	p := bc.ProofOfWork(lastProof)
	bc.NewBlock(p)

	ok := bc.ValidChain(bc.Chain)
	if !ok {
		t.Fatalf("validation failed")
	}

	lb = bc.LastBlock()
	lastProof = lb.Proof
	p = bc.ProofOfWork(lastProof)
	bc.NewBlock(p)

	ok = bc.ValidChain(bc.Chain)
	if !ok {
		t.Fatalf("validation failed")
	}

	bc2 := NewBlockchain()
	randamProof := Proof(1234)
	bc2.NewBlock(randamProof)
	ok = bc2.ValidChain(bc2.Chain)
	if ok {
		t.Fatalf("validation should be failed")
	}

}

func TestResolveConflict(t *testing.T) {
	bc := NewBlockchain()
	bc.RegisterNode("http://localhost:8082")
	fmt.Println(bc.Chain)
	ok := bc.ResolveConflict()
	if !ok {
		t.Fatalf("should be false")
	}
	fmt.Println(bc.Chain)
}
