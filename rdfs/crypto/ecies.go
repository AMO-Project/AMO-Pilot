package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func ECIESEncrypt(pubKey *ecdsa.PublicKey, plain *[]byte) *[]byte {
	eciesPubKey := ecies.ImportECDSAPublic(pubKey)

	ct, err := ecies.Encrypt(rand.Reader, eciesPubKey, *plain, nil, nil)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return &ct
}

func ECIESDecrypt(privKey *ecdsa.PrivateKey, ciphertext *[]byte) *[]byte {
	eciesPrivKey := ecies.ImportECDSA(privKey)

	pt, err := eciesPrivKey.Decrypt(*ciphertext, nil, nil)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return &pt
}
