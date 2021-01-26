package config

import (
	"github.com/BurntSushi/toml"
)

type conf struct {
	Common struct {
		JwtKey string `toml:"jwtKey"`
	} `toml:"common"`
	WebSocket struct {
		GateWay struct {
			WsPort int `toml:"wsPort"`
		} `toml:"gateWay"`
		Container struct {
			WsPort int `toml:"wsPort"`
		} `toml:"container"`
	} `toml:"webSocket"`
	Etcd struct {
		Addr     []string `toml:"addr"`
		UserName string   `toml:"userName"`
		PassWord string   `toml:"passWord"`
	} `toml:"etcd"`
	Web struct {
		RunMode string `toml:"runMode"`
		WebAddr string `toml:"webAddr"`
		SSHAddr string `toml:"sshAddr"`
	} `toml:"web"`
	GateWay struct {
		AliSms struct {
			AccessID  string `toml:"accessId"`
			AccessKey string `toml:"accessKey"`
			SignName  string `toml:"signName"`
			Template  string `toml:"template"`
		} `toml:"aliSms"`
		Mysql struct {
			Addr     string `toml:"addr"`
			Port     int    `toml:"port"`
			UserName string `toml:"userName"`
			PassWord string `toml:"passWord"`
			DbName   string `toml:"dbName"`
		} `toml:"mysql"`
	} `toml:"gateWay"`
}

var config = new(conf)

func ReadConf(file string) (err error) {
	//读取配置文件
	_, err = toml.DecodeFile(file, config)
	return err
}

func GetConf() conf {
	return *config
}