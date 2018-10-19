package core

import (
    "fmt"
    "strings"
    "time"
)

const (
    LOG_OBJ = "Bios"
    bi_get  = "Get"
)

// cache
var cache_bios = map[string]string{}

// sql
var sql_create_table = `
CREATE SEQUENCE IF NOT EXISTS bios_id_seq
INCREMENT 1
MINVALUE 1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

CREATE TABLE IF NOT EXISTS bios
(
    id INTEGER DEFAULT nextval('bios_id_seq'::regclass),
    datetime TIMESTAMP DEFAULT now(),
    name TEXT,
    value TEXT,
    CONSTRAINT bios_id_pk PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS bios_name_index
ON bios
USING btree
(name);

CREATE OR REPLACE FUNCTION public.bios_upsert( p_name TEXT, p_value TEXT )
 RETURNS integer
 LANGUAGE plpgsql
AS $function$

DECLARE
  result INTEGER;

BEGIN
  
  SELECT
    id
  FROM
    bios
  WHERE
    name = p_name INTO result;
  
  IF result IS NULL THEN
    INSERT INTO bios (name, value) VALUES (p_name, p_value) RETURNING id INTO result;
  ELSE 
    UPDATE bios SET name = p_name, value = p_value WHERE id = result;
  END IF;
  
  RETURN result;

END;
$function$;
`

type SettingObject struct {
    Id       int
    Datetime time.Time
    Name     string
    Value    string
}

func initBios() error {
    _, err := Postgres.Exec(string(sql_create_table))
    return err
}

func SettingList() (res []SettingObject, err error) {
    sql := "SELECT id, datetime, name, value FROM bios"

    rows, err := Postgres.Query(sql)
    if err != nil {
        PublishError("bios.List: " + err.Error())
        return res, err
    }

    defer rows.Close()

    for rows.Next() {
        p := SettingObject{}
        err := rows.Scan(&p.Id, &p.Datetime, &p.Name, &p.Value)
        res = append(res, p)
        if err != nil {
            PublishError("bios.List: " + err.Error())
            return res, err
        }
    }
    return res, err
}

func SetSetting(name, value string) (int, error) {
    id, err := settingUpsert(name, value)
    if err != nil {
        return id, err
    }
    cache_bios[name] = value
    return id, nil
}

func settingUpsert(name string, value string) (_id int, err error) {
    row := Postgres.QueryRow("SELECT bios_upsert($1, $2)", name, value)
    err = row.Scan(&_id)
    if err != nil {
        PublishError("settingUpsert: " + err.Error())
        return _id, err
    }
    return _id, nil
}

func ReadSetting(id int) (t SettingObject, err error) {
    t = SettingObject{}
    row := Postgres.QueryRow("SELECT * FROM bios WHERE id = $1", id)
    err = row.Scan(&t.Id, &t.Datetime, &t.Name, &t.Value)
    if err != nil {
        PublishError("core:setting:Read: " + err.Error())
        return
    }
    return
}

func getSetting(name string) (value string) {
    if value, ok := cache_bios[name]; ok {
        return value
    }
    sql := "SELECT value FROM bios WHERE name = $1"
    row := Postgres.QueryRow(sql, name)
    err := row.Scan(&value)
    if err != nil {
        ERROR("GetSessionID.PostgresQueryRow: " + err.Error())
        return value
    }
    cache_bios[name] = value
    return value
}

func DeleteSetting(name string) error {
    if _, ok := cache_bios[name]; ok {
        delete(cache_bios, name)
    }
    sql := "DELETE FROM bios WHERE name = $1"
    _, err := Postgres.Exec(sql, name)
    if err != nil {
        return err
    }
    return nil
}

func GetSettingString(name string) (string, bool) {
    val := getSetting(name)
    if val == "" {
        return "", false
    }
    return val, true
}

func GetSettingInt(name string) (int, bool) {
    val_str := getSetting(name)
    if val_str == "" {
        return 0, false
    }
    val, err := StrToInt(val_str)
    if err != nil {
        return 0, false
    }
    return val, true
}

func GetSettingFloat(name string) (float64, bool) {
    val_str := getSetting(name)
    if val_str == "" {
        return 0, false
    }
    val, err := StrToFloat64(val_str)
    if err != nil {
        return 0, false
    }
    return val, true
}

func to_bool(bstr string) (bool, error) {
    bstr = strings.ToLower(bstr)
    if strings.Contains("true", bstr) || bstr == "t" || bstr == "y" {
        return true, nil
    }

    if strings.Contains("false", bstr) || bstr == "f" || bstr == "n" {
        return false, nil
    }

    return false, fmt.Errorf("unknow param: %v", bstr)
}

func GetSettingBool(name string) (bool, bool) {
    val_str := getSetting(name)
    if val_str == "" {
        return false, false
    }
    val, err := to_bool(val_str)
    if err != nil {
        return false, false
    }
    return val, true
}
