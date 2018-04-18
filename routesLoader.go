package main

import (
	"github.com/nethruster/linksh/controllers"
	"github.com/thehowl/fasthttprouter"
)

func LoadRoutes(env *controllers.Env, router *fasthttprouter.Router) {
	router.GET("/api/users", env.Auth(env.GetUsers))
	router.GET("/api/users/:id", env.GetUser)
	router.POST("/api/users", env.CreateUser)
	router.PUT("/api/users/:id", env.EditUser)
	router.DELETE("/api/users/:id", env.DeleteUser)

	router.POST("/session/login", env.Login)
	router.POST("/session/logout", env.Auth(env.Logout))
}
