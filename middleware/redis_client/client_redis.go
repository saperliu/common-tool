package redis_client

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	Address  string
	Password string
	Db       int
	Client   *redis.Client
	PoolSize int
	Context  context.Context
}

/**
 * 新建redis客户端
 */
func (rc *RedisClient) NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password, // no password set
		DB:       rc.Db,       // use default DB
		PoolSize: rc.PoolSize,
	})
	rc.Client = client
	rc.Context = context.Background()
	return rc
}

/**
 * 保存数据
 */
func (rc *RedisClient) SetValue(key string, value string, duration time.Duration) error {
	err := rc.Client.Set(rc.Context, key, value, duration).Err()
	return err
}

/**
 * 读取数据
 */
func (rc *RedisClient) GetValue(key string) (string, error) {
	return rc.Client.Get(rc.Context, key).Result()
}

/**
 * 删除数据
 */
func (rc *RedisClient) DelValue(key string) (int64, error) {
	return rc.Client.Del(rc.Context, key).Result()
}

/**
 * 向通道发布消息
 */
func (rc *RedisClient) Publish(channel string, message string) error {
	return rc.Client.Publish(rc.Context, channel, message).Err()
}

/**
 * 订阅消息
 */
func (rc *RedisClient) Subscribe(channel string) *redis.PubSub {
	return rc.Client.Subscribe(rc.Context, channel)
}

/**
 * 保存数据
 */
func (rc *RedisClient) SAddSet(key string, members ...string) error {
	err := rc.Client.SAdd(rc.Context, key, members).Err()
	return err
}

/**
 * 读取数据
 */
func (rc *RedisClient) SMembers(key string, ) ([]string, error) {
	return rc.Client.SMembers(rc.Context, key).Result()
}

/**
 * 删除数据
 */
func (rc *RedisClient) SDelSet(key string, members ...string) *redis.IntCmd {
	return rc.Client.SRem(rc.Context, key, members)
}

/**
 * 获取两个set里的差集
 */
func (rc *RedisClient) SDiff(key string, targetKey string) ([]string, error) {
	return rc.Client.SDiff(rc.Context, key, targetKey).Result()
}

/**
 * 获取set里的值,循环出来
 */
func (rc *RedisClient) SScanSet(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return rc.Client.SScan(rc.Context, key, cursor, match, count).Result()
}

/**
 * 获取set里的成员数量
 */
func (rc *RedisClient) SCard(key string) (int64, error) {
	return rc.Client.SCard(rc.Context, key).Result()
}
