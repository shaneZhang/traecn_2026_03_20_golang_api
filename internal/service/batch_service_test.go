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
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/casdoor/casdoor/object"
	"github.com/stretchr/testify/assert"
)

// TestBatchService_ProcessCSVRecord 测试CSV记录处理
func TestBatchService_ProcessCSVRecord(t *testing.T) {
	service := NewBatchService()

	t.Run("Valid CSV record", func(t *testing.T) {
		header := []string{"owner", "name", "display_name", "email", "phone"}
		record := []string{"test-org", "user1", "User One", "user1@example.com", "1234567890"}
		lineNum := 2

		user, err := service.processCSVRecord(header, record, lineNum)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test-org", user.Owner)
		assert.Equal(t, "user1", user.Name)
		assert.Equal(t, "User One", user.DisplayName)
		assert.Equal(t, "user1@example.com", user.Email)
		assert.Equal(t, "1234567890", user.Phone)
	})

	t.Run("With prefixed headers", func(t *testing.T) {
		header := []string{"#owner", "#name", "#display_name", "#email"}
		record := []string{"test-org", "user2", "User Two", "user2@example.com"}
		lineNum := 2

		user, err := service.processCSVRecord(header, record, lineNum)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test-org", user.Owner)
		assert.Equal(t, "user2", user.Name)
		assert.Equal(t, "User Two", user.DisplayName)
		assert.Equal(t, "user2@example.com", user.Email)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		header := []string{"owner", "display_name", "email"} // Missing name
		record := []string{"test-org", "User Three", "user3@example.com"}
		lineNum := 2

		user, err := service.processCSVRecord(header, record, lineNum)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("Field count mismatch", func(t *testing.T) {
		header := []string{"owner", "name", "email"}
		record := []string{"test-org", "user4"} // Missing email field
		lineNum := 2

		user, err := service.processCSVRecord(header, record, lineNum)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

// TestBatchService_ValidateUserForImport 测试导入用户验证
func TestBatchService_ValidateUserForImport(t *testing.T) {
	service := NewBatchService()

	t.Run("Valid user", func(t *testing.T) {
		user := &object.User{
			Owner: "test-org",
			Name:  "validuser",
			Email: "valid@example.com",
		}

		err := service.validateUserForImport(user, 1)
		assert.NoError(t, err)
	})

	t.Run("Invalid email", func(t *testing.T) {
		user := &object.User{
			Owner: "test-org",
			Name:  "invaliduser",
			Email: "invalid-email-format",
		}

		err := service.validateUserForImport(user, 1)
		assert.Error(t, err)
	})

	t.Run("Missing owner", func(t *testing.T) {
		user := &object.User{
			Name:  "nouser",
			Email: "nouser@example.com",
		}

		err := service.validateUserForImport(user, 1)
		assert.Error(t, err)
	})
}

// TestBatchService_GenerateCSVContent 测试生成CSV内容
func TestBatchService_GenerateCSVContent(t *testing.T) {
	service := NewBatchService()

	users := []*object.User{
		{
			Owner:       "test-org",
			Name:        "user1",
			DisplayName: "User One",
			Email:       "user1@example.com",
			Phone:       "1234567890",
		},
		{
			Owner:       "test-org",
			Name:        "user2",
			DisplayName: "User Two",
			Email:       "user2@example.com",
			Phone:       "0987654321",
		},
	}

	content, err := service.generateCSVContent(users)
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	// Parse the CSV content back to verify
	reader := csv.NewReader(bytes.NewReader(content))
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Len(t, records, 3) // Header + 2 users

	// Check header
	expectedHeader := []string{"owner", "name", "display_name", "email", "phone", "type", "avatar", "address", "affiliation", "title", "id", "homepage", "bio", "tag", "region", "language", "gender", "birthday", "education", "score", "karma", "ranking", "is_default_avatar", "is_online", "is_admin", "is_global_admin", "is_forbidden", "created_time", "updated_time"}
	assert.Equal(t, expectedHeader, records[0])

	// Check first user data
	assert.Equal(t, "test-org", records[1][0])
	assert.Equal(t, "user1", records[1][1])
	assert.Equal(t, "User One", records[1][2])
	assert.Equal(t, "user1@example.com", records[1][3])
}

// TestBatchService_GenerateUsersExportHeaders 测试导出表头生成
func TestBatchService_GenerateUsersExportHeaders(t *testing.T) {
	service := NewBatchService()

	headers := service.generateUsersExportHeaders()

	expectedHeaders := []string{"owner", "name", "display_name", "email", "phone", "type", "avatar", "address", "affiliation", "title", "id", "homepage", "bio", "tag", "region", "language", "gender", "birthday", "education", "score", "karma", "ranking", "is_default_avatar", "is_online", "is_admin", "is_global_admin", "is_forbidden", "created_time", "updated_time"}

	assert.Equal(t, expectedHeaders, headers)
}

// TestBatchService_GenerateUserExportRecord 测试用户导出记录生成
func TestBatchService_GenerateUserExportRecord(t *testing.T) {
	service := NewBatchService()

	user := &object.User{
		Owner:            "test-org",
		Name:             "testuser",
		DisplayName:      "Test User",
		Email:            "test@example.com",
		Phone:            "1234567890",
		Type:             "normal-user",
		Avatar:           "avatar-url",
		Address:          []string{"Address Line 1", "Address Line 2"},
		Affiliation:      "Test Org",
		Title:            "Developer",
		Id:               "test-id",
		Homepage:         "https://example.com",
		Bio:              "Test bio",
		Tag:              "test",
		Region:           "CN",
		Language:         "zh",
		Gender:           "male",
		Birthday:         "1990-01-01",
		Education:        "Bachelor",
		Score:            100,
		Karma:            50,
		Ranking:          10,
		IsDefaultAvatar:  false,
		IsOnline:         true,
		IsAdmin:          false,
		IsGlobalAdmin:    false,
		IsForbidden:      false,
		CreatedTime:      "2024-01-01T00:00:00+08:00",
		UpdatedTime:      "2024-01-15T12:00:00+08:00",
	}

	record := service.generateUserExportRecord(user)

	// Verify all fields are properly converted to string
	assert.Equal(t, "test-org", record[0])
	assert.Equal(t, "testuser", record[1])
	assert.Equal(t, "Test User", record[2])
	assert.Equal(t, "test@example.com", record[3])
	assert.Equal(t, "1234567890", record[4])
	assert.Equal(t, "normal-user", record[5])
	assert.Equal(t, "avatar-url", record[6])
	assert.Equal(t, "Address Line 1;Address Line 2", record[7]) // Address slice to string
	assert.Equal(t, "Test Org", record[8])
	assert.Equal(t, "Developer", record[9])
	assert.Equal(t, "test-id", record[10])
	assert.Equal(t, "https://example.com", record[11])
	assert.Equal(t, "Test bio", record[12])
	assert.Equal(t, "test", record[13])
	assert.Equal(t, "CN", record[14])
	assert.Equal(t, "zh", record[15])
	assert.Equal(t, "male", record[16])
	assert.Equal(t, "1990-01-01", record[17])
	assert.Equal(t, "Bachelor", record[18])
	assert.Equal(t, "100", record[19])   // Score int to string
	assert.Equal(t, "50", record[20])    // Karma int to string
	assert.Equal(t, "10", record[21])    // Ranking int to string
	assert.Equal(t, "0", record[22])     // bool to string
	assert.Equal(t, "1", record[23])     // bool to string
	assert.Equal(t, "0", record[24])     // bool to string
	assert.Equal(t, "0", record[25])     // bool to string
	assert.Equal(t, "0", record[26])     // bool to string
	assert.Equal(t, "2024-01-01T00:00:00+08:00", record[27])
	assert.Equal(t, "2024-01-15T12:00:00+08:00", record[28])
}

// TestBatchService_ConvertSliceToString 测试切片转字符串
func TestBatchService_ConvertSliceToString(t *testing.T) {
	service := NewBatchService()

	t.Run("Non-empty slice", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		result := service.convertSliceToString(slice)
		assert.Equal(t, "a;b;c", result)
	})

	t.Run("Empty slice", func(t *testing.T) {
		slice := []string{}
		result := service.convertSliceToString(slice)
		assert.Equal(t, "", result)
	})

	t.Run("Nil slice", func(t *testing.T) {
		var slice []string
		result := service.convertSliceToString(slice)
		assert.Equal(t, "", result)
	})

	t.Run("Single element slice", func(t *testing.T) {
		slice := []string{"only one"}
		result := service.convertSliceToString(slice)
		assert.Equal(t, "only one", result)
	})
}

// TestBatchService_BoolToString 测试布尔转字符串
func TestBatchService_BoolToString(t *testing.T) {
	service := NewBatchService()

	assert.Equal(t, "1", service.boolToString(true))
	assert.Equal(t, "0", service.boolToString(false))
}

// TestBatchService_IntToString 测试整数转字符串
func TestBatchService_IntToString(t *testing.T) {
	service := NewBatchService()

	assert.Equal(t, "0", service.intToString(0))
	assert.Equal(t, "42", service.intToString(42))
	assert.Equal(t, "-100", service.intToString(-100))
}
