package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/OmarTariq612/torcs-go-client/client"
	"github.com/OmarTariq612/torcs-go-client/controller"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:3001", "host:port of the server")
	flag.Parse()

	udpAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		fmt.Printf("ERR (resolve udp addr): %v\n", err)
		os.Exit(1)
	}

	client := client.New(udpAddr, controller.NewSimpleDriver())
	if err := client.Start(); err != nil {
		fmt.Printf("ERR (client): %v\n", err)
	}
}
