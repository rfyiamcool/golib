package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/karlseguin/ccache/v2"
)

var (
	defaultExpire = time.Duration(5 * time.Minute)
	defaultCache  = NewCache(2000, 200)
	mlocker       = newLockMaps(1000)
)

var (
	ErrFetchTimeout = errors.New("failed to fetch data timeout !!!")
)

type Cache struct {
	cache *ccache.Cache
}

func NewCache(maxSize int64, itmesToPrune int64) *Cache {
	cc := ccache.New(ccache.Configure().MaxSize(maxSize).ItemsToPrune(uint32(itmesToPrune)))
	return &Cache{
		cache: cc,
	}
}

func (c *Cache) GetSetWithLock(key string, lockTimeout, ttl time.Duration, fetch func() (interface{}, error)) (interface{}, error) {
	lock := mlocker.getLock(key)
	lock.Lock()

	data, expire := c.Get(key)
	if data != nil && !expire {
		lock.Unlock()
		return data, nil
	}

	timer := time.NewTimer(lockTimeout)
	defer timer.Stop()

	var (
		sig = make(chan bool, 0)
		val interface{}

		err error
	)

	go func() {
		val, err = fetch()
		sig <- true
	}()

	select {
	case <-sig:
		c.Setex(key, val, ttl)
		lock.Unlock()
		return val, err

	case <-timer.C:
		lock.Unlock()

		// allow multi worker call fetch()
		// c.GetSet(key, ttl, fetch)

		return nil, ErrFetchTimeout
	}
}

func (c *Cache) GetSet(key string, ttl time.Duration, fetch func() (interface{}, error)) (interface{}, error) {
	v, err := c.cache.Fetch(key, ttl, fetch)
	if err != nil {
		return nil, err
	}
	return v.Value(), err
}

func (c *Cache) GetMust(key string) interface{} {
	val, expired := c.Get(key)
	if expired {
		return nil
	}
	return val
}

func (c *Cache) Get(key string) (interface{}, bool) {
	val := c.cache.Get(key)
	if val == nil {
		return nil, false
	}
	if val.Expired() == true {
		return val.Value(), true
	}

	return val.Value(), false
}

func (c *Cache) Set(key string, val interface{}) {
	c.Setex(key, val, defaultExpire)
}

func (c *Cache) Setex(key string, val interface{}, ttl time.Duration) {
	c.cache.Set(key, val, ttl)
}

func (c *Cache) SetWithExpire(key string, val interface{}, ttlFlag string) {
	ttl, err := time.ParseDuration(ttlFlag)
	if err != nil {
		ttl = defaultExpire
	}

	c.Setex(key, val, ttl)
}

func (c *Cache) Delete(key string) bool {
	return c.cache.Delete(key)
}

func Get(key string) (interface{}, bool) {
	return defaultCache.Get(key)
}

func GetMust(key string) interface{} {
	return defaultCache.GetMust(key)
}

func SetWithExpire(key string, val interface{}, duration string) {
	defaultCache.SetWithExpire(key, val, duration)
}

func Set(key string, val interface{}) {
	defaultCache.Set(key, val)
}

func Setex(key string, val interface{}, duration time.Duration) {
	defaultCache.Setex(key, val, duration)
}

func Delete(key string) bool {
	return defaultCache.Delete(key)
}

func GetSet(key string, duration time.Duration, fetch func() (interface{}, error)) (interface{}, error) {
	return defaultCache.GetSet(key, duration, fetch)
}

func GetSetWithLock(key string, fetchTimeout, ttl time.Duration, fetch func() (interface{}, error)) (interface{}, error) {
	return defaultCache.GetSetWithLock(key, fetchTimeout, ttl, fetch)
}

type lockMaps struct {
	length uint32
	locks  map[int]*sync.Mutex
	mutex  sync.Mutex
}

func newLockMaps(length int) *lockMaps {
	locks := make(map[int]*sync.Mutex, length)
	for i := 0; i < length; i++ {
		locks[i] = &sync.Mutex{}
	}

	lm := &lockMaps{
		length: uint32(length),
		locks:  locks,
	}

	return lm
}

func (l *lockMaps) getLock(key string) *sync.Mutex {
	idx := getMod(key, l.length)
	return l.locks[int(idx)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func getMod(key string, length uint32) uint32 {
	num := fnv32(key)
	return num % length
}
