package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dorakueyon/gochain/blockchain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTransactionPost(t *testing.T) {
	blockchain.InitBlockChain()
	testFirstNewTransaction(t)
	testSecondNewTransaction(t)
	testNewTransactionWithoutValue(t)
}

func testFirstNewTransaction(t *testing.T) {
	b, err := json.Marshal(JsonNewTransactionRequest{
		Sender:    "myAddress",
		Recipient: "recipAddress", Amount: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/transactions/new", bytes.NewBuffer(b))
	res := httptest.NewRecorder()

	newTransaction(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	resp := TransactionResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Errorf("errpr: %#v, res: %#v", err, res)
	}

	msg := fmt.Sprint("トランザクションはブロック 2 に追加されました")

	if resp.Message != msg {
		t.Errorf("invalid response: %#v, want=%#v", resp, msg)
	}
	t.Logf("%#v", resp)
}

func testSecondNewTransaction(t *testing.T) {
	b, err := json.Marshal(JsonNewTransactionRequest{
		Sender:    "myAddress2",
		Recipient: "recipAddress2",
		Amount:    10,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/transactions/new", bytes.NewBuffer(b))
	res := httptest.NewRecorder()

	newTransaction(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	resp := TransactionResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Errorf("errpr: %#v, res: %#v", err, res)
	}

	msg := fmt.Sprint("トランザクションはブロック 2 に追加されました")

	if resp.Message != msg {
		t.Errorf("invalid response: %#v, want=%#v", resp, msg)
	}
	t.Logf("%#v", resp)
}

func testNewTransactionWithoutValue(t *testing.T) {
	b, err := json.Marshal(JsonNewTransactionRequest{
		Sender:    "",
		Recipient: "recipAddress",
		Amount:    5,
	})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/transactions/new", bytes.NewBuffer(b))
	res := httptest.NewRecorder()

	newTransaction(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("invalid code: %d", res.Code)
	}

	t.Logf("%#v", res)
}

func TestNewMining(t *testing.T) {
	blockchain.InitBlockChain()
	req := httptest.NewRequest("GET", "/mine", nil)
	res := httptest.NewRecorder()

	mining(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
	resp := MiningResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Errorf("errpr: %#v, res: %#v", err, res)
	}

	var expectedMsg string = "新しいブロックを採掘しました"
	var expectedIndex int = 2

	if resp.Message != expectedMsg {
		t.Errorf("invalid response: got=%#v, want=%#v", resp.Message, expectedMsg)
	}
	if resp.Index != expectedIndex {
		t.Errorf("invalid response: got=%#v, want=%#v", resp.Index, expectedIndex)
	}

	t.Logf("%#v", resp)

	testStoredBlockchain(t)
}

func testStoredBlockchain(t *testing.T) {
	bc := &blockchain.Blockchain{}
	err := blockchain.LoadBlockChain(bc)
	if err != nil {
		t.Fatalf("laod is wrong. error:%s", err)
	}
}

func TestChain(t *testing.T) {
	blockchain.InitBlockChain()
	testFirstChainCall(t)

	req := httptest.NewRequest("GET", "/mine", nil)
	res := httptest.NewRecorder()
	mining(res, req)

	testSecondChainCall(t)
}

func testFirstChainCall(t *testing.T) {
	req := httptest.NewRequest("GET", "/chain", nil)
	res := httptest.NewRecorder()

	fullChain(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}
	resp := ChainResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Errorf("errpr: %#v, res: %#v", err, res)
	}
	if resp.Length != 1 {
		t.Errorf("length is wrong. got=%d, want=%d", resp.Length, 1)
	}
	if resp.Chain[0].Index != 1 {
		t.Errorf("val is wrong. got=%d", resp.Chain[0].Index)
	}
	if resp.Chain[0].Proof != 1 {
		t.Errorf("val is wrong. got=%d", resp.Chain[0].Proof)
	}
	if resp.Chain[0].PreviousHash != "100" {
		t.Errorf("val is wrong. got=%v", resp.Chain[0].PreviousHash)
	}

}

func testSecondChainCall(t *testing.T) {
	req := httptest.NewRequest("GET", "/chain", nil)
	res := httptest.NewRecorder()

	fullChain(res, req)

	resp := ChainResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Errorf("errpr: %#v, res: %#v", err, res)
	}
	if resp.Length != 2 {
		t.Errorf("length is wrong. got=%d, want=%d", resp.Length, 2)
	}
	t.Logf("%#v", resp)
}

func TestRegisterNode(t *testing.T) {
	tests := []struct {
		input              []string
		expectedStatusCode int
		expectedVal        []string
		expectedMessage    string
	}{
		{[]string{"http://192.168.0.5:5000"}, 201, []string{"192.168.0.5:5000"}, "新しいノードが追加されました"},
		{[]string{"http://192.168.0.5:5000", "http://192.168.0.6:5000", "http://google.com/test"}, 201, []string{"192.168.0.5:5000", "192.168.0.6:5000", "google.com"}, "新しいノードが追加されました"},
		//{[]string{}, 400, []string{}, "有効でないノードのリストです"},
		//{[]string{"http://192.168.0.5:5000/test"}, 400, "", "有効でないノードのリストです"},
	}
	for _, tt := range tests {
		blockchain.InitBlockChain()
		b, err := json.Marshal(JsonNodeRequest{
			Node: tt.input,
		})
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest("POST", "/nodes/register", bytes.NewBuffer(b))
		res := httptest.NewRecorder()

		registerNodes(res, req)

		if res.Code != tt.expectedStatusCode {
			t.Errorf("invalid code: got=%d, want=%d", res.Code, tt.expectedStatusCode)
		}
		resp := NodeRegisterResponse{}
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			t.Fatalf("errpr: %#v, res: %#v", err, res)
		}
		if resp.Message != tt.expectedMessage {
			t.Errorf("message is wrong. got=%s, want=%s", resp.Message, tt.expectedMessage)
		}

		if len(resp.TotalNodes) != len(tt.expectedVal) {
			t.Fatalf("len is wrong. got=%d, want=%d", len(resp.TotalNodes), len(tt.expectedVal))
		}
		for i, node := range resp.TotalNodes {
			if node != tt.expectedVal[i] {
				t.Errorf("total node list is wrong. got=%s, want=%s", node, tt.expectedVal[i])
			}
		}
		t.Logf("%#v", resp)
	}
}

func TestConsensus(t *testing.T) {
	blockchain.InitBlockChain()
	req := httptest.NewRequest("GET", "/nodes/resolve", nil)
	res := httptest.NewRecorder()

	consensus(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("invalid code: got=%d", res.Code)
	}

	resp := ConsensusResponse{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Errorf("errpr: %#v, res: %#v", err, res)
	}

	if resp.Message != "チェーンが確認されました" {
		t.Errorf("Message is wrong. got=%s", resp.Message)
	}
	t.Logf("%#v", resp)

}
