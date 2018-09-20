package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"

	"rdfs/crypto"
	"rdfs/ipfs"
	"rdfs/jrpc"
	"rdfs/util"
)

/*
	As following the File Purchase Scenario,

	1. EF = ipfs.Get(EFH)
	2. Request purchase by sending token, EFH to contract
		- (Contract) Check token amount and put token into deposit
		- (Contract) Return owner's geth and ip addresses to buyer

	3-a. Send buyer's pk(publicKey) to owner's json-rpc server
		- Check the transaction
		- If valid, edk = E(dk, pk) and send edk to buyer's json-rpc server

	3-b. Buyer wait for edk to receive
		(+) if receive, approve purchase if receives edk from file owner => 4
		(-) if not, abandon purchase if not receive edk from file owner  => end

	4. Receive edk and decrypt EF
		- dk = D(edk, sk)	*sk = privateKey
		- F = D(EF, dk)

*/

func purchase(args ...string) bool {
	var filePrice *big.Int
	var in *bufio.Reader
	var input string
	var encryptedFile []byte
	var encryptedFileHashBytes []byte
	var beforeHeight int

	var ok bool
	var err error

	tempDir := util.RDFS_DOWN_DIR + ".temp/"

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}

	fromAddr := GETH_KEYS[0]
	encryptedFileHashBytes = util.MultiHashToBytes(args[0])

	// 0-1. Check if buyer has already bought the file before

	fmt.Printf("[+] Checking if client has already bought '%s' before\n", args[0])

	requested := GETH_CLIENT.IsRequested(fromAddr.Address.String(), encryptedFileHashBytes, fromAddr.Address)
	approved := GETH_CLIENT.IsApproved(fromAddr.Address.String(), encryptedFileHashBytes, fromAddr.Address)
	if requested || approved {
		fmt.Printf("[+] Purchase Record: requested='%t', approved='%t'... Go through HIPASS!\n", requested, approved)
		goto HIPASS
	}

	// 0-2. Check if buyer would like to pay for the file's price(size)
	filePrice = GETH_CLIENT.GetFileSize(fromAddr.Address.String(), encryptedFileHashBytes)

	if filePrice.Uint64() < 1 {
		fmt.Printf("[-] File doesn't exist on RDFS Contract\n")
		return false
	}

	GETH_CLIENT.TokenBalance(fromAddr.Address.String(), fromAddr.Address)
	fmt.Printf("[?] Would you like to pay '%d' token to get '%s'? (y/n) ", filePrice, args[0])

	in = bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')

	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if strings.Compare(input, "y") != 0 && strings.Compare(input, "yes") != 0 {
		return false
	}

	// 1. EF = ipfs.Get(EFH)
	if ipfs.Get(IPFS_SHELL, args[0], tempDir) == false {
		return false
	}

	// 2. Request purchase by sending token, EFH to contract
	// function purchaseRequest(bytes32 _hash)
	fmt.Printf("[+] Purchase request by sending token with EFH to contract\n")

	ok = GETH_CLIENT.PurchaseRequest(fromAddr, encryptedFileHashBytes)
	if !ok {
		fmt.Printf("[-] Purchase request Failed\n")
		return false
	}

	beforeHeight, err = GETH_CLIENT.EthBlockNumber()
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Waiting for 3 confirmations\n")

	for {
		currentHeight, err := GETH_CLIENT.EthBlockNumber()
		if err != nil {
			fmt.Printf("[-] Error occured: %s\n", err)
			return false
		}
		if currentHeight-beforeHeight > 2 {
			break
		}
		time.Sleep(2 * time.Second)
		fmt.Printf("[+] before : %d , current : %d (%d)\r", beforeHeight, currentHeight, currentHeight-beforeHeight)
	}

	fmt.Println()

	requested = GETH_CLIENT.IsRequested(fromAddr.Address.String(), encryptedFileHashBytes, fromAddr.Address)
	if !requested {
		fmt.Printf("[-] Cannot find purchase request record.\n")
		return false
	}

	// 3. If accepted, receive owner's info(addr, ip)
	fmt.Printf("[+] Successfully put the amount into your deposit for the owner\n")
	GETH_CLIENT.TokenBalance(fromAddr.Address.String(), fromAddr.Address)

