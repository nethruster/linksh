package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/nethruster/linksh/models"
	"github.com/sirupsen/logrus"
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
		env.Log.WithFields(logrus.Fields{"event": "Get session", "status": "Failed"}).Error(err.Error())
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
		env.Log.WithFields(logrus.Fields{"event": "Create session", "status": "Failed"}).Error(err.Error())
		return
	}

	ctx.Response.Header.SetStatusCode(201)
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
		env.Log.WithFields(logrus.Fields{"event": "Delete session", "status": "Failed"}).Error(err.Error())
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
		env.Log.WithFields(logrus.Fields{"event": "Delete session", "status": "Failed"}).Error(err.Error())
		return
	}

	ctx.Response.Header.SetStatusCode(204)
}
