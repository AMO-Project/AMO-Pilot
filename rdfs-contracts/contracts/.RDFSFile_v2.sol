pragma solidity ^0.4.2;

import "zeppelin-solidity/contracts/math/SafeMath.sol";

contract RDFSFile {
    using SafeMath for uint256;

    string public constant symbol = "RDFS";
    string public constant name = "RDFS File";

    /*
     * Things to consider for storing and purchasing file
     *
     * files[f_hash] => FileInfo{size, list}
     * files[f_hash].list[i] => EncryptedFile{hash, price, owner}
     * files[f_hash].list[i].owner => NodeInfo{pubKey, ipAddr}
     * files[f_hash].list[i].buyers[NodeInfo] => bool(true, false)
     *
     */

    struct NodeInfo {
        string pubKey;
        string ipAddr;
    }

    struct EncryptedFile {
        string hash;
        uint256 price;
        NodeInfo owner;

        /* to check if a client already bought the file */
        mapping(NodeInfo => bool) buyers;
    }

    struct FileInfo {
        uint32 size;
        mapping(uint32 => EncryptedFile) list;
    }

    mapping(string => FileInfo) files;

    /*
     *
     *
     * Need to implement
     * [modifier, constructor, functions]
     *
     *
     */
}
