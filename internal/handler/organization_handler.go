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

// OrganizationHandler handles organization-related HTTP requests
type OrganizationHandler struct {
	web.Controller
	orgService service.OrganizationService
}

// NewOrganizationHandler creates new organization handler
func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

// GetOrganization gets organization by ID
// @router /:id [get]
func (h *OrganizationHandler) GetOrganization() {
	id := h.Ctx.Input.Param(":id")
	
	org, err := h.orgService.GetOrganization(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(org)
}

// GetOrganizationByName gets organization by owner and name
// @router /name/:owner/:name [get]
func (h *OrganizationHandler) GetOrganizationByName() {
	owner := h.Ctx.Input.Param(":owner")
	name := h.Ctx.Input.Param(":name")
	
	org, err := h.orgService.GetOrganizationByName(h.Ctx.Request.Context(), owner, name)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(org)
}

// CreateOrganization creates new organization
// @router / [post]
func (h *OrganizationHandler) CreateOrganization() {
	var req dto.CreateOrganizationRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	org, err := h.orgService.CreateOrganization(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(org)
}

// UpdateOrganization updates organization
// @router /:id [put]
func (h *OrganizationHandler) UpdateOrganization() {
	id := h.Ctx.Input.Param(":id")
	
	var req dto.UpdateOrganizationRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	org, err := h.orgService.UpdateOrganization(h.Ctx.Request.Context(), id, &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(org)
}

// DeleteOrganization deletes organization
// @router /:id [delete]
func (h *OrganizationHandler) DeleteOrganization() {
	id := h.Ctx.Input.Param(":id")
	
	if err := h.orgService.DeleteOrganization(h.Ctx.Request.Context(), id); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("deleted")
}

// ListOrganizations lists organizations
// @router / [get]
func (h *OrganizationHandler) ListOrganizations() {
	owner := h.GetString("owner")
	pageSize, _ := h.GetInt("pageSize", 10)
	page, _ := h.GetInt("p", 1)
	field := h.GetString("field")
	value := h.GetString("value")
	sortField := h.GetString("sortField")
	sortOrder := h.GetString("sortOrder")
	
	req := &dto.ListOrganizationsRequest{
		Owner:     owner,
		PageSize:  pageSize,
		Page:      page,
		Field:     field,
		Value:     value,
		SortField: sortField,
		SortOrder: sortOrder,
	}
	
	orgs, err := h.orgService.ListOrganizations(h.Ctx.Request.Context(), req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(orgs)
}

// GetOrganizationHierarchy gets organization hierarchy
// @router /:id/hierarchy [get]
func (h *OrganizationHandler) GetOrganizationHierarchy() {
	id := h.Ctx.Input.Param(":id")
	
	hierarchy, err := h.orgService.GetOrganizationHierarchy(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(hierarchy)
}

// GetOrganizationTree gets organization tree
// @router /tree [get]
func (h *OrganizationHandler) GetOrganizationTree() {
	owner := h.GetString("owner")
	
	tree, err := h.orgService.GetOrganizationTree(h.Ctx.Request.Context(), owner)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(tree)
}

// MoveOrganization moves organization to new parent
// @router /:id/move [post]
func (h *OrganizationHandler) MoveOrganization() {
	id := h.Ctx.Input.Param(":id")
	
	var req struct {
		NewParentID string `json:"newParentId"`
	}
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.orgService.MoveOrganization(h.Ctx.Request.Context(), id, req.NewParentID); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("organization moved successfully")
}

// GetOrganizationChildren gets organization children
// @router /:id/children [get]
func (h *OrganizationHandler) GetOrganizationChildren() {
	id := h.Ctx.Input.Param(":id")
	
	children, err := h.orgService.GetOrganizationChildren(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(children)
}

// GetOrganizationDescendants gets organization descendants
// @router /:id/descendants [get]
func (h *OrganizationHandler) GetOrganizationDescendants() {
	id := h.Ctx.Input.Param(":id")
	
	descendants, err := h.orgService.GetOrganizationDescendants(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(descendants)
}

// GetOrganizationAncestors gets organization ancestors
// @router /:id/ancestors [get]
func (h *OrganizationHandler) GetOrganizationAncestors() {
	id := h.Ctx.Input.Param(":id")
	
	ancestors, err := h.orgService.GetOrganizationAncestors(h.Ctx.Request.Context(), id)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(ancestors)
}

// BatchCreateOrganizations batch creates organizations
// @router /batch [post]
func (h *OrganizationHandler) BatchCreateOrganizations() {
	var req dto.BatchCreateOrganizationsRequest
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &req); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	resp, err := h.orgService.BatchCreateOrganizations(h.Ctx.Request.Context(), &req)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(resp)
}

// BatchUpdateOrganizations batch updates organizations
// @router /batch [put]
func (h *OrganizationHandler) BatchUpdateOrganizations() {
	var operation dto.BatchOrganizationOperation
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &operation); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.orgService.BatchUpdateOrganizations(h.Ctx.Request.Context(), &operation); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("batch update successful")
}

// BatchDeleteOrganizations batch deletes organizations
// @router /batch [delete]
func (h *OrganizationHandler) BatchDeleteOrganizations() {
	var ids []string
	if err := json.Unmarshal(h.Ctx.Input.RequestBody, &ids); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	if err := h.orgService.BatchDeleteOrganizations(h.Ctx.Request.Context(), ids); err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk("batch delete successful")
}

// SearchOrganizations searches organizations
// @router /search [get]
func (h *OrganizationHandler) SearchOrganizations() {
	owner := h.GetString("owner")
	keyword := h.GetString("keyword")
	
	orgs, err := h.orgService.SearchOrganizations(h.Ctx.Request.Context(), owner, keyword)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(orgs)
}

// GetOrganizationStatistics gets organization statistics
// @router /statistics [get]
func (h *OrganizationHandler) GetOrganizationStatistics() {
	owner := h.GetString("owner")
	if owner == "" {
		h.ResponseError("owner is required")
		return
	}
	
	stats, err := h.orgService.GetOrganizationStatistics(h.Ctx.Request.Context(), owner)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}
	
	h.ResponseOk(stats)
}

// ResponseOk sends success response
func (h *OrganizationHandler) ResponseOk(data interface{}) {
	h.Data["json"] = common.Success(data)
	h.ServeJSON()
}

// ResponseError sends error response
func (h *OrganizationHandler) ResponseError(message string) {
	h.Data["json"] = common.Error(http.StatusBadRequest, message)
	h.ServeJSON()
}
