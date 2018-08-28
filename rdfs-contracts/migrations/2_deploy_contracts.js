const RDFSCoin = artifacts.require("./RDFSCoin.sol");
const RDFSFile = artifacts.require("./RDFSFile.sol");

module.exports = function(deployer, network, accounts) {
	const userOne = accounts[0];
	const owner   = accounts[1];
	const coin    = accounts[2];
	const file    = accounts[3];

	const userTwo   = "0x28742aaa4f8a4c6fb31e3a3e3fb85355e3b5926b";
	const userThree = "0x82496a989c83ccd7c58f66934992c3c54f724935";
	
	//const userTwo = accounts[4];
	//const userThree = accounts[5];

	let rdfsCoin;
	let rdfsFile;

    deployer.deploy(RDFSCoin, coin, { from: owner })
    	.then(() => {
			return RDFSCoin.deployed().then(instance => {
				rdfsCoin = instance;
			});
		}).then(() => deployer.deploy(RDFSFile, file, rdfsCoin.address, { from: owner }))
    	.then(() => {
			return RDFSFile.deployed().then(instance => {
				rdfsFile = instance;
			});
		}).then(() => {
			// Transferring some token to users for test
			//console.log("Transferring 200000000 token to userOne, Two, Three")
		}).then(() => {
			//rdfsCoin.transferFrom(owner, userOne, 200000000, { from: coin });
			//rdfsCoin.transferFrom(owner, userTwo, 200000000, { from: coin });
			//rdfsCoin.transferFrom(owner, userThree, 200000000, { from: coin });
			//console.log("GOooooood")
		});
};
