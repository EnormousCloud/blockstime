package engines

import (
	"errors"
	"os"
)

type NodeRPC struct {
	// HTTP/HTTPS/IRC Address of the source to dial to
	Addr string `json:"addr" yaml:"addr"`
	// User name - in case of HTTP/HTTPS Basic Auth protection
	HttpUsername string `json:"username" yaml:"username"`
	// User password - in case of HTTP/HTTPS Basic Auth protection
	HttpPassword string `json:"password" yaml:"password"`
}

func NewTestNode() (*NodeRPC, error) {
	addr := os.Getenv("HTTP_ADDR")
	username := os.Getenv("HTTP_USERNAME")
	password := os.Getenv("HTTP_PASSWORD")
	if len(addr) == 0 {
		return nil, errors.New("HTTP_ADDR must be specified")
	}
	return &NodeRPC{
		Addr:         addr,
		HttpUsername: username,
		HttpPassword: password,
	}, nil
}

type INetworkEngine interface {
	IsSyncing() (bool, error)
	GetHeadBlockNumber() (int64, error)
	GetBlockTime(blockNumber int64) (int64, error)
	String() string
}
