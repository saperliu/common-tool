package redis_client

import (
	"fmt"
	"testing"
)

func Test_Redis(t *testing.T) {
	redisClient := RedisClient{}
	redisClient.Address = "127.0.0.1:6379"
	redisClient.PoolSize = 5
	redisClient.Db = 1
	redisClient.Password = ""
	redisClient.NewRedisClient()

    str,err:=redisClient.GetValue("TEST")
	fmt.Printf("  save  server data to etcd  error %v", err)
	fmt.Printf("  save  server  %v", str)
}
