package main

import (
	. "bitcoin/network"

	"flag"
	"fmt"
)

func test4() {
	ipnumPtr := flag.Int("ip", 2, "take n-th discovered ip address")
	versionPtr := flag.Int("pver", 69999, "pretend to have that protocol version")
	flag.Parse()
	ipnum := *ipnumPtr
	version := uint32(*versionPtr)

	client := TestClient(ipnum)
	defer client.Close()

	vermsg := NewVersionMessage()
	vermsg.Version = version
	client.SendMessage(vermsg)

	retmsg, command := client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))

	client.SendMessage(&VerAckMessage{})

	genesisBlockHash, _ := StringToHash("000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943")
	// client.SendMessage(&GetBlocksMessage{Version: version, BlockLocHashes: []Hash{}, StopHash: Hash{}})
	client.SendMessage(&GetHeadersMessage{Version: version, BlockLocHashes: []Hash{genesisBlockHash}, StopHash: Hash{}})
	for {
		retmsg, command = client.ReceiveMessage()
		fmt.Println("Received: ", command, AsJSON(retmsg))
	}
}

func main() {
	test4()
}

// Main:
// uint256 hashGenesisBlock("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f");
// Testnet:
// hashGenesisBlock = uint256("000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943");
// block.nTime    = 1296688602;
// block.nNonce   = 414098458;

//  // Genesis block
//  const char* pszTimestamp = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks";
//  CTransaction txNew;
//  txNew.vin.resize(1);
//  txNew.vout.resize(1);
//  txNew.vin[0].scriptSig = CScript() << 486604799 << CBigNum(4) << vector<unsigned char>((const unsigned char*)pszTimestamp, (const unsigned char*)pszTimestamp + strlen(pszTimestamp));
//  txNew.vout[0].nValue = 50 * COIN;
//  txNew.vout[0].scriptPubKey = CScript() << ParseHex("04678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f") << OP_CHECKSIG;
//  CBlock block;
//  block.vtx.push_back(txNew);
//  block.hashPrevBlock = 0;
//  block.hashMerkleRoot = block.BuildMerkleTree();
//  block.nVersion = 1;
//  block.nTime    = 1231006505;
//  block.nBits    = 0x1d00ffff;
//  block.nNonce   = 2083236893;

//  if (fTestNet)
//  {
// 		 block.nTime    = 1296688602;
// 		 block.nNonce   = 414098458;
//  }

// printf("%s\n", block.GetHash().ToString().c_str());
// printf("%s\n", hashGenesisBlock.ToString().c_str());
// printf("%s\n", block.hashMerkleRoot.ToString().c_str());
// assert(block.hashMerkleRoot == uint256("0x4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"));
// block.print();
// assert(block.GetHash() == hashGenesisBlock);
