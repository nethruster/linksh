package controllers

import (
	"github.com/erikdubbelboer/fasthttp"
	"github.com/jinzhu/gorm"
	"github.com/nethruster/linksh/models"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"encoding/json"
	"time"
)

type authHeaderData struct {
    UserId    string `json:"userId"`
    SessionId string `json:"sessionId"`
}

func (env Env) Auth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		print(string(ctx.Request.Header.Cookie("auth")))

		if cookie := ctx.Request.Header.Cookie("auth"); cookie != nil {
			data := strings.Split(string(cookie), "|")
			valid, user, err := checkLoginWithSession(env.Db, data[0], data[1])

			if err != nil {
				ctx.Response.Header.SetStatusCode(500)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "Internal server error"}`)
				env.Log.WithFields(logrus.Fields{"event": "Login with session", "status": "Failed"}).Error(err.Error())
				return
			}
			if valid {
				ctx.SetUserValue("currentUser", user)
				handler(ctx)
				return
			} else {
				var cookie fasthttp.Cookie
				cookie.SetKey("auth")
				cookie.SetValue("")
				cookie.SetHTTPOnly(true)
				cookie.SetExpire(time.Unix(0, 0))
				ctx.Response.Header.SetCookie(&cookie)

				ctx.Response.Header.SetStatusCode(403)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "FORBIDDEN"}`)
			}

		} else if auth := ctx.Request.Header.Peek("auth"); auth != nil {
            var data authHeaderData
			json.Unmarshal(auth, &data)

            valid, user, err := checkLoginWithSession(env.Db, data.SessionId, data.UserId)

			if err != nil {
				ctx.Response.Header.SetStatusCode(500)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "Internal server error"}`)
				env.Log.WithFields(logrus.Fields{"event": "Login with session", "status": "Failed"}).Error(err.Error())
				return
			}
			if valid {
				ctx.SetUserValue("currentUser", user)
				handler(ctx)
				return
			} else {
				ctx.Response.Header.SetStatusCode(403)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "FORBIDDEN"}`)
			}
		} else if apikey := string(ctx.Request.Header.Peek("X-Auth-Key")); apikey != "" {
			email := string(ctx.Request.Header.Peek("X-Auth-Email"))

			valid, user, err := checkLoginWithApikey(env.Db, email, apikey)
			if err != nil {
				ctx.Response.Header.SetStatusCode(500)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "Internal server error"}`)
				env.Log.WithFields(logrus.Fields{"event": "Login with apikey", "status": "Failed"}).Error(err.Error())
				return
			}
			if valid {
				ctx.SetUserValue("currentUser", user)
				handler(ctx)
				return
			} else {
				ctx.Response.Header.SetStatusCode(403)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "FORBIDDEN"}`)
			}
		} else {
			ctx.Response.Header.SetStatusCode(401)
			ctx.SetContentType("application/json")
			fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
		}

	})
}

func checkLoginWithApikey(db *gorm.DB, email, apikey string) (bool, models.User, error) {
	var user models.User

	if apikey == "" || email == "" {
		return false, models.User{}, nil
	}

	err := db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, models.User{}, nil
		} else {
			return false, models.User{}, err
		}
	}

	return user.Apikey == apikey, user, nil
}

func checkLoginWithSession(db *gorm.DB, sessionId, userId string) (bool, models.User, error) {
	if sessionId == "" || userId == "" {
		return false, models.User{}, nil
	}
	var session models.Session

	err := db.Where("id = ?", sessionId).Take(&session).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, models.User{}, nil
		} else {
			return false, models.User{}, err
		}
	}

	if session.ExpiresAt.Before(time.Now()) {
		err = db.Delete(&session).Error
		if err != nil {
			return false, models.User{}, err
		}
		return false, models.User{}, nil
	}

	if session.UserId == userId {
		var user models.User

		err := db.Where("id = ?", userId).Take(&user).Error
		if err != nil {
			return false, models.User{}, err
		}

        err = models.UpdateSessionLastUsed(db, session)
        if err != nil {
            return false, models.User{}, err
        }

		return true, user, nil
	} else {
		return false, models.User{}, nil
	}
}
