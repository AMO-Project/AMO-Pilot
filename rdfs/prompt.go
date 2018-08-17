package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"rdfs/crypto"
	"rdfs/ipfs"
	"rdfs/util"
	"strconv"
	"strings"
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
		case util.CMD_TEST:

			if strings.Compare(cmd_args[0], "-k") == 0 {
				f, err := ioutil.ReadFile("/pss/tmp/go-ipfs_v0.4.15_linux-amd64.tar.gz")
				//f, err := ioutil.ReadFile("/home/h0n9/dev/tmp/go-ipfs_v0.4.15_linux-amd64.tar.gz")
				if err != nil {
					fmt.Printf("[-] Error occured: %s\n", err)
					return
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

			} else if strings.Compare(cmd_args[0], "-j") == 0 {
				/*
					if len(cmd_args) != 4 {
						break
					}
					jrpc.TestB(cmd_args[1], cmd_args[2], cmd_args[3])
				*/
			} else if strings.Compare(cmd_args[0], "-g") == 0 {

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
		case util.CMD_BALANCE:
			balance := GETH_CLIENT.EthGetBalance(cmd_args[0])
			fmt.Println(balance)
		case util.CMD_SENDTX:
			fmt.Println(ADDR_ACCOUNT["ps1"])
			fmt.Println(CTRC_COIN)
			txhash, err := GETH_CLIENT.EthSendTransaction()

			if err != nil {
				fmt.Println(err, "at prompt")
			}
			fmt.Println(txhash)
		case util.CMD_UNLOCK:
			time, _ := strconv.Atoi(cmd_args[0])
			response, err := GETH_CLIENT.PersonalUnlockAccount("0x2074fa38f08facdf47f08b8051f9a6aff6033607", "ps1", time)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(response)
		default:
			println("Unsupported Command")
		}
	}
}
