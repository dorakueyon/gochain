package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

const (
	DATA_DIR = "../data/blockchain.data"
)

func InitBlockChain() {
	data := NewBlockchain()
	err := StoreBlockChain(data)
	if err != nil {
		panic(err)
	}
}

func StoreBlockChain(data *Blockchain) error {
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(data)
	if err != nil {
		return err
	}
	filename := fmt.Sprint(DATA_DIR)

	err = ioutil.WriteFile(filename, buffer.Bytes(), 0655)
	if err != nil {
		return err
	}
	return nil

}

func LoadBlockChain(data *Blockchain) error {
	filename := fmt.Sprint(DATA_DIR)
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
