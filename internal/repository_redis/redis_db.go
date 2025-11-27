package repository_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type RedisDB struct {
	*redis.Client
}

func (db *RedisDB) ConnectToRedis() error {
	db.Client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
		DB:   0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := db.Client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	log.Println("Service connected to Redis: " + pong)
	return nil
}

func (db *RedisDB) SetData(key string, value []byte, timeLive time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.Client.Set(ctx, key, value, timeLive).Err()
	if err != nil {
		return err
	}
	return nil
}

func (db *RedisDB) GetData(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Client.Get(ctx, key).Bytes()
}

func (db *RedisDB) DeleteData(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
