package core

import (
    "time"
    "errors"
    "net"
    "bufio"
)

var duration_sleep = time.Second
var apiClientStorage = map[string]*APIClient{}
var count = 1

func NewAPIClient(service_name string) *APIClient {
    service_address := "services." + service_name + ".address"
    cl := &APIClient{address: service_address, service_name: service_name}
    cl.connect()
    return cl
}

type APIClient struct {
    service_name string
    address      string
    conn         net.Conn
    api          []string
}

func (this *APIClient) getMethodId(name string) int {
    for idx, n := range this.api {
        if n == name {
            return idx
        }
    }
    return -1
}

func (this *APIClient) connect() error {
    conn, err := net.Dial("tcp", this.address)
    if err != nil {
        ERROR("Failed to DBConnect to test server")
    }
    reader := bufio.NewReader(conn)
    msg, err := reader.ReadSlice(delim)
    if err != nil {
        conn.Close()
        return err
    }

    var re []interface{}
    MsgUnpack(msg, &re)
    api_list := []string{}
    for _, i := range re {
        api_list = append(api_list, string(i.([]uint8)))
    }
    this.conn = conn
    this.api = api_list
    return nil
}

func (this *APIClient) Call(method string, args ...interface{}) (result interface{}, err error) {
    method_id := this.getMethodId(method)
    param := []interface{}{method_id}
    param = append(param, args...)
    b := MsgPack(args)
    b = append(b, delim)

    this.conn.Write(b)

    reader := bufio.NewReader(this.conn)
    msg, err := reader.ReadSlice(delim)
    if err != nil {
        if ok := this.reconnect(); !ok {
            return nil, errors.New("core: APIClient: Call: connection is closed: " + err.Error())
        }
        this.conn.Write(b)
        reader := bufio.NewReader(this.conn)
        msg, err = reader.ReadSlice(delim)
        if err != nil {
            return nil, errors.New("core: APIClient: Call: " + err.Error())
        }
    }

    MsgUnpackScalar(msg, &result)
    return result, nil
}

func (this *APIClient) reconnect() bool {
    iter := count
    for iter > 0 {
        WARNING("reconnect into service: " + this.service_name)
        this.conn.Close()
        err := this.connect()
        if err == nil {
            return true
        }
        ERROR("core: api: APIClient: reconnect: err:" + err.Error())
        time.Sleep(duration_sleep)
        iter--
    }
    return false
}

func GetAPIClient(service_name string) *APIClient {
    if client, ok := apiClientStorage[service_name]; ok {
        return client
    }
    client := NewAPIClient(service_name)
    apiClientStorage[service_name] = client
    return client
}

func Invoke(srv_name string, method string, args ...interface{}) (result interface{}, err error) {
    client := GetAPIClient(srv_name)
    return client.Call(method, args...)
}
