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

// ApiResponse 统一API响应格式
type ApiResponse struct {
	Status string      `json:"status"` // "ok" or "error"
	Msg    string      `json:"msg,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Data2  interface{} `json:"data2,omitempty"`
}

// PageResponse 分页响应格式
type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// ResponseSuccess 成功响应
func ResponseSuccess(data ...interface{}) *ApiResponse {
	resp := &ApiResponse{
		Status: "ok",
	}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}
	return resp
}

// ResponseSuccessWithPage 带分页信息的成功响应
func ResponseSuccessWithPage(list interface{}, total int64, page, size int) *ApiResponse {
	return ResponseSuccess(&PageResponse{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}

// ResponseError 错误响应
func ResponseError(message string, data ...interface{}) *ApiResponse {
	resp := &ApiResponse{
		Status: "error",
		Msg:    message,
	}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}
	return resp
}

// ResponseErrorWithCode 带错误码的错误响应
func ResponseErrorWithCode(err *BusinessError, data ...interface{}) *ApiResponse {
	resp := &ApiResponse{
		Status: "error",
		Msg:    err.Message,
	}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}
	return resp
}

// ResponseFromAction 从操作结果生成响应
func ResponseFromAction(affected bool, err ...error) *ApiResponse {
	if len(err) != 0 && err[0] != nil {
		return ResponseError(err[0].Error())
	}
	if affected {
		return ResponseSuccess("Affected")
	}
	return ResponseSuccess("Unaffected")
}
