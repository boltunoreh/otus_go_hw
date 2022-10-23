package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (client *Client) Connect() error {
	conn, err := net.DialTimeout("tcp", client.address, client.timeout)
	if err != nil {
		return err
	}

	client.conn = conn

	return nil
}

func (client *Client) Close() error {
	if client.conn == nil {
		return errors.New("connection is null")
	}

	if err := client.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (client *Client) Send() error {
	if client.conn == nil {
		return errors.New("connection is null")
	}

	err := client.handleMessage(client.in, client.conn)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) Receive() error {
	if client.conn == nil {
		return errors.New("connection is null")
	}

	err := client.handleMessage(client.conn, client.out)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleMessage(from io.Reader, to io.Writer) error {
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		_, err := to.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
