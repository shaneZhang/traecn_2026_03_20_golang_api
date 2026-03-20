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
	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/object"
	"github.com/xorm-io/builder"
)

// OrganizationService 组织业务接口
type OrganizationService interface {
	// 基础CRUD操作
	GetOrganization(name string) (*object.Organization, error)
	GetOrganizationByUserID(userID string) (*object.Organization, error)
	ListOrganizations(page, pageSize int, field, value, sortField, sortOrder string) ([]*object.Organization, int64, error)
	CreateOrganization(org *object.Organization) (bool, error)
	UpdateOrganization(org *object.Organization, columns ...string) (bool, error)
	DeleteOrganization(org *object.Organization) (bool, error)

	// 高级查询
	GetAllOrganizations() ([]*object.Organization, error)
	SearchOrganizations(cond builder.Cond) ([]*object.Organization, error)
	GetOrganizationByAccount(account string) (*object.Organization, error)
	GetOrganizationsByOwner(owner string) ([]*object.Organization, error)

	// 层级关系
	GetParentOrganizations(orgName string) ([]*object.Organization, error)
	GetChildOrganizations(orgName string) ([]*object.Organization, error)
	GetAllChildOrganizations(orgName string) ([]*object.Organization, error)
	GetOrganizationHierarchy(orgName string) (*OrganizationNode, error)

	// 统计信息
	GetOrganizationStats(orgName string) (*OrganizationStats, error)

	// 业务操作
	EnableOrganization(name string) (bool, error)
	DisableOrganization(name string) (bool, error)
	UpdateOrganizationTheme(name, theme string) (bool, error)
}

// OrganizationNode 组织层级节点
type OrganizationNode struct {
	Organization *object.Organization `json:"organization"`
	Children     []*OrganizationNode  `json:"children"`
}

// OrganizationStats 组织统计信息
type OrganizationStats struct {
	UserCount        int64 `json:"userCount"`
	ApplicationCount int64 `json:"applicationCount"`
	RoleCount        int64 `json:"roleCount"`
	PermissionCount  int64 `json:"permissionCount"`
}

type organizationService struct {
	orgRepo  repository.OrganizationRepository
	userRepo repository.UserRepository
	appRepo  repository.ApplicationRepository
}

// NewOrganizationService 创建组织Service实例
func NewOrganizationService() OrganizationService {
	return &organizationService{
		orgRepo:  repository.NewOrganizationRepository(),
		userRepo: repository.NewUserRepository(),
		appRepo:  repository.NewApplicationRepository(),
	}
}

// GetOrganization 获取组织
func (s *organizationService) GetOrganization(name string) (*object.Organization, error) {
	return s.orgRepo.GetByID(name)
}

// GetOrganizationByUserID 根据用户ID获取组织
func (s *organizationService) GetOrganizationByUserID(userID string) (*object.Organization, error) {
	return s.orgRepo.GetByUserID(userID)
}

