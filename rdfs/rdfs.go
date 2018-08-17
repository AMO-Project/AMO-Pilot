package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"rdfs/cryp"
	"rdfs/geth"
	"rdfs/ipfs"
	"rdfs/jrpc"
	"rdfs/util"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/go-ipfs-api"
)

var (
	geth_pid    int = -1
	ipfs_pid    int = -1
	geth_client *geth.GethRPC
	geth_keys   []*keystore.Key
	ipfs_shell  *shell.Shell
)

func purchase(args ...string) bool {
	/*
		As following the File Purchase Scenario,

		1. EF = ipfs.Get(EFH)
		2. Send token, EF to contract
			- (Contract) Check token amount and transfer token to File Owner
			- (Contract) Return true to buyer
		3. If accepted, receive true from owner with his info
		4. Send buyer's pk(publicKey) to owner's json-rpc(?) server
			- Check the transaction
			- If valid, edk = E(dk, pk) and send edk to buyer's json-rpc server
		5. Receive edk and decrypt EF
			- dk = D(edk, sk)	*sk = privateKey
			- F = D(EF, dk)

	*/
	var dir_path string = util.RDFS_DOWN_DIR
	if len(args) == 2 {
		dir_path = args[1]
	}

	// 1. EF = ipfs.Get(EFH)
	if !ipfs.Get(ipfs_shell, args[0], dir_path) {
		return false
	}

	// 2. Send token, EF to contract
	fmt.Printf("[+] Send token, EF to Ethereum through GETH\n")

	// 3. If accepted, receive true from owner with his info
	fmt.Printf("[+] Successfully made transaction with the owner\n")
	ownerIP := "127.0.0.1"

	// 4. Send buyer's pk(publickKey) to owner's json-rpc(?) server
	var edk jrpc.EDK
	if !jrpc.RequestE(ownerIP, jrpc.Key("BUYER KEY"), &edk) {
		fmt.Printf("[-] Couldn't receive EDK from owner\n")
		return false
	}

	// 5. Receive edk and decrypt EF
	fmt.Printf("[+] RPC Client: Successfully received EDK from the owner\n"+
		":%s\n", edk.Raw)

	return true
}

func store(args ...string) bool {
	/*
		As following the File Storage Scenario,

		1. ek, dk = hash(pk, sk, F)
		2. EF = E(F, ek)
		3. EFH = ipfs.Set(EF)
		4. Write ownership(EFH, nodeID) and information(EF size, Owner IP) on contract
			- Data Structure

			type Node struct {
				ID	string
				IP	string
			}

			type EF struct {
				Hash  string
				Size  int
				Owner Node
			}

	*/

	// 1. ek, dk = hash(pk, sk, F)

	// 2. EF = E(F, ek)

	// 3. EFH = ipfs.Set(EF)
	var hash string
	if !ipfs.Set(ipfs_shell, args[0], &hash) {
		return false
	}

	// 4. Write ownership(EFH, odeID) and information(EF size, Owner IP) on contract

	return true
}

func help() {
	fmt.Printf("\nRDFS, version %s. Type 'help' to see this list.\n\n", util.RDFS_VER)

	for i := 0; i < len(util.CMD); i++ {
		str := strings.Split(util.CMD[i], ":")
		fmt.Printf(" - %-30s : %s\n", str[0], str[1])
	}

	fmt.Println()
}

func prompt() {
	in := bufio.NewReader(os.Stdin)

	for {
		print(">> ")

		input, _ := in.ReadString('\n')
		cmd, cmd_args := util.Shell(input)

		switch cmd {
		case util.CMD_EXIT:
			rdfsClose()
			os.Exit(0)
		case util.CMD_STORE:
			store(cmd_args...)
		case util.CMD_PURCHASE:
			purchase(cmd_args...)
		case util.CMD_TEST:

			if strings.Compare(cmd_args[0], "-k") == 0 {

				/*
					if privKey == nil || pubKey == nil {
						privKey, pubKey = cryp.GenerateKeySet()
						fmt.Printf("[+] Successfully generated key(private, public) set\n")

						cryp.SavePrivKeyPEMTo(util.RDFS_CONFIG_DIR+"private.pem", privKey)
						cryp.SavePubKeyPEMTo(util.RDFS_CONFIG_DIR+"public.pem", pubKey)

						cryp.Hash(privKey, pubKey, "/pss/tmp/test.txt")
					}

					cryp.Hash(privKey, pubKey, "/pss/tmp/t.py")
				*/

			} else if strings.Compare(cmd_args[0], "-j") == 0 {

				if len(cmd_args) != 4 {
					break
				}

				jrpc.TestB(cmd_args[1], cmd_args[2], cmd_args[3])

			} else if strings.Compare(cmd_args[0], "-g") == 0 {

				//	geth.Test()

			}

		case util.CMD_HELP:
			help()
		case util.CMD_LIST:
			if strings.Compare(cmd_args[0], "-h") == 0 {
				ipfs.List(ipfs_shell, cmd_args[1], ".")
			} else if strings.Compare(cmd_args[0], "-f") == 0 {
				util.List(cmd_args[1])
			}
		case util.CMD_NETVERSION:
			netVersion, err := geth_client.NetVersion()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(netVersion)
		case util.CMD_COINBASE:
			address, err := geth_client.EthCoinBase()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(address)
		case util.CMD_ISMINING:
			mining, err := geth_client.EthMining()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(mining)
		case util.CMD_BLOCKNUMBER:
			num, err := geth_client.EthBlockNumber()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(num)
		case util.CMD_ACCOUNTS:
			accounts, err := geth_client.EthAccounts()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(accounts)
		case util.CMD_BALANCE:
			balance := geth_client.EthGetBalance(cmd_args[0])
			fmt.Println(balance)
		default:
			println("Unsupported Command")
		}
	}
}

func rdfsInit() {
	fmt.Println("[+] Initializing RDFS")

	/*
		Need to implement pre-processing
		for IPFS, GETH to get initialized
	*/

	geth_keys = cryp.GetKey()

	if len(geth_keys) == 0 {
		fmt.Printf("[-] Couldn't initialize RDFS\n")
		os.Exit(1)
	}

	for _, key := range geth_keys {
		fmt.Printf("[+] Processed the key set of address: %x\n", key.Address.Bytes())
	}

	ipfs_pid = ipfs.Open()
	ipfs_shell = shell.NewShell("localhost:5001")
	geth_pid, geth_client = geth.Open()

	go jrpc.InitServer()
}

func rdfsClose() {
	fmt.Println("\n[+] Closing RDFS")

	ipfs.Close(ipfs_pid)
	geth.Close(geth_pid)
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
