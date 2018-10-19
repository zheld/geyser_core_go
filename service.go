package core

import (
	"os/exec"
	"fmt"
	"strings"
	"errors"
	"time"
)

var service_info ServiceInfo
var core_api = map[string]HandlerFoo{}

type ServiceInfo struct {
	Name     string
	Version  string
	TCP_Port int
}

func ServiceName() string {
	return service_info.Name
}

func ServiceVersion() string {
	return service_info.Version
}

func GetServiceInfo() *ServiceInfo {
	return &service_info
}

func GetAddress() (string, error) {
	out, err := exec.Command("ifconfig").Output()
	if err != nil {
		message := fmt.Sprintf("core: GetAddress: %v", err.Error())
		ERROR(message)
		return "", errors.New(message)
	}
	doc := string(out)
	item_list := strings.Split(doc, " ")
	for _, item := range item_list {
		if strings.Contains(item, "192.168.1.") && !strings.Contains(item, "255") {
			def := strings.Split(item, ":")
			return def[1], nil
		}
	}
	message := fmt.Sprintf("core: GetAddress: no address in ifconfig command")
	ERROR(message)
	return "", errors.New(message)
}

func serviceInit(name, version string, port int) (si ServiceInfo, err error) {
	si.Name = name
	si.Version = version
	si.TCP_Port = port
	return si, err
}

func ServiceStart(api map[string]HandlerFoo) {
	// api methods
	setAPI(api)
	// api server
	if service_info.TCP_Port != 0 {
		api_server_address := ":" + IToStr(service_info.TCP_Port)
		start_api_server(api_server_address)
	}

	// always pause
	for {
		time.Sleep(time.Second * 1000000)
	}
}
