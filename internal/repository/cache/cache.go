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

func (c *Cache) Save(o models.Order) error {
	if ok := c.Check(o.Uid); !ok {
		c.Lock()
		c.m[o.Uid] = o.Details
		c.Unlock()
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
		return nil, fmt.Errorf("order with uid: %s doesn't exist in cache", uid)
	}
	return msg, nil
}
