package geth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"reflect"
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
		ID:      67,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	fmt.Println("params :", params)
	fmt.Println("typeof :")
	for i := range params {
		fmt.Println(reflect.TypeOf(params[i]), params[i])
	}
	fmt.Println("request :", request)
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
	fmt.Println("data : ", data)

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	fmt.Println("resp : ", resp)
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

func (rpc *GethRPC) PersonalUnlockAccount(address string, passphrase string, time int) (bool, error) {
	var response bool
	fmt.Println(address, passphrase, time)
	err := rpc.call("personal_unlockAccount", &response, address, passphrase, time)
	return response, err
}

func (rpc *GethRPC) EthSendTransaction() (string, error) {
	var response string
	var mParams = map[string]string{
		"from": "0x2074fa38f08facdf47f08b8051f9a6aff6033607",
		"to":   "0x58c62f2d8ce3d90d9c61b1117680ac0651a774fa",
		"data": "0x23e3fbd50000000000000000000000002074fa38f08facdf47f08b8051f9a6aff6033607",
	}
	err := rpc.call("eth_sendTransaction", &response, mParams)
	if err != nil {
		fmt.Println(err, "at grpc")
	}
	return response, err
}
