package main

import (
	"context"
	"fmt"
	"strings"

	"rdfs/ipfs"
)

func getPeerFile() {
	ctx := context.Background()

	swarmInfos, err := IPFS_SHELL.SwarmPeers(ctx)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return
	}

	for _, swarmInfo := range swarmInfos.Peers {
		fmt.Printf("[*] Peer : %s\n", swarmInfo.Peer)
		fmt.Printf("---------------------------------------------------------\n")
		hash := ipfs.Resolve(IPFS_SHELL, swarmInfo.Peer)
		if strings.Compare(hash, "") == 0 {
			continue
		}
		ipfs.List(IPFS_SHELL, hash, "")
		fmt.Printf("\n\n")
	}
}
