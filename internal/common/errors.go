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

import "fmt"

// ErrorCode 错误码定义
type ErrorCode string

const (
	// 通用错误码
	ErrCodeSuccess         ErrorCode = "SUCCESS"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidationError ErrorCode = "VALIDATION_ERROR"
	ErrCodeDuplicateEntry  ErrorCode = "DUPLICATE_ENTRY"

	// 用户模块错误码
	ErrCodeUserNotFound        ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists   ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidPassword     ErrorCode = "INVALID_PASSWORD"
	ErrCodePasswordComplexity  ErrorCode = "PASSWORD_COMPLEXITY"
	ErrCodeEmailAlreadyExists  ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrCodePhoneAlreadyExists  ErrorCode = "PHONE_ALREADY_EXISTS"
	ErrCodeUsernameAlreadyExists ErrorCode = "USERNAME_ALREADY_EXISTS"

	// 组织模块错误码
	ErrCodeOrgNotFound      ErrorCode = "ORG_NOT_FOUND"
	ErrCodeOrgAlreadyExists ErrorCode = "ORG_ALREADY_EXISTS"

	// 应用模块错误码
	ErrCodeAppNotFound      ErrorCode = "APP_NOT_FOUND"
	ErrCodeAppAlreadyExists ErrorCode = "APP_ALREADY_EXISTS"

	// MFA模块错误码
	ErrCodeMfaNotEnabled    ErrorCode = "MFA_NOT_ENABLED"
	ErrCodeMfaInvalidCode   ErrorCode = "MFA_INVALID_CODE"
	ErrCodeMfaAlreadyEnabled ErrorCode = "MFA_ALREADY_ENABLED"
)

// BusinessError 业务错误
type BusinessError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *BusinessError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *BusinessError) Unwrap() error {
	return e.Err
}

// NewBusinessError 创建业务错误
func NewBusinessError(code ErrorCode, message string, err ...error) *BusinessError {
	bizErr := &BusinessError{
		Code:    code,
		Message: message,
	}
	if len(err) > 0 && err[0] != nil {
		bizErr.Err = err[0]
	}
	return bizErr
}

// 预定义错误
var (
	// 通用错误
	ErrSuccess         = NewBusinessError(ErrCodeSuccess, "Success")
	ErrBadRequest      = NewBusinessError(ErrCodeBadRequest, "Invalid request parameters")
	ErrUnauthorized    = NewBusinessError(ErrCodeUnauthorized, "Unauthorized access")
	ErrForbidden       = NewBusinessError(ErrCodeForbidden, "Permission denied")
	ErrNotFound        = NewBusinessError(ErrCodeNotFound, "Resource not found")
	ErrInternalError   = NewBusinessError(ErrCodeInternalError, "Internal server error")
	ErrValidationError = NewBusinessError(ErrCodeValidationError, "Validation failed")
	ErrDuplicateEntry  = NewBusinessError(ErrCodeDuplicateEntry, "Duplicate entry")

	// 用户错误
	ErrUserNotFound        = NewBusinessError(ErrCodeUserNotFound, "User not found")
	ErrUserAlreadyExists   = NewBusinessError(ErrCodeUserAlreadyExists, "User already exists")
	ErrInvalidPassword     = NewBusinessError(ErrCodeInvalidPassword, "Invalid password")
	ErrPasswordComplexity  = NewBusinessError(ErrCodePasswordComplexity, "Password does not meet complexity requirements")
	ErrEmailAlreadyExists  = NewBusinessError(ErrCodeEmailAlreadyExists, "Email already exists")
	ErrPhoneAlreadyExists  = NewBusinessError(ErrCodePhoneAlreadyExists, "Phone already exists")
	ErrUsernameAlreadyExists = NewBusinessError(ErrCodeUsernameAlreadyExists, "Username already exists")

	// 组织错误
	ErrOrgNotFound      = NewBusinessError(ErrCodeOrgNotFound, "Organization not found")
	ErrOrgAlreadyExists = NewBusinessError(ErrCodeOrgAlreadyExists, "Organization already exists")

	// 应用错误
	ErrAppNotFound      = NewBusinessError(ErrCodeAppNotFound, "Application not found")
	ErrAppAlreadyExists = NewBusinessError(ErrCodeAppAlreadyExists, "Application already exists")

	// MFA错误
	ErrMfaNotEnabled    = NewBusinessError(ErrCodeMfaNotEnabled, "MFA not enabled")
	ErrMfaInvalidCode   = NewBusinessError(ErrCodeMfaInvalidCode, "Invalid MFA code")
	ErrMfaAlreadyEnabled = NewBusinessError(ErrCodeMfaAlreadyEnabled, "MFA already enabled")
)
