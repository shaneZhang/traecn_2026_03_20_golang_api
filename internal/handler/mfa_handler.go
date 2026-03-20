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
	"encoding/json"

	"github.com/casdoor/casdoor/controllers"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/service"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

// MfaHandler MFA认证API处理器
type MfaHandler struct {
	controllers.BaseController
	mfaService service.MfaService
}

// NewMfaHandler 创建MFA处理器实例
func NewMfaHandler() *MfaHandler {
	return &MfaHandler{
		mfaService: service.NewMfaService(),
	}
}

// MfaSetupInitiate 初始化MFA设置
// @Title MfaSetupInitiate
// @Description 初始化MFA设置，获取配置信息
// @Param	mfaType	query	string	true	"MFA类型：app(TOTP), sms, email, radius, push"
// @Success 200 {object} common.ApiResponse "成功返回MFA配置信息"
// @router /mfa/setup/initiate [post]
func (h *MfaHandler) MfaSetupInitiate() {
	mfaType := h.Input().Get("mfaType")
	if mfaType == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	props, err := h.mfaService.InitiateMfaSetup(userID, mfaType)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(props)
}

// MfaSetupVerify 验证MFA设置
// @Title MfaSetupVerify
// @Description 验证MFA验证码（启用前的验证）
// @Param	mfaType		query	string	true	"MFA类型"
// @Param	passcode	query	string	true	"验证码"
// @Success 200 {object} common.ApiResponse "成功返回验证结果"
// @router /mfa/setup/verify [post]
func (h *MfaHandler) MfaSetupVerify() {
	mfaType := h.Input().Get("mfaType")
	passcode := h.Input().Get("passcode")

	if mfaType == "" || passcode == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	err := h.mfaService.VerifyMfaSetup(userID, mfaType, passcode)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess("Verification successful")
}

// MfaSetupEnable 启用MFA
// @Title MfaSetupEnable
// @Description 启用指定类型的MFA
// @Param	mfaType		query	string	true	"MFA类型"
// @Param	passcode	query	string	true	"验证码"
// @Success 200 {object} common.ApiResponse "成功返回启用结果"
// @router /mfa/setup/enable [post]
func (h *MfaHandler) MfaSetupEnable() {
	mfaType := h.Input().Get("mfaType")
	passcode := h.Input().Get("passcode")

	if mfaType == "" || passcode == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	err := h.mfaService.EnableMfa(userID, mfaType, passcode)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	// 获取恢复码
	recoveryCodes, err := h.mfaService.GetRecoveryCodes(userID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{
		"status":        "enabled",
		"recoveryCodes": recoveryCodes,
	})
}

// GetMfaStatus 获取MFA状态
// @Title GetMfaStatus
// @Description 获取用户的MFA配置状态
// @Param	masked	query	bool	false	"是否掩码处理敏感信息"
// @Success 200 {object} common.ApiResponse "成功返回MFA状态"
// @router /mfa/status [get]
func (h *MfaHandler) GetMfaStatus() {
	masked := h.GetBoolFromInput("masked", true)

	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	configs, err := h.mfaService.GetMfaConfigurations(userID, masked)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	preferredMfa, err := h.mfaService.GetPreferredMfa(userID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{
		"configurations": configs,
		"preferredMfa":   preferredMfa,
	})
}

// DeleteMfa 禁用MFA
// @Title DeleteMfa
// @Description 禁用指定类型的MFA
// @Param	mfaType	query	string	true	"MFA类型"
// @Success 200 {object} common.ApiResponse "成功返回禁用结果"
// @router /mfa/delete [post]
func (h *MfaHandler) DeleteMfa() {
	mfaType := h.Input().Get("mfaType")
	if mfaType == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	err := h.mfaService.DisableMfa(userID, mfaType)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess("MFA disabled successfully")
}

