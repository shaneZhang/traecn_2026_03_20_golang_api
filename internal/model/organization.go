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

// Organization represents an organization entity
type Organization struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName            string     `xorm:"varchar(100)" json:"displayName"`
	WebsiteUrl             string     `xorm:"varchar(100)" json:"websiteUrl"`
	Logo                   string     `xorm:"varchar(200)" json:"logo"`
	LogoDark               string     `xorm:"varchar(200)" json:"logoDark"`
	Favicon                string     `xorm:"varchar(200)" json:"favicon"`
	HasPrivilegeConsent    bool       `xorm:"bool" json:"hasPrivilegeConsent"`
	PasswordType           string     `xorm:"varchar(100)" json:"passwordType"`
	PasswordSalt           string     `xorm:"varchar(100)" json:"passwordSalt"`
	PasswordOptions        []string   `xorm:"varchar(100)" json:"passwordOptions"`
	PasswordObfuscatorType string     `xorm:"varchar(100)" json:"passwordObfuscatorType"`
	PasswordObfuscatorKey  string     `xorm:"varchar(100)" json:"passwordObfuscatorKey"`
	PasswordExpireDays     int        `json:"passwordExpireDays"`
	CountryCodes           []string   `xorm:"mediumtext" json:"countryCodes"`
	DefaultAvatar          string     `xorm:"varchar(200)" json:"defaultAvatar"`
	DefaultApplication     string     `xorm:"varchar(100)" json:"defaultApplication"`
	UserTypes              []string   `xorm:"mediumtext" json:"userTypes"`
	Tags                   []string   `xorm:"mediumtext" json:"tags"`
	Languages              []string   `xorm:"varchar(255)" json:"languages"`
	ThemeData              *ThemeData `xorm:"json" json:"themeData"`
	MasterPassword         string     `xorm:"varchar(200)" json:"masterPassword"`
	DefaultPassword        string     `xorm:"varchar(200)" json:"defaultPassword"`
	MasterVerificationCode string     `xorm:"varchar(100)" json:"masterVerificationCode"`
	IpWhitelist            string     `xorm:"varchar(200)" json:"ipWhitelist"`
	InitScore              int        `json:"initScore"`
	EnableSoftDeletion     bool       `json:"enableSoftDeletion"`
	IsProfilePublic        bool       `json:"isProfilePublic"`
	UseEmailAsUsername     bool       `json:"useEmailAsUsername"`
	EnableTour             bool       `json:"enableTour"`
	DisableSignin          bool       `json:"disableSignin"`
	IpRestriction          string     `json:"ipRestriction"`
	NavItems               []string   `xorm:"mediumtext" json:"navItems"`
	UserNavItems           []string   `xorm:"mediumtext" json:"userNavItems"`
	WidgetItems            []string   `xorm:"mediumtext" json:"widgetItems"`

	MfaItems           []*MfaItem     `xorm:"varchar(300)" json:"mfaItems"`
	MfaRememberInHours int            `json:"mfaRememberInHours"`
	AccountMenu        string         `xorm:"varchar(20)" json:"accountMenu"`
	AccountItems       []*AccountItem `xorm:"mediumtext" json:"accountItems"`

	DcrPolicy string `xorm:"varchar(100)" json:"dcrPolicy"`

	LdapAttributes []string `xorm:"mediumtext" json:"ldapAttributes"`

	KerberosRealm       string `xorm:"varchar(200)" json:"kerberosRealm"`
	KerberosKdcHost     string `xorm:"varchar(200)" json:"kerberosKdcHost"`
	KerberosKeytab      string `xorm:"mediumtext" json:"kerberosKeytab"`
	KerberosServiceName string `xorm:"varchar(100)" json:"kerberosServiceName"`

	OrgBalance      float64 `json:"orgBalance"`
	UserBalance     float64 `json:"userBalance"`
	BalanceCredit   float64 `json:"balanceCredit"`
	BalanceCurrency string  `xorm:"varchar(100)" json:"balanceCurrency"`

	// Hierarchy fields
	ParentID string `xorm:"varchar(100)" json:"parentId"`

	// UI Customization
	HeaderHtml string `xorm:"mediumtext" json:"headerHtml"`
	FooterHtml string `xorm:"mediumtext" json:"footerHtml"`
	SigninHtml string `xorm:"mediumtext" json:"signinHtml"`
	SignupHtml string `xorm:"mediumtext" json:"signupHtml"`

	// URLs
	ForgetUrl      string `xorm:"varchar(200)" json:"forgetUrl"`
	AffiliationUrl string `xorm:"varchar(200)" json:"affiliationUrl"`
	TermsOfUse     string `xorm:"mediumtext" json:"termsOfUse"`
	SignupUrl      string `xorm:"varchar(200)" json:"signupUrl"`
	SigninUrl      string `xorm:"varchar(200)" json:"signinUrl"`

	// Phone
	PhonePrefix string `xorm:"varchar(10)" json:"phonePrefix"`

	// SAML
	EnableSamlC14n10 bool `json:"enableSamlC14n10"`
	SamlReplyLimit   int  `json:"samlReplyLimit"`

	// OAuth
	ClientId     string `xorm:"varchar(100)" json:"clientId"`
	ClientSecret string `xorm:"varchar(100)" json:"clientSecret"`
}

// ThemeData represents theme configuration
type ThemeData struct {
	ThemeType    string `xorm:"varchar(30)" json:"themeType"`
	ColorPrimary string `xorm:"varchar(10)" json:"colorPrimary"`
	BorderRadius int    `xorm:"int" json:"borderRadius"`
	IsCompact    bool   `xorm:"bool" json:"isCompact"`
	IsEnabled    bool   `xorm:"bool" json:"isEnabled"`
}

// AccountItem represents account configuration item
type AccountItem struct {
	Name       string `json:"name"`
	Visible    bool   `json:"visible"`
	ViewRule   string `json:"viewRule"`
	ModifyRule string `json:"modifyRule"`
	Regex      string `json:"regex"`
	Tab        string `json:"tab"`
}

// GetId returns the unique identifier for the organization
func (o *Organization) GetId() string {
	return o.Owner + "/" + o.Name
}

// TableName returns the table name
func (o *Organization) TableName() string {
	return "organization"
}

// IsBuiltIn checks if this is the built-in organization
func (o *Organization) IsBuiltIn() bool {
	return o.Name == "built-in"
}
