// Copyright 2024 The Casdoor Authors. All Rights Reserved.
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

package common

// 通用常量
const (
	// 默认分页参数
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 1000

	// 排序方向
	SortAsc  = "asc"
	SortDesc = "desc"

	// MFA类型
	MfaTypeEmail  = "email"
	MfaTypeSms    = "sms"
	MfaTypeTotp   = "app"
	MfaTypeRadius = "radius"
	MfaTypePush   = "push"

	// MFA会话键
	MfaSessionUserId = "MfaSessionUserId"
)

// 状态常量
const (
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"
)

// 用户相关常量
const (
	UserTypeNormal = "normal-user"
	UserTypeGuest  = "guest-user"
	UserTypeAdmin  = "admin"
)

// 组织相关常量
const (
	BuiltInOrg = "built-in"
	AdminUser  = "admin"
)

// 应用相关常量
const (
	BuiltInApp = "app-built-in"
)
