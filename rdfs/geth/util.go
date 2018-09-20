package geth

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// Address for geth
const (
	CTRC_COIN = "0xe4b23355e08f71c1a551a3cdb977a2e20d8cc0da"
	CTRC_FILE = "0x9F81F6324B0a2B7A12E45951E17a4E329bac1DED"
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

func (rpc *GethRPC) PersonalUnlockAccount(address string, passphrase string, time uint64) (bool, error) {
	var response bool

	fmt.Printf("%s\n", address)
	err := rpc.call("personal_unlockAccount", &response, address, passphrase, time)
	return response, err
}

func (rpc *GethRPC) EthSendTransaction(from *keystore.Key, to string, data *[]byte, ask bool) (string, error) {
	var response string

	nonce := rpc.getNonce(from.Address.String())
	if nonce == -1 {
		return "", errors.New("couldn't get nonce of current user")
	}

	gasPrice := rpc.getGasPrice()
	if gasPrice == nil {
		return "", errors.New("couldn't get gasPrice of current network")
	}

	gasLimit := rpc.estimateGas(from.Address.String(), data)
	if gasLimit == nil {
		return "", errors.New("couldn't estimate gasLimit. Check the hash again")
	}

	balance := rpc.EthGetBalance(from.Address.String())
	fmt.Printf("[+] %s's ethereum balance: %s\n", from.Address.String(), balance)

	if ask {
		fmt.Printf("[?] Would your like to pay gasLimit: %d, gasPrice: %s? (y/n) ", gasLimit, gasPrice)

		in := bufio.NewReader(os.Stdin)
		input, _ := in.ReadString('\n')

		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		if strings.Compare(input, "y") != 0 && strings.Compare(input, "yes") != 0 {
			return "", errors.New("client canceled sending a transaction")
		}
	}

	rawTx := types.NewTransaction(uint64(nonce), common.HexToAddress(to), big.NewInt(0), gasLimit.Uint64(), gasPrice, *data)
	signTx, err := types.SignTx(rawTx, types.NewEIP155Signer(big.NewInt(4)), from.PrivateKey)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return "", nil
	}

	signTxData, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return "", nil
	}

	err = rpc.call("eth_sendRawTransaction", &response, common.ToHex(signTxData))
	if err != nil {
		fmt.Println(err, "at eth_sendRawTransaction")
	}

	return response, err
}

func (rpc *GethRPC) EthCall(from, to, data string) (json.RawMessage, error) {
	var mParams = map[string]string{
		"from": from,
		"to":   to,
		"data": data,
	}

	response, err := rpc.Call("eth_call", mParams, "latest")

	/*
		err := rpc.call("eth_call", &response, mParams)
		if err != nil {
			fmt.Println(err, "at grpc")
		}
	*/

	return response, err
}

func (rpc *GethRPC) EstimateGas(from string, data *[]byte) *big.Int {
	var mParams = map[string]string{
		"from": from,
		"to":   CTRC_FILE,
		"data": "0x" + hex.EncodeToString(*data),
	}

	var response string

	err := rpc.call("eth_estimateGas", &response, mParams)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return ParseBigInt(response)
}

func (rpc *GethRPC) estimateGas(from string, data *[]byte) *big.Int {
	var mParams = map[string]string{
		"from": from,
		"to":   CTRC_FILE,
		"data": "0x" + hex.EncodeToString(*data),
	}

	var response string

	err := rpc.call("eth_estimateGas", &response, mParams)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return ParseBigInt(response)
}

func (rpc *GethRPC) getNonce(target string) int64 {
	var response string
	err := rpc.call("eth_getTransactionCount", &response, target, "latest")
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return -1
	}

	count, _ := ParseInt(response)

	return int64(count)
}

func (rpc *GethRPC) getGasPrice() *big.Int {
	var response string
	err := rpc.call("eth_gasPrice", &response)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return ParseBigInt(response)
}
