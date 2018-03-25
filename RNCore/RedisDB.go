package RNCore

import (
	"github.com/go-redis/redis"
)

type RedisDB struct {
	MinNode

	Client *redis.Client
}

func NewRedisDB(name, url, password, db, c string, indexKeys ...string) RedisDB {
	rdb := RedisDB{NewMinNode(name), nil}

	rdb.Client = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	_, err := rdb.Client.Ping().Result()
	if err != nil {
		rdb.Error(err.Error())
	}
	/*rdb.Client.Pipelined(func(pipe redis.Pipeliner) error {
		_ := pipe.Ping()

		return nil
	})*/
	/*_, err := rdb.Client.Get("key").Result()
	if err != nil {
		panic(err)
	}*/
	return rdb
}

func (this *RedisDB) Close() {
	this.Client.Close()
}
