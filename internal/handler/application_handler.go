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

	"github.com/casdoor/casdoor/controllers"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/service"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

// ApplicationHandler 应用API处理器
type ApplicationHandler struct {
	controllers.BaseController
	appService service.ApplicationService
}

// NewApplicationHandler 创建应用处理器实例
func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		appService: service.NewApplicationService(),
	}
}

// GetApplications 获取应用列表
// @Title GetApplications
// @Description 获取应用列表（分页）
// @Param	owner	query	string	true	"组织名称"
// @Param	page	query	int	false	"页码"
// @Param	pageSize	query	int	false	"每页大小"
// @Param	field	query	string	false	"搜索字段"
// @Param	value	query	string	false	"搜索值"
// @Param	sortField	query	string	false	"排序字段"
// @Param	sortOrder	query	string	false	"排序方向"
// @Success 200 {object} common.ApiResponse "成功返回应用列表"
// @router /get-applications [get]
func (h *ApplicationHandler) GetApplications() {
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

	apps, total, err := h.appService.ListApplications(owner, page, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccessWithPage(apps, total, page, pageSize)
}

// GetApplication 获取单个应用
// @Title GetApplication
// @Description 获取单个应用详情
// @Param	id	query	string	true	"应用ID，格式：owner/name"
// @Success 200 {object} common.ApiResponse "成功返回应用详情"
// @router /get-application [get]
func (h *ApplicationHandler) GetApplication() {
	id := h.Input().Get("id")
	if id == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	owner, name := util.GetOwnerAndNameFromId(id)
	app, err := h.appService.GetApplication(owner, name)
	if err != nil {
		if err == common.ErrAppNotFound {
			h.ResponseError("Application not found")
		} else {
			h.ResponseError(err.Error())
		}
		return
	}

	h.ResponseSuccess(app)
}

// GetApplicationByClientId 根据ClientID获取应用
// @Title GetApplicationByClientId
// @Description 根据ClientID获取应用
// @Param	clientId	query	string	true	"客户端ID"
// @Success 200 {object} common.ApiResponse "成功返回应用详情"
// @router /get-application-by-clientid [get]
func (h *ApplicationHandler) GetApplicationByClientId() {
	clientID := h.Input().Get("clientId")
	if clientID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	app, err := h.appService.GetApplicationByClientID(clientID)
	if err != nil {
		if err == common.ErrAppNotFound {
			h.ResponseError("Application not found")
		} else {
			h.ResponseError(err.Error())
		}
		return
	}

	h.ResponseSuccess(app)
}

// UpdateApplication 更新应用
// @Title UpdateApplication
// @Description 更新应用信息
// @Param	body	body	object.Application	true	"应用信息"
// @Success 200 {object} common.ApiResponse "成功返回更新结果"
// @router /update-application [post]
func (h *ApplicationHandler) UpdateApplication() {
	var app object.Application
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &app); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if app.Owner == "" || app.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.appService.UpdateApplication(&app)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// AddApplication 添加应用
// @Title AddApplication
// @Description 添加新应用
// @Param	body	body	object.Application	true	"应用信息"
// @Success 200 {object} common.ApiResponse "成功返回添加结果"
// @router /add-application [post]
func (h *ApplicationHandler) AddApplication() {
	var app object.Application
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &app); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if app.Owner == "" || app.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.appService.CreateApplication(&app)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// DeleteApplication 删除应用
// @Title DeleteApplication
// @Description 删除应用
// @Param	body	body	object.Application	true	"应用信息（只需要owner和name）"
// @Success 200 {object} common.ApiResponse "成功返回删除结果"
// @router /delete-application [post]
func (h *ApplicationHandler) DeleteApplication() {
	var app object.Application
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &app); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if app.Owner == "" || app.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.appService.DeleteApplication(&app)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// GetAllApplications 获取所有应用
// @Title GetAllApplications
// @Description 获取组织下的所有应用列表
// @Param	owner	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回所有应用列表"
// @router /get-all-applications [get]
func (h *ApplicationHandler) GetAllApplications() {
	owner := h.Input().Get("owner")
	if owner == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	apps, err := h.appService.GetAllApplications(owner)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(apps)
}

// GetApplicationsByOrganization 获取组织下的应用
// @Title GetApplicationsByOrganization
// @Description 获取组织下的所有应用列表
// @Param	orgName	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回应用列表"
// @router /get-applications-by-organization [get]
func (h *ApplicationHandler) GetApplicationsByOrganization() {
	orgName := h.Input().Get("orgName")
	if orgName == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	apps, err := h.appService.GetApplicationsByOrganization(orgName)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(apps)
}

// GetApplicationPermissions 获取应用权限列表
// @Title GetApplicationPermissions
// @Description 获取应用的权限配置列表
// @Param	appID	query	string	true	"应用ID"
// @Success 200 {object} common.ApiResponse "成功返回权限列表"
// @router /get-application-permissions [get]
func (h *ApplicationHandler) GetApplicationPermissions() {
	appID := h.Input().Get("appID")
	if appID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	permissions, err := h.appService.GetApplicationPermissions(appID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(permissions)
}

// GetApplicationRoles 获取应用角色列表
// @Title GetApplicationRoles
// @Description 获取应用的角色配置列表
// @Param	appID	query	string	true	"应用ID"
// @Success 200 {object} common.ApiResponse "成功返回角色列表"
// @router /get-application-roles [get]
func (h *ApplicationHandler) GetApplicationRoles() {
	appID := h.Input().Get("appID")
	if appID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	roles, err := h.appService.GetApplicationRoles(appID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(roles)
}

// GetApplicationUserCount 获取应用用户数量
// @Title GetApplicationUserCount
// @Description 获取注册到该应用的用户数量
// @Param	appID	query	string	true	"应用ID"
// @Success 200 {object} common.ApiResponse "成功返回用户数量"
// @router /get-application-user-count [get]
func (h *ApplicationHandler) GetApplicationUserCount() {
	appID := h.Input().Get("appID")
	if appID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	count, err := h.appService.GetApplicationUserCount(appID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{"userCount": count})
}

// EnableApplication 启用应用
// @Title EnableApplication
// @Description 启用应用
// @Param	owner	query	string	true	"组织名称"
// @Param	name	query	string	true	"应用名称"
// @Success 200 {object} common.ApiResponse "成功返回启用结果"
// @router /enable-application [post]
func (h *ApplicationHandler) EnableApplication() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")

	if owner == "" || name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.appService.EnableApplication(owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// DisableApplication 禁用应用
// @Title DisableApplication
// @Description 禁用应用
// @Param	owner	query	string	true	"组织名称"
// @Param	name	query	string	true	"应用名称"
// @Success 200 {object} common.ApiResponse "成功返回禁用结果"
// @router /disable-application [post]
func (h *ApplicationHandler) DisableApplication() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")

	if owner == "" || name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.appService.DisableApplication(owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// UpdateApplicationTheme 更新应用主题
// @Title UpdateApplicationTheme
// @Description 更新应用主题配置
// @Param	owner	formData	string	true	"组织名称"
// @Param	name	formData	string	true	"应用名称"
// @Param	theme	formData	string	true	"主题名称"
// @Success 200 {object} common.ApiResponse "成功返回更新结果"
// @router /update-application-theme [post]
func (h *ApplicationHandler) UpdateApplicationTheme() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")
	theme := h.Input().Get("theme")

	if owner == "" || name == "" || theme == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.appService.UpdateApplicationTheme(owner, name, theme)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// RotateClientCredentials 轮换客户端凭证
// @Title RotateClientCredentials
// @Description 轮换应用的客户端ID和客户端密钥
// @Param	owner	query	string	true	"组织名称"
// @Param	name	query	string	true	"应用名称"
// @Success 200 {object} common.ApiResponse "成功返回新的客户端凭证"
// @router /rotate-client-credentials [post]
func (h *ApplicationHandler) RotateClientCredentials() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")

	if owner == "" || name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	app, err := h.appService.RotateClientCredentials(owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{
		"clientId":     app.ClientId,
		"clientSecret": app.ClientSecret,
	})
}

// RegenerateClientSecret 重新生成客户端密钥
// @Title RegenerateClientSecret
// @Description 重新生成应用的客户端密钥（保留客户端ID）
// @Param	owner	query	string	true	"组织名称"
// @Param	name	query	string	true	"应用名称"
// @Success 200 {object} common.ApiResponse "成功返回新的客户端密钥"
// @router /regenerate-client-secret [post]
func (h *ApplicationHandler) RegenerateClientSecret() {
	owner := h.Input().Get("owner")
	name := h.Input().Get("name")

	if owner == "" || name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	newSecret, err := h.appService.RegenerateClientSecret(owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{"clientSecret": newSecret})
}

// ValidateRedirectURI 验证重定向URI
// @Title ValidateRedirectURI
// @Description 验证重定向URI是否合法
// @Param	clientId	query	string	true	"客户端ID"
// @Param	redirectUri	query	string	true	"重定向URI"
// @Success 200 {object} common.ApiResponse "成功返回验证结果"
// @router /validate-redirect-uri [get]
func (h *ApplicationHandler) ValidateRedirectURI() {
	clientID := h.Input().Get("clientId")
	redirectURI := h.Input().Get("redirectUri")

	if clientID == "" || redirectURI == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	app, err := h.appService.GetApplicationByClientID(clientID)
	if err != nil {
		h.ResponseError("Application not found")
		return
	}

	isValid := h.appService.ValidateRedirectURI(app, redirectURI)
	h.ResponseSuccess(map[string]interface{}{"isValid": isValid})
}

// 辅助方法

// GetIntFromInput 从输入中获取整数值
func (h *ApplicationHandler) GetIntFromInput(key string, defaultValue int) int {
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
func (h *ApplicationHandler) ResponseSuccess(data ...interface{}) {
	h.Data["json"] = common.ResponseSuccess(data...)
	h.ServeJSON()
}

// ResponseSuccessWithPage 返回带分页的成功响应
func (h *ApplicationHandler) ResponseSuccessWithPage(list interface{}, total int64, page, size int) {
	h.Data["json"] = common.ResponseSuccessWithPage(list, total, page, size)
	h.ServeJSON()
}

// ResponseError 返回错误响应
func (h *ApplicationHandler) ResponseError(message string, data ...interface{}) {
	h.Data["json"] = common.ResponseError(i18n.Translate(h.GetAcceptLanguage(), message), data...)
	h.ServeJSON()
}

// ResponseFromAction 根据操作结果返回响应
func (h *ApplicationHandler) ResponseFromAction(affected bool, err ...error) {
	h.Data["json"] = common.ResponseFromAction(affected, err...)
	h.ServeJSON()
}

// GetAcceptLanguage 获取接受语言
func (h *ApplicationHandler) GetAcceptLanguage() string {
	return h.Ctx.Input.Header("Accept-Language")
}
