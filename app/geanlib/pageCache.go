package geanlib

import (
	"sync"
)

type pageCache struct {
	sync.RWMutex
	m map[string][][2]Pages
}

func newPageCache() *pageCache {
	return &pageCache{m: make(map[string][][2]Pages)}
}

// get gets a Pages slice from the cache matching the given key and Pages slice.
// If none found in cache, a copy of the supplied slice is created.
//
// If an apply func is provided, that func is applied to the newly created copy.
//
// The cache and the execution of the apply func is protected by a RWMutex.
func (c *pageCache) get(key string, p Pages, apply func(p Pages)) (Pages, bool) {
	c.RLock()
	if cached, ok := c.m[key]; ok {
		for _, ps := range cached {
			if fastEqualPages(p, ps[0]) {
				c.RUnlock()
				return ps[1], true
			}
		}

	}
	c.RUnlock()

	c.Lock()
	defer c.Unlock()

	// double-check
	if cached, ok := c.m[key]; ok {
		for _, ps := range cached {
			if fastEqualPages(p, ps[0]) {
				return ps[1], true
			}
		}
	}

	pagesCopy := append(Pages(nil), p...)

	if apply != nil {
		apply(pagesCopy)
	}

	if v, ok := c.m[key]; ok {
		c.m[key] = append(v, [2]Pages{p, pagesCopy})
	} else {
		c.m[key] = [][2]Pages{{p, pagesCopy}}
	}

	return pagesCopy, false

}

// "fast" as in: we do not compare every element for big slices, but that is
// good enough for our use cases.
// TODO(bep) there is a similar method in pagination.go. DRY.
func fastEqualPages(p1, p2 Pages) bool {
	if p1 == nil && p2 == nil {
		return true
	}

	if p1 == nil || p2 == nil {
		return false
	}

	if p1.Len() != p2.Len() {
		return false
	}

	if p1.Len() == 0 {
		return true
	}

	step := 1

	if len(p1) >= 50 {
		step = len(p1) / 10
	}

	for i := 0; i < len(p1); i += step {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}
