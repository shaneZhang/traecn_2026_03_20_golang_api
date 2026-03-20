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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
)

// UserService 用户业务接口
type UserService interface {
	// 基础CRUD操作
	GetUser(owner, name string) (*object.User, error)
	GetUserByUserID(owner, userID string) (*object.User, error)
	GetUserByEmail(owner, email string) (*object.User, error)
	GetUserByPhone(owner, phone string) (*object.User, error)
	GetUserByAccessKey(accessKey string) (*object.User, error)
	ListUsers(owner string, page, pageSize int, field, value, sortField, sortOrder string) ([]*object.User, int64, error)
	CreateUser(user *object.User) (bool, error)
	UpdateUser(user *object.User, columns ...string) (bool, error)
	DeleteUser(user *object.User) (bool, error)

	// 批量操作
	BatchCreateUsers(users []*object.User) (bool, error)
	BatchUpdateUsers(users []*object.User) (bool, error)

	// 高级查询
	ListGlobalUsers(page, pageSize int, field, value, sortField, sortOrder string) ([]*object.User, int64, error)
	SearchUsers(owner string, cond builder.Cond) ([]*object.User, error)
	GetSortedUsers(owner string, sorter string, limit int) ([]*object.User, error)
	GetOnlineUserCount(owner string) (int64, error)

	// 业务操作
	UpdateUserPassword(owner, name, oldPassword, newPassword string, lang string) error
	ResetUserPassword(owner, name, newPassword string) error
	DisableUser(owner, name string) (bool, error)
	EnableUser(owner, name string) (bool, error)
	UpdateUserLastSignin(owner, name, ip, city string) error

	// 权限检查
	CheckUserPermission(userId, action, object string) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户Service实例
func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
	}
}

// GetUser 获取用户
func (s *userService) GetUser(owner, name string) (*object.User, error) {
	return s.userRepo.GetByID(owner, name)
}

