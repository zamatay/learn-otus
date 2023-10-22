package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

var (
	ErrorConnectIsEmpty   = errors.New("ConnectionIsEmpty")
	ErrorEOF              = errors.New("EOF")
	ErrorConnectionClosed = errors.New("ConnectionClose")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &tc{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type tc struct {
	address    string
	timeout    time.Duration
	conn       net.Conn
	connection bool
	scan       *bufio.Scanner
	inScan     *bufio.Scanner
	in         io.Reader
	out        io.Writer
}

func (t *tc) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		log.Fatal("connect", err)
		return err
	}
	t.scan = bufio.NewScanner(t.conn)
	t.inScan = bufio.NewScanner(t.in)
	return nil
}

func (t *tc) Close() error {
	if t.conn != nil {
		if err := t.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (t *tc) Send() error {
	if t.conn == nil {
		return ErrorConnectIsEmpty
	}
	if !t.inScan.Scan() {
		return ErrorEOF
	}
	if _, err := t.conn.Write(append(t.inScan.Bytes(), '\n')); err != nil {
		return err
	}
	return nil
}

func (t *tc) Receive() error {
	if t.conn == nil {
		return ErrorConnectIsEmpty
	}
	if !t.scan.Scan() {
		return ErrorConnectionClosed
	}
	var err error
	if _, err = t.out.Write(append(t.scan.Bytes(), '\n')); err != nil {
		return err
	}
	return nil
}
