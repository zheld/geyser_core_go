package core

import (
    "fmt"
    "time"
    "strings"
)

func handle(item int, prefix string) (res string) {
    res = fmt.Sprint(item)
    if len(res) < 2 {
        res = prefix + res
    }
    return
}

var publog_publisher *Publisher

func publog_message(msg string, _type string) string {
    now := time.Now()
    hour, min, sec := now.Clock()

    lname := 8
    srv_name := strings.ToUpper(ServiceName())
    rs := lname - len(srv_name)
    if rs > 0 {
        srv_name = srv_name + strings.Repeat(" ", rs)
    }

    return fmt.Sprintf(
        "%v:%v:%v  %v  %v: %v",
        handle(hour, " "),
        handle(min, "0"),
        handle(sec, "0"),
        srv_name,
        strings.ToUpper(_type),
        msg,
    )
}

func send_with_check(msg string) {
    if publog_publisher != nil {
        publog_publisher.Send(msg)
    }
    fmt.Println(msg)
}

func PublishInfo(msg string) {
    send_with_check(publog_message(msg, "    inf "))
}

func PublishWarning(msg string) {
    send_with_check(publog_message(msg, "warning "))
}

func PublishError(msg string) {
    send_with_check(publog_message(msg, "  error "))
}

func PublishCrush(msg string) {
    send_with_check(publog_message(msg, "  crush "))
}

func INFO(msg string) {
    send_with_check(publog_message(msg, "    inf "))
}

func WARNING(msg string) {
    send_with_check(publog_message(msg, "warning "))
}

func ERROR(msg string) {
    send_with_check(publog_message(msg, "  error "))
}

func DEBUG(msg string) {
    send_with_check(publog_message(msg, "  debug "))
}

