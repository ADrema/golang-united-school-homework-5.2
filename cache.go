package cache

import "time"

type Element struct {
	Key            string
	Value          string
	ExpirationDate time.Time
}

type Cache struct {
	cacheMap []Element
}

func NewCache() Cache {
	return Cache{}
}

func (c *Cache) Get(key string) (string, bool) {
	deadline := time.Now()
	c.CleanCache(deadline)
	for i := range c.cacheMap {
		if c.cacheMap[i].Key == key && c.cacheMap[i].ExpirationDate.After(deadline) {
			return c.cacheMap[i].Value, true
		}
	}
	return "", false
}

func (c *Cache) Put(key, value string) {
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
		c.cacheMap = append(c.cacheMap, Element{key, value, infiniteDate})
	}

}

func (c *Cache) Keys() []string {
	var keys []string
	for i := range c.cacheMap {
		keys = append(keys, c.cacheMap[i].Key)
	}
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
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
			c.cacheMap = append(c.cacheMap, Element{key, value, deadline})
		}
	}
}

func (c *Cache) CleanCache(deadline time.Time) {
	var cleanedMap Cache
	for i := range c.cacheMap {
		if c.cacheMap[i].ExpirationDate.After(deadline) {
			cleanedMap.cacheMap = append(cleanedMap.cacheMap, c.cacheMap[i])
		}
	}
	c.cacheMap = cleanedMap.cacheMap
}
