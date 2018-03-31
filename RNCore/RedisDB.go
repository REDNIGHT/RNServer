package RNCore

import (
	"github.com/gomodule/redigo/redis"
)

type RedisDB struct {
	MNode
	Conn redis.Conn
}

func NewRedisDB(name, url, password string, db int) RedisDB {
	rdb := RedisDB{NewMNode(name), nil}

	conn, err := redis.Dial("tcp", url, redis.DialPassword(password), redis.DialDatabase(db))
	if err != nil {
		rdb.Panic(err, "err != nil")
	}
	rdb.Conn = conn

	return rdb
}

func (this *RedisDB) Close() {
	this.MNode.Close()

	this.Conn.Close()
}
