package gcache

import (
	"GoFocusMicroService/conf"
	"encoding/json"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

var Rds *Cache

//
//  @Description: 初始化redis功能配置
//
func SetUp() {
	Rds = NewCache(conf.Conf.Redis.Url, conf.Conf.Redis.Prefix)
}

type Cache struct {
	Redis  *redis.Client
	Prefix string
}

//
//  @Description: 新建对象
//  @param uri:
//  @param prefix:
//  @return *Cache:
//
func NewCache(uri string, prefix string) *Cache {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	obj := Cache{
		Redis:  redis.NewClient(opt),
		Prefix: prefix,
	}
	return &obj
}

// fmtKey 格式化缓存key
func (c Cache) fmtKey(key string) string {
	if strings.HasPrefix(key, "-") {
		return key
	}
	if c.Prefix == "" {
		return key
	}
	return c.Prefix + key
}

// Get 读取Redis数据
func (c Cache) Get(key string) (value string, err error) {
	value, err = c.Redis.Get(c.fmtKey(key)).Result()
	// 其它故障类异常
	if err != nil && err != redis.Nil {
		panic(err)
	}
	return value, err
}

// GetObj 读取Redis数据并自动转换为对象
func (c Cache) GetObj(key string, obj interface{}) error {
	value, err := c.Get(key)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(value), &obj)
	if err != nil {
		return err
	}

	return nil
}

// Set 缓存数据到Redis中
func (c Cache) Set(key string, value interface{}, expiration time.Duration) (err error) {
	err = c.Redis.Set(c.fmtKey(key), value, expiration).Err()
	return
}

// SetObj 缓存数据时，自动将对象转换为JSON
func (c Cache) SetObj(key string, obj interface{}, expiration time.Duration) (err error) {
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Set(key, value, expiration)
}

// Expire 设置key的过期时间
func (c Cache) Expire(key string, expiration time.Duration) (err error) {
	err = c.Redis.Expire(c.fmtKey(key), expiration).Err()
	return
}

// ExpireAt 设置key在指定日期过期
func (c Cache) ExpireAt(key string, tm time.Time) (err error) {
	err = c.Redis.ExpireAt(c.fmtKey(key), tm).Err()
	return
}

// Del 删除所有Key
func (c Cache) Del(keys ...string) (err error) {
	fmtKeys := make([]string, len(keys))
	for idx, v := range keys {
		fmtKeys[idx] = c.fmtKey(v)
	}
	err = c.Redis.Del(fmtKeys...).Err()
	return
}

// Incr 增长数据
func (c Cache) Incr(key string) (err error) {
	err = c.Redis.Incr(c.fmtKey(key)).Err()
	return
}
