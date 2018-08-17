package crypto

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func ECDSAEncode(pubKey *ecdsa.PublicKey) []byte {
	encodedPubKey := crypto.FromECDSAPub(pubKey)
	return encodedPubKey
}

func ECDSADecode(encodedPubKey []byte) *ecdsa.PublicKey {
	pubKey, err := crypto.UnmarshalPubkey(encodedPubKey)
	if err != nil {
		fmt.Printf("[-] Error occured: %s", err)
		return nil
	}
	return pubKey
}
