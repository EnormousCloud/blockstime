package evmrpcclient

import (
	"encoding/json"
	"fmt"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseResult(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":66,"result":"0xbd1834"}`
	var response rpcresponsevalue
	if err := json.Unmarshal([]byte(input), &response); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, response.ID, int64(66))
	assert.Equal(t, response.Jsonrpc, "2.0")
	assert.Equal(t, response.Result, "0xbd1834")
	assert.Equal(t, response.ToInt64(), int64(12392500))
}

func TestResponseBlocktime(t *testing.T) {
	input := `{"jsonrpc":"2.0","id":66,"result":{"difficulty":"0xcdd7bbccca","extraData":"0x476574682f76312e302e302f6c696e75782f676f312e342e32","gasLimit":"0x1388","gasUsed":"0x0","hash":"0x79b4457c50a8d7bac6370f1ce4e83916b4bfffa5dead578c69afb30d81ddb8af","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd8ae46ba60513f7811dcd26882405f6be01c730f","mixHash":"0x10595c1e193959b302270a711e0b6ade71f499de74b9be48e61f9fc6a2e3eb43","nonce":"0xd251eb38e924a648","number":"0x4349","parentHash":"0x9feed2c3f297609077b2205c4809b28bf82c3305a3f01ee9d7da48ccca364116","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x21b","stateRoot":"0xc580115c8c56e1dc4cbb09af84ed806dcd20e6c88aa5c656c7b653dea7365560","timestamp":"0x55bced4b","totalDifficulty":"0x1c4beefa4f2a21","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}}`
	var response rpcresponseblocktime
	if err := json.Unmarshal([]byte(input), &response); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, response.ID, int64(66))
	assert.Equal(t, response.Jsonrpc, "2.0")
	assert.Equal(t, response.Number(), int64(17225))
	assert.Equal(t, response.Timestamp(), int64(1438444875))

	fmt.Println("time", time.Unix(response.Timestamp(), 0).UTC())
}
