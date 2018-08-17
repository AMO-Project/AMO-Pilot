package jrpc

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"

	"rdfs/util"
)

func RequestEDK(destIP []byte, file FileToRequest, info *InfoToReturn) bool {
	destAddr := net.IP(destIP).String()
	client, err := net.Dial("tcp", destAddr+":"+util.JSON_RPC_PORT)

	if err != nil {
		fmt.Printf("[-] RPC Client: Dial error: %s\n", err)
		return false
	}

	c := jsonrpc.NewClient(client)
	defer c.Close()

	err = c.Call("Crypto.Encrypt", file, info)

	if err != nil {
		fmt.Printf("[-] RPC Client: Encryption error: %s\n", err)
		return false
	}

	return true
}

/*
func TestB(addr string, a, b string) bool {
	client, err := net.Dial("tcp", addr+":2085")

	if err != nil {
		fmt.Printf("[-] RPC Client: Dial error: %s\n", err)
		return false
	}

	var key Key = Key(a)
	var result EDK

	c := jsonrpc.NewClient(client)
	defer c.Close()

	err = c.Call("Crypto.Encrypt", key, &result)

	if err != nil {
		fmt.Printf("[-] RPC Client: Encryption error: %s\n", err)
		return false
	}

	//fmt.Printf("[+] RPC Client: %d * %d = %d\n", args.A, args.B, result)

	fmt.Printf("[+] RPC Client: %v\n", result.Raw)

	return true
}

func TestA(addr string, a, b int) bool {
	client, err := net.Dial("tcp", addr+":2085")

	if err != nil {
		fmt.Printf("[-] RPC Client: Dial error: %s\n", err)
		return false
	}

	args := Args{a, b}
	var result int

	c := jsonrpc.NewClient(client)
	defer c.Close()

	err = c.Call("Arith.Multiply", args, &result)

	if err != nil {
		fmt.Printf("[-] RPC Client: Arith error: %s\n", err)
		return false
	}

	fmt.Printf("[+] RPC Client: %d * %d = %d\n", args.A, args.B, result)

	return true
}
*/
