package cache

import (
	"sync"
)

// Partition represents a cache partition where Load is the callback
// for when the partition is needed.
type Partition struct {
	Key  string
	Load func() (map[string]interface{}, error)
}

type lazyPartition struct {
	initSync sync.Once
	cache    map[string]interface{}
	load     func() (map[string]interface{}, error)
}

func (l *lazyPartition) init() error {
	var err error
	l.initSync.Do(func() {
		var c map[string]interface{}
		c, err = l.load()
		l.cache = c
	})

	return err
}

// PartitionedLazyCache is a lazily loaded cache paritioned by a supplied string key.
type PartitionedLazyCache struct {
	partitions map[string]*lazyPartition
}

// NewPartitionedLazyCache creates a new NewPartitionedLazyCache with the supplied
// partitions.
func NewPartitionedLazyCache(partitions ...Partition) *PartitionedLazyCache {
	lazyPartitions := make(map[string]*lazyPartition, len(partitions))
	for _, partition := range partitions {
		lazyPartitions[partition.Key] = &lazyPartition{load: partition.Load}
	}
	cache := &PartitionedLazyCache{partitions: lazyPartitions}

	return cache
}

// Get initializes the partition if not already done so, then looks up the given
// key in the given partition, returns nil if no value found.
func (c *PartitionedLazyCache) Get(partition, key string) (interface{}, error) {
	p, found := c.partitions[partition]

	if !found {
		return nil, nil
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	if v, found := p.cache[key]; found {
		return v, nil
	}

	return nil, nil

}
