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

// CreateApplicationRequest represents create application request
type CreateApplicationRequest struct {
	Owner                string   `json:"owner"`
	Name                 string   `json:"name" binding:"required"`
	DisplayName          string   `json:"displayName"`
	Logo                 string   `json:"logo"`
	HomepageUrl          string   `json:"homepageUrl"`
	Description          string   `json:"description"`
	Organization         string   `json:"organization"`
	Cert                 string   `json:"cert"`
	EnablePassword       bool     `json:"enablePassword"`
	EnableSignUp         bool     `json:"enableSignUp"`
	EnableSigninSession  bool     `json:"enableSigninSession"`
	EnableAutoSignin     bool     `json:"enableAutoSignin"`
	EnableCodeSignin     bool     `json:"enableCodeSignin"`
	EnableSamlCompress   bool     `json:"enableSamlCompress"`
	EnableWebAuthn       bool     `json:"enableWebAuthn"`
	EnableLinkWithEmail  bool     `json:"enableLinkWithEmail"`
	OrgChoiceMode        string   `json:"orgChoiceMode"`
	SamlReplyUrl         string   `json:"samlReplyUrl"`
	RedirectUris         []string `json:"redirectUris"`
	TokenFormat          string   `json:"tokenFormat"`
	ExpireInHours        int      `json:"expireInHours"`
	RefreshExpireInHours int      `json:"refreshExpireInHours"`
	SignupUrl            string   `json:"signupUrl"`
	SigninUrl            string   `json:"signinUrl"`
	ForgetUrl            string   `json:"forgetUrl"`
	AffiliationUrl       string   `json:"affiliationUrl"`
	TermsOfUse           string   `json:"termsOfUse"`
	SignupHtml           string   `json:"signupHtml"`
	SigninHtml           string   `json:"signinHtml"`
	ThemeData            string   `json:"themeData"`
	FormCss              string   `json:"formCss"`
	FormCssMobile        string   `json:"formCssMobile"`
	FormOffset           int      `json:"formOffset"`
	GrantTypes           []string `json:"grantTypes"`
	Tags                 []string `json:"tags"`
	IsShared             bool     `json:"isShared"`
}

// UpdateApplicationRequest represents update application request
type UpdateApplicationRequest struct {
	DisplayName          string   `json:"displayName"`
	Logo                 string   `json:"logo"`
	HomepageUrl          string   `json:"homepageUrl"`
	Description          string   `json:"description"`
	Organization         string   `json:"organization"`
	Cert                 string   `json:"cert"`
	EnablePassword       bool     `json:"enablePassword"`
	EnableSignUp         bool     `json:"enableSignUp"`
	EnableSigninSession  bool     `json:"enableSigninSession"`
	EnableAutoSignin     bool     `json:"enableAutoSignin"`
	EnableCodeSignin     bool     `json:"enableCodeSignin"`
	EnableSamlCompress   bool     `json:"enableSamlCompress"`
	EnableWebAuthn       bool     `json:"enableWebAuthn"`
	EnableLinkWithEmail  bool     `json:"enableLinkWithEmail"`
	OrgChoiceMode        string   `json:"orgChoiceMode"`
	SamlReplyUrl         string   `json:"samlReplyUrl"`
	RedirectUris         []string `json:"redirectUris"`
	TokenFormat          string   `json:"tokenFormat"`
	ExpireInHours        int      `json:"expireInHours"`
	RefreshExpireInHours int      `json:"refreshExpireInHours"`
	SignupUrl            string   `json:"signupUrl"`
	SigninUrl            string   `json:"signinUrl"`
	ForgetUrl            string   `json:"forgetUrl"`
	AffiliationUrl       string   `json:"affiliationUrl"`
	TermsOfUse           string   `json:"termsOfUse"`
	SignupHtml           string   `json:"signupHtml"`
	SigninHtml           string   `json:"signinHtml"`
	ThemeData            string   `json:"themeData"`
	FormCss              string   `json:"formCss"`
	FormCssMobile        string   `json:"formCssMobile"`
	FormOffset           int      `json:"formOffset"`
	GrantTypes           []string `json:"grantTypes"`
	Tags                 []string `json:"tags"`
	IsShared             bool     `json:"isShared"`
}

