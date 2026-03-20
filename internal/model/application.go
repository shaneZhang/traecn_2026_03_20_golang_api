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

package model

// Application represents an OAuth application
type Application struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName                  string          `xorm:"varchar(100)" json:"displayName"`
	Category                     string          `xorm:"varchar(20)" json:"category"`
	Type                         string          `xorm:"varchar(20)" json:"type"`
	Scopes                       []*ScopeItem    `xorm:"mediumtext" json:"scopes"`
	Logo                         string          `xorm:"varchar(200)" json:"logo"`
	Title                        string          `xorm:"varchar(100)" json:"title"`
	Favicon                      string          `xorm:"varchar(200)" json:"favicon"`
	Order                        int             `json:"order"`
	HomepageUrl                  string          `xorm:"varchar(100)" json:"homepageUrl"`
	Description                  string          `xorm:"varchar(100)" json:"description"`
	Organization                 string          `xorm:"varchar(100)" json:"organization"`
	Cert                         string          `xorm:"varchar(100)" json:"cert"`
	DefaultGroup                 string          `xorm:"varchar(100)" json:"defaultGroup"`
	HeaderHtml                   string          `xorm:"mediumtext" json:"headerHtml"`
	EnablePassword               bool            `json:"enablePassword"`
	EnableSignUp                 bool            `json:"enableSignUp"`
	DisableSignin                bool            `json:"disableSignin"`
	EnableSigninSession          bool            `json:"enableSigninSession"`
	EnableAutoSignin             bool            `json:"enableAutoSignin"`
	EnableCodeSignin             bool            `json:"enableCodeSignin"`
	EnableExclusiveSignin        bool            `json:"enableExclusiveSignin"`
	EnableSamlCompress           bool            `json:"enableSamlCompress"`
	EnableSamlC14n10             bool            `json:"enableSamlC14n10"`
	EnableSamlPostBinding        bool            `json:"enableSamlPostBinding"`
	DisableSamlAttributes        bool            `json:"disableSamlAttributes"`
	EnableSamlAssertionSignature bool            `json:"enableSamlAssertionSignature"`
	UseEmailAsSamlNameId         bool            `json:"useEmailAsSamlNameId"`
	EnableWebAuthn               bool            `json:"enableWebAuthn"`
	EnableLinkWithEmail          bool            `json:"enableLinkWithEmail"`
	OrgChoiceMode                string          `json:"orgChoiceMode"`
	SamlReplyUrl                 string          `xorm:"varchar(500)" json:"samlReplyUrl"`
	Providers                    []*ProviderItem `xorm:"mediumtext" json:"providers"`
	SigninMethods                []*SigninMethod `xorm:"varchar(2000)" json:"signinMethods"`
	SignupItems                  []*SignupItem   `xorm:"varchar(3000)" json:"signupItems"`
	SigninItems                  []*SigninItem   `xorm:"mediumtext" json:"signinItems"`
	GrantTypes                   []string        `xorm:"varchar(1000)" json:"grantTypes"`
	OrganizationObj              *Organization   `xorm:"-" json:"organizationObj"`
	CertPublicKey                string          `xorm:"-" json:"certPublicKey"`
	Tags                         []string        `xorm:"mediumtext" json:"tags"`
	SamlAttributes               []*SamlItem     `xorm:"varchar(1000)" json:"samlAttributes"`
	SamlHashAlgorithm            string          `xorm:"varchar(20)" json:"samlHashAlgorithm"`
	IsShared                     bool            `json:"isShared"`
	IpRestriction                string          `json:"ipRestriction"`

	// OAuth Configuration
	ClientId                string     `xorm:"varchar(100)" json:"clientId"`
	ClientSecret            string     `xorm:"varchar(100)" json:"clientSecret"`
	ClientCert              string     `xorm:"varchar(100)" json:"clientCert"`
	RedirectUris            []string   `xorm:"varchar(1000)" json:"redirectUris"`
	ForcedRedirectOrigin    string     `xorm:"varchar(100)" json:"forcedRedirectOrigin"`
	TokenFormat             string     `xorm:"varchar(100)" json:"tokenFormat"`
	TokenSigningMethod      string     `xorm:"varchar(100)" json:"tokenSigningMethod"`
	TokenFields             []string   `xorm:"varchar(1000)" json:"tokenFields"`
	TokenAttributes         []*JwtItem `xorm:"mediumtext" json:"tokenAttributes"`
	ExpireInHours           float64    `json:"expireInHours"`
	RefreshExpireInHours    float64    `json:"refreshExpireInHours"`
	CookieExpireInHours     int64      `json:"cookieExpireInHours"`
	SignupUrl               string     `xorm:"varchar(200)" json:"signupUrl"`
	SigninUrl               string     `xorm:"varchar(200)" json:"signinUrl"`
	ForgetUrl               string     `xorm:"varchar(200)" json:"forgetUrl"`
	AffiliationUrl          string     `xorm:"varchar(100)" json:"affiliationUrl"`
	IpWhitelist             string     `xorm:"varchar(200)" json:"ipWhitelist"`
	TermsOfUse              string     `xorm:"varchar(200)" json:"termsOfUse"`
	SignupHtml              string     `xorm:"mediumtext" json:"signupHtml"`
	SigninHtml              string     `xorm:"mediumtext" json:"signinHtml"`
	ThemeData               *ThemeData `xorm:"json" json:"themeData"`
	FooterHtml              string     `xorm:"mediumtext" json:"footerHtml"`
	FormCss                 string     `xorm:"text" json:"formCss"`
	FormCssMobile           string     `xorm:"text" json:"formCssMobile"`
	FormOffset              int        `json:"formOffset"`
	FormSideHtml            string     `xorm:"mediumtext" json:"formSideHtml"`
	FormBackgroundUrl       string     `xorm:"varchar(200)" json:"formBackgroundUrl"`
	FormBackgroundUrlMobile string     `xorm:"varchar(200)" json:"formBackgroundUrlMobile"`

	FailedSigninLimit      int `json:"failedSigninLimit"`
	FailedSigninFrozenTime int `json:"failedSigninFrozenTime"`
	CodeResendTimeout      int `json:"codeResendTimeout"`

	CustomScopes []*ScopeDescription `xorm:"mediumtext" json:"customScopes"`

	// Reverse proxy fields
	Domain       string   `xorm:"varchar(100)" json:"domain"`
	OtherDomains []string `xorm:"varchar(1000)" json:"otherDomains"`
	UpstreamHost string   `xorm:"varchar(100)" json:"upstreamHost"`
	SslMode      string   `xorm:"varchar(100)" json:"sslMode"`
	SslCert      string   `xorm:"varchar(100)" json:"sslCert"`

	CertObj *Cert `xorm:"-"`
}

