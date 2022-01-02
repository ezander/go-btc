package main

import (
	"bitcoin/network"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

func test3() {
	// https://bitcoin.stackexchange.com/questions/49634/testnet-peers-list-with-ip-addresses
	// dig A testnet-seed.bitcoin.jonasschnelli.ch

	ips, _ := net.LookupIP("testnet-seed.bitcoin.jonasschnelli.ch")

	tcp := net.TCPAddr{IP: ips[0], Port: 18333, Zone: ""}
	conn, err := net.Dial("tcp", tcp.String())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println(conn, err)
	conn.SetDeadline(time.Now().Add(time.Second * 2))

	vermsg := network.NewVersionMessage()
	var msg = network.CreatePacket(network.MAGIC_testnet3, "version", &vermsg)

	out := network.MarshalPacket(nil, msg)
	fmt.Printf("%+v\n", msg.Message)
	// fmt.Println(out)

	_, err = conn.Write(out)
	if err != nil {
		panic(err)
	}

	in, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println(err) // i/o timeout is okay...
	}

	retmsg, _ := network.UnmarshalPacket(in)
	fmt.Printf("%+v\n", retmsg.Message)

}

func main() {
	test3()
}
