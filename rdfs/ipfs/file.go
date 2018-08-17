package ipfs

/*
 * The information of files stored on IPFS
 * is saved in JSON File located in RDFS_CONFIG_DIR
 */

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"rdfs/crypto"
	"rdfs/util"
)

const filePath = util.RDFS_CONFIG_DIR + "files.json"

type FileList struct {
	Files map[string]File `json:"files"`
}

type File struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	DecryptionKey []byte `json:"decryptionKey"`
	Exists        bool   `json:"exists"`
}

func (jf *FileList) read() bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.OpenFile(filePath, os.O_CREATE, 0600)
		jf.Files = make(map[string]File)
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	jf.Files = make(map[string]File)
	json.Unmarshal(file, &(jf.Files))

	return true
}

func (jf *FileList) write() bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.OpenFile(filePath, os.O_CREATE, 0600)
	}

	jsonBytes, err := json.Marshal(jf.Files)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	return true
}

func (jf *FileList) SetFileInfo(hash []byte, name string, pubKey *ecdsa.PublicKey, dk []byte) bool {
	if !jf.read() {
		fmt.Printf("[-] Couldn't read file's information\n")
		return false
	}

	hashHEX := hex.EncodeToString(hash)

	if _, ok := jf.Files[hashHEX]; ok {
		fmt.Printf("[-] Already saved in file's list: %s\n", hash)
		return false
	}

	edk := crypto.ECIESEncrypt(pubKey, &dk)
	file := File{
		Name:          name,
		Path:          util.RDFS_UP_DIR + name,
		DecryptionKey: *edk,
		Exists:        true}

	jf.Files[hashHEX] = file

	if !jf.write() {
		return false
	}

	return true
}

func GetFileDecryptionKey(hash []byte, privKey *ecdsa.PrivateKey) *[]byte {
	var jf FileList

	hashHEX := hex.EncodeToString(hash)

	if !jf.read() {
		fmt.Printf("[-] Couldn't read file's information\n")
		return nil
	}

	if _, ok := jf.Files[hashHEX]; !ok {
		fmt.Printf("[-] Doesn't exist in file's list: %s\n", hashHEX)
		return nil
	}

	file := jf.Files[hashHEX]

	return crypto.ECIESDecrypt(privKey, &(file.DecryptionKey))
}

func GetFileName(hash []byte) string {
	var jf FileList

	hashHEX := hex.EncodeToString(hash)

	if !jf.read() {
		fmt.Printf("[-] Couldn't read file's information\n")
		return ""
	}

	if _, ok := jf.Files[hashHEX]; !ok {
		fmt.Printf("[-] Doesn't exist in file's list: %s\n", hashHEX)
		return ""
	}

	return jf.Files[hashHEX].Name
}
