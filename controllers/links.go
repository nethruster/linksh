package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/nethruster/linksh/models"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (env Env) GetLinks(ctx *fasthttp.RequestCtx) {
	currentUser := ctx.UserValue("currentUser").(models.User)
	args := ctx.QueryArgs()
	var links []models.Link
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

	err := query.Find(&links).Error
	if err != nil {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Lists links", "status": "Failed"}).Error(err.Error())
		return
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(&links)
}

func (env Env) GetLink(ctx *fasthttp.RequestCtx) {
	var link models.Link
	currentUser := ctx.UserValue("currentUser").(models.User)
	id := ctx.UserValue("id")

	ctx.SetContentType("application/json")

	err := env.Db.Where("id = ?", id).Take(&link).Error
	if err != nil {
		if err.Error() == "record not found" {
			ctx.Response.Header.SetStatusCode(404)
			fmt.Fprint(ctx, `{"error": "Link not found"}`)
			return
		}
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Get link", "status": "Failed"}).Error(err.Error())
		return
	}
	if !currentUser.IsAdmin && link.UserId != currentUser.Id {
		ctx.Response.Header.SetStatusCode(401)
		ctx.SetContentType("application/json")
		fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
		return
	}

	json.NewEncoder(ctx).Encode(&link)
}

type createLinkRequest struct {
	CustomId string `json:"customId"`
	Content  string `json:"content"`
}

func (env Env) CreateLink(ctx *fasthttp.RequestCtx) {
	currentUser := ctx.UserValue("currentUser").(models.User)
	var data createLinkRequest
	json.Unmarshal(ctx.Request.Body(), &data)

	link, err := models.CreateLink(env.Db, currentUser.Id, data.CustomId, data.Content)

	ctx.SetContentType("application/json")

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			ctx.Response.Header.SetStatusCode(400)
			fmt.Fprint(ctx, `{"error": "Link already exists"}`)
			return
		}
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Create link", "status": "Failed"}).Error(err.Error())
		return
	}
	ctx.Response.Header.SetStatusCode(201)
	json.NewEncoder(ctx).Encode(&link)
}

type editLinkRequest struct {
	Content string `json:"content"`
}

func (env Env) EditLink(ctx *fasthttp.RequestCtx) {
	currentUser := ctx.UserValue("currentUser").(models.User)
	var data editLinkRequest
	var link models.Link
	id := ctx.UserValue("id")
	ctx.SetContentType("application/json")
	json.Unmarshal(ctx.Request.Body(), &data)

	err := env.Db.Where("id = ?", id).Take(&link).Error
	if err != nil {
		if err.Error() == "record not found" {
			ctx.Response.Header.SetStatusCode(404)
			fmt.Fprint(ctx, `{"error": "Link not found"}`)
			return
		}
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Edit link", "status": "Failed"}).Error(err.Error())
		return
	}

	if !currentUser.IsAdmin && link.UserId != currentUser.Id {
		ctx.Response.Header.SetStatusCode(401)
		ctx.SetContentType("application/json")
		fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
		return
	}

	err = env.Db.Model(&link).Update("Content", data.Content).Error
	if err != nil {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Edit link", "status": "Failed"}).Error(err.Error())
		return
	}

	json.NewEncoder(ctx).Encode(&link)
}

func (env Env) DeleteLink(ctx *fasthttp.RequestCtx) {
	var link models.Link
	currentUser := ctx.UserValue("currentUser").(models.User)
	id := ctx.UserValue("id")

	ctx.SetContentType("application/json")

	err := env.Db.Where("id = ?", id).Take(&link).Error
	if err != nil {
		if err.Error() == "record not found" {
			ctx.Response.Header.SetStatusCode(404)
			fmt.Fprint(ctx, `{"error": "Link not found"}`)
			return
		}
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Delete link", "status": "Failed"}).Error(err.Error())
		return
	}
	if !currentUser.IsAdmin && link.UserId != currentUser.Id {
		ctx.Response.Header.SetStatusCode(401)
		ctx.SetContentType("application/json")
		fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
		return
	}

	err = env.Db.Delete(&link).Error
	if err != nil {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Delete link", "status": "Failed"}).Error(err.Error())
		return
	}

	ctx.Response.Header.SetStatusCode(204)
}
