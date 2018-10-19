package core

import (
    "database/sql"
    _ "github.com/lib/pq"
)

var Postgres *sql.DB

func DBConnect(conn_param string) (*sql.DB, error) {
    var err error
    dbi, err := sql.Open("postgres", conn_param)
    if err != nil {
        return dbi, err
    }
    return dbi, nil
}

func SqlStatement(query string) *sql.Stmt {
    stmt, err := Postgres.Prepare(query)
    if err != nil {
        ERROR("db.Prepare: " + err.Error())
        panic(err)
    }
    return stmt
}

func SqlQuery(sql string, args ...interface{}) error {
    _, err := Postgres.Exec(sql, args...)
    return err
}

// get all rows from the database
func SqlAll(sql string, args ...interface{}) (res [][]interface{}, err error) {
    rows, err := Postgres.Query(sql, args...)
    if err != nil {
        return res, err
    }
    defer rows.Close()

    col_ls, err := rows.Columns()
    if err != nil {
        return res, err
    }

    l := len(col_ls)

    for rows.Next() {
        var row = make([]interface{}, l)
        var prow = make([]interface{}, l)
        for idx := range row {
            prow[idx] = &row[idx]
        }
        err := rows.Scan(prow...)
        if err != nil {
            return res, err
        }
        res = append(res, row)
    }

    return res, err
}

// get one row from the database
func SqlOne(sql string, args ...interface{}) (res []interface{}, err error) {
    set, err := SqlAll(sql, args...)
    if err != nil {
        return res,err
    }

    if len(set) > 0 {
        res = set[0]
    }

    return res, err
}

// get one result from the database
func SqlScalar(sql string, args ...interface{}) (res interface{}, err error) {
    set, err := SqlOne(sql, args...)
    if err != nil {
        return res,err
    }

    if len(set) > 0 {
        res = set[0]
    }

    return res, err
}

