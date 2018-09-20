package geth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

/*
 * [RDFSCoin Contract]
 * - BalanceOf
 * - AvailableBalanceOf
 * - DepositOf
 *
 */

func (rpc *GethRPC) simpleCall(funcName string, from string, target common.Address) *big.Int {
	data, err := rpc.CoinABI.Pack(funcName, target)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	rawResponse, err := rpc.EthCall(from, CTRC_COIN, "0x"+hex.EncodeToString(data))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	var preResponse string
	json.Unmarshal(rawResponse, &preResponse)

	return ParseBigInt(preResponse)
}

func (rpc *GethRPC) TokenBalance(from string, target common.Address) {
	balance := rpc.BalanceOf(from, target)
	availableBalance := rpc.AvailableBalanceOf(from, target)
	deposit := rpc.DepositOf(from, target)

	if balance != nil && availableBalance != nil && deposit != nil {
		fmt.Printf("[+] %s's  %s(balance) - %s(deposit) = %s(availableBalance)\n", from, balance, deposit, availableBalance)
	}
}

func (rpc *GethRPC) BalanceOf(from string, target common.Address) *big.Int {
	return rpc.simpleCall("balanceOf", from, target)
}

func (rpc *GethRPC) AvailableBalanceOf(from string, target common.Address) *big.Int {
	return rpc.simpleCall("availableBalanceOf", from, target)
}

func (rpc *GethRPC) DepositOf(from string, target common.Address) *big.Int {
	return rpc.simpleCall("depositOf", from, target)
}
