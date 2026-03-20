package api

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/utils/pagination"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/service"
	"github.com/casdoor/casdoor/util"
)

type OrganizationAPI struct {
	BaseAPI
	orgService   service.OrganizationService
	groupService service.GroupService
}

func NewOrganizationAPI(orgService service.OrganizationService, groupService service.GroupService) *OrganizationAPI {
	return &OrganizationAPI{
		orgService:   orgService,
		groupService: groupService,
	}
}

func (api *OrganizationAPI) GetOrganizations(c *BaseController) {
	owner := c.Ctx.Input.Query("owner")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")
	organizationName := c.Ctx.Input.Query("organizationName")

	isGlobalAdmin := c.IsGlobalAdmin()
	if limit == "" || page == "" {
		var organizations []*object.Organization
		var err error
		if isGlobalAdmin {
			organizations, err = api.orgService.GetOrganizations(owner)
		} else {
			organizations, err = api.orgService.GetOrganizations(owner, c.getCurrentUser().Owner)
		}

		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		maskedOrgs, err := api.orgService.GetMaskedOrganizations(organizations)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(maskedOrgs)
	} else {
		if !isGlobalAdmin {
			organizations, err := api.orgService.GetOrganizations(owner, c.getCurrentUser().Owner)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
			maskedOrgs, err := api.orgService.GetMaskedOrganizations(organizations)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
			c.ResponseOk(maskedOrgs)
		} else {
			limitInt := util.ParseInt(limit)
			count, err := api.orgService.GetOrganizationCount(owner, organizationName, field, value)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}

			paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
			organizations, _, err := api.orgService.GetPaginationOrganizations(owner, organizationName, paginator.Offset(), limitInt, field, value, sortField, sortOrder)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}

			maskedOrgs, err := api.orgService.GetMaskedOrganizations(organizations)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}

			c.ResponseOk(maskedOrgs, paginator.Nums())
		}
	}
}

func (api *OrganizationAPI) GetOrganization(c *BaseController) {
	id := c.Ctx.Input.Query("id")
	organization, err := api.orgService.GetOrganization(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	maskedOrg, err := api.orgService.GetMaskedOrganization(organization)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if maskedOrg != nil && maskedOrg.MfaRememberInHours == 0 {
		maskedOrg.MfaRememberInHours = 12
	}

	c.ResponseOk(maskedOrg)
}

func (api *OrganizationAPI) UpdateOrganization(c *BaseController) {
	id := c.Ctx.Input.Query("id")

	var organization object.Organization
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &organization)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	isGlobalAdmin, _ := c.isGlobalAdmin()

	if organization.BalanceCurrency == "" {
		organization.BalanceCurrency = "USD"
	}

	affected, err := api.orgService.UpdateOrganization(id, &organization, isGlobalAdmin, c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *OrganizationAPI) AddOrganization(c *BaseController) {
	var organization object.Organization
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &organization)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if organization.BalanceCurrency == "" {
		organization.BalanceCurrency = "USD"
	}

	affected, err := api.orgService.CreateOrganization(&organization, c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *OrganizationAPI) DeleteOrganization(c *BaseController) {
	var organization object.Organization
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &organization)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.orgService.DeleteOrganization(&organization)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *OrganizationAPI) GetDefaultApplication(c *BaseController) {
	id := c.Ctx.Input.Query("id")

	application, err := api.orgService.GetDefaultApplication(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	userId := c.GetSessionUsername()
	application = object.GetMaskedApplication(application, userId)
	c.ResponseOk(application)
}

func (api *OrganizationAPI) GetOrganizationNames(c *BaseController) {
	owner := c.Ctx.Input.Query("owner")
	organizationNames, err := api.orgService.GetOrganizationsByFields(owner, "name", "display_name")
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(organizationNames)
}

type GroupAPI struct {
	BaseAPI
	groupService service.GroupService
}

func NewGroupAPI(groupService service.GroupService) *GroupAPI {
	return &GroupAPI{
		groupService: groupService,
	}
}

func (api *GroupAPI) GetGroups(c *BaseController) {
	owner := c.Ctx.Input.Query("owner")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")
	withTree := c.Ctx.Input.Query("withTree")

	if limit == "" || page == "" {
		groups, err := api.groupService.GetGroups(owner)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		err = api.groupService.ExtendGroupsWithUsers(groups)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		if withTree == "true" {
			c.ResponseOk(api.groupService.ConvertToTreeData(groups, owner))
			return
		}

		c.ResponseOk(groups)
	} else {
		limitInt := util.ParseInt(limit)
		count, err := api.groupService.GetGroupCount(owner, field, value)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
		groups, _, err := api.groupService.GetPaginationGroups(owner, paginator.Offset(), limitInt, field, value, sortField, sortOrder)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		err = api.groupService.ExtendGroupsWithUsers(groups)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(groups, paginator.Nums())
	}
}

func (api *GroupAPI) GetGroup(c *BaseController) {
	id := c.Ctx.Input.Query("id")

	group, err := api.groupService.GetGroup(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	err = api.groupService.ExtendGroupWithUsers(group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(group)
}

func (api *GroupAPI) UpdateGroup(c *BaseController) {
	id := c.Ctx.Input.Query("id")

	var group object.Group
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.groupService.UpdateGroup(id, &group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *GroupAPI) AddGroup(c *BaseController) {
	var group object.Group
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.groupService.CreateGroup(&group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *GroupAPI) DeleteGroup(c *BaseController) {
	var group object.Group
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.groupService.DeleteGroup(&group)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}
