package models

import "github.com/go-redis/redis/v7"

var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