// ScopeItem represents OAuth scope
type ScopeItem struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	Tools       []string `json:"tools"`
}

// SigninMethod represents signin method configuration
type SigninMethod struct {
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`
	Rule        string `json:"rule"`
}

// SignupItem represents signup form field configuration
type SignupItem struct {
	Name        string   `json:"name"`
	Visible     bool     `json:"visible"`
	Required    bool     `json:"required"`
	Prompted    bool     `json:"prompted"`
	Type        string   `json:"type"`
	CustomCss   string   `json:"customCss"`
	Label       string   `json:"label"`
	Placeholder string   `json:"placeholder"`
	Options     []string `json:"options"`
	Regex       string   `json:"regex"`
	Rule        string   `json:"rule"`
}

// SigninItem represents signin form field configuration
type SigninItem struct {
	Name        string `json:"name"`
	Visible     bool   `json:"visible"`
	Label       string `json:"label"`
	CustomCss   string `json:"customCss"`
	Placeholder string `json:"placeholder"`
	Rule        string `json:"rule"`
	IsCustom    bool   `json:"isCustom"`
}

// SamlItem represents SAML attribute mapping
type SamlItem struct {
	Name       string `json:"name"`
	NameFormat string `json:"nameFormat"`
	Value      string `json:"value"`
}

// JwtItem represents JWT token attribute
type JwtItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Value    string `json:"value"`
	Type     string `json:"type"`
}

// ScopeDescription represents custom scope description
type ScopeDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProviderItem represents provider configuration in application
type ProviderItem struct {
	Name       string    `json:"name"`
	CanSignUp  bool      `json:"canSignUp"`
	CanSignIn  bool      `json:"canSignIn"`
	CanUnlink  bool      `json:"canUnlink"`
	Prompted   bool      `json:"prompted"`
	AlertType  string    `json:"alertType"`
	Rule       string    `json:"rule"`
	Provider   *Provider `json:"provider"`
}

// Cert represents certificate
type Cert struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`
	Scope       string `xorm:"varchar(100)" json:"scope"`
	Type        string `xorm:"varchar(100)" json:"type"`
	CryptoAlgorithm string `xorm:"varchar(100)" json:"cryptoAlgorithm"`
	BitSize     int    `json:"bitSize"`
	ExpireInYears int  `json:"expireInYears"`
	Certificate string `xorm:"mediumtext" json:"certificate"`
	PrivateKey  string `xorm:"mediumtext" json:"privateKey"`
	AuthorityPublicKey string `xorm:"mediumtext" json:"authorityPublicKey"`
	AuthorityRootPublicKey string `xorm:"mediumtext" json:"authorityRootPublicKey"`
}

