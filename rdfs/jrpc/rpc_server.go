package jrpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
)

type Args struct{ A, B int }
type Result int

type Arith int

func (t *Arith) Multiply(args Args, result *Result) error {
	fmt.Printf("[+] RPC Server: Calculating %d * %d\n", args.A, args.B)
	*result = Result(args.A * args.B)
	return nil
}

type Key string
type EDK struct{ Raw string }

type Crypto int

func (t *Crypto) Encrypt(pk Key, edk *EDK) error {
	fmt.Printf("[+] RPC Server: Encrypting DK with PK\n %s\n", pk)
	/*
		Checking transaction process here
	*/

	/*
		Encrpyting decryption key with publick key here
	*/

	edk.Raw = "ENCRYPTED DECRYPTION KEY HERE"

	return nil
}

func InitServer() {
	arith := new(Arith)
	crypto := new(Crypto)

	server := rpc.NewServer()

	server.Register(arith)
	server.Register(crypto)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	listener, err := net.Listen("tcp", ":2085")

	if err != nil {
		fmt.Printf("[-] RPC Server: Listen error: %s\n", err)
	}

	for {
		if conn, err := listener.Accept(); err != nil {
			fmt.Printf("[-] RPC Server: Accept error: %s\n", err.Error())
		} else {
			addr := strings.Split(conn.RemoteAddr().String(), ":")[0]
			fmt.Printf("[+] RPC Server: New connection (%s) established\n", addr)

			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
