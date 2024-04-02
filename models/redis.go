package models

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error){
    redisClient := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // No password set
        DB:       0,  // Use default DB
    })


   if err := redisClient.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    return redisClient, nil
}
