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

package service

import (
	"context"
	"fmt"

	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/dto"
	"github.com/casdoor/casdoor/internal/model"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/util"
)

// ApplicationService defines application service interface
type ApplicationService interface {
	// CRUD operations
	GetApplication(ctx context.Context, id string) (*dto.ApplicationResponse, error)
	GetApplicationByClientID(ctx context.Context, clientID string) (*dto.ApplicationResponse, error)
	CreateApplication(ctx context.Context, req *dto.CreateApplicationRequest) (*dto.ApplicationResponse, error)
	UpdateApplication(ctx context.Context, id string, req *dto.UpdateApplicationRequest) (*dto.ApplicationResponse, error)
	DeleteApplication(ctx context.Context, id string) error
	
	// List operations
	ListApplications(ctx context.Context, req *dto.ListApplicationsRequest) (*dto.ListApplicationsResponse, error)
	GetApplicationsByOrganization(ctx context.Context, owner, organization string, page, pageSize int) (*dto.ListApplicationsResponse, error)
	
	// OAuth operations
	ValidateOAuthRequest(ctx context.Context, clientID, redirectURI, responseType string) (*dto.ApplicationResponse, error)
	ValidateClientCredentials(ctx context.Context, clientID, clientSecret string) (*dto.ApplicationResponse, error)
	GenerateToken(ctx context.Context, req *dto.TokenRequest) (*dto.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
	RevokeToken(ctx context.Context, token string) error
	
	// Permission operations
	GrantPermission(ctx context.Context, appID string, req *dto.GrantPermissionRequest) error
	RevokePermission(ctx context.Context, appID string, req *dto.RevokePermissionRequest) error
	GetPermissions(ctx context.Context, appID string) ([]*dto.PermissionInfo, error)
	
	// Batch operations
	BatchCreateApplications(ctx context.Context, req *dto.BatchCreateApplicationsRequest) (*dto.BatchOperationResponse, error)
	BatchUpdateApplications(ctx context.Context, operation *dto.BatchApplicationOperation) error
	BatchDeleteApplications(ctx context.Context, ids []string) error
	
	// Search
	SearchApplications(ctx context.Context, owner, keyword string) ([]*dto.ApplicationResponse, error)
	
	// Statistics
	GetApplicationStatistics(ctx context.Context, owner string) (*repository.ApplicationStatistics, error)
}

// applicationService implements ApplicationService
type applicationService struct {
	appRepo repository.ApplicationRepository
}

// NewApplicationService creates new application service
func NewApplicationService(appRepo repository.ApplicationRepository) ApplicationService {
	return &applicationService{appRepo: appRepo}
}

// GetApplication gets application by ID
func (s *applicationService) GetApplication(ctx context.Context, id string) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toApplicationResponse(app), nil
}

// GetApplicationByClientID gets application by client ID
func (s *applicationService) GetApplicationByClientID(ctx context.Context, clientID string) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return s.toApplicationResponse(app), nil
}

// CreateApplication creates a new application
func (s *applicationService) CreateApplication(ctx context.Context, req *dto.CreateApplicationRequest) (*dto.ApplicationResponse, error) {
	// Check if application already exists
	existingApp, _ := s.appRepo.GetByOwnerAndName(ctx, req.Owner, req.Name)
	if existingApp != nil {
		return nil, common.ErrApplicationAlreadyExists
	}
	
	app := &model.Application{
		Owner:              req.Owner,
		Name:               req.Name,
		CreatedTime:        util.GetCurrentTime(),
		DisplayName:        req.DisplayName,
		Logo:               req.Logo,
		HomepageUrl:        req.HomepageUrl,
		Description:        req.Description,
		Organization:       req.Organization,
		Cert:               req.Cert,
		EnablePassword:     req.EnablePassword,
		EnableSignUp:       req.EnableSignUp,
		EnableSigninSession: req.EnableSigninSession,
		EnableAutoSignin:   req.EnableAutoSignin,
		EnableCodeSignin:   req.EnableCodeSignin,
		EnableSamlCompress: req.EnableSamlCompress,
		EnableWebAuthn:     req.EnableWebAuthn,
		EnableLinkWithEmail: req.EnableLinkWithEmail,
		OrgChoiceMode:      req.OrgChoiceMode,
		SamlReplyUrl:       req.SamlReplyUrl,
		ClientId:           util.GenerateId(),
		ClientSecret:       util.GenerateId(),
		RedirectUris:       req.RedirectUris,
		TokenFormat:        req.TokenFormat,
		ExpireInHours:      req.ExpireInHours,
		RefreshExpireInHours: req.RefreshExpireInHours,
		SignupUrl:          req.SignupUrl,
		SigninUrl:          req.SigninUrl,
		ForgetUrl:          req.ForgetUrl,
		AffiliationUrl:     req.AffiliationUrl,
		TermsOfUse:         req.TermsOfUse,
		SignupHtml:         req.SignupHtml,
		SigninHtml:         req.SigninHtml,
		ThemeData:          req.ThemeData,
		FormCss:            req.FormCss,
		FormCssMobile:      req.FormCssMobile,
		FormOffset:         req.FormOffset,
		GrantTypes:         req.GrantTypes,
		Tags:               req.Tags,
		IsShared:           req.IsShared,
	}
	
	// Set defaults
	if app.DisplayName == "" {
		app.DisplayName = app.Name
	}
	if app.ExpireInHours == 0 {
		app.ExpireInHours = 168 // 7 days
	}
	if app.RefreshExpireInHours == 0 {
		app.RefreshExpireInHours = 720 // 30 days
	}
	
	err := s.appRepo.Create(ctx, app)
	if err != nil {
		return nil, err
	}
	
	return s.toApplicationResponse(app), nil
}

