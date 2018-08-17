const RDFSCoin = artifacts.require("./RDFSCoin.sol");
const RDFSFile = artifacts.require("./RDFSFile.sol");

contract("RDFS File Basic Test", (accounts) => {
    const owner   = accounts[0];
    const coin    = accounts[1];
    const file    = accounts[2];
    const userOne = accounts[3];
    const userTwo = accounts[4];

    let rdfsCoin = null;
    let rdfsFile = null;

    beforeEach("setup contract for each test", async() => {
        rdfsCoin = await RDFSCoin.new(coin, { from: owner });
        rdfsFile = await RDFSFile.new(file, rdfsCoin.address, { from: owner });
    });

    it("contract admin address should be set correctly", async() => {
        let checkAddr = await rdfsFile.adminAddr();
        assert.equal(checkAddr, file);
    });
});

contract("RDFS File Store/Purchase Test", (accounts) => {
    const owner   = accounts[0];
    const coin    = accounts[1];
    const file    = accounts[2];
    const userOne = accounts[3];
    const userTwo = accounts[4];

    let rdfsCoin = null;
    let rdfsFile = null;

    const hash = "0xE2BCA61C5EDC23C114CC34F923013AE98D6F70CF046D7B692AAEFB81EFEC51C8";
    const name = "test";
    const size = 123456789;
    const ip   = web3.fromUtf8("1234");

    beforeEach("setup contract for each test", async() => {
        rdfsCoin = await RDFSCoin.new(coin, { from: owner });
        rdfsFile = await RDFSFile.new(file, rdfsCoin.address, { from: owner });

        await rdfsCoin.transferFrom(owner, userOne, 200000000, { from: coin });
        await rdfsCoin.transferFrom(owner, userTwo, 200000000, { from: coin });
    });

    it("calling store request is available only when file is not in the map",
        async() => {
        try {
            await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });
            assert(true);
        } catch(err) {
            assert(err);
        }

        let fileOwner = await rdfsFile.isOwnedBy.call(hash);

        try {
            await rdfsFile.storeRequest(hash, name, size, ip, { from: userTwo });
            assert(false);
        } catch(err) {
            assert(err);
        }

        assert.equal(fileOwner, await rdfsFile.isOwnedBy.call(hash));
    });

    it("calling purchase request is available only when file is in the map",
        async() => {
        try {
            await rdfsFile.purchaseRequest(hash, { from: userTwo });
            assert(false);
        } catch(err) {
            assert(err);
        }

        await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });

        try {
            await rdfsFile.purchaseRequest(hash, { from: userTwo });
            assert(true);
        } catch(err) {
            assert(err);
        }

        assert.equal(true, await rdfsFile.isRequested.call(hash, userTwo), { from: userTwo });
    });

    it("calling purchase related functions are available "
        + "only when msg.sender is not file's owner",
        async() => {
        await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });

        try {
            await rdfsFile.purchaseRequest(hash, { from: userOne });
            assert(false);
        } catch(err) {
            assert(err);
        }

        try {
            await rdfsFile.purchaseRequest(hash, { from: userTwo });
            assert(true);
        } catch(err) {
            assert(err);
        }

        assert.equal(true, await rdfsFile.isRequested.call(hash, userTwo), { from: userTwo });
    });

    it("buyer cannot repurchase the file he requested or approved as purchase before",
        async() => {
        await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });

        await rdfsFile.purchaseRequest(hash, { from: userTwo });

        try {
            await rdfsFile.purchaseRequest(hash, { from: userTwo });
            assert(false);
        } catch(err) {
            assert(err);
        }

        await rdfsFile.purchaseApprove(hash, { from: userTwo });

        try {
            await rdfsFile.purchaseRequest(hash, { from: userTwo });
            assert(false);
        } catch(err) {
            assert(err);
        }
    });

    it("deposit amount should be same as a file's size after purchase request",
        async() => {
        await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });

        await rdfsFile.purchaseRequest(hash, { from: userTwo });

        let depoOfuserOne = await rdfsCoin.depositOf.call(userTwo, { from: userTwo });

        assert.equal(size, depoOfuserOne.toNumber());
    });

    it("after purchae approve, balance of buyer should be correct",
        async() => {
        const initBalOne = (await rdfsCoin.balanceOf(userOne)).toNumber();
        const initBalTwo = (await rdfsCoin.balanceOf(userTwo)).toNumber();

        await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });
        await rdfsFile.purchaseRequest(hash, { from: userTwo });
        await rdfsFile.purchaseApprove(hash, { from: userTwo });

        const finBalOne = (await rdfsCoin.balanceOf(userOne)).toNumber();
        const finBalTwo = (await rdfsCoin.balanceOf(userTwo)).toNumber();

        assert.equal(initBalOne + size, finBalOne);
        assert.equal(initBalTwo - size, finBalTwo);
    });

    it("after purchase abandon, balance of buyer shoud be correct",
        async() => {
        const initBalOne = (await rdfsCoin.balanceOf(userOne)).toNumber();
        const initBalTwo = (await rdfsCoin.balanceOf(userTwo)).toNumber();

        await rdfsFile.storeRequest(hash, name, size, ip, { from: userOne });
        await rdfsFile.purchaseRequest(hash, { from: userTwo });
        await rdfsFile.purchaseAbandon(hash, { from: userTwo });

        const finBalOne = (await rdfsCoin.balanceOf(userOne)).toNumber();
        const finBalTwo = (await rdfsCoin.balanceOf(userTwo)).toNumber();

        assert.equal(initBalOne, finBalOne);
        assert.equal(initBalTwo, finBalTwo);

        assert.equal(false, await rdfsFile.isRequested.call(hash, userTwo));
        assert.equal(false, await rdfsFile.isApproved.call(hash, userTwo));
    });
});
