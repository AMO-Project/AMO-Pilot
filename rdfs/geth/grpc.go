package geth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/big"
)

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (rpc *GethRPC) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// Call returns raw response of method call
func (rpc *GethRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.Client.Post(rpc.Url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	//	if rpc.Debug
	//		rpc.log.Println(fmt.Sprintf("%s\nRequest: %s\nResponse: %s\n", method, body, data))
	//	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp.Result, nil

}
func (rpc *GethRPC) NetVersion() (string, error) {
	var netVersion string

	err := rpc.call("net_version", &netVersion)
	return netVersion, err
}

func (rpc *GethRPC) EthProtocolVersion() (string, error) {
	var protocolVersion string

	err := rpc.call("eth_protocolVersion", &protocolVersion)
	return protocolVersion, err
}

func (rpc *GethRPC) EthCoinBase() (string, error) {
	var address string

	err := rpc.call("eth_coinbase", &address)
	return address, err
}

func (rpc *GethRPC) EthMining() (bool, error) {
	var mining bool

	err := rpc.call("eth_mining", &mining)
	return mining, err
}

func (rpc *GethRPC) EthBlockNumber() (int, error) {
	var response string
	if err := rpc.call("eth_blockNumber", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

func (rpc *GethRPC) EthAccounts() ([]string, error) {
	accounts := []string{}

	err := rpc.call("eth_accounts", &accounts)
	return accounts, err
}

func (rpc *GethRPC) EthGetBalance(address string) *big.Int {
	var response string
	if err := rpc.call("eth_getBalance", &response, address, "latest"); err != nil {
		return big.NewInt(0)
	}
	return ParseBigInt(response)
}
