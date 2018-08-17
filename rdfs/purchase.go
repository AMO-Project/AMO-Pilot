package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

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

	temp_dir := util.RDFS_DOWN_DIR + ".temp/"

	if _, err := os.Stat(temp_dir); os.IsNotExist(err) {
		os.Mkdir(temp_dir, 0755)
	}

	// 0. Check first if buyer would like to pay for the file's price(size)

	// filePrice := geth.GetFileSize(args[0])
	filePrice := 10209630
	fmt.Printf("[?] Would you like to pay '%d' token to get '%s'? (y/n) ", filePrice, args[0])

	in := bufio.NewReader(os.Stdin)
	input, _ := in.ReadString('\n')

	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if strings.Compare(input, "y") != 0 && strings.Compare(input, "yes") != 0 {
		return false
	}

	// 1. EF = ipfs.Get(EFH)
	if ipfs.Get(IPFS_SHELL, args[0], temp_dir) == false {
		return false
	}

	encryptedFile, err := ioutil.ReadFile(temp_dir + args[0])
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	// 2. Request purchase by sending token, EFH to contract
	// function purchaseRequest(bytes32 _hash)
	fmt.Printf("[+] Request purchase by sending token with EFH to contract\n")

	// ownerAddr, ownerIP := geth.PurchaseRequest(args[0])
	//if ownerAddr == nil {
	//	return fasle
	//}

	// 3. If accepted, receive owner's info(addr, ip)
	fmt.Printf("[+] Successfully put the amount into your deposit for the owner\n")

	ownerIP := "127.0.0.1"
	if len(args) == 2 {
		ownerIP = args[1]
	}

	// 4. Send buyer's pk(publickKey) to owner's json-rpc(?) server
	encryptedFileHash := util.MultiHashToBytes(args[0])

	privKey := GETH_KEYS[0].PrivateKey
	pubKey := &(privKey.PublicKey)

	fileToRequest := jrpc.FileToRequest{
		Hash:    encryptedFileHash,
		PubKey:  crypto.ECDSAEncode(pubKey),
		Address: GETH_KEYS[0].Address}

	var infoToReturn jrpc.InfoToReturn

	if jrpc.RequestEDK(net.ParseIP(ownerIP).To4(), fileToRequest, &infoToReturn) == false {
		fmt.Printf("[-] Couldn't receive EDK from owner\n")
		return false
	}

	// 5. Receive edk and decrypt EF
	fmt.Printf("[+] Successfully received EDK: %x\n", infoToReturn.EncryptedDecryptionKey[:31])

	decryptionKey := crypto.ECIESDecrypt(privKey, &infoToReturn.EncryptedDecryptionKey)
	decryptedFile := crypto.AESDecrypt(&encryptedFile, decryptionKey)

	err = ioutil.WriteFile(util.RDFS_DOWN_DIR+infoToReturn.Name, *decryptedFile, 0666)
	if err != nil {
		fmt.Printf("[-] Couldnt't save decrypted file\n")
		return false
	}

	err = os.Remove(temp_dir + args[0])
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	return true
}
