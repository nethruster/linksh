package main

import (
	"github.com/nethruster/linksh/controllers"
	"github.com/thehowl/fasthttprouter"
)

func LoadRoutes(env *controllers.Env, router *fasthttprouter.Router) {
	router.GET("/api/users", env.Auth(env.GetUsers))
    router.GET("/api/users/:id", env.Auth(env.GetUser))
    router.POST("/api/users", env.Auth(env.CreateUser))
    router.PUT("/api/users/:id", env.Auth(env.EditUser))
    router.DELETE("/api/users/:id", env.Auth(env.DeleteUser))

    router.GET("/api/links", env.Auth(env.GetLinks))
    router.GET("/api/links/:id", env.Auth(env.GetLink))
    router.POST("/api/links", env.Auth(env.CreateLink))
    router.PUT("/api/links/:id", env.Auth(env.EditLink))
    router.DELETE("/api/links/:id", env.Auth(env.DeleteLink))

    router.GET("/api/sessions", env.Auth(env.GetSessions))
    router.GET("/api/sessions/:id", env.Auth(env.GetSession))
    router.POST("/api/sessions", env.Auth(env.CreateSession))
    router.DELETE("/api/sessions/:id", env.Auth(env.DeleteSession))

    router.POST("/auth/login", env.Login)
    router.POST("/auth/logout", env.Auth(env.Logout))
}