HIPASS:

	if requested || approved {
		ok = ipfs.Get(IPFS_SHELL, args[0], tempDir)
		if !ok {
			return false
		}
	}

	ownerAddr, ownerIP := GETH_CLIENT.GetOwnerInfo(fromAddr.Address.String(), encryptedFileHashBytes)
	if strings.Compare(ownerAddr, "") == 0 {
		fmt.Printf("[-] Cannot resolve owner's address.\n")
		return false
	}

	fmt.Printf("[+] Got owner's address : %s, %s\n", ownerAddr, util.EncodeIP(ownerIP))

	// Check whether Owner is Stayin' alive
	var ownerStatus []byte
	ok = jrpc.IsAlive(ownerIP, &ownerStatus)
	if !ok {
		fmt.Printf("[-] Owner is not Stayin' alive.\n")

		if requested {
			// purchaseAbandon
			ok = GETH_CLIENT.PurchaseAbandon(fromAddr, encryptedFileHashBytes)
			for {
				if !ok {
					ok = GETH_CLIENT.PurchaseAbandon(fromAddr, encryptedFileHashBytes)
				} else {
					GETH_CLIENT.TokenBalance(fromAddr.Address.String(), fromAddr.Address)
					break
				}
			}
		}

		return false
	}

	// 4. Send buyer's pk(publickKey) to owner's json-rpc(?) server
	privKey := fromAddr.PrivateKey
	pubKey := &(privKey.PublicKey)

	fileToRequest := jrpc.FileToRequest{
		Hash:    encryptedFileHashBytes,
		PubKey:  crypto.ECDSAEncode(pubKey),
		Address: fromAddr.Address,
	}

	var infoToReturn jrpc.InfoToReturn
	ok = jrpc.RequestEDK(ownerIP, fileToRequest, &infoToReturn)

	if !ok {
		fmt.Printf("[-] Couldn't receive EDK from owner\n")

		// purchaseAbandon
		ok = GETH_CLIENT.PurchaseAbandon(fromAddr, encryptedFileHashBytes)
		for {
			if !ok {
				ok = GETH_CLIENT.PurchaseAbandon(fromAddr, encryptedFileHashBytes)
			} else {
				GETH_CLIENT.TokenBalance(fromAddr.Address.String(), fromAddr.Address)
				break
			}
		}

		return false
	}

	// 5. Receive edk and decrypt EF
	fmt.Printf("[+] Successfully received EDK: %x\n", infoToReturn.EncryptedDecryptionKey[:31])

	encryptedFile, err = ioutil.ReadFile(tempDir + args[0])
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	decryptionKey := crypto.ECIESDecrypt(privKey, &infoToReturn.EncryptedDecryptionKey)
	if decryptionKey == nil {
		return false
	}

	decryptedFile := crypto.AESDecrypt(&encryptedFile, decryptionKey)
	if decryptedFile == nil {
		return false
	}

	err = ioutil.WriteFile(util.RDFS_DOWN_DIR+infoToReturn.Name, *decryptedFile, 0666)
	if err != nil {
		fmt.Printf("[-] Couldnt't save decrypted file\n")
		return false
	}

	err = os.Remove(tempDir + args[0])
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	if !approved {
		ok = GETH_CLIENT.PurchaseApprove(fromAddr, encryptedFileHashBytes)
		for {
			if !ok {
				ok = GETH_CLIENT.PurchaseApprove(fromAddr, encryptedFileHashBytes)
			} else {
				GETH_CLIENT.TokenBalance(fromAddr.Address.String(), fromAddr.Address)
				break
			}
		}

		fmt.Printf("[+] Successfully approved purchase\n")
	}

	fmt.Printf("[+] Successfully purchased file '%s' at '%s'\n", args[0], util.RDFS_DOWN_DIR+infoToReturn.Name)

	return true
}
