package db

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" || port == "" {
		return nil, fmt.Errorf("As variaveis REDIS_HOST e REDIS_PORT estão vazias.")
	}

	// Cria cliente Redis
	client := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	// Verifica a conexão com o Redis
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
