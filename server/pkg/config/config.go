package config

import (
	"sync"

	"github.com/spf13/viper"
)

const (
	EnvProd = "prod"

	S = "68B385283381EBA1FD9C6C8EF3ECDAD2"
)

func IsProd() bool {
	return conf.Env == EnvProd
}

type ServerConfig struct {
	Env               string      `mapstructure:"env"`
	PaymentHttpConfig HttpConfig  `mapstructure:"payment_http"`
	WebHttpConfig     HttpConfig  `mapstructure:"web_http"`
	MysqlConfig       MysqlConfig `mapstructure:"mysql"`
	RedisConfig       RedisConfig `mapstructure:"redis"`
}

type HttpConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MysqlConfig struct {
	Uri string `mapstructure:"uri"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

var (
	conf *ServerConfig
	once sync.Once
)

func New(path string) *ServerConfig {
	once.Do(func() {
		vp := getConfig(path)
		err := vp.Unmarshal(&conf)
		if err != nil {
			panic(err)
		}
	})

	return conf
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)

	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return conf
}

func Get() *ServerConfig {
	return conf
}
