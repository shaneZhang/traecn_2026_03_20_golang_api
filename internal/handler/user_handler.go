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
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/dto"
	"github.com/casdoor/casdoor/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	web.Controller
	userService service.UserService
}

// NewUserHandler creates new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUser gets user by ID
// @router /:id [get]
func (h *UserHandler) GetUser() {
	id := h.Ctx.Input.Param(":id")
	
	user, err := h.userService.GetUser(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(user)
}

// CreateUser creates new user
// @router / [post]
func (h *UserHandler) CreateUser() {
	var req dto.CreateUserRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	user, err := h.userService.CreateUser(h.Ctx.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			h.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
		}
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(user)
}

// UpdateUser updates user
// @router /:id [put]
func (h *UserHandler) UpdateUser() {
	id := h.Ctx.Input.Param(":id")
	
	var req dto.UpdateUserRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	user, err := h.userService.UpdateUser(h.Ctx.Request.Context(), id, &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(user)
}

// DeleteUser deletes user
// @router /:id [delete]
func (h *UserHandler) DeleteUser() {
	id := h.Ctx.Input.Param(":id")
	
	if err := h.userService.DeleteUser(h.Ctx.Request.Context(), id); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("deleted")
}

// ListUsers lists users
// @router / [get]
func (h *UserHandler) ListUsers() {
	owner := h.GetString("owner")
	pageSize, _ := h.GetInt("pageSize", 10)
	page, _ := h.GetInt("p", 1)
	field := h.GetString("field")
	value := h.GetString("value")
	sortField := h.GetString("sortField")
	sortOrder := h.GetString("sortOrder")
	
	req := &dto.ListUsersRequest{
		Owner:     owner,
		PageSize:  pageSize,
		Page:      page,
		Field:     field,
		Value:     value,
		SortField: sortField,
		SortOrder: sortOrder,
	}
	
	users, err := h.userService.ListUsers(h.Ctx.Request.Context(), req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(users)
}

// ListGlobalUsers lists global users
// @router /global [get]
func (h *UserHandler) ListGlobalUsers() {
	page, _ := h.GetInt("p", 1)
	pageSize, _ := h.GetInt("pageSize", 10)
	field := h.GetString("field")
	value := h.GetString("value")
	sortField := h.GetString("sortField")
	sortOrder := h.GetString("sortOrder")
	
	users, err := h.userService.ListGlobalUsers(h.Ctx.Request.Context(), page, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(users)
}

// BatchCreateUsers batch creates users
// @router /batch [post]
func (h *UserHandler) BatchCreateUsers() {
	var req dto.ImportUsersRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	resp, err := h.userService.BatchCreateUsers(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(resp)
}

// BatchUpdateUsers batch updates users
// @router /batch [put]
func (h *UserHandler) BatchUpdateUsers() {
	var operation dto.BatchUserOperation
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &operation); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.userService.BatchUpdateUsers(h.Ctx.Request.Context(), &operation); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("batch update successful")
}

// BatchDeleteUsers batch deletes users
// @router /batch [delete]
func (h *UserHandler) BatchDeleteUsers() {
	var ids []string
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &ids); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.userService.BatchDeleteUsers(h.Ctx.Request.Context(), ids); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("batch delete successful")
}

// ImportUsers imports users from file
// @router /import [post]
func (h *UserHandler) ImportUsers() {
	owner := h.GetString("owner")
	fileType := h.GetString("fileType")
	
	file, _, err := h.GetFile("file")
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	defer file.Close()
	
	resp, err := h.userService.ImportUsers(h.Ctx.Request.Context(), owner, file, fileType)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(resp)
}

// ExportUsers exports users to file
// @router /export [get]
func (h *UserHandler) ExportUsers() {
	var req dto.ExportUsersRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	data, contentType, err := h.userService.ExportUsers(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	ext := "txt"
	switch req.Format {
	case "xlsx":
		ext = "xlsx"
	case "csv":
		ext = "csv"
	case "json":
		ext = "json"
	}
	
	h.Ctx.Output.Header("Content-Type", contentType)
	h.Ctx.Output.Header("Content-Disposition", "attachment; filename=users."+ext)
	h.Ctx.Output.Body(data)
}

// SetupMFA sets up MFA
// @router /:id/mfa/setup [post]
func (h *UserHandler) SetupMFA() {
	id := h.Ctx.Input.Param(":id")
	
	var req dto.MFASetupRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	resp, err := h.userService.SetupMFA(h.Ctx.Request.Context(), id, &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(resp)
}

// VerifyMFASetup verifies MFA setup
// @router /:id/mfa/verify [post]
func (h *UserHandler) VerifyMFASetup() {
	id := h.Ctx.Input.Param(":id")
	
	var req dto.MFASetupRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.userService.VerifyMFASetup(h.Ctx.Request.Context(), id, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("MFA verified successfully")
}

// EnableMFA enables MFA
// @router /:id/mfa/enable [post]
func (h *UserHandler) EnableMFA() {
	id := h.Ctx.Input.Param(":id")
	
	var req dto.MFASetupRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.userService.EnableMFA(h.Ctx.Request.Context(), id, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("MFA enabled successfully")
}

// DisableMFA disables MFA
// @router /:id/mfa/disable [post]
func (h *UserHandler) DisableMFA() {
	id := h.Ctx.Input.Param(":id")
	
	if err := h.userService.DisableMFA(h.Ctx.Request.Context(), id); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("MFA disabled successfully")
}

// GetMFAStatus gets MFA status
// @router /:id/mfa/status [get]
func (h *UserHandler) GetMFAStatus() {
	id := h.Ctx.Input.Param(":id")
	
	status, err := h.userService.GetMFAStatus(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(status)
}

// RecoverMFA recovers MFA
// @router /:id/mfa/recover [post]
func (h *UserHandler) RecoverMFA() {
	id := h.Ctx.Input.Param(":id")
	
	var req struct {
		RecoveryCode string `json:"recoveryCode"`
	}
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.userService.RecoverMFA(h.Ctx.Request.Context(), id, req.RecoveryCode); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("MFA recovered successfully")
}

// AddUserToGroup adds user to group
// @router /:id/groups/:groupId [post]
func (h *UserHandler) AddUserToGroup() {
	id := h.Ctx.Input.Param(":id")
	groupID := h.Ctx.Input.Param(":groupId")
	
	if err := h.userService.AddUserToGroup(h.Ctx.Request.Context(), id, groupID); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("user added to group successfully")
}

// RemoveUserFromGroup removes user from group
// @router /:id/groups/:groupId [delete]
func (h *UserHandler) RemoveUserFromGroup() {
	id := h.Ctx.Input.Param(":id")
	groupID := h.Ctx.Input.Param(":groupId")
	
	if err := h.userService.RemoveUserFromGroup(h.Ctx.Request.Context(), id, groupID); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("user removed from group successfully")
}

// GetUserStatistics gets user statistics
// @router /statistics [get]
func (h *UserHandler) GetUserStatistics() {
	owner := h.GetString("owner")
	if owner == "" {
		h.ResponseError("owner is required")
		return
	}
	
	stats, err := h.userService.GetUserStatistics(h.Ctx.Request.Context(), owner)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(stats)
}

// ResponseOk sends success response
func (h *UserHandler) ResponseOk(data interface{}) {
	h.Data["json"] = common.Success(data)
	h.ServeJSON()
}

// ResponseError sends error response
func (h *UserHandler) ResponseError(message string) {
	h.Data["json"] = common.Error(http.StatusBadRequest, message)
	h.ServeJSON()
}
