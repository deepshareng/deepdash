package storage

import (
	"time"

	"strings"

	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/redis.v3"
)

type redisClusterKV struct {
	cli *redis.ClusterClient
}

func NewRedisClusterSimpleKV(url string, password string, poolSize int) SimpleKV {
	urls := queryClusterNodes(url, password)
	opt := &redis.ClusterOptions{
		Addrs:    urls,
		Password: password,
		PoolSize: poolSize,
	}
	cli := redis.NewClusterClient(opt)

	if err := cli.Ping().Err(); err != nil {
		log.Fatal("Connect to redis cluster failed, err:", err)
	} else {
		log.Debug("[init]Redis Cluster config:", opt)
		if err := cli.Set("foo", "bar", 0).Err(); err != nil {
			log.Fatal("redis custer set failed:", err)
		} else {
			log.Debug("[init]redis cluster set foo = bar")
		}
		if err := cli.Del("foo").Err(); err != nil {
			log.Fatal("redis cluster del failed:", err)
		} else {
			log.Debug("[init]redis cluster del foo")
		}
	}
	return &redisClusterKV{cli: cli}
}

func queryClusterNodes(url, password string) []string {
	cli := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
	})
	defer cli.Close()
	var urls []string
	lines := strings.Split(cli.ClusterNodes().String(), "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		for i, word := range parts {
			if i > 1 && strings.HasSuffix(word, "master") {
				urls = append(urls, parts[i-1])
				break
			}
		}
	}
	log.Info("Redis cluster urls:", urls)
	return urls
}

func (redisKV *redisClusterKV) Get(k []byte) ([]byte, error) {
	b, err := redisKV.cli.Get(string(k)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	log.Debugf("RedisClusterKV; Get: k = %s; v = %s", string(k), string(b))
	return b, err
}

func (redisKV *redisClusterKV) Delete(k []byte) error {
	return redisKV.cli.Del(string(k)).Err()
}

func (redisKV *redisClusterKV) Set(k []byte, v []byte) error {
	return redisKV.cli.Set(string(k), v, 0).Err()
}

func (redisKV *redisClusterKV) SetEx(k []byte, v []byte, expiration time.Duration) error {
	log.Debugf("RedisClusterKV; SetEx: k = %s, v = %s, expiration = %v\n", string(k), string(v), expiration)
	return redisKV.cli.Set(string(k), v, expiration).Err()
}

func (redisKV *redisClusterKV) HSet(k []byte, hk string, v []byte) error {
	return redisKV.cli.HSet(string(k), hk, string(v)).Err()
}

func (redisKV *redisClusterKV) HGet(k []byte, hk string) ([]byte, error) {
	b, err := redisKV.cli.HGet(string(k), hk).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return b, err
}

func (redisKV *redisClusterKV) HDel(k []byte, hk string) error {
	return redisKV.cli.HDel(string(k), hk).Err()
}

func (redisKV *redisClusterKV) HIncrBy(k []byte, hk string, n int) error {
	return redisKV.cli.HIncrBy(string(k), hk, int64(n)).Err()
}

func (redisKV *redisClusterKV) Exists(k []byte) bool {
	return redisKV.cli.Exists(string(k)).Val()
}
