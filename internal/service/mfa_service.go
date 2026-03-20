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
	"fmt"

	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
)

// MfaService MFA认证服务接口
type MfaService interface {
	// MFA配置管理
	GetMfaConfigurations(userID string, masked bool) ([]*object.MfaProps, error)
	GetMfaConfiguration(userID string, mfaType string, masked bool) (*object.MfaProps, error)

	// MFA设置流程
	InitiateMfaSetup(userID string, mfaType string) (*object.MfaProps, error)
	VerifyMfaSetup(userID string, mfaType string, passcode string) error
	EnableMfa(userID string, mfaType string, passcode string) error
	DisableMfa(userID string, mfaType string) error
	DisableAllMfa(userID string) error

	// MFA验证流程
	InitiateMfaVerification(userID string, mfaType string) error
	VerifyMfa(userID string, mfaType string, passcode string) (bool, error)
	VerifyRecoveryCode(userID string, recoveryCode string) (bool, error)

	// 偏好设置
	SetPreferredMfa(userID string, mfaType string) error
	GetPreferredMfa(userID string) (string, error)

	// 恢复码管理
	GenerateRecoveryCodes(userID string) ([]string, error)
	GetRecoveryCodes(userID string) ([]string, error)
	RegenerateRecoveryCodes(userID string) ([]string, error)
}

type mfaService struct {
	userRepo repository.UserRepository
}

// NewMfaService 创建MFA服务实例
func NewMfaService() MfaService {
	return &mfaService{
		userRepo: repository.NewUserRepository(),
	}
}

// GetMfaConfigurations 获取用户所有MFA配置
func (s *mfaService) GetMfaConfigurations(userID string, masked bool) ([]*object.MfaProps, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return nil, err
	}

	return object.GetAllMfaProps(user, masked), nil
}

// GetMfaConfiguration 获取用户指定类型的MFA配置
func (s *mfaService) GetMfaConfiguration(userID string, mfaType string, masked bool) (*object.MfaProps, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return nil, err
	}

	return user.GetMfaProps(mfaType, masked), nil
}

// InitiateMfaSetup 初始化MFA设置
func (s *mfaService) InitiateMfaSetup(userID string, mfaType string) (*object.MfaProps, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return nil, err
	}

	// 检查该类型MFA是否已启用
	mfaProps := user.GetMfaProps(mfaType, false)
	if mfaProps.Enabled {
		return nil, common.ErrMfaAlreadyEnabled
	}

	// 获取MFA工具
	issuer := owner // 使用组织名称作为颁发者
	mfaUtil := object.GetMfaUtil(mfaType, nil)
	if mfaUtil == nil {
		return nil, fmt.Errorf("unsupported MFA type: %s", mfaType)
	}

	// 初始化MFA设置
	props, err := mfaUtil.Initiate(userID, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate MFA setup: %v", err)
	}

	// 临时存储配置（这里可以考虑缓存）
	// 注意：实际应用中需要将这些临时配置存储在会话或缓存中
	return props, nil
}

// VerifyMfaSetup 验证MFA设置（启用前的验证）
func (s *mfaService) VerifyMfaSetup(userID string, mfaType string, passcode string) error {
	mfaUtil := object.GetMfaUtil(mfaType, nil)
	if mfaUtil == nil {
		return fmt.Errorf("unsupported MFA type: %s", mfaType)
	}

	return mfaUtil.SetupVerify(passcode)
}

// EnableMfa 启用MFA
func (s *mfaService) EnableMfa(userID string, mfaType string, passcode string) error {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	// 先验证验证码
	if err := s.VerifyMfaSetup(userID, mfaType, passcode); err != nil {
		return common.ErrMfaInvalidCode
	}

	// 获取MFA工具并启用
	mfaUtil := object.GetMfaUtil(mfaType, nil)
	if mfaUtil == nil {
		return fmt.Errorf("unsupported MFA type: %s", mfaType)
	}

	if err := mfaUtil.Enable(user); err != nil {
		return fmt.Errorf("failed to enable MFA: %v", err)
	}

	// 如果是第一个启用的MFA，设置为首选
	if user.PreferredMfaType == "" {
		user.PreferredMfaType = mfaType
	}

	// 生成恢复码
	recoveryCodes, err := s.GenerateRecoveryCodes(userID)
	if err != nil {
		return fmt.Errorf("failed to generate recovery codes: %v", err)
	}
	user.RecoveryCodes = recoveryCodes

	// 更新用户信息
	user.UpdateHash()
	_, err = s.userRepo.Update(user, s.getMfaUpdateColumns()...)
	return err
}

// DisableMfa 禁用指定类型的MFA
func (s *mfaService) DisableMfa(userID string, mfaType string) error {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	// 禁用指定类型的MFA
	switch mfaType {
	case object.SmsType:
		user.MfaPhoneEnabled = false
	case object.EmailType:
		user.MfaEmailEnabled = false
	case object.TotpType:
		user.TotpSecret = ""
	case object.RadiusType:
		user.MfaRadiusEnabled = false
		user.MfaRadiusUsername = ""
		user.MfaRadiusProvider = ""
	case object.PushType:
		user.MfaPushEnabled = false
		user.MfaPushReceiver = ""
		user.MfaPushProvider = ""
	default:
		return fmt.Errorf("unsupported MFA type: %s", mfaType)
	}

	// 如果禁用的是首选MFA，清除首选设置
	if user.PreferredMfaType == mfaType {
		user.PreferredMfaType = ""
		// 尝试设置其他可用的MFA为首选
		user.PreferredMfaType = s.findAvailableMfaType(user)
	}

	// 如果没有启用的MFA，清除恢复码
	if !s.hasAnyMfaEnabled(user) {
		user.RecoveryCodes = []string{}
	}

	user.UpdateHash()
	_, err = s.userRepo.Update(user, s.getMfaUpdateColumns()...)
	return err
}

