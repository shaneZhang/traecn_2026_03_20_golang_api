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

package dto

// CreateOrganizationRequest represents create organization request
type CreateOrganizationRequest struct {
	Owner              string            `json:"owner"`
	Name               string            `json:"name" binding:"required"`
	DisplayName        string            `json:"displayName"`
	WebsiteUrl         string            `json:"websiteUrl"`
	Logo               string            `json:"logo"`
	LogoDark           string            `json:"logoDark"`
	Favicon            string            `json:"favicon"`
	PasswordType       string            `json:"passwordType"`
	PasswordSalt       string            `json:"passwordSalt"`
	PasswordOptions    []string          `json:"passwordOptions"`
	PasswordExpireDays int               `json:"passwordExpireDays"`
	CountryCodes       []string          `json:"countryCodes"`
	DefaultAvatar      string            `json:"defaultAvatar"`
	DefaultApplication string            `json:"defaultApplication"`
	UserTypes          []string          `json:"userTypes"`
	Tags               []string          `json:"tags"`
	Languages          []string          `json:"languages"`
	MasterPassword     string            `json:"masterPassword"`
	DefaultPassword    string            `json:"defaultPassword"`
	IpWhitelist        string            `json:"ipWhitelist"`
	InitScore          int               `json:"initScore"`
	EnableSoftDeletion bool              `json:"enableSoftDeletion"`
	IsProfilePublic    bool              `json:"isProfilePublic"`
	UseEmailAsUsername bool              `json:"useEmailAsUsername"`
	EnableTour         bool              `json:"enableTour"`
	DisableSignin      bool              `json:"disableSignin"`
	MfaItems           []*MfaItemDTO     `json:"mfaItems"`
	MfaRememberInHours int               `json:"mfaRememberInHours"`
	AccountItems       []*AccountItemDTO `json:"accountItems"`
	ParentID           string            `json:"parentId"`
	HeaderHtml         string            `json:"headerHtml"`
	FooterHtml         string            `json:"footerHtml"`
	SigninHtml         string            `json:"signinHtml"`
	SignupHtml         string            `json:"signupHtml"`
	ForgetUrl          string            `json:"forgetUrl"`
	AffiliationUrl     string            `json:"affiliationUrl"`
	TermsOfUse         string            `json:"termsOfUse"`
	SignupUrl          string            `json:"signupUrl"`
	SigninUrl          string            `json:"signinUrl"`
	ThemeData          string            `json:"themeData"`
	PhonePrefix        string            `json:"phonePrefix"`
	EnableSamlC14n10   bool              `json:"enableSamlC14n10"`
	SamlReplyLimit     int               `json:"samlReplyLimit"`
}

// UpdateOrganizationRequest represents update organization request
type UpdateOrganizationRequest struct {
	DisplayName        string            `json:"displayName"`
	WebsiteUrl         string            `json:"websiteUrl"`
	Logo               string            `json:"logo"`
	LogoDark           string            `json:"logoDark"`
	Favicon            string            `json:"favicon"`
	PasswordType       string            `json:"passwordType"`
	PasswordSalt       string            `json:"passwordSalt"`
	PasswordOptions    []string          `json:"passwordOptions"`
	PasswordExpireDays int               `json:"passwordExpireDays"`
	CountryCodes       []string          `json:"countryCodes"`
	DefaultAvatar      string            `json:"defaultAvatar"`
	DefaultApplication string            `json:"defaultApplication"`
	UserTypes          []string          `json:"userTypes"`
	Tags               []string          `json:"tags"`
	Languages          []string          `json:"languages"`
	MasterPassword     string            `json:"masterPassword"`
	DefaultPassword    string            `json:"defaultPassword"`
	IpWhitelist        string            `json:"ipWhitelist"`
	InitScore          int               `json:"initScore"`
	EnableSoftDeletion bool              `json:"enableSoftDeletion"`
	IsProfilePublic    bool              `json:"isProfilePublic"`
	UseEmailAsUsername bool              `json:"useEmailAsUsername"`
	EnableTour         bool              `json:"enableTour"`
	DisableSignin      bool              `json:"disableSignin"`
	MfaItems           []*MfaItemDTO     `json:"mfaItems"`
	MfaRememberInHours int               `json:"mfaRememberInHours"`
	AccountItems       []*AccountItemDTO `json:"accountItems"`
	HeaderHtml         string            `json:"headerHtml"`
	FooterHtml         string            `json:"footerHtml"`
	SigninHtml         string            `json:"signinHtml"`
	SignupHtml         string            `json:"signupHtml"`
	ForgetUrl          string            `json:"forgetUrl"`
	AffiliationUrl     string            `json:"affiliationUrl"`
	TermsOfUse         string            `json:"termsOfUse"`
	SignupUrl          string            `json:"signupUrl"`
	SigninUrl          string            `json:"signinUrl"`
	ThemeData          string            `json:"themeData"`
	PhonePrefix        string            `json:"phonePrefix"`
	EnableSamlC14n10   bool              `json:"enableSamlC14n10"`
	SamlReplyLimit     int               `json:"samlReplyLimit"`
}

