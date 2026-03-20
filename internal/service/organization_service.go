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

// OrganizationService defines organization service interface
type OrganizationService interface {
	// CRUD operations
	GetOrganization(ctx context.Context, id string) (*dto.OrganizationResponse, error)
	GetOrganizationByName(ctx context.Context, owner, name string) (*dto.OrganizationResponse, error)
	CreateOrganization(ctx context.Context, req *dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error)
	UpdateOrganization(ctx context.Context, id string, req *dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error)
	DeleteOrganization(ctx context.Context, id string) error

	// List operations
	ListOrganizations(ctx context.Context, req *dto.ListOrganizationsRequest) (*dto.ListOrganizationsResponse, error)
	GetOrganizationsByOwner(ctx context.Context, owner string, page, pageSize int) (*dto.ListOrganizationsResponse, error)

	// Hierarchy operations
	GetOrganizationHierarchy(ctx context.Context, orgID string) (*dto.OrganizationHierarchyResponse, error)
	GetOrganizationTree(ctx context.Context, owner string) ([]*dto.OrganizationResponse, error)
	GetOrganizationChildren(ctx context.Context, orgID string) ([]*dto.OrganizationResponse, error)
	GetOrganizationDescendants(ctx context.Context, orgID string) ([]*dto.OrganizationResponse, error)
	GetOrganizationAncestors(ctx context.Context, orgID string) ([]*dto.OrganizationResponse, error)
	MoveOrganization(ctx context.Context, orgID, newParentID string) error

	// Batch operations
	BatchCreateOrganizations(ctx context.Context, req *dto.BatchCreateOrganizationsRequest) (*dto.BatchOperationResponse, error)
	BatchUpdateOrganizations(ctx context.Context, operation *dto.BatchOrganizationOperation) error
	BatchDeleteOrganizations(ctx context.Context, ids []string) error

	// Search
	SearchOrganizations(ctx context.Context, owner, keyword string) ([]*dto.OrganizationResponse, error)

	// Statistics
	GetOrganizationStatistics(ctx context.Context, owner string) (*repository.OrganizationStatistics, error)

	// Applications
	GetOrganizationApplications(ctx context.Context, owner, orgName string) ([]*dto.ApplicationResponse, error)
}

// organizationService implements OrganizationService
type organizationService struct {
	orgRepo repository.OrganizationRepository
}

// NewOrganizationService creates new organization service
func NewOrganizationService(orgRepo repository.OrganizationRepository) OrganizationService {
	return &organizationService{orgRepo: orgRepo}
}

// GetOrganization gets organization by ID
func (s *organizationService) GetOrganization(ctx context.Context, id string) (*dto.OrganizationResponse, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toOrganizationResponse(org), nil
}

// GetOrganizationByName gets organization by owner and name
func (s *organizationService) GetOrganizationByName(ctx context.Context, owner, name string) (*dto.OrganizationResponse, error) {
	org, err := s.orgRepo.GetByOwnerAndName(ctx, owner, name)
	if err != nil {
		return nil, err
	}
	return s.toOrganizationResponse(org), nil
}

