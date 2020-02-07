package 개발관련파일들
var simple = simpleContract.new(42, {from:web3.eth.accounts[0], data: bytecode, gas: 0x47b760, privateFor: ["998c6d047b55f9d7237bd440f53a9266bf77c7d2656d4bcc2fe020f1e7e9bf09"]}, function(e, contract) {
if (e) {
console.log("err creating contract", e);
} else {
if (!contract.address) {
console.log("Contract transaction send: TransactionHash: " + contract.transactionHash + " waiting to be mined...");
} else {
console.log("Contract mined! Address: " + contract.address);
console.log(contract);
}
}
});