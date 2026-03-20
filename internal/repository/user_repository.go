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

package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/core"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	GetByID(owner, name string) (*object.User, error)
	GetByUserID(owner, userID string) (*object.User, error)
	GetByEmail(owner, email string) (*object.User, error)
	GetByPhone(owner, phone string) (*object.User, error)
	GetByAccessKey(accessKey string) (*object.User, error)
	List(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error)
	Count(owner, field, value string) (int64, error)
	Create(user *object.User) (bool, error)
	Update(user *object.User, columns ...string) (bool, error)
	Delete(user *object.User) (bool, error)
	CreateBatch(users []*object.User) (bool, error)
	GetGlobalList(offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error)
	GetGlobalCount(field, value string) (int64, error)
	GetWithFilter(owner string, cond builder.Cond) ([]*object.User, error)
	GetSorted(owner string, sorter string, limit int) ([]*object.User, error)
	GetOnlineCount(owner string, isOnline int) (int64, error)
	CheckDuplicate(user *object.User) error
}

type userRepository struct {
}

// NewUserRepository 创建用户Repository实例
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(owner, name string) (*object.User, error) {
	if owner == "" || name == "" {
		return nil, common.ErrBadRequest
	}

	user := &object.User{Owner: owner, Name: name}
	existed, err := common.GetByID(owner, name, user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// GetByUserID 根据用户ID获取用户
func (r *userRepository) GetByUserID(owner, userID string) (*object.User, error) {
	if owner == "" || userID == "" {
		return nil, common.ErrBadRequest
	}

	user := &object.User{Owner: owner, Id: userID}
	existed, err := object.GetEngine().Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(owner, email string) (*object.User, error) {
	if email == "" {
		return nil, common.ErrBadRequest
	}

	user := &object.User{Email: email}
	if owner != "" {
		user.Owner = owner
	}

	existed, err := object.GetEngine().Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// GetByPhone 根据手机号获取用户
func (r *userRepository) GetByPhone(owner, phone string) (*object.User, error) {
	if phone == "" {
		return nil, common.ErrBadRequest
	}

	phone = util.GetSeperatedPhone(phone)
	user := &object.User{Phone: phone}
	if owner != "" {
		user.Owner = owner
	}

	existed, err := object.GetEngine().Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// List 获取用户列表
func (r *userRepository) List(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error) {
	users := make([]*object.User, 0)

	sb := common.NewSessionBuilder(owner).
		SetPagination(offset, limit).
		SetFilter(field, value).
		SetSort(sortField, sortOrder)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Count 获取用户数量
func (r *userRepository) Count(owner, field, value string) (int64, error) {
	sb := common.NewSessionBuilder(owner).
		SetFilter(field, value)

	session := sb.Build()
	defer session.Close()

	return session.Count(&object.User{})
}

// Create 创建用户
func (r *userRepository) Create(user *object.User) (bool, error) {
	if err := r.CheckDuplicate(user); err != nil {
		return false, err
	}

	isUsernameLowered := conf.GetConfigBool("isUsernameLowered")
	if isUsernameLowered {
		user.Name = strings.ToLower(user.Name)
	}

	affected, err := object.GetEngine().Insert(user)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// Update 更新用户
func (r *userRepository) Update(user *object.User, columns ...string) (bool, error) {
	if user.Owner == "" || user.Name == "" {
		return false, common.ErrBadRequest
	}

	session := object.GetEngine().ID(core.PK{user.Owner, user.Name})
	if len(columns) > 0 {
		if !util.InSlice(columns, "hash") {
			columns = append(columns, "hash")
		}
		session = session.Cols(columns...)
	}

	affected, err := session.Update(user)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// Delete 删除用户
func (r *userRepository) Delete(user *object.User) (bool, error) {
	if user.Owner == "" || user.Name == "" {
		return false, common.ErrBadRequest
	}

	affected, err := object.GetEngine().ID(core.PK{user.Owner, user.Name}).Delete(&object.User{})
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// CreateBatch 批量创建用户
func (r *userRepository) CreateBatch(users []*object.User) (bool, error) {
	if len(users) == 0 {
		return false, errors.New("no users provided")
	}

	batchSize := conf.GetConfigBatchSize()
	affected := false

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		count, err := object.GetEngine().Insert(batch)
		if err != nil {
			return false, err
		}
		if count > 0 {
			affected = true
		}
	}

	return affected, nil
}

// GetGlobalList 获取全局用户列表
func (r *userRepository) GetGlobalList(offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error) {
	users := make([]*object.User, 0)

	sb := common.NewSessionBuilder("").
		SetPagination(offset, limit).
		SetFilter(field, value).
		SetSort(sortField, sortOrder)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetGlobalCount 获取全局用户数量
func (r *userRepository) GetGlobalCount(field, value string) (int64, error) {
	sb := common.NewSessionBuilder("").
		SetFilter(field, value)

	session := sb.Build()
	defer session.Close()

	return session.Count(&object.User{})
}

// GetWithFilter 带条件查询用户列表
func (r *userRepository) GetWithFilter(owner string, cond builder.Cond) ([]*object.User, error) {
	users := make([]*object.User, 0)

	sb := common.NewSessionBuilder(owner).
		AddCondition(cond)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetSorted 获取排序后的用户列表
func (r *userRepository) GetSorted(owner string, sorter string, limit int) ([]*object.User, error) {
	users := make([]*object.User, 0)

	session := object.GetEngine().Desc(sorter).Limit(limit, 0)
	if owner != "" {
		session = session.Where("owner = ?", owner)
	}
	defer session.Close()

	err := session.Find(&users, &object.User{Owner: owner})
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetOnlineCount 获取在线用户数量
func (r *userRepository) GetOnlineCount(owner string, isOnline int) (int64, error) {
	session := object.GetEngine().Where("is_online = ?", isOnline)
	if owner != "" {
		session = session.And("owner = ?", owner)
	}
	defer session.Close()

	return session.Count(&object.User{})
}

// CheckDuplicate 检查用户唯一性
func (r *userRepository) CheckDuplicate(user *object.User) error {
	existedUser := &object.User{Owner: user.Owner, Name: user.Name}
	existed, err := object.GetEngine().Get(existedUser)
	if err != nil {
		return err
	}
	if existed {
		return common.NewBusinessError(common.ErrCodeUsernameAlreadyExists,
			fmt.Sprintf("The username %s is already registered", user.Name))
	}

	if user.Email != "" {
		existedUser := &object.User{Owner: user.Owner, Email: user.Email}
		existed, err := object.GetEngine().Get(existedUser)
		if err != nil {
			return err
		}
		if existed {
			return common.NewBusinessError(common.ErrCodeEmailAlreadyExists,
				fmt.Sprintf("The email %s is already registered", user.Email))
		}
	}

	if user.Phone != "" {
		existedUser := &object.User{Owner: user.Owner, Phone: user.Phone}
		existed, err := object.GetEngine().Get(existedUser)
		if err != nil {
			return err
		}
		if existed {
			return common.NewBusinessError(common.ErrCodePhoneAlreadyExists,
				fmt.Sprintf("The phone %s is already registered", user.Phone))
		}
	}

	return nil
}
