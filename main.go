package main

import (
	"github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thehowl/fasthttprouter"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/nethruster/linksh/utils"
	"os"
	"github.com/nethruster/linksh/controllers"
	"github.com/nethruster/linksh/models"
)

func main() {
	var log = logrus.New()
	router := fasthttprouter.New()
	conf, err := utils.ParseConfigFile("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open("mysql", conf.Database.GetConnectionString())
	if err != nil {
		println("Error while connecting to the database")
		os.Exit(1009)
	}
	defer db.Close()

	conf.Server.AdjustLogSettings(log, db)

	db.AutoMigrate(&models.User{}, &models.Link{}, &models.Session{})

	env := controllers.Env{
		Config: conf,
		Db: db,
		Log: log,
	}

	LoadRoutes(&env, router)

	log.WithField("event", "Start").Info("Server listening at: ", conf.Server.GetListenString())
	log.Fatal(fasthttp.ListenAndServe(conf.Server.GetListenString(), router.Handler))
}