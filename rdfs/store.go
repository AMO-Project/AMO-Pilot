package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"rdfs/crypto"
	"rdfs/ipfs"
	"rdfs/util"
	"strings"
)

/*
	As following the File Storage Scenario,

	1. ek = hash(pk, sk, F)
	2. EF = E(F, ek)
	3. EFH = ipfs.Set(EF)
	4. Write ownership(EFH, nodeID) and information(EF size, Owner IP) on contract
		- Data Structure

		struct Node {
			address addr;
			bytes4 ip;
		}

		struct File {
			string name;
			uint256 size;
			Node owner;

			mapping(address => purchaseState) buyers;
			bool exists;
		}

*/

func store(args ...string) bool {

	fileInfo, err := os.Stat(args[0])
	if err != nil {
		fmt.Println("[-] Could not check the file info. " +
			"Check the file path, please")
		return false
	}

	if fileInfo.IsDir() {
		fmt.Printf("[-] Adding a directory is not supported yet\n")
		return false
	}

	privKey := GETH_KEYS[0].PrivateKey
	pubKey := &(privKey.PublicKey)

	f, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	filePath := strings.Split(args[0], "/")
	fileName := filePath[len(filePath)-1]
	fileSize := big.NewInt(int64(len(f)))

	/*
		if ok := util.Copy(args[0], util.RDFS_UP_DIR+fileName); !ok {
			fmt.Printf("[-] Couldn't copy file to %s(RDFS_UP_DIR)\n", util.RDFS_UP_DIR)
			return false
		}
	*/

	// 1. ek = hash(pk, sk, F)
	encryptionKey := crypto.GenerateHashKey(privKey, pubKey, &f)

	// 2. EF = E(F, ek)
	encryptedFile := crypto.AESEncrypt(&f, &encryptionKey)

	// 3. EFH = ipfs.Set(EF)
	encyptedFileHash, err := ipfs.Add(IPFS_SHELL, encryptedFile)
	if err != nil {
		fmt.Printf("[-] Couldn't add file '%s'\n", args[0])
		return false
	}

	fmt.Printf("[+] Added to IPFS '%s': '%s'\n", args[0], encyptedFileHash)

	// 4. Write ownership(EFH, nodeID) and information(EF size, Owner IP) on contract
	encyptedFileHashBytes := util.MultiHashToBytes(encyptedFileHash)

	nodeIP := util.GetPublicIP()
	nodeAddr := GETH_KEYS[0].Address.Hex()

	reqOk := GETH_CLIENT.StoreRequest(nodeAddr, encyptedFileHashBytes, fileName, fileSize, nodeIP)
	if reqOk == false {
		fmt.Printf("[-] Couldn't request storing file\n")
		return false
	}

	setOk := IPFS_FILES.SetFileInfo(encyptedFileHashBytes, fileName, pubKey, encryptionKey)
	if setOk == false {
		fmt.Printf("[-] Couldn't save file's information\n")
		return false
	}

	err = ioutil.WriteFile(util.RDFS_UP_DIR+fileName, *encryptedFile, 0666)
	if err != nil {
		fmt.Printf("[-] Couldn't copy file to '%s' (RDFS_UP_DIR)\n", util.RDFS_UP_DIR)
		return false
	}

	fmt.Printf("[+] Copied file to '%s' (RDFS_UP_DIR)\n", util.RDFS_UP_DIR)

	// ipfs add -r /pss/rdfs/up => hash
	// ipfs name publish hash
	ok := ipfs.PublishDefaultDir(IPFS_SHELL)
	if ok == false {
		return false
	}

	return true
}
