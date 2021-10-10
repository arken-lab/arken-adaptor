package terraadaptor

import (
	"encoding/json"
	"net/http"

	"github.com/arken-lab/arken-adaptor/universaladaptor"
)

const (
	CHAIN_NAME = "TERRA"
)

type ChainAdaptorImpl struct {
	LCDUrl string
}

type ContractResponse interface{}

type LogFilter struct {
}

type TransactionFilter struct {
	ContractAddress string
}

type TxEvent struct {
	Type       string `json:"type"`
	Attributes []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"attributes"`
}

type Log struct {
	MsgIndex int64     `json:"msg_index"`
	Log      string    `json:"log"`
	Events   []TxEvent `json:"events"`
}

type Transaction struct {
	Height    int64  `json:"height"`
	TxHash    string `json:"txhash"`
	Codespace string `json:"codespace"`
	Code      int64  `json:"code"`
	Data      string `json:"data"`
	RawLog    string `json:"raw_log"`
	Logs      []Log  `json:"logs"`
	Info      string `json:"info"`
	GasWanted string `json:"gas_wanted"`
	GasUsed   string `json:"gas_used"`
	Tx        struct {
		TypeUrl string `json:"type_url"`
		Value   string `json:"value"`
	} `json:"tx"`
	Timestamp string `json:"timestamp"`
}

type CosmosBlocksLatestResponse struct {
	Block struct {
		Header struct {
			Height int64 `json:"height"`
		} `json:"header"`
	} `json:"block"`
}

type CosmosTxsResponse struct {
	TxResponses []Transaction `json:"tx_responses"`
}

func NewChainAdaptor(lcdUrl string) universaladaptor.ChainAdaptor {
	c := &ChainAdaptorImpl{
		LCDUrl: lcdUrl,
	}
	return c
}

func (c *ChainAdaptorImpl) GetTransactions(fromBlock int64, toBlock int64, transactionFilters []universaladaptor.TransactionFilter) ([]universaladaptor.Transaction, error) {
	resp, err := http.Get(c.LCDUrl + "/cosmos/tx/v1beta1/txs")
	if err != nil {
		return nil, err
	}
	var data CosmosTxsResponse
	json.NewDecoder(resp.Body).Decode(&data)
	txs := []universaladaptor.Transaction{}
	for _, dataTx := range data.TxResponses {
		txs = append(txs, universaladaptor.Transaction(dataTx))
	}
	return txs, nil
}

func (c *ChainAdaptorImpl) GetLatestBlock() (universaladaptor.BlockNumber, error) {
	resp, err := http.Get(c.LCDUrl + "/blocks/latest")
	if err != nil {
		return 0, err
	}
	var data CosmosBlocksLatestResponse
	json.NewDecoder(resp.Body).Decode(&data)
	return universaladaptor.BlockNumber(data.Block.Header.Height), nil
}

func (c *ChainAdaptorImpl) GetLogs(fromBlock int64, toBlock int64, logFilters []universaladaptor.LogFilter) ([]universaladaptor.Log, error) {
	return []universaladaptor.Log{}, nil
}

func (c *ChainAdaptorImpl) QueryContract(contractQueryData interface{}) (universaladaptor.ContractResponse, error) {
	return nil, nil
}
