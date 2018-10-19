package core

import (
    "net"
    "bufio"
)

var delim = byte('\n')

// Client holds info about connection
type Client struct {
    conn   net.Conn
    Server *TCPServer
}

// Read client data from channel
func (c *Client) listen() {
    reader := bufio.NewReader(c.conn)
    for {
        message, err := reader.ReadSlice(delim)
        if err != nil {
            c.conn.Close()
            c.Server.onClientConnectionClosed(c, err)
            return
        }
        c.Server.onNewMessage(c, message)
    }
}

// Send text message to client
func (c *Client) SendString(message string) error {
    _, err := c.conn.Write([]byte(message))
    return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
    _, err := c.conn.Write(b)
    return err
}

func (c *Client) Conn() net.Conn {
    return c.conn
}

func (c *Client) Close() error {
    return c.conn.Close()
}

// TCP TCPServer
type TCPServer struct {
    address                  string
    onNewClientCallback      func(c *Client)
    onClientConnectionClosed func(c *Client, err error)
    onNewMessage             func(c *Client, message []byte)
}

// Called right after TCPServer starts listening new client
func (s *TCPServer) OnNewClient(callback func(c *Client)) {
    s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *TCPServer) OnClientConnectionClosed(callback func(c *Client, err error)) {
    s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *TCPServer) OnNewMessage(callback func(c *Client, message []byte)) {
    s.onNewMessage = callback
}

// StartByAdmin network TCPServer
func (s *TCPServer) Listen() {
    listener, err := net.Listen("tcp", s.address)
    if err != nil {

    }
    defer listener.Close()

    for {
        conn, _ := listener.Accept()
        client := &Client{
            conn:   conn,
            Server: s,
        }
        go client.listen()
        s.onNewClientCallback(client)
    }
}

// Creates new tcp TCPServer instance
func NewTCPServer(address string) *TCPServer {
    server := &TCPServer{
        address: address,
    }

    server.OnNewClient(func(c *Client) {})
    server.OnNewMessage(func(c *Client, message []byte) {})
    server.OnClientConnectionClosed(func(c *Client, err error) {})

    return server
}
