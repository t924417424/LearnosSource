package etcd

import "time"

type options struct {
	Addr []string
	Username string
	Password string
	DialTimeout time.Duration
}

type option func(*options)

func Addr(addr []string) option {
	return func(o *options) {
		o.Addr = addr
	}
}

func TimeOut(duration time.Duration) option{
	return func(o *options) {
		o.DialTimeout = duration
	}
}

func Auth(username,password string) option{
	return func(o *options) {
		o.Username = username
		o.Password = password
	}
}