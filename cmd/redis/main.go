package main

import (
	"github.com/codecrafters-io/redis-starter-go/pkg/server"
)

func main() {
	redis := server.NewRedis("0.0.0.0", 6379)
	redis.Run()
}
