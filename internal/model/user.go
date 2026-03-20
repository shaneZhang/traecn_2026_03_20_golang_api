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

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

// User represents the user entity in the system
type User struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(255) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100) index" json:"createdTime"`
	UpdatedTime string `xorm:"varchar(100)" json:"updatedTime"`
	DeletedTime string `xorm:"varchar(100)" json:"deletedTime"`

	Id                   string  `xorm:"varchar(100) index" json:"id"`
	ExternalId           string  `xorm:"varchar(100) index" json:"externalId"`
	Type                 string  `xorm:"varchar(100)" json:"type"`
	Password             string  `xorm:"varchar(150)" json:"password"`
	PasswordSalt         string  `xorm:"varchar(100)" json:"passwordSalt"`
	PasswordType         string  `xorm:"varchar(100)" json:"passwordType"`
	DisplayName          string  `xorm:"varchar(100)" json:"displayName"`
	FirstName            string  `xorm:"varchar(100)" json:"firstName"`
	LastName             string  `xorm:"varchar(100)" json:"lastName"`
	Avatar               string  `xorm:"text" json:"avatar"`
	AvatarType           string  `xorm:"varchar(100)" json:"avatarType"`
	PermanentAvatar      string  `xorm:"varchar(500)" json:"permanentAvatar"`
	Email                string  `xorm:"varchar(100) index" json:"email"`
	EmailVerified        bool    `json:"emailVerified"`
	Phone                string  `xorm:"varchar(100) index" json:"phone"`
	CountryCode          string  `xorm:"varchar(6)" json:"countryCode"`
	Region               string  `xorm:"varchar(100)" json:"region"`
	Location             string  `xorm:"varchar(100)" json:"location"`
	Address              []string `json:"address"`
	Addresses            []*Address `xorm:"addresses blob" json:"addresses"`
	Affiliation          string  `xorm:"varchar(100)" json:"affiliation"`
	Title                string  `xorm:"varchar(100)" json:"title"`
	IdCardType           string  `xorm:"varchar(100)" json:"idCardType"`
	IdCard               string  `xorm:"varchar(100) index" json:"idCard"`
	RealName             string  `xorm:"varchar(100)" json:"realName"`
	IsVerified           bool    `json:"isVerified"`
	Homepage             string  `xorm:"varchar(100)" json:"homepage"`
	Bio                  string  `xorm:"varchar(100)" json:"bio"`
	Tag                  string  `xorm:"varchar(100)" json:"tag"`
	Language             string  `xorm:"varchar(100)" json:"language"`
	Gender               string  `xorm:"varchar(100)" json:"gender"`
	Birthday             string  `xorm:"varchar(100)" json:"birthday"`
	Education            string  `xorm:"varchar(100)" json:"education"`
	Score                int     `json:"score"`
	Karma                int     `json:"karma"`
	Ranking              int     `json:"ranking"`
	Balance              float64 `json:"balance"`
	BalanceCredit        float64 `json:"balanceCredit"`
	Currency             string  `xorm:"varchar(100)" json:"currency"`
	BalanceCurrency      string  `xorm:"varchar(100)" json:"balanceCurrency"`
	IsDefaultAvatar      bool    `json:"isDefaultAvatar"`
	IsOnline             bool    `json:"isOnline"`
	IsAdmin              bool    `json:"isAdmin"`
	IsForbidden          bool    `json:"isForbidden"`
	IsDeleted            bool    `json:"isDeleted"`
	SignupApplication    string  `xorm:"varchar(100)" json:"signupApplication"`
	Hash                 string  `xorm:"varchar(100)" json:"hash"`
	PreHash              string  `xorm:"varchar(100)" json:"preHash"`
	RegisterType         string  `xorm:"varchar(100)" json:"registerType"`
	RegisterSource       string  `xorm:"varchar(100)" json:"registerSource"`
	AccessKey            string  `xorm:"varchar(100)" json:"accessKey"`
	AccessSecret         string  `xorm:"varchar(100)" json:"accessSecret"`
	AccessToken          string  `xorm:"mediumtext" json:"accessToken"`
	OriginalToken        string  `xorm:"mediumtext" json:"originalToken"`
	OriginalRefreshToken string  `xorm:"mediumtext" json:"originalRefreshToken"`

	CreatedIp      string `xorm:"varchar(100)" json:"createdIp"`
	LastSigninTime string `xorm:"varchar(100)" json:"lastSigninTime"`
	LastSigninIp   string `xorm:"varchar(100)" json:"lastSigninIp"`

	// OAuth providers
	GitHub          string `xorm:"github varchar(100)" json:"github"`
	Google          string `xorm:"varchar(100)" json:"google"`
	QQ              string `xorm:"qq varchar(100)" json:"qq"`
	WeChat          string `xorm:"wechat varchar(100)" json:"wechat"`
	Facebook        string `xorm:"facebook varchar(100)" json:"facebook"`
	DingTalk        string `xorm:"dingtalk varchar(100)" json:"dingtalk"`
	Weibo           string `xorm:"weibo varchar(100)" json:"weibo"`
	Gitee           string `xorm:"gitee varchar(100)" json:"gitee"`
	LinkedIn        string `xorm:"linkedin varchar(100)" json:"linkedin"`
	Wecom           string `xorm:"wecom varchar(100)" json:"wecom"`
	Lark            string `xorm:"lark varchar(100)" json:"lark"`
	Gitlab          string `xorm:"gitlab varchar(100)" json:"gitlab"`
	Adfs            string `xorm:"adfs varchar(100)" json:"adfs"`
	Baidu           string `xorm:"baidu varchar(100)" json:"baidu"`
	Alipay          string `xorm:"alipay varchar(100)" json:"alipay"`
	Casdoor         string `xorm:"casdoor varchar(100)" json:"casdoor"`
	Infoflow        string `xorm:"infoflow varchar(100)" json:"infoflow"`
	Apple           string `xorm:"apple varchar(100)" json:"apple"`
	AzureAD         string `xorm:"azuread varchar(100)" json:"azuread"`
	AzureADB2c      string `xorm:"azureadb2c varchar(100)" json:"azureadb2c"`
	Slack           string `xorm:"slack varchar(100)" json:"slack"`
	Steam           string `xorm:"steam varchar(100)" json:"steam"`
	Bilibili        string `xorm:"bilibili varchar(100)" json:"bilibili"`
	Okta            string `xorm:"okta varchar(100)" json:"okta"`
	Douyin          string `xorm:"douyin varchar(100)" json:"douyin"`
	Kwai            string `xorm:"kwai varchar(100)" json:"kwai"`
	Line            string `xorm:"line varchar(100)" json:"line"`
	Amazon          string `xorm:"amazon varchar(100)" json:"amazon"`
	Auth0           string `xorm:"auth0 varchar(100)" json:"auth0"`
	MetaMask        string `xorm:"metamask varchar(100)" json:"metamask"`
	Web3Onboard     string `xorm:"web3onboard varchar(100)" json:"web3onboard"`
	Custom          string `xorm:"custom varchar(100)" json:"custom"`
	Custom2         string `xorm:"custom2 text" json:"custom2"`
	Custom3         string `xorm:"custom3 text" json:"custom3"`
	Custom4         string `xorm:"custom4 text" json:"custom4"`
	Custom5         string `xorm:"custom5 text" json:"custom5"`
	Custom6         string `xorm:"custom6 text" json:"custom6"`
	Custom7         string `xorm:"custom7 text" json:"custom7"`
	Custom8         string `xorm:"custom8 text" json:"custom8"`
	Custom9         string `xorm:"custom9 text" json:"custom9"`
	Custom10        string `xorm:"custom10 text" json:"custom10"`

	// MFA fields
	WebauthnCredentials []webauthn.Credential `xorm:"webauthnCredentials blob" json:"webauthnCredentials"`
	PreferredMfaType    string                `xorm:"varchar(100)" json:"preferredMfaType"`
	RecoveryCodes       []string              `xorm:"mediumtext" json:"recoveryCodes"`
	TotpSecret          string                `xorm:"varchar(100)" json:"totpSecret"`
	MfaPhoneEnabled     bool                  `json:"mfaPhoneEnabled"`
	MfaEmailEnabled     bool                  `json:"mfaEmailEnabled"`
	MfaRadiusEnabled    bool                  `json:"mfaRadiusEnabled"`
	MfaRadiusUsername   string                `xorm:"varchar(100)" json:"mfaRadiusUsername"`
	MfaRadiusProvider   string                `xorm:"varchar(100)" json:"mfaRadiusProvider"`
	MfaPushEnabled      bool                  `json:"mfaPushEnabled"`
	MfaPushReceiver     string                `xorm:"varchar(100)" json:"mfaPushReceiver"`
	MfaPushProvider     string                `xorm:"varchar(100)" json:"mfaPushProvider"`

	Invitation          string                `xorm:"varchar(100) index" json:"invitation"`
	InvitationCode      string                `xorm:"varchar(100) index" json:"invitationCode"`
	FaceIds             []*FaceId             `json:"faceIds"`
	Cart                []ProductInfo         `xorm:"mediumtext" json:"cart"`

	Ldap       string            `xorm:"ldap varchar(100)" json:"ldap"`
	Properties map[string]string `json:"properties"`

	Roles       []*Role       `json:"roles"`
	Permissions []*Permission `json:"permissions"`
	Groups      []string      `xorm:"mediumtext" json:"groups"`

	LastChangePasswordTime string `xorm:"varchar(100)" json:"lastChangePasswordTime"`
	LastSigninWrongTime    string `xorm:"varchar(100)" json:"lastSigninWrongTime"`
	SigninWrongTimes       int    `json:"signinWrongTimes"`

	ManagedAccounts     []ManagedAccount `xorm:"managedAccounts blob" json:"managedAccounts"`
	MfaAccounts         []MfaAccount     `xorm:"mfaAccounts blob" json:"mfaAccounts"`
	MfaItems            []*MfaItem       `xorm:"varchar(300)" json:"mfaItems"`
	MfaRememberDeadline string           `xorm:"varchar(100)" json:"mfaRememberDeadline"`
	NeedUpdatePassword  bool             `json:"needUpdatePassword"`
	IpWhitelist         string           `xorm:"varchar(200)" json:"ipWhitelist"`
	ApplicationScopes   []ConsentRecord  `xorm:"mediumtext" json:"applicationScopes"`
}