// UpdateApplication updates application
func (s *applicationService) UpdateApplication(ctx context.Context, id string, req *dto.UpdateApplicationRequest) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields
	if req.DisplayName != "" {
		app.DisplayName = req.DisplayName
	}
	if req.Logo != "" {
		app.Logo = req.Logo
	}
	if req.HomepageUrl != "" {
		app.HomepageUrl = req.HomepageUrl
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	if req.Organization != "" {
		app.Organization = req.Organization
	}
	if req.Cert != "" {
		app.Cert = req.Cert
	}
	if len(req.RedirectUris) > 0 {
		app.RedirectUris = req.RedirectUris
	}
	if req.TokenFormat != "" {
		app.TokenFormat = req.TokenFormat
	}
	if req.ExpireInHours != 0 {
		app.ExpireInHours = req.ExpireInHours
	}
	if req.RefreshExpireInHours != 0 {
		app.RefreshExpireInHours = req.RefreshExpireInHours
	}
	if req.SignupUrl != "" {
		app.SignupUrl = req.SignupUrl
	}
	if req.SigninUrl != "" {
		app.SigninUrl = req.SigninUrl
	}
	if req.ForgetUrl != "" {
		app.ForgetUrl = req.ForgetUrl
	}
	if req.AffiliationUrl != "" {
		app.AffiliationUrl = req.AffiliationUrl
	}
	if req.TermsOfUse != "" {
		app.TermsOfUse = req.TermsOfUse
	}
	if req.SignupHtml != "" {
		app.SignupHtml = req.SignupHtml
	}
	if req.SigninHtml != "" {
		app.SigninHtml = req.SigninHtml
	}
	if req.ThemeData != "" {
		app.ThemeData = req.ThemeData
	}
	if req.FormCss != "" {
		app.FormCss = req.FormCss
	}
	if req.FormCssMobile != "" {
		app.FormCssMobile = req.FormCssMobile
	}
	if req.FormOffset != 0 {
		app.FormOffset = req.FormOffset
	}
	if len(req.GrantTypes) > 0 {
		app.GrantTypes = req.GrantTypes
	}
	if len(req.Tags) > 0 {
		app.Tags = req.Tags
	}
	
	app.EnablePassword = req.EnablePassword
	app.EnableSignUp = req.EnableSignUp
	app.EnableSigninSession = req.EnableSigninSession
	app.EnableAutoSignin = req.EnableAutoSignin
	app.EnableCodeSignin = req.EnableCodeSignin
	app.EnableSamlCompress = req.EnableSamlCompress
	app.EnableWebAuthn = req.EnableWebAuthn
	app.EnableLinkWithEmail = req.EnableLinkWithEmail
	app.IsShared = req.IsShared
	
	err = s.appRepo.Update(ctx, app, nil)
	if err != nil {
		return nil, err
	}
	
	return s.toApplicationResponse(app), nil
}

// DeleteApplication deletes application
func (s *applicationService) DeleteApplication(ctx context.Context, id string) error {
	return s.appRepo.Delete(ctx, id)
}

// ListApplications lists applications
func (s *applicationService) ListApplications(ctx context.Context, req *dto.ListApplicationsRequest) (*dto.ListApplicationsResponse, error) {
	filter := repository.ApplicationFilter{
		Owner:     req.Owner,
		Field:     req.Field,
		Value:     req.Value,
		SortField: req.SortField,
		SortOrder: req.SortOrder,
	}
	
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	
	offset := (req.Page - 1) * req.PageSize
	
	apps, total, err := s.appRepo.ListWithPagination(ctx, filter, offset, req.PageSize)
	if err != nil {
		return nil, err
	}
	
	appResponses := make([]*dto.ApplicationResponse, len(apps))
	for i, app := range apps {
		appResponses[i] = s.toApplicationResponse(app)
	}
	
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}
	
	return &dto.ListApplicationsResponse{
		Applications: appResponses,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
		TotalPages:   totalPages,
	}, nil
}