// ApplicationResponse represents application response
type ApplicationResponse struct {
	Owner                string   `json:"owner"`
	Name                 string   `json:"name"`
	CreatedTime          string   `json:"createdTime"`
	DisplayName          string   `json:"displayName"`
	Logo                 string   `json:"logo"`
	HomepageUrl          string   `json:"homepageUrl"`
	Description          string   `json:"description"`
	Organization         string   `json:"organization"`
	Cert                 string   `json:"cert"`
	EnablePassword       bool     `json:"enablePassword"`
	EnableSignUp         bool     `json:"enableSignUp"`
	EnableSigninSession  bool     `json:"enableSigninSession"`
	EnableAutoSignin     bool     `json:"enableAutoSignin"`
	EnableCodeSignin     bool     `json:"enableCodeSignin"`
	EnableSamlCompress   bool     `json:"enableSamlCompress"`
	EnableWebAuthn       bool     `json:"enableWebAuthn"`
	EnableLinkWithEmail  bool     `json:"enableLinkWithEmail"`
	OrgChoiceMode        string   `json:"orgChoiceMode"`
	SamlReplyUrl         string   `json:"samlReplyUrl"`
	ClientId             string   `json:"clientId"`
	RedirectUris         []string `json:"redirectUris"`
	TokenFormat          string   `json:"tokenFormat"`
	ExpireInHours        int      `json:"expireInHours"`
	RefreshExpireInHours int      `json:"refreshExpireInHours"`
	SignupUrl            string   `json:"signupUrl"`
	SigninUrl            string   `json:"signinUrl"`
	ForgetUrl            string   `json:"forgetUrl"`
	AffiliationUrl       string   `json:"affiliationUrl"`
	TermsOfUse           string   `json:"termsOfUse"`
	SignupHtml           string   `json:"signupHtml"`
	SigninHtml           string   `json:"signinHtml"`
	ThemeData            string   `json:"themeData"`
	FormCss              string   `json:"formCss"`
	FormCssMobile        string   `json:"formCssMobile"`
	FormOffset           int      `json:"formOffset"`
	GrantTypes           []string `json:"grantTypes"`
	Tags                 []string `json:"tags"`
	IsShared             bool     `json:"isShared"`
}

// ListApplicationsRequest represents list applications request
type ListApplicationsRequest struct {
	Owner     string `form:"owner"`
	PageSize  int    `form:"pageSize"`
	Page      int    `form:"p"`
	Field     string `form:"field"`
	Value     string `form:"value"`
	SortField string `form:"sortField"`
	SortOrder string `form:"sortOrder"`
}

// ListApplicationsResponse represents list applications response
type ListApplicationsResponse struct {
	Applications []*ApplicationResponse `json:"applications"`
	Total        int64                  `json:"total"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"pageSize"`
	TotalPages   int                    `json:"totalPages"`
}

// OAuthAuthorizeRequest represents OAuth authorize request
type OAuthAuthorizeRequest struct {
	ClientID     string `form:"client_id" binding:"required"`
	RedirectURI  string `form:"redirect_uri"`
	ResponseType string `form:"response_type" binding:"required"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
}

// OAuthAuthorizeResponse represents OAuth authorize response
type OAuthAuthorizeResponse struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
	Error string `json:"error,omitempty"`
}

// TokenRequest represents token request
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	RefreshToken string `form:"refresh_token"`
	Username     string `form:"username"`
	Password     string `form:"password"`
	Scope        string `form:"scope"`
}

// TokenResponse represents token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// GrantPermissionRequest represents grant permission request
type GrantPermissionRequest struct {
	UserID     string   `json:"userId" binding:"required"`
	Role       string   `json:"role"`
	Scopes     []string `json:"scopes"`
	ExpireDays int      `json:"expireDays"`
}

// RevokePermissionRequest represents revoke permission request
type RevokePermissionRequest struct {
	UserID string   `json:"userId" binding:"required"`
	Scopes []string `json:"scopes"`
}

// PermissionInfo represents permission info
type PermissionInfo struct {
	UserID    string   `json:"userId"`
	Role      string   `json:"role"`
	Scopes    []string `json:"scopes"`
	GrantedAt string   `json:"grantedAt"`
	ExpireAt  string   `json:"expireAt"`
}

// BatchCreateApplicationsRequest represents batch create applications request
type BatchCreateApplicationsRequest struct {
	Applications []*CreateApplicationRequest `json:"applications" binding:"required"`
}

// BatchApplicationOperation represents batch application operation
type BatchApplicationOperation struct {
	ApplicationIds []string `json:"applicationIds" binding:"required"`
	Operation      string   `json:"operation" binding:"required"` // enable_signup, disable_signup, enable_password, disable_password, make_shared, make_private
}
