package geth

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const rdfsFileABI = `
[
	{ 
		"type"    : "function",
		"name"    : "storeRequest",
		"constant": false,
		"inputs"  : [ 
						{ "name" : "_hash", "type" : "bytes32" },
					 	{ "name" : "_name", "type" : "string" },
					 	{ "name" : "_size", "type" : "uint256" },
					 	{ "name" : "_ip", "type" : "bytes4" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "purchaseRequest",
		"constant": false,
		"inputs"  : [ 
						{ "name" : "_hash", "type" : "bytes32" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "purchaseApprove",
		"constant": false,
		"inputs"  : [ 
						{ "name" : "_hash", "type" : "bytes32" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "purchaseAbandon",
		"constant": false,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "getFileName",
		"constant": true,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "string" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "getFileSize",
		"constant": true,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "uint256" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "isOwnedBy",
		"constant": true,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" },
						{ "name" : "buyer", "type" : "address" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "isRequested",
		"constant": true,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" },
						{ "name" : "buyer", "type" : "address" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "isApproved",
		"constant": true,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" },
						{ "name" : "buyer", "type" : "address" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "bool" }
					]
	},
	{ 
		"type"    : "function",
		"name"    : "getOwnerInfo",
		"constant": true,
		"inputs"  : [
						{ "name" : "_hash", "type" : "bytes32" }
				    ],
		"outputs" : [
						{ "name" : "", "type" : "address" },
						{ "name" : "", "type" : "bytes4" }
					]
	}
]`

func CallRDFSFileABI() *abi.ABI {
	abi, err := abi.JSON(strings.NewReader(rdfsFileABI))
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return &abi
}

/*
func code:		0x35c213eb
hash: 			122009771e2e093a910212d0b8afe5cfc35fef5983b0ee957df51f26ebb0d24e
(name_loca):	0000000000000000000000000000000000000000000000000000000000000080
size: 			0000000000000000000000000000000000000000000000000de0b6b3a7640000
ip:	  			a5c223d200000000000000000000000000000000000000000000000000000000
(name_size):	0000000000000000000000000000000000000000000000000000000000000008
name: 			7465737446696c65000000000000000000000000000000000000000000000000
*/
