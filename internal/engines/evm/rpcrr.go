package evmrpcclient

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strings"
	"time"
)

type rpcrequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int64         `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewRequest(method string, params []interface{}) *rpcrequest {
	return &rpcrequest{
		Jsonrpc: "2.0",
		ID:      time.Now().UTC().Unix(),
		Method:  method,
		Params:  params,
	}
}

func (r *rpcrequest) Buffer() *bytes.Buffer {
	payload, _ := json.Marshal(r)
	return bytes.NewBuffer(payload)
}

type rpcresponsebool struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  bool   `json:"result"`
}

func (s *rpcresponsebool) Value() bool {
	return s.Result
}

type rpcresponsevalue struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  string `json:"result"`
}

func (s *rpcresponsevalue) ToInt64() int64 {
	n := new(big.Int)
	n.SetString(strings.Replace(s.Result, "0x", "", -1), 16)
	return n.Int64()
}

type rpcresponseblocktime struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  struct {
		Number    string `json:"number"`
		Timestamp string `json:"timestamp"`
	} `json:"result"`
}

func (v *rpcresponseblocktime) Number() int64 {
	n := new(big.Int)
	n.SetString(strings.Replace(v.Result.Number, "0x", "", -1), 16)
	return n.Int64()
}

func (v *rpcresponseblocktime) Timestamp() int64 {
	n := new(big.Int)
	n.SetString(strings.Replace(v.Result.Timestamp, "0x", "", -1), 16)
	return n.Int64()
}
