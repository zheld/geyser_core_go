package core

import (
	"fmt"
	"strings"
)

type Config struct {
	Name                   string
	Version                string
	APIPort                int
	RedisConn              string
	RabbitConn             [][]string
	ExtDBConn              map[string]string
	RabbitPublogServerName string
	PostgreSQLConn         string
	BiosEnabled            bool
	LoggConn               string
	LoggDebug              bool
	SessionEnabled         bool
}

func Init(conf Config) (err error) {
	// SERVICE INIT
	service_info, err = serviceInit(conf.Name, conf.Version, conf.APIPort)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	INFO("======= START " + strings.ToUpper(ServiceName()) + " SERVICE ======================================")

	// REDIS
	if conf.RedisConn != "" {
		redis_main, err = NewReddisClient(conf.RedisConn)
		if err != nil {
			errmsg := fmt.Sprintf("core: init: getRedisClient: main: %v", err.Error())
			fmt.Println(errmsg)
			return fmt.Errorf(errmsg)
		}
		INFO("REDIS CLIENT ........... [OK]")
	} else {
		INFO("REDIS CLIENT ........... [NOT]")
	}

	// AMQP MAIN
	if conf.RabbitConn != nil && len(conf.RabbitConn) != 0 {
		for _, dconn := range conf.RabbitConn {
			rabbit_server_name := dconn[0]
			INFO("Init Rabbit server [" + rabbit_server_name + "]")
			_, err := SetRabbitServer(rabbit_server_name, dconn[1])
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

		}
		INFO("RABBIT CLIENT .......... [OK]")

		// PUBLOG
		if conf.RabbitPublogServerName != "" {
			publog_queue_name := "publog"

			if srv, err := GetRabbitServer(conf.RabbitPublogServerName); err == nil {
				publog_publisher = srv.GetPublisher(publog_queue_name)
			} else {
				fmt.Println("publog init error: ", err.Error())
				return err
			}
			INFO("PUBLOG ................. [OK]")
		} else {
			INFO("PUBLOG ................. [NOT]")
		}
	} else {
		INFO("RABBIT CLIENT .......... [NOT]")
	}

	// DATA BASE
	if conf.PostgreSQLConn != "" {
		Postgres, err = DBConnect(conf.PostgreSQLConn)
		if err != nil {
			ERROR(fmt.Sprintf("data base client is fail: %v", err.Error()))
			return err
		}
		INFO("POSTGRES CLIENT ........ [OK]")
	} else {
		INFO("POSTGRES CLIENT ........ [NOT]")
	}

	// LOG ENGINE
	if conf.LoggConn != "" {
		initLogging(conf.LoggConn, conf.LoggDebug)
		initSession()
		INFO("LOG ENGINE ............. [OK]")
	} else {
		INFO("LOG ENGINE ............. [NOT]")
	}

	// BIOS
	if conf.BiosEnabled {
		err := initBios()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		INFO("BIOS ENGINE ............ [OK]")

	} else {
		INFO("BIOS ENGINE ............ [NOT]")
	}

	// SESSION
	if conf.SessionEnabled {
		initSession()
		INFO("SESSION ENGINE ......... [OK]")
	} else {
		INFO("SESSION ENGINE ......... [NOT]")
	}

	return nil
}
