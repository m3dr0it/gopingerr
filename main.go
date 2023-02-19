package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	var ipDb string = "192.168.1.1"
	for {
		if pingWithIcmp(ipDb) {
			exec.Command("ls")
			break
		}
		log.Println("not connected")
	}
}

func pingWithIcmp(ip string) bool {
	isConnected := false
	packetCon, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err.Error())
	}

	icmpMessage := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  0,
			Data: []byte("hello"),
		},
	}

	wb, err := icmpMessage.Marshal(nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	dst, err := net.ResolveIPAddr("ip", "192.168.1.102")
	if err != nil {
		log.Fatal(err.Error())
	}

	timeout, _ := time.ParseDuration("1s")
	err = packetCon.SetDeadline(time.Now().Add(timeout))
	_, err = packetCon.WriteTo(wb, dst)

	if err != nil {
		fmt.Println(err.Error())
	}

	rb := make([]byte, 1500)

	n, peer, err := packetCon.ReadFrom(rb)

	if err != nil {
		return false
	}

	rm, err := icmp.ParseMessage(1, rb[:n])

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rm)

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("received from %v", peer)
		isConnected = true
	default:
		fmt.Printf("Failed: %+v\n", rm)
		isConnected = false
	}
	return isConnected

}
