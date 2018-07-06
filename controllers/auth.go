package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/jinzhu/gorm"
	"github.com/nethruster/linksh/models"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (env Env) Auth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		if cookie := ctx.Request.Header.Cookie("linksh-auth"); cookie != nil {
			data := string(cookie)
			cookie = nil
			valid, user, err := checkLoginWithSession(env.Db, data[:36], data[36:])

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
				ctx.Response.Header.SetCookie(createSessionCookie("", time.Unix(0, 0)))

				ctx.Response.Header.SetStatusCode(403)
				ctx.SetContentType("application/json")
				fmt.Fprint(ctx, `{"error": "FORBIDDEN"}`)
			}

		} else if auth := ctx.Request.Header.Peek("X-Auth"); auth != nil {
			data := string(auth)
			valid, user, err := checkLoginWithSession(env.Db, data[36:], data[:36])

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

type loginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	NoExpire  bool   `json:"noExpire"`
	UseCookie bool   `json:"useCookie"`
}

//Register

func (env Env) Login(ctx *fasthttp.RequestCtx) {
	var data loginRequest
	var user models.User
	ctx.SetContentType("application/json")
	json.Unmarshal(ctx.Request.Body(), &data)
	if data.Email == "" || data.Password == "" {
		ctx.Response.Header.SetStatusCode(400)
		fmt.Fprint(ctx, `{"error": "Missing email or password"}`)
	}

	err := env.Db.Where("email = ?", data.Email).Take(&user).Error
	if err != nil && err.Error() != "record not found" {
		ctx.Response.Header.SetStatusCode(500)
		fmt.Fprint(ctx, `{"error": "Internal server error"}`)
		env.Log.WithFields(logrus.Fields{"event": "Login", "status": "Failed"}).Error(err.Error())
		return
	}

	if user.Id == "" || !user.CheckIfCorrectPassword([]byte(data.Password)) {
		ctx.Response.Header.SetStatusCode(400)
		fmt.Fprint(ctx, `{"error": "The email or the password are invalid"}`)
		return
	}

	var expires time.Time

	if data.NoExpire {
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

	if data.UseCookie {
		ctx.Response.Header.SetCookie(createSessionCookie(
			fmt.Sprintf("%v%v", session.Id, user.Id),
			session.ExpiresAt.AddDate(0, 0, 2),
		))
	}

	fmt.Fprintf(ctx, `{"sessionId": "%v", "userId": "%v", "expiresAt": "%v"}`, session.Id, user.Id, session.ExpiresAt)
}

type registerRequest struct {
	loginRequest
	Username string `json:"username"`
}

func (env Env) Register(ctx *fasthttp.RequestCtx) {
	if !env.Config.Server.AllowRegister {
		ctx.Response.Header.SetStatusCode(403)
		ctx.SetContentType("application/json")
		fmt.Fprint(ctx, `{"error": "FORBIDDEN"}`)
		return
	}
	var data registerRequest
	ctx.SetContentType("application/json")

	json.Unmarshal(ctx.Request.Body(), &data)

	user := models.User{
		Username: data.Username,
		Email:    data.Email,
		Password: []byte(data.Password),
	}

	errs := user.ValidateUser()

	if errs != nil {
		ctx.Response.Header.SetStatusCode(400)

		fmt.Fprint(ctx, `{"error": [`)
		for i, err := range errs {
			fmt.Fprintf(ctx, `"%v"`, err.Error())
			if i != len(errs)-1 {
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

	var expires time.Time

	if data.NoExpire {
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

	if data.UseCookie {
		ctx.Response.Header.SetCookie(createSessionCookie(
			fmt.Sprintf("%v%v", session.Id, user.Id),
			session.ExpiresAt.AddDate(0, 0, 2),
		))
	}

	ctx.Response.Header.SetStatusCode(201)
	fmt.Fprintf(ctx, `{"sessionId": "%v", "userId": "%v", "expiresAt": "%v"}`, session.Id, user.Id, session.ExpiresAt)
	env.Log.WithFields(logrus.Fields{"event": "Register user", "status": "successful"}).Info(fmt.Sprintf(`A user was created with Id = '%v' and Email = '%v'`, user.Id, user.Email))
}

func (env Env) Logout(ctx *fasthttp.RequestCtx) {
	if cookie := ctx.Request.Header.Cookie("auth"); cookie != nil {
		data := string(cookie)
		err := env.Db.Delete(models.Session{}, "id = ?", data[:36]).Error
		if err != nil {
			ctx.Response.Header.SetStatusCode(500)
			fmt.Fprint(ctx, `{"error": "Internal server error"}`)
			env.Log.WithFields(logrus.Fields{"event": "Logout", "status": "Failed"}).Error(err.Error())
			return
		}

		ctx.Response.Header.SetCookie(createSessionCookie("", time.Unix(0, 0)))

		ctx.Response.Header.SetStatusCode(204)
	} else if auth := ctx.Request.Header.Peek("auth"); auth != nil {
		data := string(auth)

		err := env.Db.Delete(models.Session{}, "id = ?", data[:36]).Error
		if err != nil {
			ctx.Response.Header.SetStatusCode(500)
			fmt.Fprint(ctx, `{"error": "Internal server error"}`)
			env.Log.WithFields(logrus.Fields{"event": "Logout", "status": "Failed"}).Error(err.Error())
			return
		}
		ctx.Response.Header.SetStatusCode(204)
	} else {
		ctx.Response.Header.SetStatusCode(400)
		fmt.Fprint(ctx, `{"error": "Bad request"}`)
	}
}

// Helpers
func createSessionCookie(value string, expires time.Time) *fasthttp.Cookie {
	var cookie fasthttp.Cookie
	cookie.SetKey("linksh-auth")
	cookie.SetValue(value)
	cookie.SetHTTPOnly(true)
	cookie.SetExpire(expires)

	return &cookie
}
