package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/nethruster/linksh/utils"
	"github.com/sirupsen/logrus"
)

type Env struct {
	Config *utils.Config
	Db *gorm.DB
	Log *logrus.Logger
}