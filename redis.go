package core

import (
    "github.com/go-redis/redis"
    "fmt"
)

var redis_main *RedisClient

type RedisClient struct {
    host   string
    client *redis.Client
}

func NewReddisClient(address string) (*RedisClient, error) {
    client := &RedisClient{
        host:address,
    }
    red_client, err := client.redis_connect()
    if err != nil {
        return client, err
    }

    client.client = red_client

    return client, nil
}

func (this *RedisClient) redis_connect() (*redis.Client, error) {
    red_client := redis.NewClient(&redis.Options{
        Addr:     this.host + ":6379",
        Password: "", // no password set
        DB:       0, // use _default DB
    })

    _, err := red_client.Ping().Result()

    if err != nil {
        return nil, err
    }

    return red_client, nil
}

func (this *RedisClient) Set(collection string, key string, value interface{}) error {
    full_key := collection + "." + key
    err := this.client.Set(full_key, value, 0).Err()
    if err != nil {
        cl, err := this.redis_connect()
        if err != nil {
            return err
        }

        this.client = cl
        cl.Set(full_key, value, 0)
        return nil
    }

    return nil
}

func (this *RedisClient) GetString(collection string, key string) (result string, ok bool) {
    full_key := collection + "." + key
    res := this.client.Get(full_key)
    if err := res.Err(); err != nil {
        if err.Error() == "redis: nil" {
            return "", false
        } else {
            cl, err := this.redis_connect()
            if err != nil {
                return "", false
            }

            this.client = cl
            res := cl.Get(full_key)
            if err := res.Err(); err != nil {
                if err.Error() == "redis: nil" {
                    return "", false
                } else {
                    ERROR(fmt.Sprintf("core: redis: GetString: error: %v", err.Error()))
                    return "", false
                }
            }

            result, _ = res.Result()
            return result, true
        }
    }
    result, _ = res.Result()
    return result, true
}

func (this *RedisClient) GetInt(collection string, key string) (result int64, ok bool) {
    full_key := collection + "." + key
    res := this.client.Get(full_key)
    if err := res.Err(); err != nil {
        if err.Error() == "redis: nil" {
            return 0, false
        } else {
            cl, err := this.redis_connect()
            if err != nil {
                return 0, false
            }

            this.client = cl
            res := cl.Get(full_key)
            if err := res.Err(); err != nil {
                if err.Error() == "redis: nil" {
                    return 0, false
                } else {
                    ERROR(fmt.Sprintf("core: redis: GetString: error: %v", err.Error()))
                    return 0, false
                }
            }

            result, _ = res.Int64()
            return result, true
        }
    }
    result, _ = res.Int64()
    return result, true
}