// GetApplicationsByOrganization gets applications by organization
func (s *applicationService) GetApplicationsByOrganization(ctx context.Context, owner, organization string, page, pageSize int) (*dto.ListApplicationsResponse, error) {
	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}
	
	offset := (page - 1) * pageSize
	
	apps, total, err := s.appRepo.GetByOrganizationWithPagination(ctx, owner, organization, offset, pageSize)
	if err != nil {
		return nil, err
	}
	
	appResponses := make([]*dto.ApplicationResponse, len(apps))
	for i, app := range apps {
		appResponses[i] = s.toApplicationResponse(app)
	}
	
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return &dto.ListApplicationsResponse{
		Applications: appResponses,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}, nil
}

// ValidateOAuthRequest validates OAuth authorization request
func (s *applicationService) ValidateOAuthRequest(ctx context.Context, clientID, redirectURI, responseType string) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	
	// Validate redirect URI
	err = s.appRepo.ValidateRedirectURI(ctx, app.GetId(), redirectURI)
	if err != nil {
		return nil, err
	}
	
	// Validate response type
	validResponseType := false
	for _, rt := range app.GrantTypes {
		if rt == responseType {
			validResponseType = true
			break
		}
	}
	if !validResponseType && responseType != "code" {
		return nil, fmt.Errorf("unsupported response type: %s", responseType)
	}
	
	return s.toApplicationResponse(app), nil
}

// ValidateClientCredentials validates OAuth client credentials
func (s *applicationService) ValidateClientCredentials(ctx context.Context, clientID, clientSecret string) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.ValidateClientCredentials(ctx, clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	return s.toApplicationResponse(app), nil
}

// GenerateToken generates access token
func (s *applicationService) GenerateToken(ctx context.Context, req *dto.TokenRequest) (*dto.TokenResponse, error) {
	// Validate client credentials
	app, err := s.appRepo.ValidateClientCredentials(ctx, req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}
	
	// Generate tokens
	accessToken := util.GenerateId()
	refreshToken := util.GenerateId()
	
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    app.ExpireInHours * 3600,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken refreshes access token
func (s *applicationService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	// In real implementation, validate refresh token and generate new tokens
	accessToken := util.GenerateId()
	newRefreshToken := util.GenerateId()
	
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    168 * 3600, // 7 days
		RefreshToken: newRefreshToken,
	}, nil
}

// RevokeToken revokes access token
func (s *applicationService) RevokeToken(ctx context.Context, token string) error {
	// In real implementation, add token to blacklist
	return nil
}

// GrantPermission grants permission to application
func (s *applicationService) GrantPermission(ctx context.Context, appID string, req *dto.GrantPermissionRequest) error {
	// In real implementation, store permission grant
	return nil
}

// RevokePermission revokes permission from application
func (s *applicationService) RevokePermission(ctx context.Context, appID string, req *dto.RevokePermissionRequest) error {
	// In real implementation, remove permission grant
	return nil
}

// GetPermissions gets application permissions
func (s *applicationService) GetPermissions(ctx context.Context, appID string) ([]*dto.PermissionInfo, error) {
	// In real implementation, retrieve permissions from database
	return []*dto.PermissionInfo{}, nil
}

// BatchCreateApplications creates multiple applications
func (s *applicationService) BatchCreateApplications(ctx context.Context, req *dto.BatchCreateApplicationsRequest) (*dto.BatchOperationResponse, error) {
	resp := &dto.BatchOperationResponse{
		Total: len(req.Applications),
	}
	
	apps := make([]*model.Application, len(req.Applications))
	for i, reqApp := range req.Applications {
		// Check if exists
		existingApp, _ := s.appRepo.GetByOwnerAndName(ctx, reqApp.Owner, reqApp.Name)
		if existingApp != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, fmt.Sprintf("application %s already exists", reqApp.Name))
			continue
		}
		
		apps[i] = &model.Application{
			Owner:              reqApp.Owner,
			Name:               reqApp.Name,
			CreatedTime:        util.GetCurrentTime(),
			DisplayName:        reqApp.DisplayName,
			Logo:               reqApp.Logo,
			HomepageUrl:        reqApp.HomepageUrl,
			Description:        reqApp.Description,
			Organization:       reqApp.Organization,
			ClientId:           util.GenerateId(),
			ClientSecret:       util.GenerateId(),
			RedirectUris:       reqApp.RedirectUris,
			ExpireInHours:      reqApp.ExpireInHours,
			RefreshExpireInHours: reqApp.RefreshExpireInHours,
			GrantTypes:         reqApp.GrantTypes,
			Tags:               reqApp.Tags,
			IsShared:           reqApp.IsShared,
		}
		
		if apps[i].DisplayName == "" {
			apps[i].DisplayName = apps[i].Name
		}
		if apps[i].ExpireInHours == 0 {
			apps[i].ExpireInHours = 168
		}
		if apps[i].RefreshExpireInHours == 0 {
			apps[i].RefreshExpireInHours = 720
		}
	}
	
	err := s.appRepo.BatchCreate(ctx, apps)
	if err != nil {
		resp.Failed = len(req.Applications)
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}
	
	resp.Success = len(req.Applications)
	return resp, nil
}

