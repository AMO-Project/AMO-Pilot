package geth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

/*
 * [RDFSFile Contract]
 * - StoreRequest
 * - PurchaseRequest
 * - PurchaseApprove
 * - PurchaseAbandon
 *
 * - GetFileName
 * - GetFileSize
 * - IsOwnedBy
 * - IsRequested
 * - IsApproved
 * - GetOwnerInfo
 *
 */

// function storeRequest(bytes32 _hash, string _name, uint256 _size, bytes4 _ip) returns (bool)
func (rpc *GethRPC) StoreRequest(from *keystore.Key, hash []byte, name string, size *big.Int, ip [4]byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("storeRequest", hash32, name, &size, ip)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, &data, true)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully requested storing file '%s'\n", name)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

// function purchaseRequest(bytes32 _hash) returns (bool)
func (rpc *GethRPC) PurchaseRequest(from *keystore.Key, hash []byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("purchaseRequest", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, &data, true)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully requested purchase '%x'\n", hash32)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

// function purchaseApprove(bytes32 _hash) returns (bool)
func (rpc *GethRPC) PurchaseApprove(from *keystore.Key, hash []byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("purchaseApprove", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, &data, false)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully approved purchase '%x'\n", hash32)
	fmt.Printf("[+] Transaction Receipt : '%s'\n", receipt)

	return true
}

// function purchaseAbandon(bytes32 _hash) returns (bool)
func (rpc *GethRPC) PurchaseAbandon(from *keystore.Key, hash []byte) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("purchaseAbandon", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	receipt, err := rpc.EthSendTransaction(from, CTRC_FILE, &data, false)
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

	data, err := rpc.FileABI.Pack("getFileName", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	rawName, err := rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	var preName string
	json.Unmarshal(rawName, &preName)

	decodedName, err := hexutil.Decode(preName)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	return string(decodedName)
}

//function getFileSize(bytes32 _hash) public view returns (uint256)
func (rpc *GethRPC) GetFileSize(from string, hash []byte) *big.Int {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("getFileSize", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	rawName, err := rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	var preSize string
	json.Unmarshal(rawName, &preSize)

	return ParseBigInt(preSize)
}

//function isOwnedBy(bytes32 _hash) public view returns (address)
func (rpc *GethRPC) IsOwnedBy(from string, hash []byte) string {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("isOwnedBy", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	rawAddr, err := rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return ""
	}

	var preAddr string
	json.Unmarshal(rawAddr, &preAddr)

	return "0x" + preAddr[26:]
}

//function isRequested(bytes32 _hash, address buyer) returns (bool)
func (rpc *GethRPC) IsRequested(from string, hash []byte, buyer common.Address) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("isRequested", hash32, buyer)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	response, err := rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	var preResult string
	json.Unmarshal(response, &preResult)

	if preResult[len(preResult)-1] != 0x31 {
		return false
	}

	return true
}

//function isApproved(bytes32 _hash, address buyer) returns (bool)
func (rpc *GethRPC) IsApproved(from string, hash []byte, buyer common.Address) bool {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("isApproved", hash32, buyer)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	response, err := rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	var preResult string
	json.Unmarshal(response, &preResult)

	if preResult[len(preResult)-1] != 0x31 {
		return false
	}

	return true
}

//function getOwnerInfo(bytes32 _hash) returns (address, bytes4)
func (rpc *GethRPC) GetOwnerInfo(from string, hash []byte) (string, [4]byte) {
	hash32 := [32]byte{}
	copy(hash32[:], hash)

	data, err := rpc.FileABI.Pack("getOwnerInfo", hash32)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return "", [4]byte{}
	}

	response, err := rpc.EthCall(from, CTRC_FILE, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return "", [4]byte{}
	}

	var preResponse string

	json.Unmarshal(response, &preResponse)

	if len(preResponse) != 130 {
		fmt.Printf("[-] Error occured: %s\n", "no record for purchase request.")
		return "", [4]byte{}
	}

	ownerAddr := string("0x" + preResponse[26:66])
	ownerIP, err := hex.DecodeString(preResponse[66:74])
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return "", [4]byte{}
	}

	ownerIP4 := [4]byte{}
	copy(ownerIP4[:], ownerIP)

	return ownerAddr, ownerIP4
}