// CreateOrganization creates a new organization
func (s *organizationService) CreateOrganization(ctx context.Context, req *dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error) {
	// Check if organization already exists
	existingOrg, _ := s.orgRepo.GetByOwnerAndName(ctx, req.Owner, req.Name)
	if existingOrg != nil {
		return nil, common.ErrOrganizationAlreadyExists
	}

	// Validate parent if provided
	if req.ParentID != "" {
		_, err := s.orgRepo.GetByID(ctx, req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent organization not found: %w", err)
		}
	}

	org := &model.Organization{
		Owner:              req.Owner,
		Name:               req.Name,
		CreatedTime:        util.GetCurrentTime(),
		DisplayName:        req.DisplayName,
		WebsiteUrl:         req.WebsiteUrl,
		Favicon:            req.Favicon,
		Logo:               req.Logo,
		LogoDark:           req.LogoDark,
		HeaderHtml:         req.HeaderHtml,
		FooterHtml:         req.FooterHtml,
		SigninHtml:         req.SigninHtml,
		SignupHtml:         req.SignupHtml,
		ForgetUrl:          req.ForgetUrl,
		AffiliationUrl:     req.AffiliationUrl,
		TermsOfUse:         req.TermsOfUse,
		SignupUrl:          req.SignupUrl,
		SigninUrl:          req.SigninUrl,
		ClientId:           util.GenerateId(),
		ClientSecret:       util.GenerateId(),
		DefaultAvatar:      req.DefaultAvatar,
		DefaultApplication: req.DefaultApplication,
		Tags:               req.Tags,
		Languages:          req.Languages,
		MasterPassword:     req.MasterPassword,
		EnableSoftDeletion: req.EnableSoftDeletion,
		IsProfilePublic:    req.IsProfilePublic,
		DefaultPassword:    req.DefaultPassword,
		PasswordType:       req.PasswordType,
		PasswordSalt:       req.PasswordSalt,
		PasswordOptions:    req.PasswordOptions,
		CountryCodes:       req.CountryCodes,
		PhonePrefix:        req.PhonePrefix,
		InitScore:          req.InitScore,
		EnableSamlC14n10:   req.EnableSamlC14n10,
		SamlReplyLimit:     req.SamlReplyLimit,
		ParentID:           req.ParentID,
		UseEmailAsUsername: req.UseEmailAsUsername,
		EnableTour:         req.EnableTour,
		DisableSignin:      req.DisableSignin,
		IpWhitelist:        req.IpWhitelist,
		PasswordExpireDays: req.PasswordExpireDays,
	}

	// Set defaults
	if org.DisplayName == "" {
		org.DisplayName = org.Name
	}
	if org.DefaultAvatar == "" {
		org.DefaultAvatar = "https://cdn.casbin.org/img/casbin.svg"
	}
	if org.InitScore == 0 {
		org.InitScore = 2000
	}

	err := s.orgRepo.Create(ctx, org)
	if err != nil {
		return nil, err
	}

	return s.toOrganizationResponse(org), nil
}

// UpdateOrganization updates organization
func (s *organizationService) UpdateOrganization(ctx context.Context, id string, req *dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.DisplayName != "" {
		org.DisplayName = req.DisplayName
	}
	if req.WebsiteUrl != "" {
		org.WebsiteUrl = req.WebsiteUrl
	}
	if req.Favicon != "" {
		org.Favicon = req.Favicon
	}
	if req.Logo != "" {
		org.Logo = req.Logo
	}
	if req.LogoDark != "" {
		org.LogoDark = req.LogoDark
	}
	if req.HeaderHtml != "" {
		org.HeaderHtml = req.HeaderHtml
	}
	if req.FooterHtml != "" {
		org.FooterHtml = req.FooterHtml
	}
	if req.SigninHtml != "" {
		org.SigninHtml = req.SigninHtml
	}
	if req.SignupHtml != "" {
		org.SignupHtml = req.SignupHtml
	}
	if req.ForgetUrl != "" {
		org.ForgetUrl = req.ForgetUrl
	}
	if req.AffiliationUrl != "" {
		org.AffiliationUrl = req.AffiliationUrl
	}
	if req.TermsOfUse != "" {
		org.TermsOfUse = req.TermsOfUse
	}
	if req.SignupUrl != "" {
		org.SignupUrl = req.SignupUrl
	}
	if req.SigninUrl != "" {
		org.SigninUrl = req.SigninUrl
	}
	if req.DefaultAvatar != "" {
		org.DefaultAvatar = req.DefaultAvatar
	}
	if req.DefaultApplication != "" {
		org.DefaultApplication = req.DefaultApplication
	}
	if len(req.Tags) > 0 {
		org.Tags = req.Tags
	}
	if len(req.Languages) > 0 {
		org.Languages = req.Languages
	}
	if req.MasterPassword != "" {
		org.MasterPassword = req.MasterPassword
	}
	if req.DefaultPassword != "" {
		org.DefaultPassword = req.DefaultPassword
	}
	if req.PasswordType != "" {
		org.PasswordType = req.PasswordType
	}
	if req.PasswordSalt != "" {
		org.PasswordSalt = req.PasswordSalt
	}
	if len(req.PasswordOptions) > 0 {
		org.PasswordOptions = req.PasswordOptions
	}
	if len(req.CountryCodes) > 0 {
		org.CountryCodes = req.CountryCodes
	}
	if req.PhonePrefix != "" {
		org.PhonePrefix = req.PhonePrefix
	}
	if req.InitScore != 0 {
		org.InitScore = req.InitScore
	}
	if req.SamlReplyLimit != 0 {
		org.SamlReplyLimit = req.SamlReplyLimit
	}

	org.EnableSoftDeletion = req.EnableSoftDeletion
	org.IsProfilePublic = req.IsProfilePublic
	org.EnableSamlC14n10 = req.EnableSamlC14n10
	org.UseEmailAsUsername = req.UseEmailAsUsername
	org.EnableTour = req.EnableTour
	org.DisableSignin = req.DisableSignin

	err = s.orgRepo.Update(ctx, org, nil)
	if err != nil {
		return nil, err
	}

	return s.toOrganizationResponse(org), nil
}

