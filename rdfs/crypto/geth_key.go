package crypto

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"rdfs/util"
)

func GetKey() []*keystore.Key {
	ks := keystore.NewPlaintextKeyStore(util.GETH_DATA_DIR + "keystore")
	accounts := ks.Accounts()

	if len(accounts) == 0 {
		return nil
	}

	var keys []*keystore.Key

	in := bufio.NewReader(os.Stdin)

	for _, account := range accounts {
		fmt.Printf("[+] Processing %x's key\n", account.Address.Bytes())

		print("[*] Passphrase >> ")
		auth, _ := in.ReadString('\n')
		auth = strings.TrimSpace(auth)

		key := getPrivKey(account, auth)

		if key != nil {
			keys = append(keys, key)
		}
	}

	return keys
}

func getPrivKey(account accounts.Account, auth string) *keystore.Key {
	keyjson, err := ioutil.ReadFile(account.URL.Path)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return key
}
