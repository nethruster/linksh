package main

import (
	"github.com/erikdubbelboer/fasthttp"
	"github.com/nethruster/linksh/controllers"
	"github.com/thehowl/fasthttprouter"
)

type Action struct {
	Method  string
	Handler fasthttp.RequestHandler
}

type routeStruct struct {
	Path    string
	Actions []Action
}

func LoadRoutes(env *controllers.Env, router *fasthttprouter.Router) {
	routes := []routeStruct{
		{"/api/users", []Action{
			{"GET", env.Auth(env.GetUsers)},
			{"POST", env.Auth(env.CreateUser)},
		}},
		{"/api/users/:id", []Action{
			{"GET", env.Auth(env.GetUser)},
			{"PUT", env.Auth(env.EditUser)},
			{"DELETE", env.Auth(env.DeleteUser)},
		}},
		{"/api/links", []Action{
			{"GET", env.Auth(env.GetLinks)},
			{"POST", env.Auth(env.CreateLink)},
		}},
		{"/api/links/:id", []Action{
			{"GET", env.Auth(env.GetLink)},
			{"PUT", env.Auth(env.EditLink)},
			{"DELETE", env.Auth(env.DeleteLink)},
		}},
		{"/api/sessions", []Action{
			{"GET", env.Auth(env.GetSessions)},
			{"POST", env.Auth(env.CreateSession)},
		}},
		{"/api/sessions/:id", []Action{
			{"GET", env.Auth(env.GetSession)},
			{"DELETE", env.Auth(env.DeleteSession)}}},
		{"/api/auth/login", []Action{
			{"POST", env.Login},
		}},
		{"/api/auth/register", []Action{{"POST", env.Register}}},
		{"/api/auth/logout", []Action{{"POST", env.Auth(env.Logout)}}},
	}

	for _, route := range routes {
		if env.Config.Server.EnableCORS {
			for _, action := range route.Actions {
				router.Handle(action.Method, route.Path, env.CorsMiddleware(action.Handler))
			}
			router.OPTIONS(route.Path, env.Cors)
		} else {
			for _, action := range route.Actions {
				router.Handle(action.Method, route.Path, action.Handler)
			}
		}
	}
}
