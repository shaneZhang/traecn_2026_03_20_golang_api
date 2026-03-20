package api

import (
	"encoding/json"
	"fmt"

	"github.com/beego/beego/v2/core/utils/pagination"
	"github.com/casdoor/casdoor/dto"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/service"
	"github.com/casdoor/casdoor/util"
)

type ApplicationAPI struct {
	BaseAPI
	appService service.ApplicationService
}

func NewApplicationAPI(appService service.ApplicationService) *ApplicationAPI {
	return &ApplicationAPI{
		appService: appService,
	}
}

func (api *ApplicationAPI) GetApplications(c *BaseController) {
	userId := c.GetSessionUsername()
	owner := c.Ctx.Input.Query("owner")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")
	organization := c.Ctx.Input.Query("organization")

	if limit == "" || page == "" {
		var applications []*object.Application
		var err error
		if organization == "" {
			applications, err = api.appService.GetApplications(owner)
		} else {
			applications, err = api.appService.GetOrganizationApplications(owner, organization)
		}
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(api.appService.GetMaskedApplications(applications, userId))
	} else {
		limitInt := util.ParseInt(limit)
		count, err := api.appService.GetApplicationCount(owner, field, value)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
		applications, _, err := api.appService.GetPaginationApplications(owner, paginator.Offset(), limitInt, field, value, sortField, sortOrder)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		maskedApps := api.appService.GetMaskedApplications(applications, userId)
		c.ResponseOk(maskedApps, paginator.Nums())
	}
}

func (api *ApplicationAPI) GetApplication(c *BaseController) {
	userId := c.GetSessionUsername()
	id := c.Ctx.Input.Query("id")

	application, err := api.appService.GetApplication(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	if c.Ctx.Input.Query("withKey") != "" && application != nil && application.Cert != "" {
		cert, err := object.GetCert(util.GetId(application.Owner, application.Cert))
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		if cert == nil {
			cert, err = object.GetCert(util.GetId(application.Organization, application.Cert))
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
		}

		if cert != nil {
			application.CertPublicKey = cert.Certificate
		}
	}

	clientIp := util.GetClientIpFromRequest(c.Ctx.Request)
	object.CheckEntryIp(clientIp, nil, application, nil, c.GetAcceptLanguage())

	c.ResponseOk(api.appService.GetMaskedApplication(application, userId))
}

func (api *ApplicationAPI) GetUserApplication(c *BaseController) {
	userId := c.GetSessionUsername()
	id := c.Ctx.Input.Query("id")

	user, err := object.GetUser(id)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	if user == nil {
		c.ResponseError(fmt.Sprintf(c.T("general:The user: %s doesn't exist"), id))
		return
	}

	application, err := api.appService.GetApplicationByUser(user)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	if application == nil {
		c.ResponseError(fmt.Sprintf(c.T("general:The organization: %s should have one application at least"), user.Owner))
		return
	}

	c.ResponseOk(api.appService.GetMaskedApplication(application, userId))
}

func (api *ApplicationAPI) GetOrganizationApplications(c *BaseController) {
	userId := c.GetSessionUsername()
	organization := c.Ctx.Input.Query("organization")
	owner := c.Ctx.Input.Query("owner")
	limit := c.Ctx.Input.Query("pageSize")
	page := c.Ctx.Input.Query("p")
	field := c.Ctx.Input.Query("field")
	value := c.Ctx.Input.Query("value")
	sortField := c.Ctx.Input.Query("sortField")
	sortOrder := c.Ctx.Input.Query("sortOrder")

	if organization == "" {
		c.ResponseError(c.T("general:Missing parameter") + ": organization")
		return
	}

	if limit == "" || page == "" {
		applications, err := api.appService.GetOrganizationApplications(owner, organization)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		applications, err = api.appService.GetAllowedApplications(applications, userId, c.GetAcceptLanguage())
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(api.appService.GetMaskedApplications(applications, userId))
	} else {
		limitInt := util.ParseInt(limit)

		count, err := api.appService.GetOrganizationApplicationCount(owner, organization, field, value)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		paginator := pagination.NewPaginator(c.Ctx.Request, limitInt, count)
		applications, _, err := api.appService.GetPaginationOrganizationApplications(owner, organization, paginator.Offset(), limitInt, field, value, sortField, sortOrder)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		applications, err = api.appService.GetAllowedApplications(applications, userId, c.GetAcceptLanguage())
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		applications = api.appService.GetMaskedApplications(applications, userId)
		c.ResponseOk(applications, paginator.Nums())
	}
}

func (api *ApplicationAPI) UpdateApplication(c *BaseController) {
	id := c.Ctx.Input.Query("id")

	var application object.Application
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &application)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.appService.UpdateApplication(id, &application, c.IsGlobalAdmin(), c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *ApplicationAPI) AddApplication(c *BaseController) {
	var application object.Application
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &application)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.appService.CreateApplication(&application, c.GetAcceptLanguage())
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

func (api *ApplicationAPI) DeleteApplication(c *BaseController) {
	var application object.Application
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &application)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	affected, err := api.appService.DeleteApplication(&application)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.Data["json"] = wrapActionResponse(affected)
	c.ServeJSON()
}

