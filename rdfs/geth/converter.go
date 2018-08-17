package geth

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// ParseInt parse hex string value to int
func ParseInt(value string) (int, error) {
	i, err := strconv.ParseInt(strings.TrimPrefix(value, "0x"), 16, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// ParseBigInt parse hex string value to big.Int
func ParseBigInt(value string) *big.Int {
	//	cleaned := strings.Replace(value, "0x", "", -1)
	//	result, _ := strconv.ParseUint(cleaned, 16, 64)
	//	i := new(big.Int).SetUint64(result)
	//	return *i
	i := new(big.Int)
	_, err := fmt.Sscan(value, i)
	if err != nil {
		fmt.Println("error scanning value:", err)
	}
	return i
}

// IntToHex convert int to hexadecimal representation
func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}

// BigToHex covert big.Int to hexadecimal representation
func BigToHex(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return "0x0"
	}

	return "0x" + strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0")
}
