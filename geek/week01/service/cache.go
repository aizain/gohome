package service

import (
	"log"
	"time"
)

// Data 数据
type Data struct {
	Name string
	Age  int8
}

// CacheData 缓存数据，带有超时时间
type CacheData struct {
	*Data
	expire int64
}

// Cache 缓存
type Cache map[string]*CacheData

// DB 数据库
type DB map[string]*Data

// SyncCache 同步缓存，key过期时刷入数据库
type SyncCache struct {
	Cache
	DB
}

// Get 获取数据，如果过期了刷入数据库
func (s *SyncCache) Get(key string) (*Data, bool) {
	log.Printf("获取缓存 %v\n", key)
	val, ok := s.Cache[key]
	if ok {
		if val.expire > time.Now().UnixMilli() {
			s.Sync(key, val.Data)
		}
		return val.Data, ok
	}
	return nil, ok
}

// Set 缓存数据，并增加过期时间
func (s *SyncCache) Set(key string, data *Data) {
	log.Printf("设置缓存 %v\n", key)
	s.Cache[key] = &CacheData{
		data,
		time.Now().UnixMilli() + time.Minute.Milliseconds()*10,
	}
}

// Sync 同步数据
func (s *SyncCache) Sync(key string, data *Data) {
	log.Printf("同步缓存 %v\n", key)
	s.DB[key] = data
}

// SyncAll 同步所有书
func (s *SyncCache) SyncAll() int {
	log.Printf("正在同步缓存，缓存数量 %v\n", len(s.Cache))
	for k, v := range s.Cache {
		s.Sync(k, v.Data)
	}
	log.Printf("同步缓存结束，DB数量 %v", len(s.DB))
	return len(s.Cache)
}