// BatchUpdateApplications updates multiple applications
func (s *applicationService) BatchUpdateApplications(ctx context.Context, operation *dto.BatchApplicationOperation) error {
	apps := make([]*model.Application, len(operation.ApplicationIds))
	for i, id := range operation.ApplicationIds {
		app, err := s.appRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		
		switch operation.Operation {
		case "enable_signup":
			app.EnableSignUp = true
		case "disable_signup":
			app.EnableSignUp = false
		case "enable_password":
			app.EnablePassword = true
		case "disable_password":
			app.EnablePassword = false
		case "make_shared":
			app.IsShared = true
		case "make_private":
			app.IsShared = false
		}
		
		apps[i] = app
	}
	
	return s.appRepo.BatchUpdate(ctx, apps, nil)
}

// BatchDeleteApplications deletes multiple applications
func (s *applicationService) BatchDeleteApplications(ctx context.Context, ids []string) error {
	return s.appRepo.BatchDelete(ctx, ids)
}

// SearchApplications searches applications
func (s *applicationService) SearchApplications(ctx context.Context, owner, keyword string) ([]*dto.ApplicationResponse, error) {
	fields := []string{"name", "display_name", "description"}
	apps, err := s.appRepo.Search(ctx, owner, keyword, fields)
	if err != nil {
		return nil, err
	}
	
	responses := make([]*dto.ApplicationResponse, len(apps))
	for i, app := range apps {
		responses[i] = s.toApplicationResponse(app)
	}
	
	return responses, nil
}

// GetApplicationStatistics gets application statistics
func (s *applicationService) GetApplicationStatistics(ctx context.Context, owner string) (*repository.ApplicationStatistics, error) {
	return s.appRepo.GetStatistics(ctx, owner)
}

// Helper functions

func (s *applicationService) toApplicationResponse(app *model.Application) *dto.ApplicationResponse {
	return &dto.ApplicationResponse{
		Owner:              app.Owner,
		Name:               app.Name,
		CreatedTime:        app.CreatedTime,
		DisplayName:        app.DisplayName,
		Logo:               app.Logo,
		HomepageUrl:        app.HomepageUrl,
		Description:        app.Description,
		Organization:       app.Organization,
		Cert:               app.Cert,
		EnablePassword:     app.EnablePassword,
		EnableSignUp:       app.EnableSignUp,
		EnableSigninSession: app.EnableSigninSession,
		EnableAutoSignin:   app.EnableAutoSignin,
		EnableCodeSignin:   app.EnableCodeSignin,
		EnableSamlCompress: app.EnableSamlCompress,
		EnableWebAuthn:     app.EnableWebAuthn,
		EnableLinkWithEmail: app.EnableLinkWithEmail,
		OrgChoiceMode:      app.OrgChoiceMode,
		SamlReplyUrl:       app.SamlReplyUrl,
		ClientId:           app.ClientId,
		RedirectUris:       app.RedirectUris,
		TokenFormat:        app.TokenFormat,
		ExpireInHours:      app.ExpireInHours,
		RefreshExpireInHours: app.RefreshExpireInHours,
		SignupUrl:          app.SignupUrl,
		SigninUrl:          app.SigninUrl,
		ForgetUrl:          app.ForgetUrl,
		AffiliationUrl:     app.AffiliationUrl,
		TermsOfUse:         app.TermsOfUse,
		SignupHtml:         app.SignupHtml,
		SigninHtml:         app.SigninHtml,
		ThemeData:          app.ThemeData,
		FormCss:            app.FormCss,
		FormCssMobile:      app.FormCssMobile,
		FormOffset:         app.FormOffset,
		GrantTypes:         app.GrantTypes,
		Tags:               app.Tags,
		IsShared:           app.IsShared,
	}
}
