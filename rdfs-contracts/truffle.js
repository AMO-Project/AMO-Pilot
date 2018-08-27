/*
 * NB: since truffle-hdwallet-provider 0.0.5 you must wrap HDWallet providers in a
 * function when declaring them. Failure to do so will cause commands to hang. ex:
 * ```
 * mainnet: {
 *     provider: function() {
 *       return new HDWalletProvider(mnemonic, 'https://mainnet.infura.io/<infura-key>')
 *     },
 *     network_id: '1',
 *     gas: 4500000,
 *     gasPrice: 10000000000,
 *   },
 */

module.exports = {
	networks: {
		development: {
			host: "127.0.0.1",
			port: 9545,
			network_id: "*",
			gas: 3000000,
			gasPrice : 10000000
		},
		rdfs: {
			host: "127.0.0.1",
			port: 8545,
			network_id: 208518,
			from: "0x1d8582d7d5c85a9b65d9beee856f5408aae96215",
			gas: 3000000,
			gasPrice : 10000000
		}
	}
};
