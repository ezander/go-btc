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

	vermsg := network.NewVersionMessage()
	var packet = network.CreatePacket(network.MAGIC_testnet3, "version", &vermsg)
	fmt.Printf("%+v\n", packet)
	fmt.Println(AsJSON(packet))

	conn := network.GetTestConn(1)
	client := network.Client(conn)
	defer client.Close()

	client.SendPacket(packet)
	retmsg := client.ReadPacket()

	fmt.Printf("%+v\n", retmsg)
	fmt.Println(AsJSON(retmsg))
}

func main() {
	test4()
}
