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

package service

import (
	"testing"

	"github.com/casdoor/casdoor/object"
	"github.com/stretchr/testify/assert"
)

// TestPreprocessUser 测试用户数据预处理
func TestPreprocessUser(t *testing.T) {
	service := NewUserService()

	t.Run("Normalize user data", func(t *testing.T) {
		user := &object.User{
			Name:  "UserName",
			Email: "User@Example.COM",
			Phone: " 123-456-7890 ",
		}

		service.preprocessUser(user)

		assert.Equal(t, "username", user.Name)
		assert.Equal(t, "user@example.com", user.Email)
		assert.Equal(t, "1234567890", user.Phone)
	})

	t.Run("Set default values", func(t *testing.T) {
		user := &object.User{
			Name: "testuser",
		}

		service.preprocessUser(user)

		assert.Equal(t, "testuser", user.Name)
		// ID should be set
		assert.NotEmpty(t, user.Id)
	})

	t.Run("Hash password", func(t *testing.T) {
		user := &object.User{
			Name:     "testuser",
			Password: "plainpassword",
		}

		originalPassword := user.Password
		service.preprocessUser(user)

		assert.NotEqual(t, originalPassword, user.Password)
		assert.NotEmpty(t, user.Password)
		assert.Empty(t, user.Salt)
	})
}

// TestValidateUser 测试用户数据验证
func TestValidateUser(t *testing.T) {
	service := NewUserService()

	t.Run("Valid user", func(t *testing.T) {
		user := &object.User{
			Owner: "test-org",
			Name:  "testuser",
			Email: "test@example.com",
		}

		err := service.validateUser(user)
		assert.NoError(t, err)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		testCases := []struct {
			name string
			user *object.User
		}{
			{
				name: "Missing owner",
				user: &object.User{
					Name: "testuser",
				},
			},
			{
				name: "Missing name",
				user: &object.User{
					Owner: "test-org",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := service.validateUser(tc.user)
				assert.Error(t, err)
			})
		}
	})

	t.Run("Invalid email format", func(t *testing.T) {
		user := &object.User{
			Owner: "test-org",
			Name:  "testuser",
			Email: "invalid-email",
		}

		err := service.validateUser(user)
		assert.Error(t, err)
	})
}

// TestUserService_UpdateUser 测试更新用户
func TestUserService_UpdateUser(t *testing.T) {
	service := NewUserService()

	t.Run("Update without password change", func(t *testing.T) {
		oldUser := &object.User{
			Owner:    "test-org",
			Name:     "testuser",
			Password: "hashedpassword",
			Email:    "old@example.com",
		}

		newUser := &object.User{
			Owner: "test-org",
			Name:  "testuser",
			Email: "new@example.com",
			// Password is empty, should not change
		}

		cols := service.determineUpdateColumns(oldUser, newUser)

		// Password should not be in the update columns
		assert.NotContains(t, cols, "password")
		assert.Contains(t, cols, "email")
	})

	t.Run("Update with password change", func(t *testing.T) {
		oldUser := &object.User{
			Owner:    "test-org",
			Name:     "testuser",
			Password: "hashedpassword",
		}

		newUser := &object.User{
			Owner:    "test-org",
			Name:     "testuser",
			Password: "newpassword",
		}

		cols := service.determineUpdateColumns(oldUser, newUser)

		// Password should be in the update columns
		assert.Contains(t, cols, "password")
	})
}

// TestUserService_CalculateUpdateColumns 测试计算更新列
func TestCalculateUpdateColumns(t *testing.T) {
	service := NewUserService()

	oldUser := &object.User{
		Owner:             "test-org",
		Name:              "testuser",
		DisplayName:       "Old Name",
		Email:             "old@example.com",
		Phone:             "1234567890",
		IsAdmin:           false,
		IsGlobalAdmin:     false,
		IsForbidden:       false,
		IsDeleted:         false,
		MfaPhoneEnabled:   false,
		MfaEmailEnabled:   false,
		MfaRadiusEnabled:  false,
		MfaPushEnabled:    false,
		PreferredMfaType:  "",
		RecoveryCodes:     nil,
		CountryCode:       "US",
		Address:           []string{"Old Address"},
		Avatar:            "old-avatar-url",
		Type:              "normal-user",
	}

	newUser := &object.User{
		Owner:            "test-org",
		Name:             "testuser",
		DisplayName:      "New Name",
		Email:            "new@example.com",
		Phone:            "0987654321",
		IsAdmin:          true,
		IsGlobalAdmin:    true,
		IsForbidden:      true,
		IsDeleted:        true,
		MfaPhoneEnabled:  true,
		MfaEmailEnabled:  true,
		MfaRadiusEnabled: true,
		MfaPushEnabled:   true,
		PreferredMfaType: "app",
		RecoveryCodes:    []string{"code1", "code2"},
		CountryCode:      "CN",
		Address:          []string{"New Address"},
		Avatar:           "new-avatar-url",
		Type:             "admin-user",
	}

	cols := service.calculateUpdateColumns(oldUser, newUser)

	// Should include all changed fields
	expectedCols := []string{
		"display_name", "email", "phone", "is_admin", "is_global_admin",
		"is_forbidden", "is_deleted", "mfa_phone_enabled", "mfa_email_enabled",
		"mfa_radius_enabled", "mfa_push_enabled", "preferred_mfa_type",
		"recovery_codes", "country_code", "address", "avatar", "type", "hash",
	}

	for _, col := range expectedCols {
		assert.Contains(t, cols, col, "Column %s should be included", col)
	}
}

// TestUserService_GetMaskedUser 测试获取掩码用户
func TestUserService_GetMaskedUser(t *testing.T) {
	service := NewUserService()

	user := &object.User{
		Owner:       "test-org",
		Name:        "testuser",
		Password:    "hashedpassword",
		Phone:       "1234567890",
		CountryCode: "CN",
		Email:       "user@example.com",
	}

	maskedUser := service.getMaskedUser(user)

	// Sensitive fields should be masked
	assert.Empty(t, maskedUser.Password)
	assert.Empty(t, maskedUser.Phone)
	assert.Empty(t, maskedUser.CountryCode)

	// Non-sensitive fields should be preserved
	assert.Equal(t, user.Owner, maskedUser.Owner)
	assert.Equal(t, user.Name, maskedUser.Name)
	assert.Equal(t, user.Email, maskedUser.Email)
}
