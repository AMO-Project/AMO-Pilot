const RDFSCoin = artifacts.require("./RDFSCoin.sol");
const RDFSFile = artifacts.require("./RDFSFile.sol");

module.exports = async(deployer, network, accounts) => {
    const owner   = accounts[0];
    const coin    = accounts[1];
    const file    = accounts[2];
    const userOne = accounts[3];
    const userTwo = accounts[4];

    await deployer.deploy(RDFSCoin, coin, { from: owner });
    const rdfsCoin = await RDFSCoin.deployed();

    await deployer.deploy(RDFSFile, file, rdfsCoin.address, { from: owner });
    await RDFSFile.deployed();
};
