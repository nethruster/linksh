package controllers

import (
    "fmt"
	"github.com/nethruster/linksh/pkg/models"
    "github.com/nethruster/linksh/pkg/models/sessions"
    "github.com/valyala/fasthttp"
)

//GetFasthttpSessionMiddleware returns a middleware for fasyhttp
//if strict is set to true the middleware will cancel the connection if a valid session is not provided
func GetFasthttpSessionMiddleware(db *models.DB,manager *sessions.SessionManager,strict bool, cookiePrefix string,handler fasthttp.RequestHandler) fasthttp.RequestHandler {
    cookieName := fmt.Sprintf("%v%v",cookiePrefix, "auth")
	if strict {
		return func(ctx *fasthttp.RequestCtx) {
			if cookie := ctx.Request.Header.Cookie("linksh-auth"); cookie != nil {
				data := string(cookie)
				ok, session, err := manager.CheckValidSession(data[:36], data[36:])

				if err != nil {
					ctx.Error("Internal server error", 500)
					return
				}

				if ok {
                    ctx.SetUserValue("session", session)
                    handler(ctx)
                    return
                }
                ctx.Error("Invalid session", 400)
			} else if auth := ctx.Request.Header.Peek("X-Auth"); auth != nil {
                data := string(auth)
                ok, session, err := manager.CheckValidSession(data[:36], data[36:])
                if err != nil {
					ctx.Error("Internal server error", 500)
					return
				}

                if ok {
                    ctx.SetUserValue("session", session)
                    handler(ctx)
                    return
                }
                ctx.Error("Invalid session", 400)
            } else if apikeyBytes := ctx.Request.Header.Peek("X-Auth-Key"); apikeyBytes != nil {
                emailBytes := ctx.Request.Header.Peek("X-Auth-Email")
                if emailBytes == nil {
                    ctx.Error("Missing email", 400)
                    return
                }
                var user models.User
                err := db.Where("email = ? AND apikey = ?", string(emailBytes), string(apikeyBytes)).Take(&user).Error
                if err != nil {
                    if err.Error() == "record not found" {
                        ctx.Error("Bad apikey or email", 400)
                    }
                    handler(ctx)
                    return
                }

                ctx.SetUserValue("user", user)
                ctx.SetUserValue("session", sessions.Session{ID: "Temporary", OwnerID: user.ID})

            } else {
                ctx.Error("Unauthorized", 401)
            }
		}
	}

	return func(ctx *fasthttp.RequestCtx) {
		if cookie := ctx.Request.Header.Cookie(cookieName); cookie != nil {
			data := string(cookie)
			ok, session, _ := manager.CheckValidSession(data[:36], data[36:])

			if ok {
				ctx.SetUserValue("session", session)
            }
		} else if auth := ctx.Request.Header.Peek("X-Auth"); auth != nil {
            data := string(auth)
            ok, session, _ := manager.CheckValidSession(data[:36], data[36:])
            if ok {
                ctx.SetUserValue("session", session)
            }
        } else if apikeyBytes := ctx.Request.Header.Peek("X-Auth-Key"); apikeyBytes != nil {
            emailBytes := ctx.Request.Header.Peek("X-Auth-Email")
            if emailBytes == nil {
                handler(ctx)
                return
            }
            var user models.User
            err := db.Where("email = ? AND apikey = ?", string(emailBytes), string(apikeyBytes)).Take(&user).Error
            if err != nil {
                handler(ctx)
                return
            }

            ctx.SetUserValue("user", user)
            ctx.SetUserValue("session", sessions.Session{ID: "Temporary", OwnerID: user.ID})
        }
        handler(ctx)
	}
}
