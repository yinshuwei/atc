package base

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	configPath = ""

	// Config 配置内容
	Config = &AtcConfig{}

	// RefCache 引用模块内存
	RefCache = map[string]*string{}

	// RefCacheMutex 引用模块内存锁
	RefCacheMutex sync.Mutex

	redisPool        *redis.Pool
	memoryCache      = map[string]*AtcPageCache{}
	memoryCacheMutex sync.Mutex
)

// AtcConfig AtcConfig
type AtcConfig struct {
	WebPath      string
	Port         string
	RedisAddress string
	RedisAuth    string
	Envs         map[string]string
	IsDev        bool
	Page404      string
	Caches       map[string]int64
}

// AtcStatic Static
type AtcStatic struct {
	// Dir is the directory to serve static files from
	Dir http.FileSystem
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
}

// AtcPageCache 内存缓存
type AtcPageCache struct {
	Body        []byte
	ExpiredTime time.Time
}

// ReadConfig 读取配置
func ReadConfig() error {
	RefCacheMutex.Lock()
	defer RefCacheMutex.Unlock()
	RefCache = map[string]*string{}

	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	err = json.Unmarshal(b, Config)
	if err != nil {
		return err
	}

	if Config.RedisAddress != "" {
		redisPool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", Config.RedisAddress)
				if err != nil {
					return nil, err
				}
				if Config.RedisAuth != "" {
					if _, err := c.Do("AUTH", Config.RedisAuth); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
	}

	return nil
}

// ConfigPath 配置文件路径
func ConfigPath() {
	c := flag.String("c", "./config/atc.json", "config path")
	flag.Parse()
	configPath = *(c)
}

// ReadPageCache 读取缓存
func ReadPageCache(name, query string) ([]byte, bool, bool) {
	if Config.Caches == nil {
		return nil, false, false
	}
	ttl, ok := Config.Caches[name]
	if !ok || ttl <= 0 {
		return nil, false, false
	}

	redisKey := fmt.Sprintf("actcache:%s_%s", name, query)
	if redisPool != nil {
		redisKeyTo := fmt.Sprintf("actcache:to:%s_%s", name, query)

		conn := redisPool.Get()
		defer conn.Close()
		conn.Send("SELECT", "8")
		conn.Send("EXISTS", redisKeyTo)
		conn.Send("GET", redisKey)
		conn.Flush()
		conn.Receive()
		exist, err := redis.Bool(conn.Receive())
		if err != nil {
			log.Println(err)
		}
		cache, err := redis.Bytes(conn.Receive())
		if err != nil {
			log.Println(err)
		}
		return cache, !exist, true
	}
	memoryCacheMutex.Lock()
	defer memoryCacheMutex.Unlock()
	mc, ok := memoryCache[redisKey]
	if mc != nil && ok {
		return mc.Body, mc.ExpiredTime.Before(time.Now()), true
	}
	return nil, true, true
}

// WritePageCache 写入缓存
func WritePageCache(name, query string, cache []byte) {
	if Config.Caches == nil {
		return
	}
	ttl, ok := Config.Caches[name]
	if !ok || ttl <= 0 {
		return
	}

	redisKey := fmt.Sprintf("actcache:%s_%s", name, query)
	if redisPool != nil {
		redisKeyTo := fmt.Sprintf("actcache:to:%s_%s", name, query)

		conn := redisPool.Get()
		defer conn.Close()
		conn.Send("SELECT", "8")
		conn.Send("SET", redisKeyTo, "true")
		conn.Send("EXPIRE", redisKeyTo, ttl)
		conn.Send("SET", redisKey, cache)
		conn.Send("EXPIRE", redisKey, 60*60*24*3)
		conn.Flush()
		conn.Receive()
		conn.Receive()
		conn.Receive()
		conn.Receive()
		conn.Receive()
	} else {
		memoryCacheMutex.Lock()
		defer memoryCacheMutex.Unlock()

		// 内存缓存需要新区域，不可使用bufferPool中的内存
		cacheNew := make([]byte, len(cache))
		copy(cacheNew, cache)

		memoryCache[redisKey] = &AtcPageCache{
			Body:        cacheNew,
			ExpiredTime: time.Now().Add(time.Duration(ttl) * time.Second),
		}
	}
}
