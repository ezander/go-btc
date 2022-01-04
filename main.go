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
	// client := network.TestClient(6)
	client := network.TestClient(5)
	defer client.Close()

	var vermsg network.Message = network.NewVersionMessage()
	client.SendMessage(vermsg)
	retmsg, command := client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))

	client.SendMessage(&network.VerAckMessage{})
	retmsg, command = client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))

	retmsg, command = client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))

	retmsg, command = client.ReceiveMessage()
	fmt.Println("Received: ", command, AsJSON(retmsg))
}

func main() {
	test4()
}
