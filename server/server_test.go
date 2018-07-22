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
	b, err := json.Marshal(JsonNewTransactionRequest{
		sender:    "myAddress",
		recipient: "recipAddress",
		amount:    5,
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

	msg := fmt.Sprint("トランザクションはブロック 1 に追加されました")

	if resp.Message != msg {
		t.Errorf("invalid response: %#v", resp)
	}
	t.Logf("%#v", resp)
}

func TestNewMining(t *testing.T) {
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
	//	"message":       "新しいブロックを採掘しました",
	//		"index":         "2",
	//		"transactions":  "transaction",
	//		"proof":         "1",
	//		"previous_hash": "preHash",
	//	}`

	if resp.Message != expectedMsg {
		t.Errorf("invalid response: got=%#v, want=%#v", resp.Message, expectedMsg)
	}
	t.Logf("%#v", resp)
}

func TestHelloWorld(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	env := Env{bc: blockchain.NewBlockchain()}
	HelloWorld(env, res, req)
}
