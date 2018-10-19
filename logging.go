package core

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "time"
)

type logg struct {
    level  string
    object string
    method string
    params string
    body   string
}

var service_id = 0
var cache_level = make(map[string]int)
var cache_object = make(map[string]int)
var cache_method = make(map[string]int)
var is_debug = true
var log_db *sql.DB
var log_queue = make(chan *logg, 100)
var sql_logs = "INSERT INTO logs(service, level, object, method, params, body, session) VALUES($1, $2, $3, $4, $5, $6, $7)"

func initLogging(conn string, debug bool) {
    go run_log_engine(conn)
    is_debug = debug
}

func get_service_id(name string) (_id int) {
    row := log_db.QueryRow("SELECT service_upsert($1)", name)
    err := row.Scan(&_id)
    if err != nil {
        PublishError("logg.service_upsert: " + err.Error())
        return _id
    }
    PublishWarning("Get service id from DB, name: " + name)
    service_id = _id
    return _id
}

func get_level_id(name string) (_id int) {
    if id, ok := cache_level[name]; ok {
        return id
    }
    row := log_db.QueryRow("SELECT level_upsert($1)", name)
    err := row.Scan(&_id)
    if err != nil {
        PublishError("logg.level_upsert: " + err.Error())
        return _id
    }
    PublishWarning("Get level id from DB, name: " + name)
    cache_level[name] = _id
    return _id
}

func get_object_id(name string) (_id int) {
    if id, ok := cache_object[name]; ok {
        return id
    }
    row := log_db.QueryRow("SELECT object_upsert($1)", name)
    err := row.Scan(&_id)
    if err != nil {
        PublishError("logg.object_upsert: " + err.Error())
        return _id
    }
    PublishWarning("Get object id from DB, name: " + name)
    cache_object[name] = _id
    return _id
}

func get_method_id(name string) (_id int) {
    if id, ok := cache_method[name]; ok {
        return id
    }
    row := log_db.QueryRow("SELECT method_upsert($1)", name)
    err := row.Scan(&_id)

    if err != nil {
        PublishError("logg.method_upsert: " + err.Error())
        return _id
    }
    PublishWarning("Get method id from DB, name: " + name)
    cache_method[name] = _id
    return _id
}

func (this *logg) save() error {
    var srv_id = service_id
    if srv_id == 0 {
        service_id = get_service_id(ServiceName())
        srv_id = service_id
    }
    level_id := get_level_id(this.level)
    object_id := get_object_id(this.object)
    method_id := get_method_id(this.method)
    sid := GetSessionID()

    row, err := log_db.Query(sql_logs, srv_id, level_id, object_id, method_id, this.params, this.body, sid)
    defer row.Close()

    if err != nil {
        PublishError("core::logs save: " + err.Error())
        return err
    }
    return nil
}

func run_log_engine(log_db_conn string) {
    var err error
    for {
        time.Sleep(time.Second * 1)

        log_db, err = DBConnect(log_db_conn)
        if err != nil {
            ERROR(err.Error())
            continue
        }

        INFO("Log engine is initializer, conn: " + log_db_conn)
        INFO("Log engine start...")

        var log *logg
        for {
            log = <-log_queue
            err = log.save()
            if err != nil {
                ERROR(err.Error())
                break
            }
        }

    }
}

func log(level, object, method, params string, body string) {
    msg := fmt.Sprintf("%v__%v__%v__%v", object, method, params, body)
    switch level {
    case "INFO":
        INFO(msg)
    case "ERROR":
        ERROR(msg)
    case "WARNING":
        WARNING(msg)
    }

    l := &logg{level: level, object: object, method: method, params: params, body: body}
    select {
    case log_queue <- l:
    default:
    }

}

func Info(object, method, params string, body string) {
    if is_debug {
        log("INFO", object, method, params, body)
    }
}

func Error(object, method, params string, body string) {
    log("ERROR", object, method, params, body)
}

func TradeHistory(object, method, params string, body string) {
    log("TRADE_HISTORY", object, method, params, body)
}

func Warning(object, method, params string, body string) {
    log("WARNING", object, method, params, body)
}

func Fatal(object, method, params string, body string) {
    log("FATAL", object, method, params, body)
}