// GetUserByUserID 根据用户ID获取用户
func (s *userService) GetUserByUserID(owner, userID string) (*object.User, error) {
	return s.userRepo.GetByUserID(owner, userID)
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(owner, email string) (*object.User, error) {
	return s.userRepo.GetByEmail(owner, email)
}

// GetUserByPhone 根据手机号获取用户
func (s *userService) GetUserByPhone(owner, phone string) (*object.User, error) {
	return s.userRepo.GetByPhone(owner, phone)
}

// GetUserByAccessKey 根据AccessKey获取用户
func (s *userService) GetUserByAccessKey(accessKey string) (*object.User, error) {
	return s.userRepo.GetByAccessKey(accessKey)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(owner string, page, pageSize int, field, value, sortField, sortOrder string) ([]*object.User, int64, error) {
	if page < 1 {
		page = common.DefaultPage
	}
	if pageSize < 1 || pageSize > common.MaxPageSize {
		pageSize = common.DefaultPageSize
	}

	offset := (page - 1) * pageSize
	users, err := s.userRepo.List(owner, offset, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.Count(owner, field, value)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(user *object.User) (bool, error) {
	// 预处理用户数据
	preprocessUser(user)

	// 密码加密
	if user.Password != "" {
		user.Password = conf.GetEncryptedPassword(user.Password, user.Owner, user.Name)
	}

	// 处理默认值
	if user.Type == "" {
		user.Type = common.UserTypeNormal
	}
	if user.Status == "" {
		user.Status = common.StatusEnabled
	}

	return s.userRepo.Create(user)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(user *object.User, columns ...string) (bool, error) {
	// 密码加密
	if user.Password != "" {
		user.Password = conf.GetEncryptedPassword(user.Password, user.Owner, user.Name)
	}

	// 更新hash
	user.UpdateHash()

	return s.userRepo.Update(user, columns...)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(user *object.User) (bool, error) {
	return s.userRepo.Delete(user)
}

// BatchCreateUsers 批量创建用户
func (s *userService) BatchCreateUsers(users []*object.User) (bool, error) {
	for _, user := range users {
		preprocessUser(user)
		if user.Password != "" {
			user.Password = conf.GetEncryptedPassword(user.Password, user.Owner, user.Name)
		}
		if user.Type == "" {
			user.Type = common.UserTypeNormal
		}
		if user.Status == "" {
			user.Status = common.StatusEnabled
		}
	}

	return s.userRepo.CreateBatch(users)
}

// BatchUpdateUsers 批量更新用户
func (s *userService) BatchUpdateUsers(users []*object.User) (bool, error) {
	return s.userRepo.UpdateBatch(users)
}

// ListGlobalUsers 获取全局用户列表
func (s *userService) ListGlobalUsers(page, pageSize int, field, value, sortField, sortOrder string) ([]*object.User, int64, error) {
	if page < 1 {
		page = common.DefaultPage
	}
	if pageSize < 1 || pageSize > common.MaxPageSize {
		pageSize = common.DefaultPageSize
	}

	offset := (page - 1) * pageSize
	users, err := s.userRepo.GetGlobalList(offset, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.GetGlobalCount(field, value)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SearchUsers 搜索用户
func (s *userService) SearchUsers(owner string, cond builder.Cond) ([]*object.User, error) {
	return s.userRepo.GetWithFilter(owner, cond)
}

// GetSortedUsers 获取排序后的用户列表
func (s *userService) GetSortedUsers(owner string, sorter string, limit int) ([]*object.User, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.userRepo.GetSorted(owner, sorter, limit)
}

// GetOnlineUserCount 获取在线用户数量
func (s *userService) GetOnlineUserCount(owner string) (int64, error) {
	return s.userRepo.GetOnlineCount(owner, 1)
}

// UpdateUserPassword 更新用户密码
func (s *userService) UpdateUserPassword(owner, name, oldPassword, newPassword string, lang string) error {
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	// 验证旧密码
	if oldPassword != "" {
		oldPasswordEncrypted := conf.GetEncryptedPassword(oldPassword, owner, name)
		if user.Password != oldPasswordEncrypted {
			return fmt.Errorf(i18n.Translate(lang, "user:Wrong password"))
		}
	}

	// 加密新密码
	newPasswordEncrypted := conf.GetEncryptedPassword(newPassword, owner, name)
	user.Password = newPasswordEncrypted
	user.UpdateHash()

	_, err = s.userRepo.Update(user, "password", "hash")
	return err
}

// ResetUserPassword 重置用户密码
func (s *userService) ResetUserPassword(owner, name, newPassword string) error {
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	user.Password = conf.GetEncryptedPassword(newPassword, owner, name)
	user.UpdateHash()

	_, err = s.userRepo.Update(user, "password", "hash")
	return err
}

// DisableUser 禁用用户
func (s *userService) DisableUser(owner, name string) (bool, error) {
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	user.IsDisabled = true
	user.UpdateHash()

	return s.userRepo.Update(user, "is_disabled", "hash")
}

// EnableUser 启用用户
func (s *userService) EnableUser(owner, name string) (bool, error) {
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	user.IsDisabled = false
	user.UpdateHash()

	return s.userRepo.Update(user, "is_disabled", "hash")
}

// UpdateUserLastSignin 更新用户最后登录信息
func (s *userService) UpdateUserLastSignin(owner, name, ip, city string) error {
	user, err := s.userRepo.GetByID(owner, name)
	if err != nil {
		return err
	}

	user.LastSigninTime = util.GetCurrentTime()
	user.LastSigninIp = ip
	user.LastSigninCity = city
	user.SigninCount += 1
	user.UpdateHash()

	_, err = s.userRepo.Update(user, "last_signin_time", "last_signin_ip", "last_signin_city", "signin_count", "hash")
	return err
}

// CheckUserPermission 检查用户权限
func (s *userService) CheckUserPermission(userId, action, object string) (bool, error) {
	// 这里可以集成Casbin等权限框架
	// 简化实现，默认返回true
	return true, nil
}

// preprocessUser 预处理用户数据
func preprocessUser(user *object.User) {
	// 用户名转小写
	isUsernameLowered := conf.GetConfigBool("isUsernameLowered")
	if isUsernameLowered {
		user.Name = strings.ToLower(user.Name)
	}

	// 邮箱转小写
	if user.Email != "" {
		user.Email = strings.ToLower(user.Email)
	}

	// 格式化手机号
	if user.Phone != "" {
		user.Phone = util.GetSeperatedPhone(user.Phone)
	}

	// 设置默认显示名称
	if user.DisplayName == "" {
		user.DisplayName = user.Name
	}

	// 解析首选项
	if user.Preferences != "" {
		var preferences map[string]interface{}
		if err := json.Unmarshal([]byte(user.Preferences), &preferences); err == nil {
			if theme, ok := preferences["theme"].(string); ok && theme != "" {
				user.PreferredTheme = theme
			}
			if language, ok := preferences["language"].(string); ok && language != "" {
				user.PreferredLanguage = language
			}
		}
	}

	// 更新hash
	user.UpdateHash()
}
