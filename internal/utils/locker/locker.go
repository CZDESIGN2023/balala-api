package locker

import (
	"go-cs/internal/utils"
	"sync"

	"github.com/spf13/cast"
)

var _sharedLocks *utils.ShardedLocks

func init() {
	_sharedLocks = utils.NewShardedLocks(1024)
}

func Lock(key string) *sync.RWMutex {
	return _sharedLocks.GetLockForKey(key)
}

func NewWorkItemLockKey(workItemId int64) string {
	return "balala:workItem:" + cast.ToString(workItemId)
}
