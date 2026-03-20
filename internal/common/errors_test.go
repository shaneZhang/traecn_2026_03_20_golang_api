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

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewBusinessError 测试创建业务错误
func TestNewBusinessError(t *testing.T) {
	code := ErrCodeUserNotFound
	message := "User not found"

	err := NewBusinessError(code, message)
	assert.NotNil(t, err)
	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Nil(t, err.Err)
}

// TestNewBusinessErrorWithErr 测试创建带嵌套错误的业务错误
func TestNewBusinessErrorWithErr(t *testing.T) {
	code := ErrCodeInternalError
	message := "Internal server error"
	innerErr := errors.New("database connection failed")

	err := NewBusinessErrorWithErr(code, message, innerErr)
	assert.NotNil(t, err)
	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, innerErr, err.Err)
	assert.Equal(t, innerErr, err.Unwrap())
}

// TestBusinessErrorError 测试Error方法
func TestBusinessErrorError(t *testing.T) {
	err := &BusinessError{
		Code:    ErrCodeBadRequest,
		Message: "Invalid parameters",
	}

	errStr := err.Error()
	assert.Contains(t, errStr, "Invalid parameters")
	assert.Contains(t, errStr, string(ErrCodeBadRequest))
}

// TestPredefinedErrors 测试预定义错误
func TestPredefinedErrors(t *testing.T) {
	// 测试通用错误
	assert.Equal(t, ErrCodeBadRequest, ErrBadRequest.Code)
	assert.Equal(t, ErrCodeUnauthorized, ErrUnauthorized.Code)
	assert.Equal(t, ErrCodeForbidden, ErrForbidden.Code)
	assert.Equal(t, ErrCodeNotFound, ErrNotFound.Code)
	assert.Equal(t, ErrCodeInternalError, ErrInternalError.Code)
	assert.Equal(t, ErrCodeValidationError, ErrValidationError.Code)
	assert.Equal(t, ErrCodeDuplicateEntry, ErrDuplicateEntry.Code)

	// 测试用户模块错误
	assert.Equal(t, ErrCodeUserNotFound, ErrUserNotFound.Code)
	assert.Equal(t, ErrCodeUserAlreadyExists, ErrUserAlreadyExists.Code)
	assert.Equal(t, ErrCodeInvalidPassword, ErrInvalidPassword.Code)

	// 测试组织模块错误
	assert.Equal(t, ErrCodeOrgNotFound, ErrOrgNotFound.Code)

	// 测试应用模块错误
	assert.Equal(t, ErrCodeAppNotFound, ErrAppNotFound.Code)

	// 测试MFA模块错误
	assert.Equal(t, ErrCodeMfaNotEnabled, ErrMfaNotEnabled.Code)
	assert.Equal(t, ErrCodeMfaInvalidCode, ErrMfaInvalidCode.Code)
}

// TestIsBusinessError 测试是否为业务错误
func TestIsBusinessError(t *testing.T) {
	businessErr := &BusinessError{
		Code:    ErrCodeBadRequest,
		Message: "test",
	}

	regularErr := errors.New("regular error")

	assert.True(t, IsBusinessError(businessErr))
	assert.False(t, IsBusinessError(regularErr))
	assert.False(t, IsBusinessError(nil))
}
