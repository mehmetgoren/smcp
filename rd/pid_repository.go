package rd

import (
	"golang.org/x/net/context"
	"os"
)

type PidRepository struct {
	*RedisOptions
}

func getKey() string {
	return "pid:smcp_service"
}

func (r *PidRepository) Add() (int64, error) {
	pid := os.Getpid()

	return r.Client.SAdd(context.Background(), getKey(), pid).Result()
}
