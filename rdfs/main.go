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
	GETH_CLIENT.CoinABI = geth.CallRDFSCoinABI()
	GETH_CLIENT.FileABI = geth.CallRDFSFileABI()

	go jrpc.InitServer(GETH_KEYS, GETH_CLIENT)
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
