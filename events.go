package core

import (
	"fmt"
	"strings"
)

func CreateHandleEvent(server_name string, event_name string, callback func([]interface{})) {
	queue_name := fmt.Sprintf("%v_%v", ServiceName(), event_name)

	server, err := GetRabbitServer(server_name)
	if err != nil {
		message := fmt.Sprintf("CreteHandleEvent: GetRabbitServer: %s", err.Error())
		ERROR(message)
		return
	}

	consumer, err := server.BindConsumer(queue_name, event_name, false)
	if err != nil {
		message := fmt.Sprintf("CreteHandleEvent: BindConsumer: %s", err.Error())
		ERROR(message)
		return
	}

	consumer.Listen(callback)
}

func CreateTemporaryListener(server_name string, event_name string, callback func([]interface{})) {
	queue_name := fmt.Sprintf("%v_%v", ServiceName(), event_name)

	server, err := GetRabbitServer(server_name)
	if err != nil {
		message := fmt.Sprintf("CreteTemporaryHandleEvent: GetRabbitServer: %s", err.Error())
		ERROR(message)
		return
	}

	consumer, err := server.BindConsumer(queue_name, event_name, true)
	if err != nil {
		message := fmt.Sprintf("CreteHandleEvent: BindConsumer: %s", err.Error())
		ERROR(message)
		return
	}

	consumer.Listen(callback)
}

func CreateEventPublisher(server_name string, name string) *Publisher {
	server, err := GetRabbitServer(server_name)
	if err != nil {
		message := fmt.Sprintf("CreteEvent: GetRabbitServer: %s", err.Error())
		ERROR(message)
		return nil
	}

	exchangeName := strings.ToUpper(ServiceName() + "_" + name)

	publ := server.GetPublisher(exchangeName)

	INFO("Create event publisher: name: " + server_name + ", exchange_name: " + exchangeName)

	return publ
}
