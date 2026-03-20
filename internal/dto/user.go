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

import (
	"time"
)

// CreateUserRequest represents create user request
type CreateUserRequest struct {
	Owner             string            `json:"owner"`
	Name              string            `json:"name" binding:"required"`
	DisplayName       string            `json:"displayName"`
	Email             string            `json:"email"`
	Phone             string            `json:"phone"`
	CountryCode       string            `json:"countryCode"`
	Password          string            `json:"password"`
	Type              string            `json:"type"`
	Avatar            string            `json:"avatar"`
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Gender            string            `json:"gender"`
	Birthday          string            `json:"birthday"`
	Location          string            `json:"location"`
	Address           []string          `json:"address"`
	Affiliation       string            `json:"affiliation"`
	Title             string            `json:"title"`
	Homepage          string            `json:"homepage"`
	Bio               string            `json:"bio"`
	Tag               string            `json:"tag"`
	Region            string            `json:"region"`
	Language          string            `json:"language"`
	Score             int               `json:"score"`
	SignupApplication string            `json:"signupApplication"`
	Properties        map[string]string `json:"properties"`
	Groups            []string          `json:"groups"`
}

// UpdateUserRequest represents update user request
type UpdateUserRequest struct {
	DisplayName string            `json:"displayName"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	CountryCode string            `json:"countryCode"`
	Password    string            `json:"password"`
	Avatar      string            `json:"avatar"`
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	Gender      string            `json:"gender"`
	Birthday    string            `json:"birthday"`
	Location    string            `json:"location"`
	Address     []string          `json:"address"`
	Affiliation string            `json:"affiliation"`
	Title       string            `json:"title"`
	Homepage    string            `json:"homepage"`
	Bio         string            `json:"bio"`
	Tag         string            `json:"tag"`
	Region      string            `json:"region"`
	Language    string            `json:"language"`
	Score       int               `json:"score"`
	IsAdmin     bool              `json:"isAdmin"`
	IsForbidden bool              `json:"isForbidden"`
	Properties  map[string]string `json:"properties"`
	Groups      []string          `json:"groups"`
	Roles       []string          `json:"roles"`
	Permissions []string          `json:"permissions"`
}

// UserResponse represents user response
type UserResponse struct {
	Owner             string            `json:"owner"`
	Name              string            `json:"name"`
	CreatedTime       string            `json:"createdTime"`
	Id                string            `json:"id"`
	Type              string            `json:"type"`
	DisplayName       string            `json:"displayName"`
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Avatar            string            `json:"avatar"`
	PermanentAvatar   string            `json:"permanentAvatar"`
	Email             string            `json:"email"`
	EmailVerified     bool              `json:"emailVerified"`
	Phone             string            `json:"phone"`
	CountryCode       string            `json:"countryCode"`
	Region            string            `json:"region"`
	Location          string            `json:"location"`
	Affiliation       string            `json:"affiliation"`
	Title             string            `json:"title"`
	Homepage          string            `json:"homepage"`
	Bio               string            `json:"bio"`
	Tag               string            `json:"tag"`
	Language          string            `json:"language"`
	Gender            string            `json:"gender"`
	Birthday          string            `json:"birthday"`
	Education         string            `json:"education"`
	Score             int               `json:"score"`
	Karma             int               `json:"karma"`
	Ranking           int               `json:"ranking"`
	Balance           float64           `json:"balance"`
	IsAdmin           bool              `json:"isAdmin"`
	IsForbidden       bool              `json:"isForbidden"`
	SignupApplication string            `json:"signupApplication"`
	Groups            []string          `json:"groups"`
	Roles             []*RoleInfo       `json:"roles"`
	Permissions       []*PermissionInfo `json:"permissions"`
	Properties        map[string]string `json:"properties"`
}

// RoleInfo represents role information in user response
type RoleInfo struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// PermissionInfo represents permission information in user response
type PermissionInfo struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// ListUsersRequest represents list users request
type ListUsersRequest struct {
	Owner     string `form:"owner" binding:"required"`
	GroupName string `form:"groupName"`
	PageSize  int    `form:"pageSize"`
	Page      int    `form:"p"`
	Field     string `form:"field"`
	Value     string `form:"value"`
	SortField string `form:"sortField"`
	SortOrder string `form:"sortOrder"`
}

// ListUsersResponse represents list users response
type ListUsersResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
}

// ImportUsersRequest represents import users request
type ImportUsersRequest struct {
	Users []*CreateUserRequest `json:"users" binding:"required"`
}

// ImportUsersResponse represents import users response
type ImportUsersResponse struct {
	Total     int      `json:"total"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	FailedIds []string `json:"failedIds"`
	Errors    []string `json:"errors"`
}

// ExportUsersRequest represents export users request
type ExportUsersRequest struct {
	Owner     string   `json:"owner"`
	Fields    []string `json:"fields"`
	Format    string   `json:"format"` // xlsx, csv, json
	GroupName string   `json:"groupName"`
}

// MFASetupRequest represents MFA setup request
type MFASetupRequest struct {
	MfaType     string `json:"mfaType" binding:"required"` // sms, email, totp, radius, push
	Secret      string `json:"secret"`
	Dest        string `json:"dest"`
	CountryCode string `json:"countryCode"`
	Passcode    string `json:"passcode"`
}

// MFASetupResponse represents MFA setup response
type MFASetupResponse struct {
	Enabled            bool     `json:"enabled"`
	IsPreferred        bool     `json:"isPreferred"`
	MfaType            string   `json:"mfaType"`
	Secret             string   `json:"secret,omitempty"`
	URL                string   `json:"url,omitempty"`
	QRCode             string   `json:"qrCode,omitempty"`
	RecoveryCodes      []string `json:"recoveryCodes,omitempty"`
	MfaRememberInHours int      `json:"mfaRememberInHours"`
}

// MFAVerifyRequest represents MFA verify request
type MFAVerifyRequest struct {
	MfaType  string `json:"mfaType" binding:"required"`
	Passcode string `json:"passcode" binding:"required"`
}

// UserFilter represents user filter options
type UserFilter struct {
	Owner       string
	GroupName   string
	Field       string
	Value       string
	SortField   string
	SortOrder   string
	IsAdmin     *bool
	IsForbidden *bool
	IsDeleted   *bool
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

// BatchUserOperation represents batch user operation
type BatchUserOperation struct {
	UserIds   []string          `json:"userIds" binding:"required"`
	Operation string            `json:"operation" binding:"required"` // enable, disable, delete, add_to_group, remove_from_group
	Params    map[string]string `json:"params"`
}

// UserActivity represents user activity log
type UserActivity struct {
	UserId    string    `json:"userId"`
	Action    string    `json:"action"`
	IpAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}
