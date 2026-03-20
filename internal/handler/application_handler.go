// Copyright 2024 The Refactored Authors. All Rights Reserved.
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
	"net/http"

	"github.com/beego/beego/v2/server/web"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/dto"
	"github.com/casdoor/casdoor/internal/service"
)

// ApplicationHandler handles application-related HTTP requests
type ApplicationHandler struct {
	web.Controller
	appService service.ApplicationService
}

// NewApplicationHandler creates new application handler
func NewApplicationHandler(appService service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{appService: appService}
}

// GetApplication gets application by ID
// @router /:id [get]
func (h *ApplicationHandler) GetApplication() {
	id := h.Ctx.Input.Param(":id")

	app, err := h.appService.GetApplication(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(app)
}

// GetApplicationByClientID gets application by client ID
// @router /client/:clientId [get]
func (h *ApplicationHandler) GetApplicationByClientID() {
	clientID := h.Ctx.Input.Param(":clientId")

	app, err := h.appService.GetApplicationByClientID(h.Ctx.Request.Context(), clientID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(app)
}

// CreateApplication creates new application
// @router / [post]
func (h *ApplicationHandler) CreateApplication() {
	var req dto.CreateApplicationRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	app, err := h.appService.CreateApplication(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(app)
}

// UpdateApplication updates application
// @router /:id [put]
func (h *ApplicationHandler) UpdateApplication() {
	id := h.Ctx.Input.Param(":id")

	var req dto.UpdateApplicationRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	app, err := h.appService.UpdateApplication(h.Ctx.Request.Context(), id, &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(app)
}

// DeleteApplication deletes application
// @router /:id [delete]
func (h *ApplicationHandler) DeleteApplication() {
	id := h.Ctx.Input.Param(":id")

	if err := h.appService.DeleteApplication(h.Ctx.Request.Context(), id); err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk("deleted")
}

// ListApplications lists applications
// @router / [get]
func (h *ApplicationHandler) ListApplications() {
	owner := h.GetString("owner")
	pageSize, _ := h.GetInt("pageSize", 10)
	page, _ := h.GetInt("p", 1)
	field := h.GetString("field")
	value := h.GetString("value")
	sortField := h.GetString("sortField")
	sortOrder := h.GetString("sortOrder")

	req := &dto.ListApplicationsRequest{
		Owner:     owner,
		PageSize:  pageSize,
		Page:      page,
		Field:     field,
		Value:     value,
		SortField: sortField,
		SortOrder: sortOrder,
	}

	apps, err := h.appService.ListApplications(h.Ctx.Request.Context(), req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(apps)
}

// GetApplicationsByOrganization gets applications by organization
// @router /organization/:owner/:organization [get]
func (h *ApplicationHandler) GetApplicationsByOrganization() {
	owner := h.Ctx.Input.Param(":owner")
	organization := h.Ctx.Input.Param(":organization")
	page, _ := h.GetInt("p", 1)
	pageSize, _ := h.GetInt("pageSize", 10)

	apps, err := h.appService.GetApplicationsByOrganization(h.Ctx.Request.Context(), owner, organization, page, pageSize)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(apps)
}

// OAuthAuthorize OAuth authorization endpoint
// @router /oauth/authorize [get]
func (h *ApplicationHandler) OAuthAuthorize() {
	var req dto.OAuthAuthorizeRequest
	if err := h.Ctx.Input.Bind(&req.ClientID, "client_id"); err != nil {
		h.ResponseError(err.Error())
		return
	}
	if err := h.Ctx.Input.Bind(&req.RedirectURI, "redirect_uri"); err != nil {
		h.ResponseError(err.Error())
		return
	}
	if err := h.Ctx.Input.Bind(&req.ResponseType, "response_type"); err != nil {
		h.ResponseError(err.Error())
		return
	}
	if err := h.Ctx.Input.Bind(&req.Scope, "scope"); err != nil {
		h.ResponseError(err.Error())
		return
	}
	if err := h.Ctx.Input.Bind(&req.State, "state"); err != nil {
		h.ResponseError(err.Error())
		return
	}

	app, err := h.appService.ValidateOAuthRequest(h.Ctx.Request.Context(), req.ClientID, req.RedirectURI, req.ResponseType)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(app)
}

// OAuthToken OAuth token endpoint
// @router /oauth/token [post]
func (h *ApplicationHandler) OAuthToken() {
	var req dto.TokenRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	token, err := h.appService.GenerateToken(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(token)
}

// OAuthRefreshToken OAuth refresh token endpoint
// @router /oauth/refresh [post]
func (h *ApplicationHandler) OAuthRefreshToken() {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	token, err := h.appService.RefreshToken(h.Ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(token)
}

// OAuthRevokeToken OAuth revoke token endpoint
// @router /oauth/revoke [post]
func (h *ApplicationHandler) OAuthRevokeToken() {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	if err := h.appService.RevokeToken(h.Ctx.Request.Context(), req.Token); err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk("token revoked")
}

// GrantPermission grants permission to user
// @router /:id/permissions [post]
func (h *ApplicationHandler) GrantPermission() {
	id := h.Ctx.Input.Param(":id")

	var req dto.GrantPermissionRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	if err := h.appService.GrantPermission(h.Ctx.Request.Context(), id, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk("permission granted")
}

// RevokePermission revokes permission from user
// @router /:id/permissions/revoke [post]
func (h *ApplicationHandler) RevokePermission() {
	id := h.Ctx.Input.Param(":id")

	var req dto.RevokePermissionRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	if err := h.appService.RevokePermission(h.Ctx.Request.Context(), id, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk("permission revoked")
}

// GetPermissions gets application permissions
// @router /:id/permissions [get]
func (h *ApplicationHandler) GetPermissions() {
	id := h.Ctx.Input.Param(":id")

	permissions, err := h.appService.GetPermissions(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(permissions)
}

// BatchCreateApplications batch creates applications
// @router /batch [post]
func (h *ApplicationHandler) BatchCreateApplications() {
	var req dto.BatchCreateApplicationsRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}

	resp, err := h.appService.BatchCreateApplications(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(resp)
}

// BatchUpdateApplications batch updates applications
// @router /batch [put]
func (h *ApplicationHandler) BatchUpdateApplications() {
	var operation dto.BatchApplicationOperation
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &operation); err != nil {
		h.ResponseError(err.Error())
		return
	}

	if err := h.appService.BatchUpdateApplications(h.Ctx.Request.Context(), &operation); err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk("batch update successful")
}

// BatchDeleteApplications batch deletes applications
// @router /batch [delete]
func (h *ApplicationHandler) BatchDeleteApplications() {
	var ids []string
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &ids); err != nil {
		h.ResponseError(err.Error())
		return
	}

	if err := h.appService.BatchDeleteApplications(h.Ctx.Request.Context(), ids); err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk("batch delete successful")
}

// SearchApplications searches applications
// @router /search [get]
func (h *ApplicationHandler) SearchApplications() {
	owner := h.GetString("owner")
	keyword := h.GetString("keyword")

	apps, err := h.appService.SearchApplications(h.Ctx.Request.Context(), owner, keyword)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(apps)
}

// GetApplicationStatistics gets application statistics
// @router /statistics [get]
func (h *ApplicationHandler) GetApplicationStatistics() {
	owner := h.GetString("owner")
	if owner == "" {
		h.ResponseError("owner is required")
		return
	}

	stats, err := h.appService.GetApplicationStatistics(h.Ctx.Request.Context(), owner)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseOk(stats)
}

// ResponseOk sends success response
func (h *ApplicationHandler) ResponseOk(data interface{}) {
	h.Data["json"] = common.Success(data)
	h.ServeJSON()
}

// ResponseError sends error response
func (h *ApplicationHandler) ResponseError(message string) {
	h.Data["json"] = common.Error(http.StatusBadRequest, message)
	h.ServeJSON()
}
