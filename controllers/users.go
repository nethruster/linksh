package controllers

import (
	"encoding/json"
	"github.com/nethruster/linksh/models"
	"fmt"
	"github.com/erikdubbelboer/fasthttp"
	"strings"
	"github.com/sirupsen/logrus"
	"strconv"
)

func (env Env) GetUsers(ctx *fasthttp.RequestCtx) {
    if !ctx.UserValue("currentUser").(models.User).IsAdmin {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

	var users []models.User
	args := ctx.QueryArgs()
	query := env.Db

	if email := string(args.Peek("email")); email != "" {
		query = query.Where("email like ?", fmt.Sprintf("%%%v%%", email))
	}
	if offset, err := strconv.Atoi(string(args.Peek("offset"))); err == nil && offset != 0 {
		query = query.Offset(offset)
	}
	if limit, err := strconv.Atoi(string(args.Peek("limit"))); err == nil && limit != 0 {
		query = query.Limit(limit)
	}

    err := query.Find(&users).Error
    if err != nil {
        ctx.Response.Header.SetStatusCode(500)
        fmt.Fprint(ctx, `{"error": "Internal server error"}`)
        env.Log.WithFields(logrus.Fields{"event": "Lists sessions", "status": "Failed"}).Error(err.Error())
        return
    }

    ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(&users)
}

func (env Env) GetUser(ctx *fasthttp.RequestCtx) {
    currentUser := ctx.UserValue("currentUser").(models.User)
	var user models.User
	args := ctx.QueryArgs()
	id := ctx.UserValue("id")
	ctx.SetContentType("application/json")
    if currentUser.Id == id {
        user = currentUser
    } else if currentUser.IsAdmin {
        err := env.Db.Where("id = ?", id).Take(&user).Error
        if err != nil && err.Error() == "record not found" {
            ctx.Response.Header.SetStatusCode(500)
            fmt.Fprint(ctx, `{"error": "Internal server error"}`)
            env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
            return
        }

        if user.Id == "" {
            ctx.Response.Header.SetStatusCode(404)
            fmt.Fprint(ctx, `{"error": "User not found"}`)
            return
        }
    } else {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

	if string(args.Peek("includeSessions")) == "true" {
		var sessions []models.Session
		query := env.Db

		if offset, err := strconv.Atoi(string(args.Peek("sessionsOffset"))); err == nil && offset != 0 {
			query = query.Offset(offset)
		}
		if limit, err := strconv.Atoi(string(args.Peek("sessionsLimit"))); err == nil && limit != 0 {
			query = query.Limit(limit)
		}

        err := query.Model(&user).Related(&sessions).Error
        if err != nil {
            ctx.Response.Header.SetStatusCode(500)
            fmt.Fprint(ctx, `{"error": "Internal server error"}`)
            env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
            return
        }

		user.Sessions = sessions
	}

	if string(args.Peek("includeLinks")) == "true" {
		var links []models.Link
		query := env.Db

		if offset, err := strconv.Atoi(string(args.Peek("linksOffset"))); err == nil && offset != 0 {
			query = query.Offset(offset)
		}
		if limit, err := strconv.Atoi(string(args.Peek("linksLimit"))); err == nil && limit != 0 {
			query = query.Limit(limit)
		}

        err := query.Model(&user).Related(&links).Error
        if err != nil {
            ctx.Response.Header.SetStatusCode(500)
            fmt.Fprint(ctx, `{"error": "Internal server error"}`)
            env.Log.WithFields(logrus.Fields{"event": "GetUser", "status": "Failed"}).Error(err.Error())
            return
        }

		user.Links = links
	}

	json.NewEncoder(ctx).Encode(&user)
}

func (env Env) CreateUser(ctx *fasthttp.RequestCtx) {
    if !ctx.UserValue("currentUser").(models.User).IsAdmin {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

	var data map[string] string
	ctx.SetContentType("application/json")

	json.Unmarshal(ctx.Request.Body(), &data)


	user := models.User{
		Username: data["username"],
		Email: data["email"],
		Password: []byte(data["password"]),
	}

	errs := user.ValidateUser()

	if errs != nil {
		ctx.Response.Header.SetStatusCode(400)

		fmt.Fprint(ctx, `{"error": [`)
		for i,err := range errs {
			fmt.Fprintf(ctx, `"%v"`, err.Error())
			if i != len(errs) -1 {
				fmt.Fprint(ctx, ",")
			}
		}
		fmt.Fprint(ctx, "]}")
		return
	}

	err := user.SaveToDatabase(env.Db)


	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			ctx.Response.Header.SetStatusCode(400)
			fmt.Fprint(ctx, `{"error": "User already exists"}`)
			return
		}
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Create user", "status": "Failed"}).Error(err.Error())
		return
	}

	ctx.Response.Header.SetStatusCode(201)
	json.NewEncoder(ctx).Encode(&user)
	env.Log.WithFields(logrus.Fields{"event": "Create user", "status": "successful"}).Info(fmt.Sprintf(`A user was created with Id = '%v' and Email = '%v'`, user.Id, user.Email))
}

func (env Env) EditUser(ctx *fasthttp.RequestCtx) {
    currentUser := ctx.UserValue("currentUser").(models.User)
    var data map[string]string
	var user models.User
	changes := make(map[string] interface{})
	id := ctx.UserValue("id")
	ctx.SetContentType("application/json")

	json.Unmarshal(ctx.Request.Body(), &data)

	if username := data["username"]; username != "" {
		if err := models.ValidateUsername(username); err != nil {
			ctx.Response.Header.SetStatusCode(400)
			fmt.Fprintf(ctx, `{"error": "%v"}`, err)
			return
		}
		changes["Username"] = username
	}
	if email := data["email"]; email != "" {
		if err := models.ValidateEmail(email); err != nil {
			ctx.Response.Header.SetStatusCode(400)
			fmt.Fprintf(ctx, `{"error": "%v"}`, err)
			return
		}
		changes["Email"] = email
	}
	if password := data["password"]; password != "" {
		passwordBytes := []byte(password)
		if err := models.ValidatePassword(passwordBytes); err != nil {
			ctx.Response.Header.SetStatusCode(400)
			fmt.Fprintf(ctx, `{"error": "%v"}`, err)
			return
		}
		passwordBytes, err := models.HashPassword(passwordBytes)
		if err != nil {
			ctx.Response.Header.SetStatusCode(500)
			fmt.Fprint(ctx, `{"error": "Internal server error"}`)
			env.Log.WithFields(logrus.Fields{"event": "Edit user", "status": "Failed"}).Error(err.Error())
			return
		}


		changes["Password"] = passwordBytes
	}
	if data["apikey"] == "true" {
		apikey, err := models.GenerateUserApiKey()
		if err != nil {
			ctx.Response.Header.SetStatusCode(500)
			fmt.Fprint(ctx, `{"error": "Internal server error"}`)
			env.Log.WithFields(logrus.Fields{"event": "Edit user", "status": "Failed"}).Error(err.Error())
			return
		}

		changes["Apikey"] = apikey
	}
    if currentUser.Id == id {
        user = currentUser
        if changes["Password"] != nil && !user.CheckIfCorrectPassword([]byte(data["originalPassword"])) {
            ctx.Response.Header.SetStatusCode(400)
            fmt.Fprint(ctx, `{"error": "Missing or invalid originalPassword"}`)
            return
        }

        err := env.Db.Model(&user).Updates(changes).Error
        if err != nil {
            ctx.Response.Header.SetStatusCode(500)
            fmt.Fprint(ctx, `{"error": "Internal server error"}`)
            env.Log.WithFields(logrus.Fields{"event": "Edit user", "status": "Failed"}).Error(err.Error())
            return
        }
    } else if currentUser.IsAdmin {
        env.Db.Where("id = ?", id).Take(&user)
        if user.Id == "" {
            ctx.Response.Header.SetStatusCode(404)
            fmt.Fprint(ctx, `{"error": "User not found"}`)
            return
        }

        err := env.Db.Model(&user).Updates(changes).Error
        if err != nil {
            ctx.Response.Header.SetStatusCode(500)
            fmt.Fprint(ctx, `{"error": "Internal server error"}`)
            env.Log.WithFields(logrus.Fields{"event": "Edit user", "status": "Failed"}).Error(err.Error())
            return
        }
    } else {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

	json.NewEncoder(ctx).Encode(&user)
}

func (env Env) DeleteUser(ctx *fasthttp.RequestCtx) {
    currentUser := ctx.UserValue("currentUser").(models.User)
	id := ctx.UserValue("id")

    if currentUser.Id != id && !currentUser.IsAdmin {
        ctx.Response.Header.SetStatusCode(401)
        ctx.SetContentType("application/json")
        fmt.Fprint(ctx, `{"error": "UNAUTHORIZED"}`)
        return
    }

	result := env.Db.Delete(models.User{}, "id = ?", id)
	if err := result.Error; err != nil {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Delete user", "status": "Failed"}).Error(err.Error())
		return
	}
	if result.RowsAffected == 0 {
		ctx.Response.Header.SetStatusCode(404)
		fmt.Fprint(ctx, `{"error": "User not found"}`)
		return
	}
	ctx.Response.Header.SetStatusCode(204)
}
