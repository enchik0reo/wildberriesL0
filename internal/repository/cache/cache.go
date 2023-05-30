package cache

import (
	"fmt"
	"sync"

	"github.com/enchik0reo/wildberriesL0/internal/models"
)

type Cache struct {
	m     map[string][]byte
	slice []string
	count int
	size  int
	sync.RWMutex
}

func New(size int) *Cache {
	return &Cache{
		m:       make(map[string][]byte, size),
		slice:   make([]string, size),
		count:   0,
		size:    size,
		RWMutex: sync.RWMutex{},
	}
}

func (c *Cache) Save(o models.Order) error {
	if ok := c.Check(o.Uid); !ok {
		c.Lock()
		c.m[o.Uid] = o.Details
		c.Unlock()

		delete(c.m, c.slice[c.count])
		c.slice[c.count] = o.Uid
		c.count++
		if c.count >= c.size {
			c.count = 0
		}
		return nil
	}
	return fmt.Errorf("order already exists")
}

func (c *Cache) Check(uid string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.m[uid]
	return ok
}

func (c *Cache) GetById(uid string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	msg, ok := c.m[uid]
	if !ok {
		return nil, fmt.Errorf("doesn't exist in cache; ")
	}
	return msg, nil
}
