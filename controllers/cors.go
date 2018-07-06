package controllers

import (
	"github.com/erikdubbelboer/fasthttp"
)

func (env Env) Cors(ctx *fasthttp.RequestCtx) {
	env.addCorsHeader(ctx)
	ctx.Response.Header.Add("Content-Type", "text/plain; charset=utf-8")
	ctx.Response.Header.Add("Content-Length", "0")
	ctx.Response.SetStatusCode(204)
}

func (env Env) CorsMiddleware(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		env.addCorsHeader(ctx)
		handler(ctx)
	})
}

func (env Env) addCorsHeader(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", env.Config.Server.CORSString)
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Response.Header.Add("Access-Control-Max-Age", "1728000")
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,X-Auth,X-Auth-Key,X-Auth-Email")
	ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
}
