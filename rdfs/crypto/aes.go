package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"fmt"
	"io"
)

func AESEncrypt(plain *[]byte, key *[]byte) *[]byte {
	block, err := aes.NewCipher(*key)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	ciphertext := make([]byte, aes.BlockSize+len(*plain))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], *plain)

	fmt.Printf("[+] Successfully encrypted file: %x\n", ciphertext[:32])

	return &ciphertext
}

func AESDecrypt(ciphertext *[]byte, key *[]byte) *[]byte {
	block, err := aes.NewCipher(*key)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	if len(*ciphertext) < aes.BlockSize {
		fmt.Printf("[-] Error occured: ciphertext is too short\n")
		return nil
	}

	iv := (*ciphertext)[:aes.BlockSize]
	encrypted := (*ciphertext)[aes.BlockSize:]

	plain := make([]byte, len(encrypted))

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plain, encrypted)

	fmt.Printf("[+] Successfully decrypted file: %x\n", plain[:32])

	return &plain
}
