// Copyright 2024 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/beego/beego/utils/pagination"
	"github.com/casdoor/casdoor/controllers"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/service"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

// UserHandler 用户API处理器
type UserHandler struct {
	controllers.BaseController
	userService service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// GetUsers 获取用户列表
// @Title GetUsers
// @Description 获取用户列表（分页）
// @Param	owner	query	string	true	"组织名称"
// @Param	page	query	int	false	"页码"
// @Param	pageSize	query	int	false	"每页大小"
// @Param	field	query	string	false	"搜索字段"
// @Param	value	query	string	false	"搜索值"
// @Param	sortField	query	string	false	"排序字段"
// @Param	sortOrder	query	string	false	"排序方向"
// @Success 200 {object} common.ApiResponse "成功返回用户列表"
// @router /get-users [get]
func (h *UserHandler) GetUsers() {
	owner := h.Input().Get("owner")
	page := h.GetIntFromInput("page", common.DefaultPage)
	pageSize := h.GetIntFromInput("pageSize", common.DefaultPageSize)
	field := h.Input().Get("field")
	value := h.Input().Get("value")
	sortField := h.Input().Get("sortField")
	sortOrder := h.Input().Get("sortOrder")

	if owner == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	users, total, err := h.userService.ListUsers(owner, page, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccessWithPage(users, total, page, pageSize)
}

// GetGlobalUsers 获取全局用户列表
// @Title GetGlobalUsers
// @Description 获取全局用户列表（跨组织）
// @Param	page	query	int	false	"页码"
// @Param	pageSize	query	int	false	"每页大小"
// @Param	field	query	string	false	"搜索字段"
// @Param	value	query	string	false	"搜索值"
// @Param	sortField	query	string	false	"排序字段"
// @Param	sortOrder	query	string	false	"排序方向"
// @Success 200 {object} common.ApiResponse "成功返回全局用户列表"
// @router /get-global-users [get]
func (h *UserHandler) GetGlobalUsers() {
	page := h.GetIntFromInput("page", common.DefaultPage)
	pageSize := h.GetIntFromInput("pageSize", common.DefaultPageSize)
	field := h.Input().Get("field")
	value := h.Input().Get("value")
	sortField := h.Input().Get("sortField")
	sortOrder := h.Input().Get("sortOrder")

	users, total, err := h.userService.ListGlobalUsers(page, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccessWithPage(users, total, page, pageSize)
}

// GetUser 获取单个用户
// @Title GetUser
// @Description 获取单个用户详情
// @Param	id	query	string	true	"用户ID，格式：owner/name"
// @Success 200 {object} common.ApiResponse "成功返回用户详情"
// @router /get-user [get]
func (h *UserHandler) GetUser() {
	id := h.Input().Get("id")
	if id == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	owner, name := util.GetOwnerAndNameFromId(id)
	user, err := h.userService.GetUser(owner, name)
	if err != nil {
		if err == common.ErrUserNotFound {
			h.ResponseError("User not found")
		} else {
			h.ResponseError(err.Error())
		}
		return
	}

	h.ResponseSuccess(user)
}

// GetUserByUserId 根据用户ID获取用户
// @Title GetUserByUserId
// @Description 根据用户ID获取用户
// @Param	userId	query	string	true	"用户ID"
// @Success 200 {object} common.ApiResponse "成功返回用户详情"
// @router /get-user-by-userId [get]
func (h *UserHandler) GetUserByUserId() {
	userId := h.Input().Get("userId")
	if userId == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 从userId中解析owner
	// 注意：这里简化处理，实际应该根据userId的格式来解析
	owner := h.GetSessionUser().Owner

	user, err := h.userService.GetUserByUserID(owner, userId)
	if err != nil {
		if err == common.ErrUserNotFound {
			h.ResponseError("User not found")
		} else {
			h.ResponseError(err.Error())
		}
		return
	}

	h.ResponseSuccess(user)
}

// UpdateUser 更新用户
// @Title UpdateUser
// @Description 更新用户信息
// @Param	body	body	object.User	true	"用户信息"
// @Success 200 {object} common.ApiResponse "成功返回更新结果"
// @router /update-user [post]
func (h *UserHandler) UpdateUser() {
	var user object.User
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &user); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if user.Owner == "" || user.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 密码处理：如果密码不为空则更新，否则不更新密码字段
	var columns []string
	if user.Password != "" {
		columns = nil // 更新所有字段
	} else {
		// 排除密码字段
		columns = []string{"display_name", "first_name", "last_name", "avatar", "email", "phone", "location", "address", "affiliation", "title", "homepage", "bio", "tag", "region", "language", "gender", "birthday", "education", "score", "ranking", "is_online", "is_admin", "is_global_admin", "is_forbidden", "access_key", "access_secret", "type", "properties"}
	}

	affected, err := h.userService.UpdateUser(&user, columns...)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// AddUser 添加用户
// @Title AddUser
// @Description 添加新用户
// @Param	body	body	object.User	true	"用户信息"
// @Success 200 {object} common.ApiResponse "成功返回添加结果"
// @router /add-user [post]
func (h *UserHandler) AddUser() {
	var user object.User
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &user); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if user.Owner == "" || user.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.userService.CreateUser(&user)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// DeleteUser 删除用户
// @Title DeleteUser
// @Description 删除用户
// @Param	body	body	object.User	true	"用户信息（只需要owner和name）"
// @Success 200 {object} common.ApiResponse "成功返回删除结果"
// @router /delete-user [post]
func (h *UserHandler) DeleteUser() {
	var user object.User
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &user); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if user.Owner == "" || user.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.userService.DeleteUser(&user)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// UpdatePassword 更新用户密码
// @Title UpdatePassword
// @Description 更新用户密码
// @Param	userOwner	formData	string	true	"用户所属组织"
// @Param	userName	formData	string	true	"用户名"
// @Param	oldPassword	formData	string	true	"旧密码"
// @Param	newPassword	formData	string	true	"新密码"
// @Success 200 {object} common.ApiResponse "成功返回更新结果"
// @router /update-password [post]
func (h *UserHandler) UpdatePassword() {
	userOwner := h.Input().Get("userOwner")
	userName := h.Input().Get("userName")
	oldPassword := h.Input().Get("oldPassword")
	newPassword := h.Input().Get("newPassword")

	if userOwner == "" || userName == "" || oldPassword == "" || newPassword == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	err := h.userService.UpdateUserPassword(userOwner, userName, oldPassword, newPassword, h.GetAcceptLanguage())
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess("Password updated successfully")
}

// ResetPassword 重置用户密码
// @Title ResetPassword
// @Description 重置用户密码（管理员操作）
// @Param	userOwner	formData	string	true	"用户所属组织"
// @Param	userName	formData	string	true	"用户名"
// @Param	newPassword	formData	string	true	"新密码"
// @Success 200 {object} common.ApiResponse "成功返回重置结果"
// @router /reset-password [post]
func (h *UserHandler) ResetPassword() {
	userOwner := h.Input().Get("userOwner")
	userName := h.Input().Get("userName")
	newPassword := h.Input().Get("newPassword")

	if userOwner == "" || userName == "" || newPassword == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	err := h.userService.ResetUserPassword(userOwner, userName, newPassword)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess("Password reset successfully")
}

// EnableUser 启用用户
// @Title EnableUser
// @Description 启用用户
// @Param	owner	query	string	true	"组织名称"
// @Param	name	query	string	true	"用户名"
// @Success 200 {object} common.ApiResponse "成功返回启用结果"
// @router /enable-user [post]
func (h *UserHandler) EnableUser() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")

	if owner == "" || name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.userService.EnableUser(owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// DisableUser 禁用用户
// @Title DisableUser
// @Description 禁用用户
// @Param	owner	query	string	true	"组织名称"
// @Param	name	query	string	true	"用户名"
// @Success 200 {object} common.ApiResponse "成功返回禁用结果"
// @router /disable-user [post]
func (h *UserHandler) DisableUser() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")

	if owner == "" || name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.userService.DisableUser(owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// GetSortedUsers 获取排序后的用户列表
// @Title GetSortedUsers
// @Description 获取排序后的用户列表（用于排行榜等场景）
// @Param	owner	query	string	true	"组织名称"
// @Param	sorter	query	string	true	"排序字段（如：score、signin_count）"
// @Param	limit	query	int	false	"返回数量限制"
// @Success 200 {object} common.ApiResponse "成功返回用户列表"
// @router /get-sorted-users [get]
func (h *UserHandler) GetSortedUsers() {
	owner := h.Input().Get("owner")
	sorter := h.Input().Get("sorter")
	limit := h.GetIntFromInput("limit", 10)

	if owner == "" || sorter == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	users, err := h.userService.GetSortedUsers(owner, sorter, limit)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(users)
}

// GetOnlineUserCount 获取在线用户数量
// @Title GetOnlineUserCount
// @Description 获取在线用户数量
// @Param	owner	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回在线用户数量"
// @router /get-online-user-count [get]
func (h *UserHandler) GetOnlineUserCount() {
	owner := h.Input().Get("owner")

	if owner == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	count, err := h.userService.GetOnlineUserCount(owner)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{"onlineCount": count})
}

// 辅助方法

// GetIntFromInput 从输入中获取整数值
func (h *UserHandler) GetIntFromInput(key string, defaultValue int) int {
	value := h.Input().Get(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

// ResponseSuccess 返回成功响应
func (h *UserHandler) ResponseSuccess(data ...interface{}) {
	h.Data["json"] = common.ResponseSuccess(data...)
	h.ServeJSON()
}

// ResponseSuccessWithPage 返回带分页的成功响应
func (h *UserHandler) ResponseSuccessWithPage(list interface{}, total int64, page, size int) {
	h.Data["json"] = common.ResponseSuccessWithPage(list, total, page, size)
	h.ServeJSON()
}

// ResponseError 返回错误响应
func (h *UserHandler) ResponseError(message string, data ...interface{}) {
	h.Data["json"] = common.ResponseError(i18n.Translate(h.GetAcceptLanguage(), message), data...)
	h.ServeJSON()
}

// ResponseFromAction 根据操作结果返回响应
func (h *UserHandler) ResponseFromAction(affected bool, err ...error) {
	h.Data["json"] = common.ResponseFromAction(affected, err...)
	h.ServeJSON()
}

// GetAcceptLanguage 获取接受语言
func (h *UserHandler) GetAcceptLanguage() string {
	return h.Ctx.Input.Header("Accept-Language")
}

// GetSessionUser 获取当前会话用户
func (h *UserHandler) GetSessionUser() *object.User {
	user, ok := h.Ctx.Input.GetData("user").(*object.User)
	if !ok {
		return nil
	}
	return user
}
