pragma solidity ^0.4.2;

import "openzeppelin-solidity/contracts/token/ERC20/StandardToken.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract RDFSCoin is StandardToken, Ownable {
    using SafeMath for uint256;

    string public constant symbol = "RDFS";
    string public constant name = "RDFS Coin";

    uint8 public constant decimals = 18;
    uint256 public constant INITIAL_SUPPLY = 20000000000 * (10 ** uint256(decimals));

    // Address of token administrator
    address public adminAddr;

    // Deposit balacnes
    mapping(address => uint256) deposit;

    /*
     * Check if token transfer destination is valid
     */
    modifier onlyValidDestination(address to) {
        require(to != address(0x0)
            && to != address(this)
            && to != owner
            && to != adminAddr);
        _;
    }

    modifier onlyAllowedAmount(address target, uint256 amount) {
        require(balances[target].sub(deposit[target]) >= amount);
        _;
    }

    /*
     * Constructor of RDFSCoin contract
     * @param _adminAddr: Address of token administrator
     */
    constructor(address _adminAddr) public {
        totalSupply_ = INITIAL_SUPPLY;

        balances[msg.sender] = totalSupply_;
        emit Transfer(address(0x0), msg.sender, totalSupply_);

        adminAddr = _adminAddr;
        approve(adminAddr, INITIAL_SUPPLY);
    }

    /*
     * Transfer token from message sender to another
     * @param to: Destination address
     * @param value: Amount of RDFS Coin to transfer
     */
    function transfer(address to, uint256 value)
        public
        onlyValidDestination(to)
        onlyAllowedAmount(tx.origin, value)
        returns (bool)
    {
        return super.transfer(to, value);
    }

    function transferFrom(address from, address to, uint256 value)
        public
        onlyValidDestination(to)
        onlyAllowedAmount(from, value)
        returns (bool)
    {
        return super.transferFrom(from, to, value);
    }

    function transferDepositTo(address to, uint256 value)
        public
        onlyValidDestination(to)
        returns (bool)
    {
        require(deposit[tx.origin].sub(value) >= 0);

        // release and transfer deposit from 'from' address to 'to' address
        balances[tx.origin] = balances[tx.origin].sub(value);
        deposit[tx.origin] = deposit[tx.origin].sub(value);
        balances[to] = balances[to].add(value);

        emit Transfer(tx.origin, to, value);
        return true;
    }

    function addDeposit(uint256 value) public returns (bool) {
        require(balances[tx.origin].sub(deposit[tx.origin]) >= value);

        deposit[tx.origin] = deposit[tx.origin].add(value);
        return true;
    }

    function subDeposit(uint256 value) public returns (bool) {
        require(deposit[tx.origin].sub(value) >= 0);

        deposit[tx.origin] = deposit[tx.origin].sub(value);
        return true;
    }

    function depositOf(address target) public view returns (uint256) {
        return deposit[target];
    }

    function availableBalanceOf(address target) public view returns (uint256) {
        return balances[target].sub(deposit[target]);
    }
}
