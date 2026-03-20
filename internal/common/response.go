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

package common

import (
	"net/http"

	"github.com/beego/beego/v2/server/web/context"
)

// Response represents standard API response
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Status     string      `json:"status"`
	Code       int         `json:"code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination represents pagination info
type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"totalPages"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Status  string            `json:"status"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// Success returns success response
func Success(data interface{}) *Response {
	return &Response{
		Status: "ok",
		Code:   http.StatusOK,
		Data:   data,
	}
}

// SuccessWithPagination returns success response with pagination
func SuccessWithPagination(data interface{}, pagination *Pagination) *PaginatedResponse {
	return &PaginatedResponse{
		Status:     "ok",
		Code:       http.StatusOK,
		Data:       data,
		Pagination: pagination,
	}
}

// Error returns error response
func Error(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: message,
	}
}

// ErrorWithDetails returns error response with details
func ErrorWithDetails(code int, message string, errors map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}

// Created returns created response
func Created(data interface{}) *Response {
	return &Response{
		Status: "ok",
		Code:   http.StatusCreated,
		Data:   data,
	}
}

// BeegoResponseWriter wraps beego context for response writing
type BeegoResponseWriter struct {
	Ctx *context.Context
}

// WriteJSON writes JSON response
func (w *BeegoResponseWriter) WriteJSON(data interface{}) {
	w.Ctx.Output.JSON(data, true, true)
}

// WriteError writes error response
func (w *BeegoResponseWriter) WriteError(code int, err error) {
	resp := Error(code, err.Error())
	w.Ctx.Output.SetStatus(code)
	w.Ctx.Output.JSON(resp, true, true)
}
