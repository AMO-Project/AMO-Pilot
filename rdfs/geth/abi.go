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
						{ "name" : "_hash", "type" : "bytes32" }
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

const rdfsCoinABI = `
[
    {
      "constant": true,
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_spender",
          "type": "address"
        },
        {
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "INITIAL_SUPPLY",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "decimals",
      "outputs": [
        {
          "name": "",
          "type": "uint8"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_spender",
          "type": "address"
        },
        {
          "name": "_subtractedValue",
          "type": "uint256"
        }
      ],
      "name": "decreaseApproval",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "_owner",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [],
      "name": "renounceOwnership",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "adminAddr",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "owner",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_spender",
          "type": "address"
        },
        {
          "name": "_addedValue",
          "type": "uint256"
        }
      ],
      "name": "increaseApproval",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "_owner",
          "type": "address"
        },
        {
          "name": "_spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_newOwner",
          "type": "address"
        }
      ],
      "name": "transferOwnership",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "name": "_adminAddr",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "previousOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipRenounced",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "previousOwner",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "OwnershipTransferred",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "from",
          "type": "address"
        },
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "from",
          "type": "address"
        },
        {
          "name": "to",
          "type": "address"
        },
        {
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "transferDepositTo",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "from",
          "type": "address"
        },
        {
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "addDeposit",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "from",
          "type": "address"
        },
        {
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "subDeposit",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "target",
          "type": "address"
        }
      ],
      "name": "depositOf",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "target",
          "type": "address"
        }
      ],
      "name": "availableBalanceOf",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
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

func CallRDFSCoinABI() *abi.ABI {
	abi, err := abi.JSON(strings.NewReader(rdfsCoinABI))
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
