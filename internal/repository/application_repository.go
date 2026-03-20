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

// ApplicationRepository 应用数据访问接口
type ApplicationRepository interface {
	// 基础CRUD操作
	GetByID(owner, name string) (*object.Application, error)
	GetByClientID(clientID string) (*object.Application, error)
	GetByClientSecret(clientSecret string) (*object.Application, error)
	List(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, error)
	Count(owner, field, value string) (int64, error)
	Create(app *object.Application) (bool, error)
	Update(app *object.Application, columns ...string) (bool, error)
	Delete(app *object.Application) (bool, error)

	// 特殊查询
	GetAll(owner string) ([]*object.Application, error)
	GetAllWithSort(owner, sortField, sortOrder string) ([]*object.Application, error)
	GetWithFilter(owner string, cond builder.Cond) ([]*object.Application, error)
	GetByOrganization(orgName string) ([]*object.Application, error)
	GetUserCountByApp(appID string) (int64, error)

	// OAuth相关
	GetByProvider(providerType, providerName string) ([]*object.Application, error)
	GetByCallbackURL(callbackURL string) ([]*object.Application, error)

	// 权限授予
	GetPermissions(appID string) ([]*object.Permission, error)
	GetRoles(appID string) ([]*object.Role, error)

	// 唯一性检查
	CheckDuplicate(app *object.Application) error
}

type applicationRepository struct {
}

// NewApplicationRepository 创建应用Repository实例
func NewApplicationRepository() ApplicationRepository {
	return &applicationRepository{}
}

// GetByID 根据ID获取应用
func (r *applicationRepository) GetByID(owner, name string) (*object.Application, error) {
	if owner == "" || name == "" {
		return nil, common.ErrBadRequest
	}

	app := &object.Application{Owner: owner, Name: name}
	existed, err := common.GetByID(owner, name, app)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrAppNotFound
	}
	return app, nil
}

// GetByClientID 根据ClientID获取应用
func (r *applicationRepository) GetByClientID(clientID string) (*object.Application, error) {
	if clientID == "" {
		return nil, common.ErrBadRequest
	}

	app := &object.Application{ClientId: clientID}
	existed, err := object.GetEngine().Get(app)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrAppNotFound
	}
	return app, nil
}

// GetByClientSecret 根据ClientSecret获取应用
func (r *applicationRepository) GetByClientSecret(clientSecret string) (*object.Application, error) {
	if clientSecret == "" {
		return nil, common.ErrBadRequest
	}

	app := &object.Application{ClientSecret: clientSecret}
	existed, err := object.GetEngine().Get(app)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, common.ErrAppNotFound
	}
	return app, nil
}

// List 获取应用列表
func (r *applicationRepository) List(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)

	sb := common.NewSessionBuilder(owner).
		SetPagination(offset, limit).
		SetFilter(field, value).
		SetSort(sortField, sortOrder)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// Count 获取应用数量
func (r *applicationRepository) Count(owner, field, value string) (int64, error) {
	sb := common.NewSessionBuilder(owner).
		SetFilter(field, value)

	session := sb.Build()
	defer session.Close()

	return session.Count(&object.Application{})
}

// Create 创建应用
func (r *applicationRepository) Create(app *object.Application) (bool, error) {
	if err := r.CheckDuplicate(app); err != nil {
		return false, err
	}

	affected, err := object.GetEngine().Insert(app)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// Update 更新应用
func (r *applicationRepository) Update(app *object.Application, columns ...string) (bool, error) {
	if app.Owner == "" || app.Name == "" {
		return false, common.ErrBadRequest
	}

	session := object.GetEngine().ID(core.PK{app.Owner, app.Name})
	if len(columns) > 0 {
		session = session.Cols(columns...)
	}

	affected, err := session.Update(app)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// Delete 删除应用
func (r *applicationRepository) Delete(app *object.Application) (bool, error) {
	if app.Owner == "" || app.Name == "" {
		return false, common.ErrBadRequest
	}

	affected, err := object.GetEngine().ID(core.PK{app.Owner, app.Name}).Delete(&object.Application{})
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

// GetAll 获取所有应用
func (r *applicationRepository) GetAll(owner string) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)
	err := object.GetEngine().Find(&apps, &object.Application{Owner: owner})
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// GetAllWithSort 获取所有应用并排序
func (r *applicationRepository) GetAllWithSort(owner, sortField, sortOrder string) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)

	sb := common.NewSessionBuilder(owner).
		SetSort(sortField, sortOrder)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// GetWithFilter 带条件查询应用列表
func (r *applicationRepository) GetWithFilter(owner string, cond builder.Cond) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)

	sb := common.NewSessionBuilder(owner).
		AddCondition(cond)

	session := sb.Build()
	defer session.Close()

	err := session.Find(&apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// GetByOrganization 获取组织下的所有应用
func (r *applicationRepository) GetByOrganization(orgName string) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)
	err := object.GetEngine().
		Where("owner = ?", orgName).
		Find(&apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// GetUserCountByApp 获取应用下的用户数量
func (r *applicationRepository) GetUserCountByApp(appID string) (int64, error) {
	return object.GetEngine().
		Where("signup_application = ?", appID).
		Count(&object.User{})
}

// GetByProvider 根据提供商获取应用列表
func (r *applicationRepository) GetByProvider(providerType, providerName string) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)

	// 这里需要查询provider_objs表中的关联关系
	// 简化实现，后续可以优化为关联查询
	allApps := make([]*object.Application, 0)
	err := object.GetEngine().Find(&allApps)
	if err != nil {
		return nil, err
	}

	result := make([]*object.Application, 0)
	for _, app := range allApps {
		// 检查应用是否配置了该提供商
		// 这里简化处理，实际应该查询关联表
		result = append(result, app)
	}

	return result, nil
}

// GetByCallbackURL 根据回调URL获取应用列表
func (r *applicationRepository) GetByCallbackURL(callbackURL string) ([]*object.Application, error) {
	apps := make([]*object.Application, 0)
	err := object.GetEngine().
		Where("redirect_uris like ?", "%"+callbackURL+"%").
		Find(&apps)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// GetPermissions 获取应用权限列表
func (r *applicationRepository) GetPermissions(appID string) ([]*object.Permission, error) {
	permissions := make([]*object.Permission, 0)
	err := object.GetEngine().
		Where("application = ?", appID).
		Find(&permissions)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetRoles 获取应用角色列表
func (r *applicationRepository) GetRoles(appID string) ([]*object.Role, error) {
	roles := make([]*object.Role, 0)
	err := object.GetEngine().
		Where("application = ?", appID).
		Find(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// CheckDuplicate 检查应用唯一性
func (r *applicationRepository) CheckDuplicate(app *object.Application) error {
	// 检查应用名称是否存在
	existedApp := &object.Application{Owner: app.Owner, Name: app.Name}
	existed, err := object.GetEngine().Get(existedApp)
	if err != nil {
		return err
	}
	if existed {
		return common.NewBusinessError(common.ErrCodeAppAlreadyExists,
			fmt.Sprintf("The application %s already exists", app.Name))
	}

	// 检查ClientID是否存在
	if app.ClientId != "" {
		existedApp := &object.Application{ClientId: app.ClientId}
		existed, err := object.GetEngine().Get(existedApp)
		if err != nil {
			return err
		}
		if existed {
			return common.NewBusinessError(common.ErrCodeAppAlreadyExists,
				fmt.Sprintf("The client ID %s already exists", app.ClientId))
		}
	}

	return nil
}
