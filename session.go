package core

import (
    "fmt"
    "strings"
    "time"
)

var sid = 0
var previous_sid int

var sql_table_session = `
CREATE SEQUENCE IF NOT EXISTS session_id_seq
INCREMENT 1
MINVALUE 1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
    

CREATE TABLE IF NOT EXISTS session
(
    id INTEGER DEFAULT nextval('session_id_seq'::regclass),
    datetime TIMESTAMP DEFAULT now(),
    description TEXT,
    params TEXT,
    CONSTRAINT session_id_pk PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS session_datetime_index
ON session
USING btree
(datetime);



`
type SessionObject struct {
    Id int
    Datetime time.Time
    Params string
    Description string
}


func initSession() error {
    _, err := Postgres.Exec(string(sql_table_session))
    if err != nil {
        ERROR("initSession.PostgresQueryRow: " + err.Error())
    }
    GetSessionID()
    return err
}

func SessionListFilter(from *int, limit int, filter_list *[]string) (res []SessionObject, count int, err error) {
    where_block_list := []string{}

    if filter_list != nil {
        for _, filter := range *filter_list {
            where_block_list = append(where_block_list, filter)
        }
    }

    if from != nil {
        where_block_list = append(where_block_list, fmt.Sprintf("id > %v", *from))
        //where_block_list = append(where_block_list, fmt.Sprintf("id <= %v", from + limit))
    }

    limit_block := ""
    if limit != 0 {
        limit_block = fmt.Sprintf(" LIMIT %v", limit)
    }

    where_block := ""
    if len(where_block_list) != 0 {
        where_block = strings.Join(where_block_list, " AND ")
        where_block = "WHERE " + where_block
    }
    sql := "SELECT id, datetime, params, description FROM session " + where_block + limit_block

    rows, err := Postgres.Query(sql)
    if err != nil {
        PublishError("session.ListFilter: " + err.Error())
        return res, count, err
    }

    defer rows.Close()

    for rows.Next() {
        p := SessionObject{}
        err := rows.Scan(&p.Id, &p.Datetime, &p.Params, &p.Description)
        res = append(res, p)
        if err != nil {
            PublishError("session.List: " + err.Error())
            return res, count, err
        }
        count++
        if *from < p.Id {
            *from = p.Id
        }
    }
    return res, count, err
}

// get the current session id
func GetSessionID() int {
    var err error
    if sid != 0 {
        return sid
    }
    // insert new session and return sid

    // todo get the list of starting parameters of the service
    var params = ""

    // ask for the identifier of the previous session
    // if this is the first session in the life of the service - return zero
    sql := "SELECT id FROM session WHERE datetime = (SELECT max(datetime) FROM session)"
    row := Postgres.QueryRow(sql)
    err = row.Scan(&previous_sid)
    if err != nil {
        ERROR("GetSessionID:PostgresQueryRow:get_previous_id " + err.Error())
        previous_sid = 0
    }

    // insert new session
    sql2 := "INSERT INTO session(description, params) VALUES ($1,$2) RETURNING id;"
    row2 := Postgres.QueryRow(sql2, "", params)
    err = row2.Scan(&sid)
    if err != nil {
        ERROR("GetSessionID:PostgresQueryRow:get_current_id " + err.Error())
        sid = 0
    }

    return sid
}

// get the previous session id
func GetPreviousSessionID() int {
    return previous_sid
}
