const RDFSCoin = artifacts.require("./RDFSCoin.sol");

contract("RDFS Coin Basic Test", (accounts) => {
    const owner   = accounts[0];
    const coin    = accounts[1];
    const file    = accounts[2];
    const userOne = accounts[3];
    const userTwo = accounts[4];

    let token = null;

    beforeEach("setup contract for each test", async() => {
        token = await RDFSCoin.new(coin, { from: owner });
    });

    it("contract admin address should be set correctly", async() => {
        let checkAddr = await token.adminAddr();
        assert.equal(checkAddr, coin);
    });
});

contract("RDFS Coin Transfer Test", (accounts) => {
    const owner   = accounts[0];
    const coin    = accounts[1];
    const file    = accounts[2];
    const userOne = accounts[3];
    const userTwo = accounts[4];

    let token = null;

    beforeEach("setup contract for each test", async() => {
        token = await RDFSCoin.new(coin, { from: owner });
    });

    it("calling transfer is available only when destination is valid",
        async() => {
        try {
            await token.transfer(0x0, 10);
            assert(false);
        } catch(err) {
            assert(err);
        }
    });

    it("calling transfer is available only when sender has more than sending value",
        async() => {
        try {
            await token.transfer(owner, 100, {from: userOne });
            assert(false);
        } catch(err) {
            assert(err);
        }
    });

    it("calling transfer is available only when sender's available balance"
        + "(balances[sender] - deposit[sender]) is bigger than sending value",
        async() => {
        await token.transfer(userOne, 100, { from: owner });
        await token.addDeposit(60, { from: userOne });

        let abOfuserOne = await token.availableBalanceOf.call(userOne, { from: userOne });

        assert(40, abOfuserOne.toNumber());

        try {
            await token.transfer(userTwo, 50, { from: userOne });
            assert(false);
        } catch(err) {
            assert(err);
        }

        await token.transferDepositTo(userTwo, 50, { from: userOne });

        let bOfuserOne = await token.balanceOf(userOne, { from: userOne });
        let bOfuserTwo = await token.balanceOf(userTwo, { from: userTwo });

        assert(50, bOfuserOne.toNumber());
        assert(50, bOfuserTwo.toNumber());

        abOfuserOne = await token.availableBalanceOf.call(userOne, { from: userOne });

        assert(50, abOfuserOne.toNumber());
    })

    it("balance of destination should be same as transferred value",
        async() => {
        const initialBalance = (await token.balanceOf(owner)).toNumber();

        await token.transfer(userOne, 100, { from: owner });

        const finalBalance = (await token.balanceOf(owner)).toNumber();
        const userBalance = (await token.balanceOf(userOne)).toNumber();

        assert.equal(initialBalance - 100, finalBalance);
        assert.equal(userBalance, 100);
    });
})
