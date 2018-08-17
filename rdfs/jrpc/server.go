package jrpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	"rdfs/crypto"
	"rdfs/ipfs"
	"rdfs/util"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

var GETH_KEYS []*keystore.Key

type Args struct{ A, B int }
type Result int

type Arith int

func (t *Arith) Multiply(args Args, result *Result) error {
	fmt.Printf("[+] RPC Server: Calculating %d * %d\n", args.A, args.B)
	*result = Result(args.A * args.B)
	return nil
}

type FileToRequest struct {
	Hash    []byte
	PubKey  []byte
	Address common.Address
}

type InfoToReturn struct {
	Name                   string
	EncryptedDecryptionKey []byte
}

type Crypto int

func (t *Crypto) Encrypt(file FileToRequest, info *InfoToReturn) error {

	fmt.Printf("[+] Validating requested purchase\n")

	/*
		Checking transaction process here
	*/

	fmt.Printf("[+] RPC Server: Encrypting DK with buyer's PK\n")

	// Encrpyting decryption key with public key here
	pubKey := crypto.ECDSADecode(file.PubKey)

	decryptionKey := ipfs.GetFileDecryptionKey(file.Hash, GETH_KEYS[0].PrivateKey)
	info.EncryptedDecryptionKey = *crypto.ECIESEncrypt(pubKey, decryptionKey)

	// Adding file's name
	info.Name = ipfs.GetFileName(file.Hash)

	return nil
}

func InitServer(keys []*keystore.Key) {
	GETH_KEYS = keys

	arith := new(Arith)
	crypto := new(Crypto)

	server := rpc.NewServer()

	server.Register(arith)
	server.Register(crypto)

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	listener, err := net.Listen("tcp", ":"+util.JSON_RPC_PORT)

	if err != nil {
		fmt.Printf("[-] RPC Server: Listen error: %s\n", err)
	}

	for {
		if conn, err := listener.Accept(); err != nil {
			fmt.Printf("[-] RPC Server: Accept error: %s\n", err.Error())
		} else {
			addr := strings.Split(conn.RemoteAddr().String(), ":")[0]
			fmt.Printf("[+] RPC Server: New connection (%s) established\n", addr)

			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
