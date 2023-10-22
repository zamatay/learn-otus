package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

var conn net.Conn

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	return nil
}

type tc struct {
}

func (t tc) Connect() error {
	connectionStr := fmt.Sprintf("%s:%s", programOptions.host, programOptions.port)
	var err error
	conn, err = net.Dial("tcp", connectionStr)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (t tc) Close() error {
	conn.Close()
}

func (t tc) Send() error {
	//TODO implement me
	panic("implement me")
}

func (t tc) Receive() error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
		return err
	}

	response := string(buf[:n])
	fmt.Println(response)

	return nil
}