// OrganizationResponse represents organization response
type OrganizationResponse struct {
	Owner              string            `json:"owner"`
	Name               string            `json:"name"`
	CreatedTime        string            `json:"createdTime"`
	DisplayName        string            `json:"displayName"`
	WebsiteUrl         string            `json:"websiteUrl"`
	Logo               string            `json:"logo"`
	LogoDark           string            `json:"logoDark"`
	Favicon            string            `json:"favicon"`
	PasswordType       string            `json:"passwordType"`
	PasswordOptions    []string          `json:"passwordOptions"`
	PasswordExpireDays int               `json:"passwordExpireDays"`
	CountryCodes       []string          `json:"countryCodes"`
	DefaultAvatar      string            `json:"defaultAvatar"`
	DefaultApplication string            `json:"defaultApplication"`
	UserTypes          []string          `json:"userTypes"`
	Tags               []string          `json:"tags"`
	Languages          []string          `json:"languages"`
	ThemeData          string            `json:"themeData"`
	InitScore          int               `json:"initScore"`
	EnableSoftDeletion bool              `json:"enableSoftDeletion"`
	IsProfilePublic    bool              `json:"isProfilePublic"`
	UseEmailAsUsername bool              `json:"useEmailAsUsername"`
	EnableTour         bool              `json:"enableTour"`
	DisableSignin      bool              `json:"disableSignin"`
	MfaItems           []*MfaItemDTO     `json:"mfaItems"`
	MfaRememberInHours int               `json:"mfaRememberInHours"`
	AccountItems       []*AccountItemDTO `json:"accountItems"`
	UserCount          int64             `json:"userCount"`
	ApplicationCount   int64             `json:"applicationCount"`
	HeaderHtml         string            `json:"headerHtml"`
	FooterHtml         string            `json:"footerHtml"`
	SigninHtml         string            `json:"signinHtml"`
	SignupHtml         string            `json:"signupHtml"`
	ForgetUrl          string            `json:"forgetUrl"`
	AffiliationUrl     string            `json:"affiliationUrl"`
	TermsOfUse         string            `json:"termsOfUse"`
	SignupUrl          string            `json:"signupUrl"`
	SigninUrl          string            `json:"signinUrl"`
	PhonePrefix        string            `json:"phonePrefix"`
	EnableSamlC14n10   bool              `json:"enableSamlC14n10"`
	SamlReplyLimit     int               `json:"samlReplyLimit"`
	ParentID           string            `json:"parentId"`
	ClientId           string            `json:"clientId"`
}

// MfaItemDTO represents MFA item DTO
type MfaItemDTO struct {
	Name string `json:"name"`
	Rule string `json:"rule"`
}

// AccountItemDTO represents account item DTO
type AccountItemDTO struct {
	Name       string `json:"name"`
	Visible    bool   `json:"visible"`
	ViewRule   string `json:"viewRule"`
	ModifyRule string `json:"modifyRule"`
	Regex      string `json:"regex"`
	Tab        string `json:"tab"`
}

// ThemeDataDTO represents theme data DTO
type ThemeDataDTO struct {
	ThemeType    string `json:"themeType"`
	ColorPrimary string `json:"colorPrimary"`
	BorderRadius int    `json:"borderRadius"`
	IsCompact    bool   `json:"isCompact"`
	IsEnabled    bool   `json:"isEnabled"`
}

// ListOrganizationsRequest represents list organizations request
type ListOrganizationsRequest struct {
	Owner            string `form:"owner"`
	PageSize         int    `form:"pageSize"`
	Page             int    `form:"p"`
	Field            string `form:"field"`
	Value            string `form:"value"`
	SortField        string `form:"sortField"`
	SortOrder        string `form:"sortOrder"`
	OrganizationName string `form:"organizationName"`
}

// ListOrganizationsResponse represents list organizations response
type ListOrganizationsResponse struct {
	Organizations []*OrganizationResponse `json:"organizations"`
	Total         int64                   `json:"total"`
	Page          int                     `json:"page"`
	PageSize      int                     `json:"pageSize"`
	TotalPages    int                     `json:"totalPages"`
}

// OrganizationTreeNode represents organization tree node
type OrganizationTreeNode struct {
	Organization *OrganizationResponse   `json:"organization"`
	Children     []*OrganizationTreeNode `json:"children"`
	Level        int                     `json:"level"`
	Path         string                  `json:"path"`
}

// OrganizationHierarchyRequest represents organization hierarchy request
type OrganizationHierarchyRequest struct {
	RootId       string `form:"rootId"`
	MaxDepth     int    `form:"maxDepth"`
	IncludeUsers bool   `form:"includeUsers"`
}

// OrganizationHierarchyResponse represents organization hierarchy response
type OrganizationHierarchyResponse struct {
	Organization *OrganizationResponse   `json:"organization"`
	Ancestors    []*OrganizationResponse `json:"ancestors"`
	Children     []*OrganizationResponse `json:"children"`
	Descendants  []*OrganizationResponse `json:"descendants"`
	Level        int                     `json:"level"`
}

// OrganizationStats represents organization statistics
type OrganizationStats struct {
	OrganizationId   string  `json:"organizationId"`
	UserCount        int64   `json:"userCount"`
	ActiveUserCount  int64   `json:"activeUserCount"`
	ApplicationCount int64   `json:"applicationCount"`
	RoleCount        int64   `json:"roleCount"`
	PermissionCount  int64   `json:"permissionCount"`
	StorageUsed      int64   `json:"storageUsed"`
	Balance          float64 `json:"balance"`
}

// BatchCreateOrganizationsRequest represents batch create organizations request
type BatchCreateOrganizationsRequest struct {
	Organizations []*CreateOrganizationRequest `json:"organizations" binding:"required"`
}

// BatchOrganizationOperation represents batch organization operation
type BatchOrganizationOperation struct {
	OrganizationIds []string `json:"organizationIds" binding:"required"`
	Operation       string   `json:"operation" binding:"required"` // enable_soft_deletion, disable_soft_deletion, make_public, make_private
}

// BatchOperationResponse represents batch operation response
type BatchOperationResponse struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}