// DeleteOrganization deletes organization
func (s *organizationService) DeleteOrganization(ctx context.Context, id string) error {
	return s.orgRepo.Delete(ctx, id)
}

// ListOrganizations lists organizations
func (s *organizationService) ListOrganizations(ctx context.Context, req *dto.ListOrganizationsRequest) (*dto.ListOrganizationsResponse, error) {
	filter := repository.OrganizationFilter{
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

	orgs, total, err := s.orgRepo.ListWithPagination(ctx, filter, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	orgResponses := make([]*dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		orgResponses[i] = s.toOrganizationResponse(org)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.ListOrganizationsResponse{
		Organizations: orgResponses,
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
		TotalPages:    totalPages,
	}, nil
}

// GetOrganizationsByOwner gets organizations by owner
func (s *organizationService) GetOrganizationsByOwner(ctx context.Context, owner string, page, pageSize int) (*dto.ListOrganizationsResponse, error) {
	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	filter := repository.OrganizationFilter{
		Owner: owner,
	}

	orgs, total, err := s.orgRepo.ListWithPagination(ctx, filter, offset, pageSize)
	if err != nil {
		return nil, err
	}

	orgResponses := make([]*dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		orgResponses[i] = s.toOrganizationResponse(org)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListOrganizationsResponse{
		Organizations: orgResponses,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}, nil
}

// GetOrganizationHierarchy gets organization hierarchy
func (s *organizationService) GetOrganizationHierarchy(ctx context.Context, orgID string) (*dto.OrganizationHierarchyResponse, error) {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Get ancestors
	ancestors, err := s.orgRepo.GetAncestors(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Get children
	children, err := s.orgRepo.GetChildren(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Get descendants (up to 3 levels)
	descendants, err := s.orgRepo.GetDescendants(ctx, orgID, 3)
	if err != nil {
		return nil, err
	}

	ancestorResponses := make([]*dto.OrganizationResponse, len(ancestors))
	for i, ancestor := range ancestors {
		ancestorResponses[i] = s.toOrganizationResponse(ancestor)
	}

	childrenResponses := make([]*dto.OrganizationResponse, len(children))
	for i, child := range children {
		childrenResponses[i] = s.toOrganizationResponse(child)
	}

	descendantResponses := make([]*dto.OrganizationResponse, len(descendants))
	for i, descendant := range descendants {
		descendantResponses[i] = s.toOrganizationResponse(descendant)
	}

	return &dto.OrganizationHierarchyResponse{
		Organization: s.toOrganizationResponse(org),
		Ancestors:    ancestorResponses,
		Children:     childrenResponses,
		Descendants:  descendantResponses,
		Level:        len(ancestors),
	}, nil
}

// GetOrganizationTree gets organization tree for owner
func (s *organizationService) GetOrganizationTree(ctx context.Context, owner string) ([]*dto.OrganizationResponse, error) {
	// Get all organizations for owner
	orgs, err := s.orgRepo.ListByOwner(ctx, owner, 0, 1000)
	if err != nil {
		return nil, err
	}

	// Filter root organizations (no parent)
	var rootOrgs []*model.Organization
	for _, org := range orgs {
		if org.ParentID == "" {
			rootOrgs = append(rootOrgs, org)
		}
	}

	responses := make([]*dto.OrganizationResponse, len(rootOrgs))
	for i, org := range rootOrgs {
		responses[i] = s.toOrganizationResponse(org)
	}

	return responses, nil
}

// GetOrganizationChildren gets organization children
func (s *organizationService) GetOrganizationChildren(ctx context.Context, orgID string) ([]*dto.OrganizationResponse, error) {
	children, err := s.orgRepo.GetChildren(ctx, orgID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrganizationResponse, len(children))
	for i, child := range children {
		responses[i] = s.toOrganizationResponse(child)
	}

	return responses, nil
}

// GetOrganizationDescendants gets organization descendants
func (s *organizationService) GetOrganizationDescendants(ctx context.Context, orgID string) ([]*dto.OrganizationResponse, error) {
	descendants, err := s.orgRepo.GetDescendants(ctx, orgID, 10)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrganizationResponse, len(descendants))
	for i, descendant := range descendants {
		responses[i] = s.toOrganizationResponse(descendant)
	}

	return responses, nil
}

// GetOrganizationAncestors gets organization ancestors
func (s *organizationService) GetOrganizationAncestors(ctx context.Context, orgID string) ([]*dto.OrganizationResponse, error) {
	ancestors, err := s.orgRepo.GetAncestors(ctx, orgID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrganizationResponse, len(ancestors))
	for i, ancestor := range ancestors {
		responses[i] = s.toOrganizationResponse(ancestor)
	}

	return responses, nil
}

// MoveOrganization moves organization to new parent
func (s *organizationService) MoveOrganization(ctx context.Context, orgID, newParentID string) error {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return err
	}

	// Validate new parent
	if newParentID != "" {
		_, err := s.orgRepo.GetByID(ctx, newParentID)
		if err != nil {
			return fmt.Errorf("new parent organization not found: %w", err)
		}

		// Check for circular reference
		descendants, err := s.orgRepo.GetDescendants(ctx, orgID, 10)
		if err != nil {
			return err
		}

		for _, descendant := range descendants {
			if descendant.Name == newParentID {
				return fmt.Errorf("cannot move organization to its own descendant")
			}
		}
	}

	org.ParentID = newParentID
	return s.orgRepo.Update(ctx, org, []string{"parent_id"})
}

// BatchCreateOrganizations creates multiple organizations
func (s *organizationService) BatchCreateOrganizations(ctx context.Context, req *dto.BatchCreateOrganizationsRequest) (*dto.BatchOperationResponse, error) {
	resp := &dto.BatchOperationResponse{
		Total: len(req.Organizations),
	}

	orgs := make([]*model.Organization, len(req.Organizations))
	for i, reqOrg := range req.Organizations {
		// Check if exists
		existingOrg, _ := s.orgRepo.GetByOwnerAndName(ctx, reqOrg.Owner, reqOrg.Name)
		if existingOrg != nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, fmt.Sprintf("organization %s already exists", reqOrg.Name))
			continue
		}

		orgs[i] = &model.Organization{
			Owner:              reqOrg.Owner,
			Name:               reqOrg.Name,
			CreatedTime:        util.GetCurrentTime(),
			DisplayName:        reqOrg.DisplayName,
			WebsiteUrl:         reqOrg.WebsiteUrl,
			Favicon:            reqOrg.Favicon,
			Logo:               reqOrg.Logo,
			LogoDark:           reqOrg.LogoDark,
			DefaultAvatar:      reqOrg.DefaultAvatar,
			DefaultApplication: reqOrg.DefaultApplication,
			Tags:               reqOrg.Tags,
			Languages:          reqOrg.Languages,
			MasterPassword:     reqOrg.MasterPassword,
			EnableSoftDeletion: reqOrg.EnableSoftDeletion,
			IsProfilePublic:    reqOrg.IsProfilePublic,
			DefaultPassword:    reqOrg.DefaultPassword,
			PasswordType:       reqOrg.PasswordType,
			PasswordSalt:       reqOrg.PasswordSalt,
			PasswordOptions:    reqOrg.PasswordOptions,
			CountryCodes:       reqOrg.CountryCodes,
			PhonePrefix:        reqOrg.PhonePrefix,
			InitScore:          reqOrg.InitScore,
			ClientId:           util.GenerateId(),
			ClientSecret:       util.GenerateId(),
			ParentID:           reqOrg.ParentID,
			HeaderHtml:         reqOrg.HeaderHtml,
			FooterHtml:         reqOrg.FooterHtml,
			SigninHtml:         reqOrg.SigninHtml,
			SignupHtml:         reqOrg.SignupHtml,
			ForgetUrl:          reqOrg.ForgetUrl,
			AffiliationUrl:     reqOrg.AffiliationUrl,
			TermsOfUse:         reqOrg.TermsOfUse,
			SignupUrl:          reqOrg.SignupUrl,
			SigninUrl:          reqOrg.SigninUrl,
			EnableSamlC14n10:   reqOrg.EnableSamlC14n10,
			SamlReplyLimit:     reqOrg.SamlReplyLimit,
			UseEmailAsUsername: reqOrg.UseEmailAsUsername,
			EnableTour:         reqOrg.EnableTour,
			DisableSignin:      reqOrg.DisableSignin,
			IpWhitelist:        reqOrg.IpWhitelist,
			PasswordExpireDays: reqOrg.PasswordExpireDays,
		}

		if orgs[i].DisplayName == "" {
			orgs[i].DisplayName = orgs[i].Name
		}
		if orgs[i].DefaultAvatar == "" {
			orgs[i].DefaultAvatar = "https://cdn.casbin.org/img/casbin.svg"
		}
		if orgs[i].InitScore == 0 {
			orgs[i].InitScore = 2000
		}
	}

	err := s.orgRepo.BatchCreate(ctx, orgs)
	if err != nil {
		resp.Failed = len(req.Organizations)
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.Success = len(req.Organizations)
	return resp, nil
}

// BatchUpdateOrganizations updates multiple organizations
func (s *organizationService) BatchUpdateOrganizations(ctx context.Context, operation *dto.BatchOrganizationOperation) error {
	orgs := make([]*model.Organization, len(operation.OrganizationIds))
	for i, id := range operation.OrganizationIds {
		org, err := s.orgRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}

		switch operation.Operation {
		case "enable_soft_deletion":
			org.EnableSoftDeletion = true
		case "disable_soft_deletion":
			org.EnableSoftDeletion = false
		case "make_public":
			org.IsProfilePublic = true
		case "make_private":
			org.IsProfilePublic = false
		}

		orgs[i] = org
	}

	return s.orgRepo.BatchUpdate(ctx, orgs, nil)
}

// BatchDeleteOrganizations deletes multiple organizations
func (s *organizationService) BatchDeleteOrganizations(ctx context.Context, ids []string) error {
	return s.orgRepo.BatchDelete(ctx, ids)
}

// SearchOrganizations searches organizations
func (s *organizationService) SearchOrganizations(ctx context.Context, owner, keyword string) ([]*dto.OrganizationResponse, error) {
	fields := []string{"name", "display_name", "website_url"}
	orgs, err := s.orgRepo.Search(ctx, keyword, fields)
	if err != nil {
		return nil, err
	}

	// Filter by owner
	var filtered []*model.Organization
	for _, org := range orgs {
		if org.Owner == owner {
			filtered = append(filtered, org)
		}
	}

	responses := make([]*dto.OrganizationResponse, len(filtered))
	for i, org := range filtered {
		responses[i] = s.toOrganizationResponse(org)
	}

	return responses, nil
}

// GetOrganizationStatistics gets organization statistics
func (s *organizationService) GetOrganizationStatistics(ctx context.Context, owner string) (*repository.OrganizationStatistics, error) {
	return s.orgRepo.GetStatistics(ctx, owner)
}

// GetOrganizationApplications gets organization applications
func (s *organizationService) GetOrganizationApplications(ctx context.Context, owner, orgName string) ([]*dto.ApplicationResponse, error) {
	// This would typically call application service
	// For now, return empty list
	return []*dto.ApplicationResponse{}, nil
}

// Helper functions

func (s *organizationService) toOrganizationResponse(org *model.Organization) *dto.OrganizationResponse {
	return &dto.OrganizationResponse{
		Owner:              org.Owner,
		Name:               org.Name,
		CreatedTime:        org.CreatedTime,
		DisplayName:        org.DisplayName,
		WebsiteUrl:         org.WebsiteUrl,
		Logo:               org.Logo,
		LogoDark:           org.LogoDark,
		Favicon:            org.Favicon,
		PasswordType:       org.PasswordType,
		PasswordOptions:    org.PasswordOptions,
		PasswordExpireDays: org.PasswordExpireDays,
		CountryCodes:       org.CountryCodes,
		DefaultAvatar:      org.DefaultAvatar,
		DefaultApplication: org.DefaultApplication,
		UserTypes:          org.UserTypes,
		Tags:               org.Tags,
		Languages:          org.Languages,
		InitScore:          org.InitScore,
		EnableSoftDeletion: org.EnableSoftDeletion,
		IsProfilePublic:    org.IsProfilePublic,
		UseEmailAsUsername: org.UseEmailAsUsername,
		EnableTour:         org.EnableTour,
		DisableSignin:      org.DisableSignin,
		MfaRememberInHours: org.MfaRememberInHours,
		HeaderHtml:         org.HeaderHtml,
		FooterHtml:         org.FooterHtml,
		SigninHtml:         org.SigninHtml,
		SignupHtml:         org.SignupHtml,
		ForgetUrl:          org.ForgetUrl,
		AffiliationUrl:     org.AffiliationUrl,
		TermsOfUse:         org.TermsOfUse,
		SignupUrl:          org.SignupUrl,
		SigninUrl:          org.SigninUrl,
		PhonePrefix:        org.PhonePrefix,
		EnableSamlC14n10:   org.EnableSamlC14n10,
		SamlReplyLimit:     org.SamlReplyLimit,
		ParentID:           org.ParentID,
		ClientId:           org.ClientId,
	}
}
