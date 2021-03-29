// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{

	// ReapChain Primary Bootnode
	"enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@116.126.85.23:30301",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{
	//"enode://6ce05930c72abc632c58e2e4324f7c7ea478cec0ed4fa2528982cf34483094e9cbc9216e7aa349691242576d552a2a56aaeae426c5303ded677ce455ba1acd9d@13.84.180.240:30303", // US-TX
	//"enode://20c9ad97c081d63397d7b685a412227a40e23c8bdc6688c6f37e97cfbc22d2b4d1db1510d8f61e6a8866ad7f0e17c02b14182d37ea7c3c8b9c2683aeb6b733a1@52.169.14.227:30303", // IE
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{
	//"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303", // IE
}

// RinkebyV5Bootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network for the experimental RLPx v5 topic-discovery network.
var RinkebyV5Bootnodes = []string{
	//"enode://a24ac7c5484ef4ed0c5eb2d36620ba4e4aa13b8c84684e1b4aab0cebea2ae45cb4d375b77eab56516d34bfbd3c1a833fc51296ff084b770b94fb9028c4d25ccf@52.169.42.101:30303?discport=30304", // IE
}

// Ottoman are the enode URLs of the P2P bootstrap nodes running on the
// Ottoman test network.
var OttomanBootnodes = []string{

}

var ReapChainBootnodes = []string{
	// ReapChain Primary Bootnode
	//"enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@116.126.85.23:30301",
	"enode://78321f807436cd3632753f44153dfef3e61bd8cb59ed366fc3c81f73df2293640e6d28fe1b54d56476b68acf7df8bf847ce7e109f904dac0ffbcb04574c1e976@192.168.0.98:30301",
	"enode://3f93f712a5e0676f2334bdfcbeff42c887b790a5b5367bbf06c5fce83c202cf682a6c00f862342545086dadd58a870925c4b8b82ab752ef912453e629400e8c2@192.168.0.98:30302",
	"enode://532079b5edd286744b02c96fbcf85827e92189174f57bbcf04f7409600b6848f87e338156a569566f703688dc7ff4cf9ba2e102bc7e04dc3f4b884f4bf54e515@192.168.0.98:30303",
}
var ReapChainQMannodes = []string{
	"enode://d69e8911b2e081fab00660b62423ffaedb84f9c6fdec3df677cb6e66e765f1439a9aed94a92dab83e8af99c4c54a53954c19c7b692e52c0a51953e9631b5e344@192.168.0.98:30500",
	"enode://df8aba7c4a1db494f77dde477e4e2792acf420607e6e7bd5a8489973ea9b537e551e5c0d5d356d43901749ffb62185f7e9b44308e7e0ddb22fd1b1084fcf7832@192.168.0.98:30409",
	"enode://3666e1c8498a081a0acb62d08fd6204b5baaa5522d8d7207273016482dde31cd3b5d698ed24de213fbbf27a164dae004001c8fbc6af429f089999eda2f07ecb2@192.168.0.98:30408",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
	"enode://5d686a07e38d2862322a2b7e829ee90c9931f119391c63328cab0d565067835808e46cb16dc2a0e920cf1a6a68806e6129b986b6b143cdb7d0752dec45a7f12c@116.126.85.23:30301",
}