// Provider represents identity provider
type Provider struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`
	Category    string `xorm:"varchar(100)" json:"category"`
	Type        string `xorm:"varchar(100)" json:"type"`
	SubType     string `xorm:"varchar(100)" json:"subType"`
	Method      string `xorm:"varchar(100)" json:"method"`
	ClientId    string `xorm:"varchar(100)" json:"clientId"`
	ClientSecret string `xorm:"varchar(2000)" json:"clientSecret"`
	ClientId2   string `xorm:"varchar(100)" json:"clientId2"`
	ClientSecret2 string `xorm:"varchar(100)" json:"clientSecret2"`
	Cert        string `xorm:"varchar(100)" json:"cert"`
	CustomAuthUrl    string `xorm:"varchar(200)" json:"customAuthUrl"`
	CustomTokenUrl   string `xorm:"varchar(200)" json:"customTokenUrl"`
	CustomUserInfoUrl string `xorm:"varchar(200)" json:"customUserInfoUrl"`
	CustomLogo       string `xorm:"varchar(200)" json:"customLogo"`
	Scopes           string `xorm:"varchar(100)" json:"scopes"`
	Domain           string `xorm:"varchar(100)" json:"domain"`
}

// Role represents role entity
type Role struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`
	Description string `xorm:"varchar(100)" json:"description"`
	Users       []string `xorm:"mediumtext" json:"users"`
	Groups      []string `xorm:"mediumtext" json:"groups"`
	Roles       []string `xorm:"mediumtext" json:"roles"`
	Domains     []string `xorm:"mediumtext" json:"domains"`
	IsEnabled   bool   `json:"isEnabled"`
}

// Permission represents permission entity
type Permission struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`
	Description string `xorm:"varchar(100)" json:"description"`
	Users       []string `xorm:"mediumtext" json:"users"`
	Groups      []string `xorm:"mediumtext" json:"groups"`
	Roles       []string `xorm:"mediumtext" json:"roles"`
	Domains     []string `xorm:"mediumtext" json:"domains"`
	Model       string `xorm:"varchar(100)" json:"model"`
	Adapter     string `xorm:"varchar(100)" json:"adapter"`
	ResourceType string `xorm:"varchar(100)" json:"resourceType"`
	Resources   []string `xorm:"mediumtext" json:"resources"`
	Actions     []string `xorm:"mediumtext" json:"actions"`
	Effect      string `xorm:"varchar(100)" json:"effect"`
	IsEnabled   bool   `json:"isEnabled"`
	Submitter   string `xorm:"varchar(100)" json:"submitter"`
	Approver    string `xorm:"varchar(100)" json:"approver"`
	ApproveTime string `xorm:"varchar(100)" json:"approveTime"`
	State       string `xorm:"varchar(100)" json:"state"`
}

// GetId returns the unique identifier for the application
func (a *Application) GetId() string {
	return a.Owner + "/" + a.Name
}

// TableName returns the table name
func (a *Application) TableName() string {
	return "application"
}

// IsOAuthApp checks if this is an OAuth application
func (a *Application) IsOAuthApp() bool {
	return a.ClientId != ""
}

// GetOrganizationId returns the organization ID
func (a *Application) GetOrganizationId() string {
	return a.Owner + "/" + a.Organization
}
