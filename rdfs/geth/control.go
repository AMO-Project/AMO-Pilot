package geth

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"rdfs/util"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type GethRPC struct {
	Url    string
	Client http.Client
	Abi    *abi.ABI
}

func Open() (int, *GethRPC) {
	cmd := exec.Command("geth",
		"--datadir="+util.GETH_DATA_DIR,
		"--networkid=208518",
		"--rpc", "--rpcaddr=0.0.0.0",
		"--rpcapi='db,eth,net,web3,personal,miner,web3'")

	err := cmd.Start()

	fmt.Printf("[+] Initializing 'GETH' ")

	for range make([]int, 3) {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
	}
	fmt.Println()

	if err != nil {
		fmt.Printf("[-] Failed to initialize 'GETH' with '%s'\n", err)
		cmd.Process.Kill()
		os.Exit(1)
	}

	fmt.Printf("[+] Successfully initialized 'GETH' with pid %d\n",
		cmd.Process.Pid)

	gclient := GethRPC{
		Url:    "http://127.0.0.1:8545",
		Client: *http.DefaultClient,
	}

	return cmd.Process.Pid, &gclient
}

func Close(pid int) {
	if pid == -1 {
		return
	}
	syscall.Kill(pid, syscall.SIGTERM)
	fmt.Println("[+] Successfully closed 'GETH'")
}
