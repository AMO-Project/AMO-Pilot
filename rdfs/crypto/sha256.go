package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"

	"fmt"
)

func GenerateHashKey(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey, file *[]byte) []byte {
	h := sha256.New()

	encodedPrivKey, _ := x509.MarshalECPrivateKey(privKey)
	encodedPubKey, _ := x509.MarshalPKIXPublicKey(pubKey)

	h.Write(encodedPrivKey)
	h.Write(encodedPubKey)
	h.Write(*file)

	hash := h.Sum(nil)
	fmt.Printf("[+] Successfully generated hash: %x\n", hash)

	return hash
}
