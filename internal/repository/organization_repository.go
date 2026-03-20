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
	"fmt"

	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/object"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/core"
)

// OrganizationRepository 组织数据访问接口
type OrganizationRepository interface {
	// 基础CRUD操作
	GetByID(name string) (*object.Organization, error)
	GetByUserID(userID string) (*object.Organization, error)
	List(offset, limit int, field, value, sortField, sortOrder string) ([]*object.Organization, error)
	Count(field, value string) (int64, error)
	Create(org *object.Organization) (bool, error)
	Update(org *object.Organization, columns ...string) (bool, error)
	Delete(org *object.Organization) (bool, error)

	// 特殊查询
	GetAll() ([]*object.Organization, error)
	GetWithFilter(cond builder.Cond) ([]*object.Organization, error)
	GetByAccount(account string) (*object.Organization, error)
	GetByOwner(owner string) ([]*object.Organization, error)

	// 层级关系
	GetParentOrganizations(orgName string) ([]*object.Organization, error)
	GetChildOrganizations(orgName string) ([]*object.Organization, error)
	GetAllChildOrganizations(orgName string) ([]*object.Organization, error)

	// 统计
	GetUserCount(orgName string) (int64, error)
	GetApplicationCount(orgName string) (int64, error)

	// 唯一性检查
	CheckDuplicate(org *object.Organization) error
}

type organizationRepository struct {
}

// NewOrganizationRepository 创建组织Repository实例
func NewOrganizationRepository() OrganizationRepository {
	return &organizationRepository{}
}

// GetByID 根据ID获取组织
func (r *organizationRepository) GetByID(name string) (*object.Organization, error) {
	if name == "" {
		return nil, common.ErrBadRequest
	}

	org := &object.Organization{Owner: common.AdminUser, Name: name}
	existed, err := common.GetByID(common.AdminUser, name, org)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrOrgNotFound
	}
	return org, nil
}

// GetByUserID 根据用户ID获取组织
func (r *organizationRepository) GetByUserID(userID string) (*object.Organization, error) {
	if userID == "" {
		return nil, common.ErrBadRequest
	}

	user := &object.User{Id: userID}
	existed, err := object.GetEngine().Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrUserNotFound
	}

	return r.GetByID(user.Owner)
}

// List 获取组织列表
func (r *organizationRepository) List(offset, limit int, field, value, sortField, sortOrder string) ([]*object.Organization, error) {
	orgs := make([]*object.Organization, 0)

	sb := common.NewSessionBuilder(common.AdminUser).
		SetPagination(offset, limit).
		SetFilter(field, value).
		SetSort(sortField, sortOrder)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&orgs)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

// Count 获取组织数量
func (r *organizationRepository) Count(field, value string) (int64, error) {
	sb := common.NewSessionBuilder(common.AdminUser).
		SetFilter(field, value)

	session := sb.Build()
	defer session.Close()

	return session.Count(&object.Organization{})
}

