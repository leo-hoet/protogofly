package main

import (
	"bytes"
	"errors"
	"net"
	"testing"
)

func TestServerReq(t *testing.T) {
	go StartServer()

	tt := []struct {
		test string
		req  []byte
		res  []byte
	}{
		{
			test: "Ok req",
			req:  []byte("{\"method\":\"isPrime\",\"number\":2}\n"),
			res:  []byte("{\"method\":\"isPrime\",\"prime\":true}\n"),
		},
		{
			test: "Ok req no prime",
			req:  []byte("{\"method\":\"isPrime\",\"number\":4}\n"),
			res:  []byte("{\"method\":\"isPrime\",\"prime\":false}\n"),
		},
		{
			test: "Ok req float",
			req:  []byte("{\"method\":\"isPrime\",\"number\":4.01}\n"),
			res:  []byte("{\"method\":\"isPrime\",\"prime\":false}\n"),
		},
		{
			test: "Malformed req, number as string",
			req:  []byte("{\"method\":\"isPrime\",\"number\":\"4\"}\n"),
			res:  []byte("Error decoding json"),
		},
		{
			test: "Malformed req. Missing number",
			req:  []byte("{\"method\": \"isPrime\"}\n"),
			res:  []byte("Error, empty number\n"),
		},
		{
			test: "Malformed req. Line break in req",
			req:  []byte("{\"method\": \n \"isPrime\"}\n"),
			res:  []byte("Error decoding json\n"),
		},
		{
			test: "Anything",
			req:  []byte("asdasd\n"),
			res:  []byte("Error decoding json\n"),
		},
		{
			test: "Big payload",
			req: []byte(
				"{\"method\":\"isPrime\",\"number\":1}\n" +
					"{\"method\":\"isPrime\",\"number\":2}\n" +
					"{\"method\":\"isPrime\",\"number\":3}\n" +
					"{\"method\":\"isPrime\",\"number\":4}\n",
			),
			res: []byte(
				"{\"method\":\"isPrime\",\"prime\":false}\n" +
					"{\"method\":\"isPrime\",\"prime\":true}\n" +
					"{\"method\":\"isPrime\",\"prime\":true}\n" +
					"{\"method\":\"isPrime\",\"prime\":false}\n",
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.test, func(t *testing.T) {
			conn, err := net.Dial("tcp", "localhost:5000")
			if err != nil {
				t.Fatal("Server is not running")
			}
			defer conn.Close()
			if _, err := conn.Write(tc.req); err != nil {
				t.Error("could not write payload to TCP server:", err)
			}

			out := make([]byte, 1024)

			if n, err := conn.Read(out); err == nil {
				read := out[:n]
				if !bytes.Contains(tc.res, read) {
					t.Errorf("Response did not match expected output: out: \n %+v\n expected: \n%+v\n", string(read), string(tc.res[:]))
				}
			} else {
				t.Error("could not read from connection")
			}

		})
	}

}

func addrF64(f64 float64) *float64 {
	return &f64
}
func TestPrimeTimeHandleReq(t *testing.T) {

	testCases := []struct {
		name    string
		req     Request
		wantRes bool
		wantErr error
	}{
		{
			name:    "Invalid method",
			req:     Request{Method: "is_prime", Number: addrF64(1)},
			wantRes: false,
			wantErr: ErrMethodUnknown,
		},
		{

			name:    "Number is float",
			req:     Request{Method: "isPrime", Number: addrF64(1.1)},
			wantRes: false,
			wantErr: nil,
		},
		{
			name:    "Number is prime",
			req:     Request{Method: "isPrime", Number: addrF64(3)},
			wantRes: true,
			wantErr: nil,
		},
		{
			name:    "Number is integer not prime",
			req:     Request{Method: "isPrime", Number: addrF64(4)},
			wantRes: false,
			wantErr: nil,
		},
		{
			name:    "Number is negative",
			req:     Request{Method: "isPrime", Number: addrF64(-1)},
			wantRes: false,
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		got, err := HandleRequestIsPrime(tc.req)

		if (got != tc.wantRes) || !errors.Is(err, tc.wantErr) {
			t.Errorf("Name: '%s' got: %t %s, expected: %t %s", tc.name, got, err, tc.wantRes, tc.wantErr)
		}
	}
}
