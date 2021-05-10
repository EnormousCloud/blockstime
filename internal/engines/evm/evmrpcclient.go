package evmrpcclient

import (
	"blockstime/internal/engines"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type rpcclient struct {
	Node     *engines.NodeRPC
	MaxBlock int64
}

func NewClient(node *engines.NodeRPC) (engines.INetworkEngine, error) {
	// initialize
	res := &rpcclient{
		Node: node,
	}
	if isSyncing, _ := res.IsSyncing(); isSyncing {
		return nil, fmt.Errorf("Node %v is syncing", node.String())
	}
	if lastBlock, err := res.GetHeadBlockNumber(); err != nil {
		return nil, err
	} else {
		res.MaxBlock = lastBlock
	}
	return res, nil
}

func (c *rpcclient) String() string {
	return c.Node.Addr
}

func (c *rpcclient) rpcGet(method string, params []interface{}) ([]byte, error) {
	req, err := http.NewRequest("POST", c.Node.Addr, NewRequest(method, params).Buffer())
	if err != nil {
		return nil, err
	}
	if len(c.Node.HttpPassword) > 0 {
		authHdr := base64.RawStdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", c.Node.HttpUsername, c.Node.HttpPassword)),
		)
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHdr))
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *rpcclient) GetHeadBlockNumber() (int64, error) {
	body, err := c.rpcGet("eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}
	var response rpcresponsevalue
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return 0, err
	}
	return response.ToInt64(), nil
}

func (c *rpcclient) IsSyncing() (bool, error) {
	body, err := c.rpcGet("eth_syncing", []interface{}{})
	if err != nil {
		return false, err
	}
	var response rpcresponsebool
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		// example result: {"currentBlock":"0x4724d0","highestBlock":"0x48976f","knownStates":"0x0","pulledStates":"0x0","startingBlock":"0x451f50"}
		// fmt.Println("body:", string(body))
		return true, err
	}
	return response.Value(), nil
}

func (c *rpcclient) GetBlockTime(blockNumber int64) (int64, error) {
	body, err := c.rpcGet("eth_getBlockByNumber", []interface{}{
		fmt.Sprintf("0x%x", blockNumber), false,
	})
	if err != nil {
		return 0, err
	}
	var response rpcresponseblocktime
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return 0, err
	}
	if response.Number() != blockNumber {
		return 0, fmt.Errorf("blocknumber in response '%d' doesn't match one in request '%d'",
			response.Number(), blockNumber)
	}
	return response.Timestamp(), nil
}
