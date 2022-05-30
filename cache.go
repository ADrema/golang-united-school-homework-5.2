package cache

import "time"

type Cache struct {
	Key            string
	Value          string
	ExpirationDate time.Time
}

type CachesMap struct {
	cacheMap []Cache
}

func NewCache() Cache {
	return Cache{}
}

func (c *CachesMap) Get(key string) (string, bool) {
	deadline := time.Now()
	c.CleanCache(deadline)
	for i := range c.cacheMap {
		if c.cacheMap[i].Key == key && c.cacheMap[i].ExpirationDate.After(deadline) {
			return c.cacheMap[i].Value, true
		}
	}
	return "", false
}

func (c *CachesMap) Put(key, value string) {
	deadline := time.Now()
	c.CleanCache(deadline)
	infiniteDate := time.Unix(1<<63-62135596801, 999999999)
	var exists = false
	for i := range c.cacheMap {
		if c.cacheMap[i].Key == key {
			c.cacheMap[i].Value = value
			c.cacheMap[i].ExpirationDate = infiniteDate
			exists = true
		}
	}

	if !exists {
		c.cacheMap = append(c.cacheMap, Cache{key, value, infiniteDate})
	}

}

func (c *CachesMap) Keys() []string {
	var keys []string
	for i := range c.cacheMap {
		keys = append(keys, c.cacheMap[i].Key)
	}
	return keys
}

func (c *CachesMap) PutTill(key, value string, deadline time.Time) {
	currentTime := time.Now()
	c.CleanCache(currentTime)
	if deadline.After(currentTime) {
		var exists = false
		for i := range c.cacheMap {
			if c.cacheMap[i].Key == key {
				c.cacheMap[i].Value = value
				c.cacheMap[i].ExpirationDate = deadline
				exists = true
			}
		}

		if !exists {
			c.cacheMap = append(c.cacheMap, Cache{key, value, deadline})
		}
	}
}

func (c *CachesMap) CleanCache(deadline time.Time) {
	var cleanedMap CachesMap
	for i := range c.cacheMap {
		if c.cacheMap[i].ExpirationDate.After(deadline) {
			cleanedMap.cacheMap = append(cleanedMap.cacheMap, c.cacheMap[i])
		}
	}
	c.cacheMap = cleanedMap.cacheMap
}
