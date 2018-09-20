package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"rdfs/crypto"
	"rdfs/geth"
	"rdfs/ipfs"
	"rdfs/util"
)

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
		case util.CMD_TOKENBAL:
			GETH_CLIENT.TokenBalance(GETH_KEYS[0].Address.String(), GETH_KEYS[0].Address)
		case util.CMD_PEERFILE:
			getPeerFile()
		case util.CMD_TEST:

			if strings.Compare(cmd_args[0], "-k") == 0 {
				f, err := ioutil.ReadFile("/pss/tmp/go-ipfs_v0.4.15_linux-amd64.tar.gz")
				//f, err := ioutil.ReadFile("/home/h0n9/dev/tmp/go-ipfs_v0.4.15_linux-amd64.tar.gz")
				if err != nil {
					fmt.Printf("[-] Error occured: %s\n", err)
					break
				}

				privKey := GETH_KEYS[0].PrivateKey
				pubKey := privKey.PublicKey

				hash := crypto.GenerateHashKey(privKey, &pubKey, &f)

				encrypted := crypto.AESEncrypt(&f, &hash)
				decrypted := crypto.AESDecrypt(encrypted, &hash)

				fmt.Printf("AES TEST: %t\n", bytes.Equal(f, *decrypted))

				edk := crypto.ECIESEncrypt(&pubKey, &hash)
				dk := crypto.ECIESDecrypt(privKey, edk)

				fmt.Printf("ECIES TEST: %t\n", bytes.Equal(hash, *dk))

			} else if strings.Compare(cmd_args[0], "-pub") == 0 {
				ipfs.PublishDefaultDir(IPFS_SHELL)
			} else if strings.Compare(cmd_args[0], "-j") == 0 {
				rdfsFileABI := geth.CallRDFSFileABI()

				rawHash := util.MultiHashToBytes("QmNybj8qNJnLL8LRKKanVbZuwV9SCbN4YRXdm7Pwb7mZ6h")
				hash := [32]byte{}
				copy(hash[:], rawHash)

				//name := "testFile"
				//size := big.NewInt(1000000000000000000)
				//ip := util.GetPublicIP()

				fmt.Printf("0x%x\n", hash)

				//data, err := rdfsFileABI.Pack("storeRequest", hash, name, &size, ip)

				data, err := rdfsFileABI.Pack("purchaseRequest", hash)
				if err != nil {
					fmt.Printf("[-] Error occured: %s\n", err)
					break
				}
				fmt.Printf("%x\n%s\n", data, hex.EncodeToString(data))
				fmt.Printf("%x\n", GETH_KEYS[0].Address)

				gasPrice := GETH_CLIENT.EstimateGas(GETH_KEYS[0].Address.Hex(), &data)

				fmt.Println("Estimated gasPrice: ", gasPrice)

			} else if strings.Compare(cmd_args[0], "-st") == 0 {
				hash := util.MultiHashToBytes("QmNybj8qNJnLL8LRKKanVbZuwV9SCbN4YRXdm7Pwb7mZ6h")

				name := "testFile"
				size := big.NewInt(10000)
				ip := util.GetPublicIP()

				ok := GETH_CLIENT.StoreRequest(GETH_KEYS[0], hash, name, size, ip)
				if !ok {
					fmt.Printf("OOps!\n")
					break
				}

			} else if strings.Compare(cmd_args[0], "-pr") == 0 {
				hash := util.MultiHashToBytes("QmNybj8qNJnLL8LRKKanVbZuwV9SCbN4YRXdm7Pwb7mZ6h")

				ok := GETH_CLIENT.PurchaseRequest(GETH_KEYS[1], hash)
				if !ok {
					fmt.Printf("OOps!\n")
					break
				}
			} else if strings.Compare(cmd_args[0], "-info") == 0 {
				rawHash := util.MultiHashToBytes("QmNybj8qNJnLL8LRKKanVbZuwV9SCbN4YRXdm7Pwb7mZ6h")

				fileName := GETH_CLIENT.GetFileName(GETH_KEYS[1].Address.String(), rawHash)
				fmt.Printf("%s\n", fileName)

				fileSize := GETH_CLIENT.GetFileSize(GETH_KEYS[1].Address.String(), rawHash)
				fmt.Printf("%s\n", fileSize)

				ownerAddr := GETH_CLIENT.IsOwnedBy(GETH_KEYS[1].Address.String(), rawHash)
				fmt.Printf("%s\n", ownerAddr)

				isRequested := GETH_CLIENT.IsRequested(GETH_KEYS[1].Address.String(), rawHash, GETH_KEYS[1].Address)
				fmt.Printf("%t\n", isRequested)

				isApproved := GETH_CLIENT.IsApproved(GETH_KEYS[1].Address.String(), rawHash, GETH_KEYS[1].Address)
				fmt.Printf("%t\n", isApproved)

				ownerAddr, ownerIP := GETH_CLIENT.GetOwnerInfo(GETH_KEYS[1].Address.String(), rawHash)
				fmt.Printf("%s, %x\n", ownerAddr, ownerIP)

				beforeHeight, err := GETH_CLIENT.EthBlockNumber()
				if err != nil {
					fmt.Printf("[-] Error occured: %s\n", err)
					break
				}

				for {
					currentHeight, err := GETH_CLIENT.EthBlockNumber()
					if err != nil {
						fmt.Printf("[-] Error occured: %s\n", err)
						break
					}
					if currentHeight-beforeHeight > 2 {
						break
					}
					time.Sleep(2 * time.Second)
					fmt.Printf("before : %d , current : %d\r", beforeHeight, currentHeight)
				}
			}

		case util.CMD_HELP:
			help()
		case util.CMD_LIST:
			if strings.Compare(cmd_args[0], "-h") == 0 {
				if len(cmd_args) != 2 {
					fmt.Printf("[-] Type the hash of file\n")
					break
				}
				ipfs.List(IPFS_SHELL, cmd_args[1], ".")
			} else if strings.Compare(cmd_args[0], "-f") == 0 {
				util.List(cmd_args[1])
			}
		case util.CMD_NETVERSION:
			netVersion, err := GETH_CLIENT.NetVersion()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(netVersion)
		case util.CMD_COINBASE:
			address, err := GETH_CLIENT.EthCoinBase()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(address)
		case util.CMD_ISMINING:
			mining, err := GETH_CLIENT.EthMining()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(mining)
		case util.CMD_BLOCKNUMBER:
			num, err := GETH_CLIENT.EthBlockNumber()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(num)
		case util.CMD_ACCOUNTS:
			accounts, err := GETH_CLIENT.EthAccounts()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(accounts)
		case util.CMD_ETHBALANCE:
			balance := GETH_CLIENT.EthGetBalance(cmd_args[0])
			fmt.Println(balance)
		case util.CMD_SENDTX:
			/*
				need to re-implement sendTx func if necessary...

				txhash, err := GETH_CLIENT.EthSendTransaction(GETH_KEYS[0].Address.Hex(), geth.ADDR_ACCOUNT["ps2"], "0x23e3fbd5d64162eb51f593bffaa11b223fe337ed5f7e34c0d0fa1a08d1f1cf9d")

				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(txhash)
			*/
		case util.CMD_UNLOCK:
			addrIndex, _ := strconv.ParseInt(cmd_args[0], 10, 64)
			pass := strings.TrimSpace(cmd_args[1])
			time, _ := strconv.ParseUint(cmd_args[2], 10, 64)

			response, err := GETH_CLIENT.PersonalUnlockAccount(GETH_KEYS[addrIndex-1].Address.Hex(), pass, time)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(response)
		default:
			println("Unsupported Command")
		}
	}
}
