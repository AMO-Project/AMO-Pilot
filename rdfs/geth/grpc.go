package geth

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
)

// Address for geth
const (
	CTRC_COIN = "0xda569df17e14ae6a865d299239646c7906d41e10"
	CTRC_FILE = "0x508be8c23bc05e8c3bb56b07a9caa0d364bbbaed"
)

var (
	ADDR_ACCOUNT = map[string]string{
		"ps1": "0x2074fa38f08facdf47f08b8051f9a6aff6033607",
		"ps2": "0x28742aaa4f8a4c6fb31e3a3e3fb85355e3b5926b",
		"ps3": "0x82496a989c83ccd7c58f66934992c3c54f724935",
	}
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
	err := rpc.call("personal_unlockAccount", &response, address, passphrase, time)
	return response, err
}

func (rpc *GethRPC) EthSendTransaction(from, to, data string) (string, error) {
	var response string
	var mParams = map[string]string{
		"from": from,
		"to":   to,
		"data": data,
		"gas":  "0xf4240", // 1,000,000
	}

	err := rpc.call("eth_sendTransaction", &response, mParams)
	if err != nil {
		fmt.Println(err, "at grpc")
	}

	return response, err
}

func (rpc *GethRPC) EthCall(from, to, data string) (json.RawMessage, error) {
	var mParams = map[string]string{
		"from": from,
		"to":   to,
		"data": data,
	}

	response, err := rpc.Call("eth_call", mParams)

	fmt.Printf("%x\n", response)

	/*
		err := rpc.call("eth_call", &response, mParams)
		if err != nil {
			fmt.Println(err, "at grpc")
		}
	*/

	return response, err
}

// function storeRequest(bytes32 _hash, string _name, uint256 _size, bytes4 _ip) returns (bool)
func (rpc *GethRPC) StoreRequest(from string, hash []byte, name string, size *big.Int, ip [4]byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.Abi.Pack("storeRequest", hash32, name, &size, ip)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}
	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully requested storing file '%s'\n", name)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

// function purchaseRequest(bytes32 _hash) returns (bool)
func (rpc *GethRPC) PurchaseRequest(from string, hash []byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.Abi.Pack("purchaseRequest", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully requested purchase '%x'\n", hash32)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

// function purchaseApprove(bytes32 _hash) returns (bool)
func (rpc *GethRPC) PurchaseApprove(from string, hash []byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.Abi.Pack("purchaseApprove", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully approved purchase '%x'\n", hash32)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

// function purchaseAbandon(bytes32 _hash) returns (bool)
func (rpc *GethRPC) PurchaseAbandon(from string, hash []byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.Abi.Pack("purchaseAbandon", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully abandoned purchase '%x'\n", hash32)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

//function getFileName(bytes32 _hash) public view returns (string)
func (rpc *GethRPC) GetFileName(from string, hash []byte) string {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.Abi.Pack("getFileName", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	_, err = rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	return "good"
}

//function getFileSize(bytes32 _hash) public view returns (uint256)
func (rpc *GethRPC) GetFileSize(hash []byte) *big.Int {
	var response string

	return ParseBigInt(response)
}

//function isOwnedBy(bytes32 _hash) public view returns (address)
func (rpc *GethRPC) IsOwnedBy(hash []byte) string {
	var address string

	return address
}

//function isRequested(bytes32 _hash, address buyer) returns (bool)
func (rpc *GethRPC) IsRequested(hash []byte, buyer string) bool {
	var response bool

	return response
}

//function isApproved(bytes32 _hash, address buyer) returns (bool)
func (rpc *GethRPC) IsApproved(hash []byte, buyer string) bool {
	var response bool

	return response
}

//function getOwnerInfo(bytes32 _hash) returns (address, bytes4)
func (rpc *GethRPC) GetOwnerInfo(hash []byte) ([20]byte, [4]byte) {

	return [20]byte{}, [4]byte{}
}