// Address represents user address
type Address struct {
	Tag     string `xorm:"varchar(100)" json:"tag"`
	Line1   string `xorm:"varchar(100)" json:"line1"`
	Line2   string `xorm:"varchar(100)" json:"line2"`
	City    string `xorm:"varchar(100)" json:"city"`
	State   string `xorm:"varchar(100)" json:"state"`
	ZipCode string `xorm:"varchar(100)" json:"zipCode"`
	Region  string `xorm:"varchar(100)" json:"region"`
}

// FaceId represents face identification data
type FaceId struct {
	Name       string    `xorm:"varchar(100) notnull pk" json:"name"`
	FaceIdData []float64 `json:"faceIdData"`
	ImageUrl   string    `json:"ImageUrl"`
}

// ManagedAccount represents a managed account
type ManagedAccount struct {
	Application string `xorm:"varchar(100)" json:"application"`
	Username    string `xorm:"varchar(100)" json:"username"`
	Password    string `xorm:"varchar(100)" json:"password"`
	SigninUrl   string `xorm:"varchar(200)" json:"signinUrl"`
}

// MfaAccount represents MFA account
type MfaAccount struct {
	AccountName string `xorm:"varchar(100)" json:"accountName"`
	Issuer      string `xorm:"varchar(100)" json:"issuer"`
	SecretKey   string `xorm:"varchar(100)" json:"secretKey"`
	Origin      string `xorm:"varchar(100)" json:"origin"`
}

// MfaItem represents MFA configuration item
type MfaItem struct {
	Name string `json:"name"`
	Rule string `json:"rule"`
}

// ProductInfo represents product in cart
type ProductInfo struct {
	ProductName string  `json:"productName"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// ConsentRecord represents OAuth consent record
type ConsentRecord struct {
	Application string   `json:"application"`
	Scopes      []string `json:"scopes"`
	GrantedAt   time.Time `json:"grantedAt"`
}

// GetId returns the unique identifier for the user
func (u *User) GetId() string {
	return u.Owner + "/" + u.Name
}

// GetCreatedAt returns the creation time
func (u *User) GetCreatedAt() time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", u.CreatedTime)
	return t
}

// IsGlobalAdmin checks if user is global admin
func (u *User) IsGlobalAdmin() bool {
	return u.Owner == "built-in" && u.Name == "admin"
}

// IsSoftDeleted checks if user is soft deleted
func (u *User) IsSoftDeleted() bool {
	return u.IsDeleted || u.DeletedTime != ""
}

// TableName returns the table name
func (u *User) TableName() string {
	return "user"
}
