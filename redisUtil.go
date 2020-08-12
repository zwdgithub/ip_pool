package main

import (
	"errors"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

const (
	ValidIpSortedSetKey = "ip_pool_ValidIpSortedSetKey"
	WaitingKey          = "ip_pool_WaitingKey"
	EndTimeStamp        = 2543902107 // 2050-08-12 15:28:27
)

type RedisUtil struct {
	conn *redis.Client
}

func NewRedis() *RedisUtil {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisUtil{conn: conn}
}

func (r *RedisUtil) GetProxy() (string, error) {
	ip, err := r.conn.ZRange(ValidIpSortedSetKey, 0, 100).Result()
	if err != nil {
		return "", err
	}
	if len(ip) == 0 {
		return "", errors.New("ip not exists")
	}
	index := rand.Intn(len(ip))
	return ip[index], nil
}

func (r *RedisUtil) AddProxy(ip string) {
	r.conn.ZAdd(ValidIpSortedSetKey, redis.Z{
		Score:  float64(EndTimeStamp - time.Now().Unix()),
		Member: ip,
	})
}

func (r *RedisUtil) DeleteProxy(ip string) {
	r.conn.ZRem(ValidIpSortedSetKey, ip)
}

func (r *RedisUtil) AddIpToWaitingList(ip string) {
	r.conn.LPush(WaitingKey, ip)
}

func (r *RedisUtil) GetWaitingIp() (string, error) {
	return r.conn.RPop(WaitingKey).Result()
}

func (r *RedisUtil) AllValidIp() ([]string, error) {
	return r.conn.ZRevRange(ValidIpSortedSetKey, 0, -1).Result()
}
