package util

import (
	"fmt"

	"github.com/multiformats/go-multihash"
)

func MultiHashToBytes(mhb58 string) []byte {
	mh, err := multihash.FromB58String(mhb58)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return mh[2:]
}
