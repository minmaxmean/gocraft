package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/m-nny/goinit/pkg/mcnet"
)

var (
	host = flag.String("host", "localhost", "ip host")
	port = flag.Uint("port", 8080, "port")
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Printf("Server started at %s:%d", *host, *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		client := mcnet.NewClient(conn)

		log.Printf("Got connection: %v\n", conn)
		if err := client.Welcome(); err != nil {
			log.Printf("err: %v", err)
		}
		client.Close()
	}
}
