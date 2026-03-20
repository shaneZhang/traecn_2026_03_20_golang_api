package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/utils/pagination"
	"github.com/casdoor/casdoor/dto"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/service"
	"github.com/casdoor/casdoor/util"
)

type UserAPI struct {
	BaseAPI
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

func (api *UserAPI) GetGlobalUsers(c *BaseController) {
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")

	if limit == "" || page == "" {
		users, _, err := api.userService.GetGlobalUsers(-1, -1, field, value, sortField, sortOrder)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		maskedUsers, err := api.userService.GetMaskedUsers(users)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(maskedUsers)
	} else {
		limitInt := util.ParseInt(limit)
		count, err := api.userService.GetUserCount("", field, value, "")
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
		users, _, err := api.userService.GetGlobalUsers(paginator.Offset(), limitInt, field, value, sortField, sortOrder)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		maskedUsers, err := api.userService.GetMaskedUsers(users)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(maskedUsers, paginator.Nums())
	}
}

func (api *UserAPI) GetUsers(c *BaseController) {
	owner := c.Ctx.Input.Query("owner")
	groupName := c.Ctx.Input.Query("groupName")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")

	if limit == "" || page == "" {
		if groupName != "" {
			users, _, err := api.userService.GetUsers(owner, -1, -1, field, value, sortField, sortOrder, groupName)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
			maskedUsers, err := api.userService.GetMaskedUsers(users)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
			c.ResponseOk(maskedUsers)
			return
		}

		users, _, err := api.userService.GetUsers(owner, -1, -1, field, value, sortField, sortOrder, "")
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		maskedUsers, err := api.userService.GetMaskedUsers(users)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(maskedUsers)
	} else {
		limitInt := util.ParseInt(limit)
		count, err := api.userService.GetUserCount(owner, field, value, groupName)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
		users, _, err := api.userService.GetUsers(owner, paginator.Offset(), limitInt, field, value, sortField, sortOrder, groupName)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		maskedUsers, err := api.userService.GetMaskedUsers(users)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(maskedUsers, paginator.Nums())
	}
}

func (api *UserAPI) GetUser(c *BaseController) {
	id := c.Ctx.Input.Query("id")
	email := c.Ctx.Input.Query("email")
	phone := c.Ctx.Input.Query("phone")
	userId := c.Ctx.Input.Query("userId")
	owner := c.Ctx.Input.Query("owner")

	var user *object.User
	var err error

	if userId != "" && owner != "" {
		user, err = api.userService.GetUserByUserId(owner, userId)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		if user == nil {
			c.ResponseOk(nil)
			return
		}
		id = util.GetId(user.Owner, user.Name)
	}

	if id == "" && owner == "" {
		switch {
		case email != "":
			user, err = api.userService.GetUserByEmail("", email)
		case phone != "":
			user, err = api.userService.GetUserByPhone("", phone)
		case userId != "":
			user, err = api.userService.GetUserByUserId("", userId)
		}
	} else {
		if owner == "" {
			owner = util.GetOwnerFromId(id)
		}

		switch {
		case email != "":
			user, err = api.userService.GetUserByEmail(owner, email)
		case phone != "":
			user, err = api.userService.GetUserByPhone(owner, phone)
		case userId != "":
		default:
			user, err = api.userService.GetUser(id)
		}
	}

	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if user != nil {
		user.MultiFactorAuths = object.GetAllMfaProps(user, true)
	}

	err = api.userService.ExtendUserWithRolesAndPermissions(user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	isAdminOrSelf := c.IsAdminOrSelf(user)
	user, err = api.userService.GetMaskedUser(user, isAdminOrSelf)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(user)
}

func (api *UserAPI) UpdateUser(c *BaseController) {
	id := c.Ctx.Input.Query("id")
	userId := c.Ctx.Input.Query("userId")
	owner := c.Ctx.Input.Query("owner")
	columnsStr := c.Ctx.Input.Query("columns")

	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if id == "" && userId == "" {
		id = c.GetSessionUsername()
		if id == "" {
			c.ResponseError(c.T("general:Missing parameter"))
			return
		}
	}

	if userId != "" && owner != "" {
		userFromUserId, err := api.userService.GetUserByUserId(owner, userId)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		if userFromUserId == nil {
			c.ResponseError(fmt.Sprintf(c.T("general:The user: %s doesn't exist"), userId))
			return
		}

		id = util.GetId(userFromUserId.Owner, userFromUserId.Name)
	}

	oldUser, err := api.userService.GetUser(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if oldUser == nil {
		c.ResponseError(fmt.Sprintf(c.T("general:The user: %s doesn't exist"), id))
		return
	}

	if oldUser.Owner == "built-in" && oldUser.Name == "admin" && (user.Owner != "built-in" || user.Name != "admin") {
		c.ResponseError(c.T("auth:Unauthorized operation"))
		return
	}

	columns := []string{}
	if columnsStr != "" {
		columns = strings.Split(columnsStr, ",")
	}

	isAdmin := c.IsAdmin()
	affected, err := api.userService.UpdateUser(id, &user, columns, isAdmin, c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *UserAPI) AddUser(c *BaseController) {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if user.RegisterType == "" {
		user.RegisterType = "Add User"
	}
	if user.RegisterSource == "" {
		currentUser := c.getCurrentUser()
		if currentUser != nil {
			user.RegisterSource = currentUser.GetId()
		}
	}

	affected, err := api.userService.CreateUser(&user, c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *UserAPI) DeleteUser(c *BaseController) {
	var user object.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.userService.DeleteUser(&user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *UserAPI) GetSortedUsers(c *BaseController) {
	owner := c.Ctx.Input.Query("owner")
	sorter := c.Ctx.Input.Query("sorter")
	limit := util.ParseInt(c.Ctx.Input.Query("limit"))

	users, err := api.userService.GetSortedUsers(owner, sorter, limit)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	maskedUsers, err := api.userService.GetMaskedUsers(users)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(maskedUsers)
}

func (api *UserAPI) GetUserCount(c *BaseController) {
	owner := c.Ctx.Input.Query("owner")
	isOnline := c.Ctx.Input.Query("isOnline")

	var count int64
	var err error
	if isOnline == "" {
		count, err = api.userService.GetUserCount(owner, "", "", "")
	} else {
		count, err = api.userService.GetOnlineUserCount(owner, util.ParseInt(isOnline))
	}
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(count)
}

func (api *UserAPI) SetPassword(c *BaseController) {
	userOwner := c.Ctx.Request.Form.Get("userOwner")
	userName := c.Ctx.Request.Form.Get("userName")
	oldPassword := c.Ctx.Request.Form.Get("oldPassword")
	newPassword := c.Ctx.Request.Form.Get("newPassword")
	code := c.Ctx.Request.Form.Get("code")

	err := api.userService.SetPassword(userOwner, userName, oldPassword, newPassword, code, c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk()
}

func wrapActionResponse(affected bool, e ...error) *dto.Response {
	if len(e) != 0 && e[0] != nil {
		return &dto.Response{Status: "error", Msg: e[0].Error()}
	} else if affected {
		return &dto.Response{Status: "ok", Msg: "", Data: "Affected"}
	} else {
		return &dto.Response{Status: "ok", Msg: "", Data: "Unaffected"}
	}
}
