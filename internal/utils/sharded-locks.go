package utils

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

// ShardedLocks 分片鎖
type ShardedLocks struct {
	locks     []*sync.RWMutex
	numShards uint64 // 分片數量
}

func NewShardedLocks(numShards uint64) *ShardedLocks {
	minNums := uint64(2)
	maxNums := uint64(1024)
	if numShards <= minNums {
		numShards = minNums
	}
	if numShards > maxNums {
		numShards = maxNums
	}

	sl := &ShardedLocks{
		locks:     make([]*sync.RWMutex, numShards),
		numShards: numShards,
	}
	for i := range sl.locks {
		sl.locks[i] = &sync.RWMutex{}
	}
	return sl
}

func (sl *ShardedLocks) GetLockForKey(key string) *sync.RWMutex {
	hash := xxhash.Sum64String(key)
	return sl.locks[hash%sl.numShards]
}
