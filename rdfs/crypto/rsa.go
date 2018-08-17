package crypto

import (
	"fmt"
	"os"

	"encoding/asn1"
	"encoding/pem"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

func GenerateKeySet() (*rsa.PrivateKey, *rsa.PublicKey) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("[-] Error occured while generating key set\n")
		return nil, nil
	}

	err = privKey.Validate()
	if err != nil {
		fmt.Printf("[-] Error occured while validating generated private key\n")
		return nil, nil
	}

	pubKey := &privKey.PublicKey
	fmt.Printf("[+] Successfully generated key set\n")

	return privKey, pubKey
}

func SavePrivKeyPEMTo(fileName string, privKey *rsa.PrivateKey) bool {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		outFile.Close()
		return false
	}
	defer outFile.Close()

	privBlock := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(privKey),
	}

	err = pem.Encode(outFile, privBlock)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	fmt.Printf("[+] Successfully saved private key to %s in PEM format\n", fileName)

	return true
}

func SavePubKeyPEMTo(fileName string, pubKey *rsa.PublicKey) bool {
	asn1Bytes, err := asn1.Marshal(*pubKey)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}

	pubBlock := &pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   asn1Bytes,
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return false
	}
	defer outFile.Close()

	err = pem.Encode(outFile, pubBlock)

	fmt.Printf("[+] Successfully saved public key to %s in PEM format\n", fileName)

	return true
}
