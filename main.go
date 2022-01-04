package main

import (
	"bitcoin/network"
	"encoding/json"
	"fmt"
)

func AsJSON(object interface{}) string {
	json, _ := json.MarshalIndent(object, "", "\t")
	return string(json)
}

func test4() {
	conn := network.GetTestConn(6)
	client := network.Client(conn)
	defer client.Close()

	var vermsg network.Message = network.NewVersionMessage()
	var packet = network.CreatePacket(network.MAGIC_testnet3, "version", vermsg)
	fmt.Println("Sending: ", AsJSON(packet))
	client.SendPacket(packet)
	retmsg := client.ReadPacket()
	fmt.Println("Received: ", AsJSON(retmsg))

	vermsg = &network.VerAckMessage{}
	packet = network.CreatePacket(network.MAGIC_testnet3, "verack", vermsg)
	fmt.Println("Sending: ", AsJSON(packet))
	client.SendPacket(packet)
	retmsg = client.ReadPacket()
	fmt.Println("Received: ", AsJSON(retmsg))

}

func main() {
	test4()
}
