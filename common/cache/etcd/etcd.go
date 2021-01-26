package etcd

import (
	"Learnos/common/cache"
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/google/uuid"
	"time"
)

type etcd struct {
	client *clientv3.Client
	mutex  *concurrency.Mutex
}

func NewCache(opt ...option) (client cache.Cache, err error) {
	conf := &options{}
	var config clientv3.Config
	e := &etcd{}
	for _, o := range opt {
		o(conf)
	}
	if len(conf.Addr) > 0 {
		config.Endpoints = conf.Addr
	}
	config.Username = conf.Username
	config.Password = conf.Password
	config.DialTimeout = conf.DialTimeout
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	session, err := concurrency.NewSession(cli)
	if err != nil {
		return nil, err
	}
	e.client = cli
	e.mutex = concurrency.NewMutex(session, "/create/lock")
	return e, nil
}

func (e etcd) Lock() error {
	return e.mutex.Lock(context.TODO())
}

func (e etcd) UnLock() error {
	return e.mutex.Unlock(context.TODO())
}

func (e etcd) Set(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := e.client.Put(ctx, key, value)
	return err
}

func (e etcd) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	v, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v.Count < 1 {
		return []byte{}, err
	}
	return v.Kvs[0].Value, err
}

func (e etcd) GetByPreFix(prefix string) (int, map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	v, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	var result = make(map[string][]byte)
	if err != nil {
		return 0, result, err
	}
	if v.Count < 1 {
		return int(v.Count), result, err
	}
	for _, v := range v.Kvs {
		result[string(v.Key)] = v.Value
	}
	return int(v.Count), result, err
}

func (e etcd) Exists(key string) (bool, []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	v, err := e.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return false, []byte{}
	}
	if v.Count < 1 {
		return false, []byte{}
	}
	return true, v.Kvs[0].Value
}

func (e etcd) SetEx(key string, value string, second int64) error {
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()
	lease, err := e.client.Grant(ctx1, second)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = e.client.Put(ctx, key, value, clientv3.WithLease(lease.ID))
	return err
}

func (e etcd) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e etcd) Push(prefix string, key string, value string) (bool, error) {
	if exists, _ := e.Exists(prefix + key); exists {
		return false, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := e.client.Put(ctx, prefix+key+uuid.New().String(), value)
	return true, err
}

func (e etcd) Pop(prefix string) (msg []byte, err error) {
	if err = e.Lock(); err != nil {
		return msg, err
	}
	defer e.UnLock()
	data, err := e.client.Get(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend))
	if err != nil {
		return msg, err
	}
	if data.Count > 0 {
		key := string(data.Kvs[0].Key)
		defer e.Del(key)
		msg = data.Kvs[0].Value
	}
	//for _, v := range data.Kvs {
	//	//log.Println(string(v.Key))
	//	msg[string(v.Key)] = v.Value
	//}
	return msg, err
	//ch := make(chan []byte,1)
	//c := <-e.client.Watch(context.Background(), key)
	//ch <- c.Events[0].Kv.Value
	//return ch
	//for {
	//	select {
	//	case c := <-e.client.Watch(context.Background(), prefix):
	//		recv <- c.Events[0].Kv.Value
	//	}
	//}
}

func (e etcd) KeepAlive(exitCtx context.Context, key string, value string, second int64, exit chan<- error) {
	lease := clientv3.NewLease(e.client)
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()
	leaseRsp, err := e.client.Grant(ctx1, second)
	if err != nil {
		exit <- err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = e.client.Put(ctx, key, value, clientv3.WithLease(leaseRsp.ID))
	if err != nil {
		exit <- err
	}
	keepRespChan, err := lease.KeepAlive(context.TODO(), leaseRsp.ID)
	if err != nil {
		exit <- err
	}
	for {
		select {
		case _ = <-keepRespChan:
			if keepRespChan == nil {
				return
			}
		case <-exitCtx.Done():
			return
		}
	}
}
