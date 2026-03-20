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

// OrganizationHandler 组织API处理器
type OrganizationHandler struct {
	controllers.BaseController
	orgService service.OrganizationService
}

// NewOrganizationHandler 创建组织处理器实例
func NewOrganizationHandler() *OrganizationHandler {
	return &OrganizationHandler{
		orgService: service.NewOrganizationService(),
	}
}

// GetOrganizations 获取组织列表
// @Title GetOrganizations
// @Description 获取组织列表（分页）
// @Param	page	query	int	false	"页码"
// @Param	pageSize	query	int	false	"每页大小"
// @Param	field	query	string	false	"搜索字段"
// @Param	value	query	string	false	"搜索值"
// @Param	sortField	query	string	false	"排序字段"
// @Param	sortOrder	query	string	false	"排序方向"
// @Success 200 {object} common.ApiResponse "成功返回组织列表"
// @router /get-organizations [get]
func (h *OrganizationHandler) GetOrganizations() {
	page := h.GetIntFromInput("page", common.DefaultPage)
	pageSize := h.GetIntFromInput("pageSize", common.DefaultPageSize)
	field := h.Input().Get("field")
	value := h.Input().Get("value")
	sortField := h.Input().Get("sortField")
	sortOrder := h.Input().Get("sortOrder")

	orgs, total, err := h.orgService.ListOrganizations(page, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccessWithPage(orgs, total, page, pageSize)
}

// GetOrganization 获取单个组织
// @Title GetOrganization
// @Description 获取单个组织详情
// @Param	id	query	string	true	"组织ID，格式：admin/name"
// @Success 200 {object} common.ApiResponse "成功返回组织详情"
// @router /get-organization [get]
func (h *OrganizationHandler) GetOrganization() {
	id := h.Input().Get("id")
	if id == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	owner, name := util.GetOwnerAndNameFromId(id)
	if owner != common.AdminUser {
		h.ResponseError("Invalid organization ID format")
		return
	}

	org, err := h.orgService.GetOrganization(name)
	if err != nil {
		if err == common.ErrOrgNotFound {
			h.ResponseError("Organization not found")
		} else {
			h.ResponseError(err.Error())
		}
		return
	}

	h.ResponseSuccess(org)
}

// GetOrganizationByUserID 根据用户ID获取组织
// @Title GetOrganizationByUserID
// @Description 根据用户ID获取所属组织
// @Param	userID	query	string	true	"用户ID"
// @Success 200 {object} common.ApiResponse "成功返回组织详情"
// @router /get-organization-by-userid [get]
func (h *OrganizationHandler) GetOrganizationByUserID() {
	userID := h.Input().Get("userID")
	if userID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	org, err := h.orgService.GetOrganizationByUserID(userID)
	if err != nil {
		if err == common.ErrOrgNotFound || err == common.ErrUserNotFound {
			h.ResponseError("Organization not found for the user")
		} else {
			h.ResponseError(err.Error())
		}
		return
	}

	h.ResponseSuccess(org)
}

// UpdateOrganization 更新组织
// @Title UpdateOrganization
// @Description 更新组织信息
// @Param	body	body	object.Organization	true	"组织信息"
// @Success 200 {object} common.ApiResponse "成功返回更新结果"
// @router /update-organization [post]
func (h *OrganizationHandler) UpdateOrganization() {
	var org object.Organization
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &org); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if org.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.orgService.UpdateOrganization(&org)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// AddOrganization 添加组织
// @Title AddOrganization
// @Description 添加新组织
// @Param	body	body	object.Organization	true	"组织信息"
// @Success 200 {object} common.ApiResponse "成功返回添加结果"
// @router /add-organization [post]
func (h *OrganizationHandler) AddOrganization() {
	var org object.Organization
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &org); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if org.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.orgService.CreateOrganization(&org)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// DeleteOrganization 删除组织
// @Title DeleteOrganization
// @Description 删除组织
// @Param	body	body	object.Organization	true	"组织信息（只需要owner和name）"
// @Success 200 {object} common.ApiResponse "成功返回删除结果"
// @router /delete-organization [post]
func (h *OrganizationHandler) DeleteOrganization() {
	var org object.Organization
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &org); err != nil {
		h.ResponseError("Invalid request body: " + err.Error())
		return
	}

	if org.Name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.orgService.DeleteOrganization(&org)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// GetAllOrganizations 获取所有组织
// @Title GetAllOrganizations
// @Description 获取所有组织列表
// @Success 200 {object} common.ApiResponse "成功返回所有组织列表"
// @router /get-all-organizations [get]
func (h *OrganizationHandler) GetAllOrganizations() {
	orgs, err := h.orgService.GetAllOrganizations()
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(orgs)
}

// GetParentOrganizations 获取父组织列表
// @Title GetParentOrganizations
// @Description 获取父组织列表
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回父组织列表"
// @router /get-parent-organizations [get]
func (h *OrganizationHandler) GetParentOrganizations() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	orgs, err := h.orgService.GetParentOrganizations(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(orgs)
}

// GetChildOrganizations 获取子组织列表
// @Title GetChildOrganizations
// @Description 获取直接子组织列表
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回子组织列表"
// @router /get-child-organizations [get]
func (h *OrganizationHandler) GetChildOrganizations() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	orgs, err := h.orgService.GetChildOrganizations(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(orgs)
}

// GetAllChildOrganizations 获取所有子组织列表
// @Title GetAllChildOrganizations
// @Description 获取所有子组织列表（递归）
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回所有子组织列表"
// @router /get-all-child-organizations [get]
func (h *OrganizationHandler) GetAllChildOrganizations() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	orgs, err := h.orgService.GetAllChildOrganizations(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(orgs)
}

// GetOrganizationHierarchy 获取组织层级结构
// @Title GetOrganizationHierarchy
// @Description 获取组织完整层级结构
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回组织层级结构"
// @router /get-organization-hierarchy [get]
func (h *OrganizationHandler) GetOrganizationHierarchy() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	hierarchy, err := h.orgService.GetOrganizationHierarchy(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(hierarchy)
}

// GetOrganizationStats 获取组织统计信息
// @Title GetOrganizationStats
// @Description 获取组织统计信息（用户数、应用数等）
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回组织统计信息"
// @router /get-organization-stats [get]
func (h *OrganizationHandler) GetOrganizationStats() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	stats, err := h.orgService.GetOrganizationStats(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(stats)
}

// EnableOrganization 启用组织
// @Title EnableOrganization
// @Description 启用组织
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回启用结果"
// @router /enable-organization [post]
func (h *OrganizationHandler) EnableOrganization() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.orgService.EnableOrganization(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// DisableOrganization 禁用组织
// @Title DisableOrganization
// @Description 禁用组织
// @Param	name	query	string	true	"组织名称"
// @Success 200 {object} common.ApiResponse "成功返回禁用结果"
// @router /disable-organization [post]
func (h *OrganizationHandler) DisableOrganization() {
	name := h.Input().Get("name")
	if name == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.orgService.DisableOrganization(name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// UpdateOrganizationTheme 更新组织主题
// @Title UpdateOrganizationTheme
// @Description 更新组织主题配置
// @Param	name	formData	string	true	"组织名称"
// @Param	theme	formData	string	true	"主题名称"
// @Success 200 {object} common.ApiResponse "成功返回更新结果"
// @router /update-organization-theme [post]
func (h *OrganizationHandler) UpdateOrganizationTheme() {
	name := h.Input().Get("name")
	theme := h.Input().Get("theme")

	if name == "" || theme == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	affected, err := h.orgService.UpdateOrganizationTheme(name, theme)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseFromAction(affected)
}

// 辅助方法

// GetIntFromInput 从输入中获取整数值
func (h *OrganizationHandler) GetIntFromInput(key string, defaultValue int) int {
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
func (h *OrganizationHandler) ResponseSuccess(data ...interface{}) {
	h.Data["json"] = common.ResponseSuccess(data...)
	h.ServeJSON()
}

// ResponseSuccessWithPage 返回带分页的成功响应
func (h *OrganizationHandler) ResponseSuccessWithPage(list interface{}, total int64, page, size int) {
	h.Data["json"] = common.ResponseSuccessWithPage(list, total, page, size)
	h.ServeJSON()
}

// ResponseError 返回错误响应
func (h *OrganizationHandler) ResponseError(message string, data ...interface{}) {
	h.Data["json"] = common.ResponseError(i18n.Translate(h.GetAcceptLanguage(), message), data...)
	h.ServeJSON()
}

// ResponseFromAction 根据操作结果返回响应
func (h *OrganizationHandler) ResponseFromAction(affected bool, err ...error) {
	h.Data["json"] = common.ResponseFromAction(affected, err...)
	h.ServeJSON()
}

// GetAcceptLanguage 获取接受语言
func (h *OrganizationHandler) GetAcceptLanguage() string {
	return h.Ctx.Input.Header("Accept-Language")
}
