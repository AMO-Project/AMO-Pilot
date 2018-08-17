package cryp

import (
	"crypto/rsa"
	"crypto/sha256"
	"fmt"

	"io/ioutil"
)

func Hash(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, filePath string) []byte {
	f, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	h := sha256.New()

	h.Write(privKey.N.Bytes())
	h.Write(pubKey.N.Bytes())
	h.Write(f)

	hash := h.Sum(nil)
	fmt.Printf("[+] Successfully generated hash: %x\n", hash)

	return hash
}