// DeleteAllMfa 禁用所有MFA
// @Title DeleteAllMfa
// @Description 禁用所有类型的MFA
// @Success 200 {object} common.ApiResponse "成功返回禁用结果"
// @router /mfa/delete-all [post]
func (h *MfaHandler) DeleteAllMfa() {
	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	err := h.mfaService.DisableAllMfa(userID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess("All MFA disabled successfully")
}

// SetPreferredMfa 设置首选MFA
// @Title SetPreferredMfa
// @Description 设置首选的MFA类型
// @Param	mfaType	query	string	true	"MFA类型"
// @Success 200 {object} common.ApiResponse "成功返回设置结果"
// @router /mfa/set-preferred [post]
func (h *MfaHandler) SetPreferredMfa() {
	mfaType := h.Input().Get("mfaType")
	if mfaType == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	err := h.mfaService.SetPreferredMfa(userID, mfaType)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess("Preferred MFA set successfully")
}

// MfaVerify 验证MFA
// @Title MfaVerify
// @Description MFA登录验证
// @Param	mfaType		query	string	true	"MFA类型"
// @Param	passcode	query	string	true	"验证码"
// @Param	userId		query	string	true	"用户ID"
// @Success 200 {object} common.ApiResponse "成功返回验证结果"
// @router /mfa/verify [post]
func (h *MfaHandler) MfaVerify() {
	mfaType := h.Input().Get("mfaType")
	passcode := h.Input().Get("passcode")
	userID := h.Input().Get("userId")

	if mfaType == "" || passcode == "" || userID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	valid, err := h.mfaService.VerifyMfa(userID, mfaType, passcode)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	if valid {
		h.ResponseSuccess("Verification successful")
	} else {
		h.ResponseError("Verification failed")
	}
}

// MfaRecover 使用恢复码验证
// @Title MfaRecover
// @Description 使用恢复码进行MFA验证
// @Param	recoveryCode	query	string	true	"恢复码"
// @Param	userId			query	string	true	"用户ID"
// @Success 200 {object} common.ApiResponse "成功返回验证结果"
// @router /mfa/recover [post]
func (h *MfaHandler) MfaRecover() {
	recoveryCode := h.Input().Get("recoveryCode")
	userID := h.Input().Get("userId")

	if recoveryCode == "" || userID == "" {
		h.ResponseError(common.ErrBadRequest.Message)
		return
	}

	valid, err := h.mfaService.VerifyRecoveryCode(userID, recoveryCode)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	if valid {
		h.ResponseSuccess("Recovery code verified successfully")
	} else {
		h.ResponseError("Invalid recovery code")
	}
}

// RegenerateRecoveryCodes 重新生成恢复码
// @Title RegenerateRecoveryCodes
// @Description 重新生成MFA恢复码
// @Success 200 {object} common.ApiResponse "成功返回新的恢复码"
// @router /mfa/regenerate-recovery-codes [post]
func (h *MfaHandler) RegenerateRecoveryCodes() {
	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	recoveryCodes, err := h.mfaService.RegenerateRecoveryCodes(userID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{
		"recoveryCodes": recoveryCodes,
	})
}

// GetRecoveryCodes 获取恢复码
// @Title GetRecoveryCodes
// @Description 获取当前的MFA恢复码
// @Success 200 {object} common.ApiResponse "成功返回恢复码"
// @router /mfa/recovery-codes [get]
func (h *MfaHandler) GetRecoveryCodes() {
	// 获取当前用户
	user := h.GetSessionUser()
	if user == nil {
		h.ResponseError(common.ErrUnauthorized.Message)
		return
	}

	userID := util.GetId(user.Owner, user.Name)
	recoveryCodes, err := h.mfaService.GetRecoveryCodes(userID)
	if err != nil {
		h.ResponseError(err.Error())
		return
	}

	h.ResponseSuccess(map[string]interface{}{
		"recoveryCodes": recoveryCodes,
	})
}

// 辅助方法

// GetBoolFromInput 从输入中获取布尔值
func (h *MfaHandler) GetBoolFromInput(key string, defaultValue bool) bool {
	value := h.Input().Get(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

// ResponseSuccess 返回成功响应
func (h *MfaHandler) ResponseSuccess(data ...interface{}) {
	h.Data["json"] = common.ResponseSuccess(data...)
	h.ServeJSON()
}

// ResponseError 返回错误响应
func (h *MfaHandler) ResponseError(message string, data ...interface{}) {
	h.Data["json"] = common.ResponseError(i18n.Translate(h.GetAcceptLanguage(), message), data...)
	h.ServeJSON()
}

// GetAcceptLanguage 获取接受语言
func (h *MfaHandler) GetAcceptLanguage() string {
	return h.Ctx.Input.Header("Accept-Language")
}

// GetSessionUser 获取当前会话用户
func (h *MfaHandler) GetSessionUser() *object.User {
	user, ok := h.Ctx.Input.GetData("user").(*object.User)
	if !ok {
		return nil
	}
	return user
}
