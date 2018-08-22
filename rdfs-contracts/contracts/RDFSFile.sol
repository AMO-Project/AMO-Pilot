pragma solidity ^0.4.2;

import "./RDFSCoin.sol";
import "zeppelin-solidity/contracts/math/SafeMath.sol";

contract RDFSFile {
    using SafeMath for uint256;

    string public constant symbol = "RDFS";
    string public constant name = "RDFS File";

    /*
     * Things to consider for storing and purchasing file
     *
     * files[ef_path] => FileInfo{name, size, owner, buyers, exists}
     * files[ef_path].owner => NodeInfo{pubKey, ipAddr}
     * files[ef_path].buyers[address] => purchaseState
     *
     * [modifiers]: checkDuplicateFile, checkDuplicatePurchase, onlyExistingFile,
     *              onlyAllowedAmount, onlyRequestedFile, onlyValidBuyer,
     *
     * [functions]: storeRequest,
     *              purchaseRequest, purchaseApprove, purchaseAbandon,
     *              isOwnedBy, isRequested, isApproved
     */

    enum purchaseState {
        None,
        Requested,
        Approved
    }

    struct Node {
        address addr;
        bytes4 ip;
    }

    struct File {
        string name;
        uint256 size;
        Node owner;

        /* to check the buyer's purchase state */
        mapping(address => purchaseState) buyers;
        bool exists;
    }

    mapping(bytes32 => File) files;

    address public adminAddr;
    RDFSCoin public token;

    modifier checkDuplicateFile(bytes32 _hash) {
        require(files[_hash].exists == false);
        _;
    }

    // check if msg.sender has already bought the file with given _hash
    modifier checkDuplicatePurchase(bytes32 _hash) {
        require(files[_hash].buyers[msg.sender] == purchaseState.None);
        _;
    }

    modifier onlyExistingFile(bytes32 _hash) {
        require(files[_hash].exists == true);
        _;
    }

    modifier onlyRequestedFile(bytes32 _hash) {
        require(files[_hash].buyers[msg.sender] == purchaseState.Requested);
        _;
    }

    modifier onlyValidBuyer(bytes32 _hash) {
        require(files[_hash].owner.addr != msg.sender);
        _;
    }

    modifier onlyAllowedAmount(bytes32 _hash, address target) {
        require(token.availableBalanceOf(target) >= files[_hash].size);
        _;
    }

    constructor(address _adminAddr, address _tokenAddr) public {
        adminAddr = _adminAddr;
        token = RDFSCoin(_tokenAddr);
    }

    function storeRequest(bytes32 _hash, string _name, uint256 _size, bytes4 _ip)
        public
        checkDuplicateFile(_hash)
        returns (bool)
    {
        /*
        files[_hash].name = _name;
        files[_hash].size = _size;
        files[_hash].owner = Node(msg.sender, _ip);
        files[_hash].exists = true;
        */

        files[_hash] = File(_name, _size, Node(msg.sender, _ip), true);

        return true;
    }

    function purchaseRequest(bytes32 _hash)
        public
        onlyExistingFile(_hash)
        onlyValidBuyer(_hash)
        onlyAllowedAmount(_hash, msg.sender)
        checkDuplicatePurchase(_hash)
        returns (bool)
    {
        token.addDeposit(files[_hash].size);
        files[_hash].buyers[msg.sender] = purchaseState.Requested;

        return true;
    }

    function purchaseApprove(bytes32 _hash)
        public
        onlyExistingFile(_hash)
        onlyValidBuyer(_hash)
        onlyRequestedFile(_hash)
        returns (bool)
    {
        token.transferDepositTo(files[_hash].owner.addr, files[_hash].size);

        files[_hash].buyers[msg.sender] = purchaseState.Approved;

        return true;
    }

    function purchaseAbandon(bytes32 _hash)
        public
        onlyExistingFile(_hash)
        onlyValidBuyer(_hash)
        onlyRequestedFile(_hash)
        returns (bool)
    {
        token.subDeposit(files[_hash].size);
        files[_hash].buyers[msg.sender] = purchaseState.None;

        return true;
    }

    function getFileName(bytes32 _hash) public view returns (string) {
        return files[_hash].name;
    }

    function getFileSize(bytes32 _hash) public view returns (uint256) {
        return files[_hash].size;
    }

    function isOwnedBy(bytes32 _hash) public view returns (address) {
        return files[_hash].owner.addr;
    }

    function isRequested(bytes32 _hash, address buyer) public view returns (bool) {
        return files[_hash].buyers[buyer] == purchaseState.Requested;
    }

    function isApproved(bytes32 _hash, address buyer) public view returns (bool) {
        return files[_hash].buyers[buyer] == purchaseState.Approved;
    }

    function getOwnerInfo(bytes32 _hash)
        public
        onlyExistingFile(_hash)
        onlyValidBuyer(_hash)
        onlyRequestedFile(_hash)
        view
        returns (address, bytes4) {
        return (files[_hash].owner.addr, files[_hash].owner.ip);
    }
}
