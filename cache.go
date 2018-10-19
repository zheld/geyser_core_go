package core

var TAG_COMMON = []string{"common"}

func GetIdStorage(collection, key string) (int, bool) {
    if id, ok := redis_main.GetInt(collection, key); ok {
        return int(id), true
    }

    return 0, false
}

func SetIdStorage(collection, key string, id int) {
    redis_main.Set(collection, key, id)
}

func GetID(service string, object string, name string) (int64, bool) {
    collection := "id." +  service + "." + object
    if id, ok := redis_main.GetInt(collection, name); ok {
        return id, true
    }

    return 0, false
}

func SetID(service string, object string, name string, id int64) error {
    collection := "id." + service + "." + object
    return redis_main.Set(collection, name, id)
}
