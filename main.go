package main

import (
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	io.Copy(conn, conn)
}

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalln("Could not open tcp server: ", err.Error())
	}
	defer listener.Close()

	log.Println("Starting TCP server in 5000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
