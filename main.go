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
	flag.Parse()

	client := network.TestClient(*numptr)
	defer client.Close()

	var vermsg network.Message = network.NewVersionMessage()
	client.SendMessage(vermsg)
	retmsg, command := client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))

	client.SendMessage(&network.VerAckMessage{})
	for {
		retmsg, command = client.ReceiveMessage()
		fmt.Println("Received: ", command, AsJSON(retmsg))
	}
}

func main() {
	test4()
}
