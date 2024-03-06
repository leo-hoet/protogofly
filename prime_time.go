package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
)

var ErrMethodUnknown = errors.New("MethodNotPrime")

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}

	upperLimit := int64(math.Ceil(math.Sqrt(float64(n))))
	for i := range upperLimit {
		if i+1 == 1 {
			continue
		}
		if (i + 1) == n {
			return true
		}
		if n%(i+1) == 0 {
			return false
		}
	}

	return true
}

func HandleRequestIsPrime(r Request) (bool, error) {
	if r.Method != "isPrime" {
		return false, ErrMethodUnknown
	}

	isInteger := math.Round(*r.Number) == *r.Number
	if !isInteger {
		return false, nil
	}

	return isPrime(int64(*r.Number)), nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	encoder := json.NewEncoder(conn)
	c := bufio.NewReader(conn)

	for {

		reqRaw, err := c.ReadBytes('\n')
		if err != nil {
			fmt.Fprintf(conn, "Error read bytes")
			return
		}

		log.Printf("Processing: %+v %+v", string(reqRaw[:]), reqRaw)

		var req Request

		err = json.Unmarshal(reqRaw[:], &req)
		if err != nil {
			fmt.Fprint(conn, "Error decoding json")
			return
		}

		if req.Number == nil {
			fmt.Fprintf(conn, "Error, empty number")
			return
		}

		isPrime, err := HandleRequestIsPrime(req)
		if err != nil {
			fmt.Fprint(conn, "Error: ", err)
			return
		}

		encoder.Encode(Response{Method: "isPrime", Prime: isPrime})
	}
}

func StartServer() {
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
