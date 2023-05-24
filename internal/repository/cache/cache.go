package cache

import (
	"fmt"
	"sync"

	"github.com/enchik0reo/wildberriesL0/internal/models"
)

type Cache struct {
	m map[string][]byte
	sync.RWMutex
}

func New() *Cache {
	return &Cache{
		m:       make(map[string][]byte),
		RWMutex: sync.RWMutex{},
	}
}

func (c *Cache) Save(o models.Order) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.m[o.Uid]; !ok {
		c.m[o.Uid] = o.Details
	}
}

func (c *Cache) Check(uid string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()

	msg, ok := c.m[uid]
	if !ok {
		return nil, false
	}

	return msg, true
}

func (c *Cache) GetById(uid string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()

	msg, ok := c.Check(uid)
	if !ok {
		return nil, fmt.Errorf("order with uid: %s doesn't exist in cache", uid)
	}

	return msg, nil
}
