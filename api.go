package core

import (
	"fmt"
	"net"
)

var api_server *TCPServer
var api_delim = "@$@"

func start_api_server(conn_str string) {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", conn_str)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

var block_size = 1408

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		ERROR(err.Error())
	}

	// handle
	ans := HandleMsg(buf)

	conn.Write(ans)

	//enc := EasyEncrypt(ans)
	//
	//if len(enc) < block_size {
	//	// encrypt
	//	enc = append(enc, delim)
	//	conn.Write(enc)
	//} else {
	//	idx_start := 0
	//	idx_stop := block_size
	//	for {
	//		if idx_stop < len(enc) {
	//			bd := enc[idx_start:idx_stop]
	//			// encrypt
	//			conn.Write(bd)
	//		} else {
	//			bd := enc[idx_start:]
	//			// encrypt
	//			enc = append(enc, delim)
	//			conn.Write(bd)
	//			break
	//		}
	//
	//		idx_start = idx_stop
	//		idx_stop += block_size
	//	}
	//
	//}
	// Close the connection when you're done with it.
	//conn.Close()
}

func HandleMsg(msg []byte) []byte {
	var args = []interface{}{}
	err := MsgUnpack(msg, &args)
	if err != nil {
		msg := fmt.Sprintf("core: api: HandleMsg: %v", err.Error())
		ERROR(msg)
		return nil
	}

	if len(args) < 1 {
		msg := fmt.Sprintf("core: api: HadleMsg: method idx not found")
		ERROR(msg)
		return nil
	}
	method_name := ToString(args[0])
	handler := getHandler(method_name)

	fmt.Println("called method [", method_name, "]", "args:", args[1:])

	var result interface{}
	if handler != nil {
		param := []interface{}{}
		if len(args) > 1 {
			param = args[1:]
		}
		result = handler(param)
		return MsgPack(result)
	} else {
		result = []interface{}{500, fmt.Sprintf("method [%v] is not found", method_name)}
		return MsgPack(result)
	}

}

type HandlerFoo func(args []interface{}) interface{}
type HandlerItem struct {
	name    string
	handler HandlerFoo
}

var api_list = []HandlerItem{}

func setAPI(a map[string]HandlerFoo) {
	for name, handler := range a {
		item := HandlerItem{name, handler}
		api_list = append(api_list, item)
	}
}

func API() (ls []string) {
	for _, item := range api_list {
		ls = append(ls, item.name)
	}
	return ls
}

func getHandler(method string) HandlerFoo {
	for _, hand := range api_list {
		if hand.name == method {
			return hand.handler
		}
	}
	return nil
}

//func onNewAPIClient(c *Client) {
//}
//
//func onNewMessageAPIClient(c *Client, message []byte) {
//   answer := HandleMsg(message)
//
//   answer = append(answer, delim)
//   _, err := c.conn.Write(answer)
//   if err != nil {
//       ERROR("core: api: onNewMessageAPIclient: " + err.Error())
//   }
//}
//
//func OnAPIClientConnectionClosed(c *Client, err error) {
//}
