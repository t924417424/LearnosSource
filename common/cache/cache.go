package cache

import "context"

type Cache interface {
	Lock() error
	UnLock() error
	Set(key string, value string) error
	Get(key string) ([]byte, error)
	GetByPreFix(prefix string) (int, map[string][]byte, error)
	Exists(key string) (bool, []byte)
	SetEx(key string, value string, second int64) error
	Del(key string) error
	Push(prefix string, key string, value string) (bool, error) //队列中不允许有重复Key的消息
	Pop(prefix string) (msg []byte, err error)
	KeepAlive(exitCtx context.Context, key string, value string, second int64, exit chan<- error)
}
