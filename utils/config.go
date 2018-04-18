package utils

import (
	"github.com/go-ini/ini"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Address            string
	Port               uint16
	AllowRegister      bool
	LogLevel           string
	LogDatabaseQueries bool
}

type DatabaseConfig struct {
	Host         string
	Port         uint16
	DatabaseName string
	User         string
	Password     string
}

func (config ServerConfig) AdjustLogSettings(log *logrus.Logger, db *gorm.DB) {
	var level logrus.Level
	switch config.LogLevel {
	case "debug":
		level = logrus.DebugLevel
		break
	case "info":
		level = logrus.InfoLevel
		break
	case "warn":
		level = logrus.WarnLevel
		break
	case "error":
		level = logrus.ErrorLevel
		break
	}
	log.SetLevel(level)
	db.LogMode(config.LogDatabaseQueries)
}

func (config ServerConfig) GetListenString() string {
	return fmt.Sprintf("%v:%v", config.Address, config.Port)
}

func (config DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName,
	)
}

func ParseConfigFile(filePath string) (*Config, error){
	config := new(Config)
	err := ini.MapTo(config, filePath)
	return config, err
}