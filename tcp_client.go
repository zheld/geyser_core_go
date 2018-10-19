package core

import (
	"net"
	"fmt"
	"github.com/fatih/pool"
	"sync"
)

type TCPClient struct {
	conn *net.TCPConn
	addr string
}

var client_stor = map[string]*TCPClient{}
var pl pool.Pool

func InitTCPClientPool(addr string) {
	var err error
	//create a factory() to be used with channel based pool
	factory := func() (net.Conn, error) {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			println("ResolveTCPAddr failed:", err.Error())
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			println("Dial failed:", err.Error())
		}

		return conn, err
	}

	//create a new channel based pool with an initial capacity of 5 and maximum
	//capacity of 30. The factory will create 5 initial connections and put it
	//into the pool.
	pl, err = pool.NewChannelPool(1, 2, factory)
	if err != nil {
		fmt.Println("error:", err.Error())
	}
	fmt.Println("init connection pool")
}

func NewTCPClient(addr string) *TCPClient {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
	}

	tcp := &TCPClient{
		conn: conn,
		addr: addr,
	}

	return tcp
}

var mtx = sync.Mutex{}

func SendFromPool(method, payload string) (resp []byte, err error) {
	req := fmt.Sprintf("{\"method\":\"%s\",\"payload\":%s}", method, payload)

	data := []byte(req)
	data = append(data, 0)

	mtx.Lock()
	conn, err := pl.Get()
	if err != nil {
		fmt.Println(err.Error())
	}
	mtx.Unlock()

	_, err = conn.Write(data)
	if err != nil {
		println("Write to server failed:", err.Error())
		return resp, err
	}

	reply := make([]byte, 1024)

	read_len, err := conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		return resp, err
	}

	reply = reply[:read_len]
	//fmt.Println(string(reply))
	conn.Close()

	return reply, nil
}

func CloseTCPPool() {
	pl.Close()
}

func CallHost(host, method string, args ... interface{}) (resp []interface{}, err error) {
	req := []interface{}{method}
	req = append(req, args...)
	breq := MsgPack(req)
	breq = append(breq, 0)

	var tcp_cli *TCPClient

	if cli, ok := client_stor[host]; ok {
		tcp_cli = cli
	} else {
		tcp_cli = NewTCPClient(host)
	}

	if tcp_cli.conn != nil {
		_, err = tcp_cli.conn.Write(breq)
		if err != nil {
			println("Write to server failed:", err.Error())
			return resp, err
		}

		reply := make([]byte, 1024)

		read_len, err := tcp_cli.conn.Read(reply)
		if err != nil {
			println("Write to server failed:", err.Error())
			return resp, err
		}

		reply = reply[:read_len]

		MsgUnpack(reply, &resp)

		code := resp[0].(uint64)
		if code == 200 {
			return resp[1:],nil
		}

		return resp,fmt.Errorf(string(resp[1].([]byte)))
	}

	return resp, fmt.Errorf("host not exist")
}

func (this *TCPClient) Send(method, payload string) (resp []byte, err error) {
	req := fmt.Sprintf("{\"method\":\"%s\",\"payload\":%s}", method, payload)

	data := []byte(req)
	data = append(data, 0)
	_, err = this.conn.Write(data)
	if err != nil {
		println("Write to server failed:", err.Error())
		return resp, err
	}

	reply := make([]byte, 1024)

	read_len, err := this.conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		return resp, err
	}

	reply = reply[:read_len]
	return reply, nil
}

func (this *TCPClient) Close() {
	this.conn.Close()
}
