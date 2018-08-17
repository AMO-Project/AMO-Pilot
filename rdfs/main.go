package main

import (
	"fmt"
	"os"
	"os/signal"

	"rdfs/crypto"
	"rdfs/geth"
	"rdfs/ipfs"
	"rdfs/jrpc"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/go-ipfs-api"
)

var (
	GETH_PID    int = -1
	GETH_CLIENT *geth.GethRPC
	GETH_KEYS   []*keystore.Key
	IPFS_PID    int = -1
	IPFS_SHELL  *shell.Shell
	IPFS_FILES  ipfs.FileList
	// Address for geth
	CTRC_COIN    string = "0x58c62f2d8ce3d90d9c61b1117680ac0651a774fa"
	CTRC_FILE    string = "0x76174989bc2026d845ec487654839973d88345dc"
	ADDR_ACCOUNT        = map[string]string{
		"ps1": "0x2074fa38f08facdf47f08b8051f9a6aff6033607",
		"ps2": "0x28742aaa4f8a4c6fb31e3a3e3fb85355e3b5926b",
		"ps3": "0x82496a989c83ccd7c58f66934992c3c54f724935",
	}
)

func rdfsInit() {
	fmt.Println("[+] Initializing RDFS")

	/*
	 *	Need to implement pre-processing
	 *	for IPFS, GETH to get initialized
	 */

	GETH_KEYS = crypto.GetKey()

	if len(GETH_KEYS) == 0 {
		fmt.Printf("[-] Couldn't initialize RDFS\n")
		os.Exit(1)
	}

	for _, key := range GETH_KEYS {
		fmt.Printf("[+] Processed the key set of address: %x\n", key.Address.Bytes())
	}

	IPFS_PID = ipfs.Open()
	IPFS_SHELL = shell.NewShell("localhost:5001")
	GETH_PID, GETH_CLIENT = geth.Open()

	go jrpc.InitServer(GETH_KEYS)
}

func rdfsClose() {
	fmt.Println("\n[+] Closing RDFS")

	ipfs.Close(IPFS_PID)
	geth.Close(GETH_PID)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	rdfsInit()

	go func() {
		<-c
		rdfsClose()
		os.Exit(0)
	}()

	prompt()
}
