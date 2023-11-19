package rpc

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/iquidus/blockspider/common"
	"github.com/iquidus/blockspider/util"
)

type Config struct {
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}

type RPCClient struct {
	client *rpc.Client
	eth    *ethclient.Client
}

func dialNewClient(cfg *Config) (*rpc.Client, error) {

	var (
		client *rpc.Client
		err    error
	)

	switch cfg.Type {
	case "http", "https":
		if client, err = rpc.DialHTTP(cfg.Endpoint); err != nil {
			return nil, err
		}
	case "unix", "ipc":
		if client, err = rpc.DialIPC(context.Background(), cfg.Endpoint); err != nil {
			return nil, err
		}
	case "ws", "websocket", "websockets":
		if client, err = rpc.DialWebsocket(context.Background(), cfg.Endpoint, ""); err != nil {
			return nil, err
		}
	default:
		fp, err := homedir.Expand(cfg.Endpoint)
		if err != nil {
			return nil, err
		}
		if client, err = rpc.DialIPC(context.Background(), fp); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func NewRPCClient(cfg *Config) *RPCClient {

	client, err := dialNewClient(cfg)
	if err != nil {
		log.Error("could not dial rpc client", "err", err)
		os.Exit(1)
	}
	eth := ethclient.NewClient(client)
	rpcClient := &RPCClient{client, eth}

	return rpcClient
}

func (r *RPCClient) getBlockBy(method string, params ...interface{}) (common.RawBlock, error) {
	var reply common.RawBlock

	err := r.client.Call(&reply, method, params...)

	if err != nil {
		return common.RawBlock{}, err
	}

	return reply, nil
}

func (r *RPCClient) GetLatestBlock() (common.RawBlock, error) {
	bn, err := r.LatestBlockNumber()

	if err != nil {
		return common.RawBlock{}, err
	}

	return r.getBlockBy("eth_getBlockByNumber", util.EncodeUint64(bn))
}

func (r *RPCClient) GetBlockByHeight(height uint64) (common.RawBlock, error) {
	return r.getBlockBy("eth_getBlockByNumber", util.EncodeUint64(height), true)
}

func (r *RPCClient) GetBlockByHash(hash string) (common.RawBlock, error) {
	return r.getBlockBy("eth_getBlockByHash", hash, true)
}

func (r *RPCClient) LatestBlockNumber() (uint64, error) {
	var bn string

	err := r.client.Call(&bn, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	return util.DecodeHex(bn), nil
}

func (r *RPCClient) GetLogs(address []string, hash string, topics []string) ([]common.TxLog, error) {
	var logs []common.TxLog
	err := r.client.Call(&logs, "eth_getLogs", &common.TxLogRequest{
		BlockHash: hash,
		Address:   address,
		Topics:    topics,
	})
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *RPCClient) Ping() (string, error) {
	var version string

	err := r.client.Call(&version, "web3_clientVersion")
	if err != nil {
		return "", err
	}

	return version, nil
}