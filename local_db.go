package core

import (
    "github.com/peterbourgon/diskv"
    "time"
    "sync"
)

func sync_cache() {
    go func() {
        for {
            select {
            case <-time.After(time.Second * 5):
                saveCache()
            }
        }
    }()
}

type Collection map[string]string

var cache_mutex = sync.Mutex{}
var base = diskv.New(diskv.Options{})
var cache = map[string]Collection{}
var splitter = " : "

func saveCache() {
    for collname, coll := range cache {
        key_value_ls := []string{}
        for key, value := range coll {
            ln := key + splitter + value
            key_value_ls = append(key_value_ls, ln)
        }
        write(collname, StrJoin(key_value_ls, "\n"))
    }
}

func write(collname, value string) error {
    return base.Write(collname, []byte(value))
}

func read(collname string) (string, error) {
    v, err := base.Read(collname)
    return string(v), err
}

func serializeCollection(body string) Collection {
    coll := Collection{}
    lines := StrSplit(body, "\n")
    for _, line := range lines {
        kv := StrSplit(line, splitter)
        coll[kv[0]] = kv[1]
    }
    return coll
}

func Insert(collname, key, value string) {
    coll, ok := cache[collname]
    if !ok {
        text, err := read(collname)
        if err != nil {
            coll = Collection{}
        } else {
            coll = serializeCollection(text)
        }
        cache_mutex.Lock()
        cache[collname] = coll
        cache_mutex.Unlock()
    }
    coll[key] = value
}

func InsertList(collname, key string, list []string) {
    value := StrJoin(list, ",")
    Insert(collname, key, value)
}

func Select(collname, key string) (string, bool) {
    coll, ok := cache[collname];
    if !ok {
        text, err := read(collname)
        if err != nil {
            return "", false
        }
        coll_map := serializeCollection(text)
        cache_mutex.Lock()
        cache[collname] = coll_map
        cache_mutex.Unlock()
        coll = coll_map
    }

    if val, ok := coll[key]; ok {
        return val, true
    }
    return "", false

}

func SelectList(collname, key string) ([]string, bool) {
    if value, ok := Select(collname, key); ok {
        return StrSplit(value, ","), true
    }
    return []string{}, false
}
