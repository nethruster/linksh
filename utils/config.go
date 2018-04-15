package utils

import (
	"github.com/go-ini/ini"
	"fmt"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host         string
	Port         uint16
	DatabaseName string
	User         string
	Password     string
}

func (self DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		self.User,
		self.Password,
		self.Host,
		self.Port,
		self.DatabaseName,
	)
}

func ParseConfigFile(filePath string) (*Config, error){
	config := new(Config)
	err := ini.MapTo(config, filePath)
	return config, err
}