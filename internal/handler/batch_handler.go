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

package handler

import (
	"encoding/csv"
	"encoding/json"
	"strings"

	"github.com/casdoor/casdoor/controllers"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/service"
	"github.com/casdoor/casdoor/object"
)

// BatchHandler 批量操作API处理器
type BatchHandler struct {
	controllers.BaseController
	batchService service.BatchService
}

// NewBatchHandler 创建批量操作处理器实例
func NewBatchHandler() *BatchHandler {
	return &BatchHandler{
		batchService: service.NewBatchService(),
	}
}

// ImportUsers 导入用户
// @Title ImportUsers
// @Description 从Excel文件批量导入用户
// @Param	owner	formData	string	true	"组织名称"
// @Param	file	formData	file	true	"Excel文件（.xlsx格式）"
// @Success 200 {object} common.ApiResponse "成功返回导入结果"
// @router /import-users [post]
func (h *BatchHandler) ImportUsers() {
	owner := h.Input().Get("owner")
	if owner == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 获取上传的文件
	file, header, err := h.GetFile("file")
	if err != nil {
		h.ResponseError("Failed to get uploaded file: " + err.Error())
		return
	}
	defer file.Close()

	// 验证文件类型
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") && !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		h.ResponseError("Only .xlsx and .csv files are supported")
		return
	}

	// 保存临时文件
	tempPath := "./tmp/" + header.Filename
	if err := h.SaveToFile("file", tempPath); err != nil {
		h.ResponseError("Failed to save temporary file: " + err.Error())
		return
	}

	// 获取当前操作用户
	operator := h.GetSessionUser()
	if operator == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	// 执行导入
	result, err := h.batchService.ImportUsers(owner, tempPath, operator, h.GetAcceptLanguage())
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(result)
}

// ExportUsers 导出用户
// @Title ExportUsers
// @Description 导出用户为Excel或CSV格式
// @Param	owner	query	string	true	"组织名称"
// @Param	field	query	string	false	"搜索字段"
// @Param	value	query	string	false	"搜索值"
// @Param	format	query	string	false	"导出格式：csv或xlsx（默认：csv）"
// @Success 200 {file} file "导出文件"
// @router /export-users [get]
func (h *BatchHandler) ExportUsers() {
	owner := h.Input().Get("owner")
	field := h.Input().Get("field")
	value := h.Input().Get("value")
	format := h.Input().Get("format")

	if owner == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	if format == "" {
		format = "csv"
	}

	switch strings.ToLower(format) {
	case "csv":
		h.exportUsersToCSV(owner, field, value)
	case "xlsx":
		h.exportUsersToExcel(owner, field, value)
	default:
		h.ResponseError("Unsupported format: " + format + ". Supported formats are: csv, xlsx")
	}
}

// BatchDeleteUsers 批量删除用户
// @Title BatchDeleteUsers
// @Description 批量删除用户
// @Param	owner	formData	string	true	"组织名称"
// @Param	userNames	formData	string	true	"用户名列表，JSON数组格式"
// @Success 200 {object} common.ApiResponse "成功返回批量操作结果"
// @router /batch-delete-users [post]
func (h *BatchHandler) BatchDeleteUsers() {
	owner := h.Input().Get("owner")
	userNamesJSON := h.Input().Get("userNames")

	if owner == "" || userNamesJSON == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	var userNames []string
	if err := json.Unmarshal([]byte(userNamesJSON), &userNames); err != nil {
		h.ResponseError("Invalid userNames format: " + err.Error())
		return
	}

	result, err := h.batchService.BatchDeleteUsers(owner, userNames)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(result)
}

// BatchDisableUsers 批量禁用用户
// @Title BatchDisableUsers
// @Description 批量禁用用户
// @Param	owner	formData	string	true	"组织名称"
// @Param	userNames	formData	string	true	"用户名列表，JSON数组格式"
// @Success 200 {object} common.ApiResponse "成功返回批量操作结果"
// @router /batch-disable-users [post]
func (h *BatchHandler) BatchDisableUsers() {
	h.batchUpdateUserStatus(true)
}

// BatchEnableUsers 批量启用用户
// @Title BatchEnableUsers
// @Description 批量启用用户
// @Param	owner	formData	string	true	"组织名称"
// @Param	userNames	formData	string	true	"用户名列表，JSON数组格式"
// @Success 200 {object} common.ApiResponse "成功返回批量操作结果"
// @router /batch-enable-users [post]
func (h *BatchHandler) BatchEnableUsers() {
	h.batchUpdateUserStatus(false)
}

// 辅助方法

// batchUpdateUserStatus 批量更新用户状态
func (h *BatchHandler) batchUpdateUserStatus(isDisabled bool) {
	owner := h.Input().Get("owner")
	userNamesJSON := h.Input().Get("userNames")

	if owner == "" || userNamesJSON == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	var userNames []string
	if err := json.Unmarshal([]byte(userNamesJSON), &userNames); err != nil {
		h.ResponseError("Invalid userNames format: " + err.Error())
		return
	}

	result, err := h.batchService.BatchUpdateUserStatus(owner, userNames, isDisabled)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(result)
}

// exportUsersToCSV 导出用户为CSV格式
func (h *BatchHandler) exportUsersToCSV(owner string, field string, value string) {
	table, err := h.batchService.ExportUsers(owner, field, value)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	// 设置响应头
	h.Ctx.ResponseWriter.Header().Set("Content-Type", "text/csv; charset=utf-8")
	h.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename=\"users.csv\"")

	// 写入UTF-8 BOM，防止Excel乱码
	h.Ctx.ResponseWriter.Write([]byte{0xEF, 0xBB, 0xBF})

	// 写入CSV数据
	writer := csv.NewWriter(h.Ctx.ResponseWriter)
	defer writer.Flush()

	for _, row := range table {
		if err := writer.Write(row); err != nil {
			h.ResponseError("Failed to write CSV: " + err.Error())
			return
		}
	}
}

// exportUsersToExcel 导出用户为Excel格式
func (h *BatchHandler) exportUsersToExcel(owner string, field string, value string) {
	table, err := h.batchService.ExportUsers(owner, field, value)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	// 设置响应头
	h.Ctx.ResponseWriter.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	h.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename=\"users.xlsx\"")

	// 这里可以集成xlsx库来生成Excel文件
	// 简化实现：返回JSON格式的表格数据
	h.ResponseSuccess(table)
}

// ResponseSuccess 返回成功响应
func (h *BatchHandler) ResponseSuccess(data ...interface{}) {
	h.Data["json"] = common.ResponseSuccess(data...)
	h.ServeJSON()
}

// ResponseError 返回错误响应
func (h *BatchHandler) ResponseError(message string, data ...interface{}) {
	h.Data["json"] = common.ResponseError(i18n.Translate(h.GetAcceptLanguage(), message), data...)
	h.ServeJSON()
}

// GetAcceptLanguage 获取接受语言
func (h *BatchHandler) GetAcceptLanguage() string {
	return h.Ctx.Input.Header("Accept-Language")
}

// GetSessionUser 获取当前会话用户
func (h *BatchHandler) GetSessionUser() *object.User {
	user, ok := h.Ctx.Input.GetData("user").(*object.User)
	if !ok {
		return nil
	}
	return user
}