type MfaAPI struct {
	BaseAPI
	mfaService service.MfaService
}

func NewMfaAPI(mfaService service.MfaService) *MfaAPI {
	return &MfaAPI{
		mfaService: mfaService,
	}
}

func (api *MfaAPI) MfaSetupInitiate(c *BaseController) {
	owner := c.Ctx.Request.Form.Get("owner")
	name := c.Ctx.Request.Form.Get("name")
	mfaType := c.Ctx.Request.Form.Get("mfaType")

	resp, err := api.mfaService.SetupInitiate(owner, name, mfaType)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(resp)
}

func (api *MfaAPI) MfaSetupVerify(c *BaseController) {
	mfaType := c.Ctx.Request.Form.Get("mfaType")
	passcode := c.Ctx.Request.Form.Get("passcode")
	secret := c.Ctx.Request.Form.Get("secret")
	dest := c.Ctx.Request.Form.Get("dest")
	countryCode := c.Ctx.Request.Form.Get("countryCode")

	err := api.mfaService.SetupVerify(mfaType, passcode, secret, dest, countryCode)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk("OK")
}

func (api *MfaAPI) MfaSetupEnable(c *BaseController) {
	owner := c.Ctx.Request.Form.Get("owner")
	name := c.Ctx.Request.Form.Get("name")
	mfaType := c.Ctx.Request.Form.Get("mfaType")
	secret := c.Ctx.Request.Form.Get("secret")
	dest := c.Ctx.Request.Form.Get("dest")
	countryCode := c.Ctx.Request.Form.Get("countryCode")
	recoveryCodes := c.Ctx.Request.Form.Get("recoveryCodes")

	err := api.mfaService.SetupEnable(owner, name, mfaType, secret, dest, countryCode, recoveryCodes)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk("OK")
}

func (api *MfaAPI) DeleteMfa(c *BaseController) {
	owner := c.Ctx.Request.Form.Get("owner")
	name := c.Ctx.Request.Form.Get("name")

	mfaProps, err := api.mfaService.DeleteMfa(owner, name)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(mfaProps)
}

func (api *MfaAPI) SetPreferredMfa(c *BaseController) {
	mfaType := c.Ctx.Request.Form.Get("mfaType")
	owner := c.Ctx.Request.Form.Get("owner")
	name := c.Ctx.Request.Form.Get("name")

	mfaProps, err := api.mfaService.SetPreferredMfa(owner, name, mfaType)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(mfaProps)
}

type OAuthAPI struct {
	BaseAPI
	oauthService service.OAuthService
	appService   service.ApplicationService
}

func NewOAuthAPI(oauthService service.OAuthService, appService service.ApplicationService) *OAuthAPI {
	return &OAuthAPI{
		oauthService: oauthService,
		appService:   appService,
	}
}

func (api *OAuthAPI) HandleGrant(c *BaseController) {
	var req dto.OAuthGrantRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	err = api.oauthService.HandleOAuthGrant(&req)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk()
}

func (api *OAuthAPI) HandleToken(c *BaseController) {
	var req dto.OAuthTokenRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	resp, err := api.oauthService.HandleOAuthToken(&req)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(resp)
}