// ListOrganizations 获取组织列表
func (s *organizationService) ListOrganizations(page, pageSize int, field, value, sortField, sortOrder string) ([]*object.Organization, int64, error) {
	if page < 1 {
		page = common.DefaultPage
	}
	if pageSize < 1 || pageSize > common.MaxPageSize {
		pageSize = common.DefaultPageSize
	}

	offset := (page - 1) * pageSize
	orgs, err := s.orgRepo.List(offset, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.orgRepo.Count(field, value)
	if err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// CreateOrganization 创建组织
func (s *organizationService) CreateOrganization(org *object.Organization) (bool, error) {
	// 设置默认值
	if org.Owner == "" {
		org.Owner = common.AdminUser
	}
	if org.Status == "" {
		org.Status = common.StatusEnabled
	}

	return s.orgRepo.Create(org)
}

// UpdateOrganization 更新组织
func (s *organizationService) UpdateOrganization(org *object.Organization, columns ...string) (bool, error) {
	return s.orgRepo.Update(org, columns...)
}

// DeleteOrganization 删除组织
func (s *organizationService) DeleteOrganization(org *object.Organization) (bool, error) {
	return s.orgRepo.Delete(org)
}

// GetAllOrganizations 获取所有组织
func (s *organizationService) GetAllOrganizations() ([]*object.Organization, error) {
	return s.orgRepo.GetAll()
}

// SearchOrganizations 搜索组织
func (s *organizationService) SearchOrganizations(cond builder.Cond) ([]*object.Organization, error) {
	return s.orgRepo.GetWithFilter(cond)
}

// GetOrganizationByAccount 根据账号获取组织
func (s *organizationService) GetOrganizationByAccount(account string) (*object.Organization, error) {
	return s.orgRepo.GetByAccount(account)
}

// GetOrganizationsByOwner 根据所有者获取组织
func (s *organizationService) GetOrganizationsByOwner(owner string) ([]*object.Organization, error) {
	return s.orgRepo.GetByOwner(owner)
}

// GetParentOrganizations 获取父组织列表
func (s *organizationService) GetParentOrganizations(orgName string) ([]*object.Organization, error) {
	return s.orgRepo.GetParentOrganizations(orgName)
}

// GetChildOrganizations 获取直接子组织列表
func (s *organizationService) GetChildOrganizations(orgName string) ([]*object.Organization, error) {
	return s.orgRepo.GetChildOrganizations(orgName)
}

// GetAllChildOrganizations 获取所有子组织列表
func (s *organizationService) GetAllChildOrganizations(orgName string) ([]*object.Organization, error) {
	return s.orgRepo.GetAllChildOrganizations(orgName)
}

// GetOrganizationHierarchy 获取组织层级结构
func (s *organizationService) GetOrganizationHierarchy(orgName string) (*OrganizationNode, error) {
	org, err := s.orgRepo.GetByID(orgName)
	if err != nil {
		return nil, err
	}

	root := &OrganizationNode{
		Organization: org,
		Children:     make([]*OrganizationNode, 0),
	}

	// 递归构建层级结构
	err = s.buildOrganizationHierarchy(root)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// buildOrganizationHierarchy 递归构建组织层级结构
func (s *organizationService) buildOrganizationHierarchy(node *OrganizationNode) error {
	children, err := s.orgRepo.GetChildOrganizations(node.Organization.Name)
	if err != nil {
		return err
	}

	for _, child := range children {
		childNode := &OrganizationNode{
			Organization: child,
			Children:     make([]*OrganizationNode, 0),
		}
		node.Children = append(node.Children, childNode)

		err = s.buildOrganizationHierarchy(childNode)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetOrganizationStats 获取组织统计信息
func (s *organizationService) GetOrganizationStats(orgName string) (*OrganizationStats, error) {
	userCount, err := s.orgRepo.GetUserCount(orgName)
	if err != nil {
		return nil, err
	}

	appCount, err := s.orgRepo.GetApplicationCount(orgName)
	if err != nil {
		return nil, err
	}

	// 这里可以继续获取角色和权限统计
	stats := &OrganizationStats{
		UserCount:        userCount,
		ApplicationCount: appCount,
		RoleCount:        0,
		PermissionCount:  0,
	}

	return stats, nil
}

// EnableOrganization 启用组织
func (s *organizationService) EnableOrganization(name string) (bool, error) {
	org, err := s.orgRepo.GetByID(name)
	if err != nil {
		return false, err
	}

	org.Status = common.StatusEnabled
	return s.orgRepo.Update(org, "status")
}

// DisableOrganization 禁用组织
func (s *organizationService) DisableOrganization(name string) (bool, error) {
	org, err := s.orgRepo.GetByID(name)
	if err != nil {
		return false, err
	}

	org.Status = common.StatusDisabled
	return s.orgRepo.Update(org, "status")
}

// UpdateOrganizationTheme 更新组织主题
func (s *organizationService) UpdateOrganizationTheme(name, theme string) (bool, error) {
	org, err := s.orgRepo.GetByID(name)
	if err != nil {
		return false, err
	}

	org.Theme = theme
	return s.orgRepo.Update(org, "theme")
}
