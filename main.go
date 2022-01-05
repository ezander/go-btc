package main

import (
	"bitcoin/network"
	"encoding/json"
	"flag"
	"fmt"
)

func AsJSON(object interface{}) string {
	json, _ := json.MarshalIndent(object, "", "\t")
	return string(json)
}

func test4() {
	numptr := flag.Int("ip", 0, "take n-th discovered ip address")
	versionptr := flag.Int("pver", 70015, "pretend to have that protocol version")
	flag.Parse()
	ipnum := *numptr
	version := uint32(*versionptr)

	client := network.TestClient(ipnum)
	defer client.Close()

	vermsg := network.NewVersionMessage()
	vermsg.Version = version
	client.SendMessage(vermsg)

	retmsg, command := client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))

	client.SendMessage(&network.VerAckMessage{})

	client.SendMessage(&network.GetBlocksMessage{Version: version, BlockLocHashes: []network.Hash{}, StopHash: network.Hash{}})
	for {
		retmsg, command = client.ReceiveMessage()
		fmt.Println("Received: ", command, AsJSON(retmsg))
	}
}

func main() {
	test4()
}