// DisableAllMfa 禁用所有MFA
func (s *mfaService) DisableAllMfa(userID string) error {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	return object.DisabledMultiFactorAuth(user)
}

// InitiateMfaVerification 初始化MFA验证
func (s *mfaService) InitiateMfaVerification(userID string, mfaType string) error {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	// 检查MFA是否已启用
	mfaProps := user.GetMfaProps(mfaType, false)
	if !mfaProps.Enabled {
		return common.ErrMfaNotEnabled
	}

	// 对于需要发送验证码的类型（短信、邮件、推送），触发发送
	switch mfaType {
	case object.SmsType, object.EmailType, object.PushType:
		mfaUtil := object.GetMfaUtil(mfaType, mfaProps)
		if mfaUtil == nil {
			return fmt.Errorf("unsupported MFA type: %s", mfaType)
		}
		// 这里可以触发发送验证码的逻辑
		// 注意：实际发送逻辑可能需要在controller层处理
	}

	return nil
}

// VerifyMfa 验证MFA验证码
func (s *mfaService) VerifyMfa(userID string, mfaType string, passcode string) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	// 检查MFA是否已启用
	mfaProps := user.GetMfaProps(mfaType, false)
	if !mfaProps.Enabled {
		return false, common.ErrMfaNotEnabled
	}

	// 使用对应MFA工具验证
	mfaUtil := object.GetMfaUtil(mfaType, mfaProps)
	if mfaUtil == nil {
		return false, fmt.Errorf("unsupported MFA type: %s", mfaType)
	}

	if err := mfaUtil.Verify(passcode); err != nil {
		return false, common.ErrMfaInvalidCode
	}

	return true, nil
}

// VerifyRecoveryCode 使用恢复码验证
func (s *mfaService) VerifyRecoveryCode(userID string, recoveryCode string) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	// 使用恢复码验证
	err = object.MfaRecover(user, recoveryCode)
	if err != nil {
		return false, err
	}

	return true, nil
}

// SetPreferredMfa 设置首选MFA类型
func (s *mfaService) SetPreferredMfa(userID string, mfaType string) error {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	// 检查该类型MFA是否已启用
	mfaProps := user.GetMfaProps(mfaType, false)
	if !mfaProps.Enabled {
		return fmt.Errorf("MFA type %s is not enabled", mfaType)
	}

	return object.SetPreferredMultiFactorAuth(user, mfaType)
}

// GetPreferredMfa 获取首选MFA类型
func (s *mfaService) GetPreferredMfa(userID string) (string, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return "", err
	}

	return user.PreferredMfaType, nil
}

// GenerateRecoveryCodes 生成恢复码
func (s *mfaService) GenerateRecoveryCodes(userID string) ([]string, error) {
	// 生成10个随机恢复码
	recoveryCodes := make([]string, 10)
	for i := range recoveryCodes {
		recoveryCodes[i] = util.GenerateRecoveryCode()
	}
	return recoveryCodes, nil
}

// GetRecoveryCodes 获取恢复码
func (s *mfaService) GetRecoveryCodes(userID string) ([]string, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return nil, err
	}

	return user.RecoveryCodes, nil
}

// RegenerateRecoveryCodes 重新生成恢复码
func (s *mfaService) RegenerateRecoveryCodes(userID string) ([]string, error) {
	owner, name := util.GetOwnerAndNameFromId(userID)
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return nil, err
	}

	// 检查是否有启用的MFA
	if !s.hasAnyMfaEnabled(user) {
		return nil, fmt.Errorf("no MFA method enabled")
	}

	// 生成新的恢复码
	recoveryCodes, err := s.GenerateRecoveryCodes(userID)
	if err != nil {
		return nil, err
	}

	// 更新用户恢复码
	user.RecoveryCodes = recoveryCodes
	user.UpdateHash()
	_, err = s.userRepo.Update(user, "recovery_codes", "hash")
	if err != nil {
		return nil, err
	}

	return recoveryCodes, nil
}

// 辅助方法

// getMfaUpdateColumns 获取MFA相关的更新列
func (s *mfaService) getMfaUpdateColumns() []string {
	return []string{
		"preferred_mfa_type",
		"recovery_codes",
		"mfa_phone_enabled",
		"mfa_email_enabled",
		"totp_secret",
		"mfa_radius_enabled",
		"mfa_radius_username",
		"mfa_radius_provider",
		"mfa_push_enabled",
		"mfa_push_receiver",
		"mfa_push_provider",
		"hash",
	}
}

// hasAnyMfaEnabled 检查是否有任何MFA已启用
func (s *mfaService) hasAnyMfaEnabled(user *object.User) bool {
	return user.MfaPhoneEnabled ||
		user.MfaEmailEnabled ||
		user.TotpSecret != "" ||
		user.MfaRadiusEnabled ||
		user.MfaPushEnabled
}

// findAvailableMfaType 查找可用的MFA类型
func (s *mfaService) findAvailableMfaType(user *object.User) string {
	mfaTypes := []string{object.TotpType, object.SmsType, object.EmailType, object.RadiusType, object.PushType}
	for _, mfaType := range mfaTypes {
		props := user.GetMfaProps(mfaType, false)
		if props.Enabled {
			return mfaType
		}
	}
	return ""
}