// Create 创建组织
func (r *organizationRepository) Create(org *object.Organization) (bool, error) {
	if err := r.CheckDuplicate(org); err != nil {
		return false, err
	}

	affected, err := object.GetEngine().Insert(org)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// Update 更新组织
func (r *organizationRepository) Update(org *object.Organization, columns ...string) (bool, error) {
	if org.Name == "" {
		return false, common.ErrBadRequest
	}

	session := object.GetEngine().ID(core.PK{common.AdminUser, org.Name})
	if len(columns) > 0 {
		session = session.Cols(columns...)
	}

	affected, err := session.Update(org)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// Delete 删除组织
func (r *organizationRepository) Delete(org *object.Organization) (bool, error) {
	if org.Name == "" {
		return false, common.ErrBadRequest
	}

	affected, err := object.GetEngine().ID(core.PK{common.AdminUser, org.Name}).Delete(&object.Organization{})
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// GetAll 获取所有组织
func (r *organizationRepository) GetAll() ([]*object.Organization, error) {
	orgs := make([]*object.Organization, 0)
	err := object.GetEngine().Find(&orgs, &object.Organization{Owner: common.AdminUser})
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetWithFilter 带条件查询组织列表
func (r *organizationRepository) GetWithFilter(cond builder.Cond) ([]*object.Organization, error) {
	orgs := make([]*object.Organization, 0)

	sb := common.NewSessionBuilder(common.AdminUser).
		AddCondition(cond)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&orgs)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

// GetByAccount 根据账号获取组织（邮箱/手机号）
func (r *organizationRepository) GetByAccount(account string) (*object.Organization, error) {
	orgs := make([]*object.Organization, 0)
	err := object.GetEngine().Find(&orgs, &object.Organization{Owner: common.AdminUser})
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		// 按优先级匹配组织
		if org.Name != "built-in" {
			return org, nil
		}
	}

	// 返回默认组织
	for _, org := range orgs {
		if org.Name == "built-in" {
			return org, nil
		}
	}

	return nil, common.ErrOrgNotFound
}

// GetByOwner 根据所有者获取组织
func (r *organizationRepository) GetByOwner(owner string) ([]*object.Organization, error) {
	orgs := make([]*object.Organization, 0)
	err := object.GetEngine().Find(&orgs, &object.Organization{Owner: owner})
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetParentOrganizations 获取父组织列表
func (r *organizationRepository) GetParentOrganizations(orgName string) ([]*object.Organization, error) {
	result := make([]*object.Organization, 0)
	visited := make(map[string]bool)

	currentOrg, err := r.GetByID(orgName)
	if err != nil {
		return nil, err
	}

	for currentOrg != nil && currentOrg.ParentId != "" {
		if visited[currentOrg.Name] {
			break
		}
		visited[currentOrg.Name] = true

		parentOrg, err := r.GetByID(currentOrg.ParentId)
		if err != nil {
			break
		}

		result = append([]*object.Organization{parentOrg}, result...)
		currentOrg = parentOrg
	}

	return result, nil
}

// GetChildOrganizations 获取直接子组织列表
func (r *organizationRepository) GetChildOrganizations(orgName string) ([]*object.Organization, error) {
	orgs := make([]*object.Organization, 0)
	err := object.GetEngine().
		Where("parent_id = ?", orgName).
		And("owner = ?", common.AdminUser).
		Find(&orgs)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetAllChildOrganizations 获取所有子组织列表（递归）
func (r *organizationRepository) GetAllChildOrganizations(orgName string) ([]*object.Organization, error) {
	result := make([]*object.Organization, 0)
	visited := make(map[string]bool)

	// BFS遍历
	queue := []string{orgName}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		children, err := r.GetChildOrganizations(current)
		if err != nil {
			return nil, err
		}

		for _, child := range children {
			if !visited[child.Name] {
				result = append(result, child)
				queue = append(queue, child.Name)
			}
		}
	}

	return result, nil
}

// GetUserCount 获取组织用户数量
func (r *organizationRepository) GetUserCount(orgName string) (int64, error) {
	return object.GetEngine().
		Where("owner = ?", orgName).
		Count(&object.User{})
}

// GetApplicationCount 获取组织应用数量
func (r *organizationRepository) GetApplicationCount(orgName string) (int64, error) {
	return object.GetEngine().
		Where("owner = ?", orgName).
		Count(&object.Application{})
}

// CheckDuplicate 检查组织唯一性
func (r *organizationRepository) CheckDuplicate(org *object.Organization) error {
	existedOrg := &object.Organization{Owner: common.AdminUser, Name: org.Name}
	existed, err := object.GetEngine().Get(existedOrg)
	if err != nil {
		return err
	}
	if existed {
		return common.NewBusinessError(common.ErrCodeOrgAlreadyExists,
			fmt.Sprintf("The organization %s already exists", org.Name))
	}

	return nil
}
