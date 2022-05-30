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
	for i := range c.cacheMap {
		if c.cacheMap[i].Key == key && c.cacheMap[i].ExpirationDate.After(deadline) {
			return c.cacheMap[i].Value, true
		}
	}
	return "", false
}

func (c *CachesMap) Put(key, value string) {
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
	if deadline.After(time.Now()) {
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

func (c *CachesMap) CleanCache() {
	currentTime := time.Now()
	var cleanedMap CachesMap
	for i := range c.cacheMap {
		if c.cacheMap[i].ExpirationDate.After(currentTime) {
			cleanedMap.cacheMap = append(cleanedMap.cacheMap, c.cacheMap[i])
		}
	}
	c.cacheMap = cleanedMap.cacheMap
}
