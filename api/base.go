package api

import (
	"context"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

type BaseAPI struct{}

type BaseController struct {
	web.Controller
}

type SessionData struct {
	ExpireTime int64
}

func (c *BaseController) IsGlobalAdmin() bool {
	isGlobalAdmin, _ := c.isGlobalAdmin()
	return isGlobalAdmin
}

func (c *BaseController) IsAdmin() bool {
	isGlobalAdmin, user := c.isGlobalAdmin()
	if !isGlobalAdmin && user == nil {
		return false
	}
	return isGlobalAdmin || user.IsAdmin
}

func (c *BaseController) IsAdminOrSelf(user2 *object.User) bool {
	isGlobalAdmin, user := c.isGlobalAdmin()
	if isGlobalAdmin || (user != nil && user.IsAdmin) {
		return true
	}

	if user == nil || user2 == nil {
		return false
	}

	if user.Owner == user2.Owner && user.Name == user2.Name {
		return true
	}
	return false
}

func (c *BaseController) isGlobalAdmin() (bool, *object.User) {
	username := c.GetSessionUsername()
	if object.IsAppUser(username) {
		return true, nil
	}

	user := c.getCurrentUser()
	if user == nil {
		return false, nil
	}

	return user.IsGlobalAdmin(), user
}

func (c *BaseController) getCurrentUser() *object.User {
	var user *object.User
	var err error
	userId := c.GetSessionUsername()
	if userId == "" {
		user = nil
	} else {
		user, err = object.GetUser(userId)
		if err != nil {
			c.ResponseError(err.Error())
			return nil
		}
	}
	return user
}

func (c *BaseController) GetSessionUsername() string {
	if ctxUser := c.Ctx.Input.GetData("currentUserId"); ctxUser != nil {
		if username, ok := ctxUser.(string); ok {
			return username
		}
	}

	sessionData := c.GetSessionData()

	if sessionData != nil &&
		sessionData.ExpireTime != 0 &&
		sessionData.ExpireTime < time.Now().Unix() {
		c.ClearUserSession()
		return ""
	}

	user := c.GetSession("username")
	if user == nil {
		return ""
	}

	return user.(string)
}

func (c *BaseController) GetSessionToken() string {
	accessToken := c.GetSession("accessToken")
	if accessToken == nil {
		return ""
	}
	return accessToken.(string)
}

func (c *BaseController) GetSessionApplication() *object.Application {
	clientId := c.GetSession("aud")
	if clientId == nil {
		return nil
	}
	application, err := object.GetApplicationByClientId(clientId.(string))
	if err != nil {
		c.ResponseError(err.Error())
		return nil
	}
	return application
}

func (c *BaseController) ClearUserSession() {
	c.SetSessionUsername("")
	c.SetSessionData(nil)
	_ = c.SessionRegenerateID()
}

func (c *BaseController) ClearTokenSession() {
	c.SetSessionToken("")
}

func (c *BaseController) SetSessionUsername(user string) {
	c.SetSession("username", user)
}

func (c *BaseController) SetSessionToken(accessToken string) {
	c.SetSession("accessToken", accessToken)
}

func (c *BaseController) GetSessionData() *SessionData {
	session := c.GetSession("SessionData")
	if session == nil {
		return nil
	}

	sessionData := &SessionData{}
	err := util.JsonToStruct(session.(string), sessionData)
	if err != nil {
		logs.Error("GetSessionData failed, error: %s", err)
		return nil
	}

	return sessionData
}

func (c *BaseController) SetSessionData(s *SessionData) {
	if s == nil {
		c.DelSession("SessionData")
		return
	}
	c.SetSession("SessionData", util.StructToJson(s))
}

func (c *BaseController) setMfaUserSession(userId string) {
	c.SetSession(object.MfaSessionUserId, userId)
}

func (c *BaseController) getMfaUserSession() string {
	userId := c.Ctx.Input.CruSession.Get(context.Background(), object.MfaSessionUserId)
	if userId == nil {
		return ""
	}
	return userId.(string)
}

func (c *BaseController) setExpireForSession(cookieExpireInHours int64) {
	timestamp := time.Now().Unix()
	if cookieExpireInHours == 0 {
		cookieExpireInHours = 720
	}
	timestamp += 3600 * cookieExpireInHours
	c.SetSessionData(&SessionData{
		ExpireTime: timestamp,
	})
}

func (c *BaseController) ResponseOk(data ...interface{}) {
	c.Data["json"] = c.buildResponse("ok", "", data...)
	c.ServeJSON()
}

func (c *BaseController) ResponseError(msg string) {
	c.Data["json"] = c.buildResponse("error", msg)
	c.ServeJSON()
}

func (c *BaseController) buildResponse(status, msg string, data ...interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"status": status,
		"msg":    msg,
	}
	if len(data) > 0 {
		resp["data"] = data[0]
	}
	if len(data) > 1 {
		resp["data2"] = data[1]
	}
	return resp
}

func (c *BaseController) T(key string) string {
	return key
}

func (c *BaseController) GetAcceptLanguage() string {
	return "en"
}

func (c *BaseController) Finish() {
	if strings.HasPrefix(c.Ctx.Input.URL(), "/api") {
		startTime := c.Ctx.Input.GetData("startTime")
		if startTime != nil {
			latency := time.Since(startTime.(time.Time)).Milliseconds()
			object.ApiLatency.WithLabelValues(c.Ctx.Input.URL(), c.Ctx.Input.Method()).Observe(float64(latency))
		}
	}
	c.Controller.Finish()
}
