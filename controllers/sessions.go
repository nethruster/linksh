package controllers

import (
	"github.com/erikdubbelboer/fasthttp"
	"encoding/json"
	"fmt"
	"github.com/nethruster/linksh/models"
	"github.com/sirupsen/logrus"
	"time"
	"strings"
    "strconv"
)

func (env Env) GetSessions(ctx *fasthttp.RequestCtx) {
    currentUser := ctx.UserValue("currentUser").(models.User)
    args := ctx.QueryArgs()
    var sessions []models.Session
    ownerId := string(args.Peek("ownerId"))
    if ownerId != currentUser.Id && !currentUser.IsAdmin {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }
    query := env.Db
    if offset, err := strconv.Atoi(string(args.Peek("offset"))); err == nil && offset != 0 {
        query = query.Offset(offset)
    }
    if limit, err := strconv.Atoi(string(args.Peek("limit"))); err == nil && limit != 0 {
        query = query.Limit(limit)
    }
    if ownerId != "" {
        query = query.Where("user_id = ?", ownerId)
    }

    err := query.Find(&sessions).Error
    if err != nil {
        ctx.Response.Header.SetStatusCode(500)
        fmt.Fprint(ctx, `{"error": "Internal server error"}`)
        env.Log.WithFields(logrus.Fields{"event": "Lists sessions", "status": "Failed"}).Error(err.Error())
        return
    }

    ctx.SetContentType("application/json")
    json.NewEncoder(ctx).Encode(&sessions)
}

func (env Env) GetSession(ctx *fasthttp.RequestCtx) {
    var session models.Session
    currentUser := ctx.UserValue("currentUser").(models.User)
    id := ctx.UserValue("id")

    ctx.SetContentType("application/json")

    err := env.Db.Where("id = ?", id).Take(&session).Error
    if err != nil {
        if err.Error() == "record not found" {
            ctx.Response.Header.SetStatusCode(404)
            fmt.Fprint(ctx, `{"error": "Session not found"}`)
            return
        }
        ctx.Response.Header.SetStatusCode(500)
        fmt.Fprint(ctx, `{"error": "Internal server error"}`)
        env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
        return
    }
    if !currentUser.IsAdmin && session.UserId != currentUser.Id {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

    json.NewEncoder(ctx).Encode(&session)
}
func (env Env) CreateSession(ctx *fasthttp.RequestCtx) {
    currentUser := ctx.UserValue("currentUser").(models.User)
    session, err := models.CreateSession(env.Db, currentUser)

    ctx.SetContentType("application/json")

    if err != nil {
        ctx.Response.Header.SetStatusCode(500)
        fmt.Fprint(ctx, `{"error": "Internal server error"}`)
        env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
        return
    }

    json.NewEncoder(ctx).Encode(&session)
}

func (env Env) DeleteSession(ctx *fasthttp.RequestCtx) {
    var session models.Session
    currentUser := ctx.UserValue("currentUser").(models.User)
    id := ctx.UserValue("id")

    ctx.SetContentType("application/json")

    err := env.Db.Where("id = ?", id).Take(&session).Error
    if err != nil {
        if err.Error() == "record not found" {
            ctx.Response.Header.SetStatusCode(404)
            fmt.Fprint(ctx, `{"error": "Session not found"}`)
            return
        }
        ctx.Response.Header.SetStatusCode(500)
        fmt.Fprint(ctx, `{"error": "Internal server error"}`)
        env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
        return
    }
    if !currentUser.IsAdmin && session.UserId != currentUser.Id {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

    err = env.Db.Delete(&session).Error
    if err != nil {
        ctx.Response.Header.SetStatusCode(500)
        fmt.Fprint(ctx, `{"error": "Internal server error"}`)
        env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
        return
    }

    ctx.Response.Header.SetStatusCode(204)
}

func (env Env) Login(ctx *fasthttp.RequestCtx) {
	var data map[string]string
	var user models.User
	ctx.SetContentType("application/json")
	json.Unmarshal(ctx.Request.Body(), &data)
	if data["email"] == "" || data["password"] == "" {
		ctx.Response.Header.SetStatusCode(400)
		fmt.Fprint(ctx, `{"error": "Missing email or password"}`)
	}

	err := env.Db.Where("email = ?", data["email"]).Take(&user).Error
	if err != nil && err.Error() != "record not found" {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Login", "status": "Failed"}).Error(err.Error())
		return
	}

	if user.Id == "" || !user.CheckIfCorrectPassword([]byte(data["password"])) {
		ctx.Response.Header.SetStatusCode(400)
		fmt.Fprint(ctx, `{"error": "The email or the password are invalid"}`)
		return
	}

	var expires time.Time

	if data["notExpire"] == "true" {
		expires = time.Now().AddDate(100, 0, 0)
	} else {
		expires = time.Now().AddDate(0, 0, 1)
	}
	id, err := models.GenerateSessionId()
	if err != nil {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Login", "status": "Failed"}).Error(err.Error())
		return
	}
	session := models.Session{
		Id:        id,
		UserId:    user.Id,
		ExpiresAt: expires,
	}
	err = env.Db.Create(&session).Error
	if err != nil {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Login", "status": "Failed"}).Error(err.Error())
		return
	}

	if data["useCookie"] == "true" {
		var cookie fasthttp.Cookie
		cookie.SetKey("auth")
		cookie.SetValue(fmt.Sprintf("%v|%v", session.Id, user.Id))
		cookie.SetHTTPOnly(true)
		cookie.SetExpire(session.ExpiresAt.AddDate(0, 0, 2))
		ctx.Response.Header.SetCookie(&cookie)
	}

	fmt.Fprintf(ctx, `{"sessionId": "%v", "userId": "%v", "expiresAt": "%v"}"`, session.Id, user.Id, session.ExpiresAt)
}

func (env Env) Logout(ctx *fasthttp.RequestCtx) {
	if cookie := ctx.Request.Header.Cookie("auth"); cookie != nil {
		data := strings.Split(string(cookie), "|")
		err := env.Db.Delete(models.Session{}, "id = ?", data[0]).Error
		if err != nil {
			ctx.Response.Header.SetStatusCode(500)
			fmt.Fprint(ctx, `{"error": "Internal server error"}`)
			env.Log.WithFields(logrus.Fields{"event": "Logout", "status": "Failed"}).Error(err.Error())
			return
		}
		var cookie fasthttp.Cookie
		cookie.SetKey("auth")
		cookie.SetValue("")
		cookie.SetHTTPOnly(true)
		cookie.SetExpire(time.Unix(0, 0))
		ctx.Response.Header.SetCookie(&cookie)

		ctx.Response.Header.SetStatusCode(204)
	} else if auth := ctx.Request.Header.Peek("auth"); auth != nil {
		var data map[string]string
		json.Unmarshal(auth, &data)

		err := env.Db.Delete(models.Session{}, "id = ?", data["sessionId"]).Error
		if err != nil {
			ctx.Response.Header.SetStatusCode(500)
			fmt.Fprint(ctx, `{"error": "Internal server error"}`)
			env.Log.WithFields(logrus.Fields{"event": "Logout", "status": "Failed"}).Error(err.Error())
			return
		}
		ctx.Response.Header.SetStatusCode(204)
	}
}
