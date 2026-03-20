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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResponseSuccess 测试成功响应
func TestResponseSuccess(t *testing.T) {
	// 测试无数据响应
	resp := ResponseSuccess()
	assert.Equal(t, StatusOk, resp.Status)
	assert.Empty(t, resp.Msg)
	assert.Nil(t, resp.Data)

	// 测试单数据响应
	testData := map[string]string{"key": "value"}
	resp = ResponseSuccess(testData)
	assert.Equal(t, StatusOk, resp.Status)
	assert.Equal(t, testData, resp.Data)

	// 测试双数据响应
	resp = ResponseSuccess("data1", "data2")
	assert.Equal(t, StatusOk, resp.Status)
	assert.Equal(t, "data1", resp.Data)
	assert.Equal(t, "data2", resp.Data2)
}

// TestResponseError 测试错误响应
func TestResponseError(t *testing.T) {
	resp := ResponseError("error message")
	assert.Equal(t, StatusError, resp.Status)
	assert.Equal(t, "error message", resp.Msg)
	assert.Nil(t, resp.Data)

	// 测试带数据的错误响应
	testData := map[string]string{"error": "details"}
	resp = ResponseError("error message", testData)
	assert.Equal(t, StatusError, resp.Status)
	assert.Equal(t, "error message", resp.Msg)
	assert.Equal(t, testData, resp.Data)
}

// TestNewPageResponse 测试分页响应
func TestNewPageResponse(t *testing.T) {
	list := []string{"item1", "item2", "item3"}
	total := int64(100)
	page := 2
	size := 10

	resp := NewPageResponse(list, total, page, size)
	assert.Equal(t, list, resp.List)
	assert.Equal(t, total, resp.Total)
	assert.Equal(t, page, resp.Page)
	assert.Equal(t, size, resp.Size)
}

// TestResponseJsonSerialization 测试响应JSON序列化
func TestResponseJsonSerialization(t *testing.T) {
	// 测试成功响应序列化
	testData := map[string]interface{}{
		"id":   "1",
		"name": "test",
	}
	resp := ResponseSuccess(testData)

	jsonBytes, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	var parsedResp ApiResponse
	err = json.Unmarshal(jsonBytes, &parsedResp)
	assert.NoError(t, err)
	assert.Equal(t, StatusOk, parsedResp.Status)

	// 测试错误响应序列化
	errResp := ResponseError("test error")
	jsonBytes, err = json.Marshal(errResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	var parsedErrResp ApiResponse
	err = json.Unmarshal(jsonBytes, &parsedErrResp)
	assert.NoError(t, err)
	assert.Equal(t, StatusError, parsedErrResp.Status)
}

// TestConstants 测试常量定义
func TestConstants(t *testing.T) {
	assert.Equal(t, "ok", StatusOk)
	assert.Equal(t, "error", StatusError)

	assert.Equal(t, 1, DefaultPage)
	assert.Equal(t, 10, DefaultPageSize)
	assert.Equal(t, 100, MaxPageSize)
}
