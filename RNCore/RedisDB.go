package RNCore

import (
	"github.com/gomodule/redigo/redis"
)

type RedisDB struct {
	Node
	Conn redis.Conn
}

func NewRedisDB(name, url, password string, db int) RedisDB {
	rdb := RedisDB{Node: NewNode(name)}

	conn, err := redis.Dial("tcp", url, redis.DialPassword(password), redis.DialDatabase(db))
	if err != nil {
		rdb.Panic(err.Error())
	}
	rdb.Conn = conn

	return rdb
}

func (this *RedisDB) Close() {
	this.Node.Close()

	this.Conn.Close()
}
