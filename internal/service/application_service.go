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
	"strings"

	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
)

// ApplicationService 应用业务接口
type ApplicationService interface {
	// 基础CRUD操作
	GetApplication(owner, name string) (*object.Application, error)
	GetApplicationByClientID(clientID string) (*object.Application, error)
	GetApplicationByClientSecret(clientSecret string) (*object.Application, error)
	ListApplications(owner string, page, pageSize int, field, value, sortField, sortOrder string) ([]*object.Application, int64, error)
	CreateApplication(app *object.Application) (bool, error)
	UpdateApplication(app *object.Application, columns ...string) (bool, error)
	DeleteApplication(app *object.Application) (bool, error)

	// 高级查询
	GetAllApplications(owner string) ([]*object.Application, error)
	SearchApplications(owner string, cond builder.Cond) ([]*object.Application, error)
	GetApplicationsByOrganization(orgName string) ([]*object.Application, error)

	// OAuth相关
	ValidateClientCredentials(clientID, clientSecret string) (*object.Application, error)
	ValidateRedirectURI(app *object.Application, redirectURI string) bool
	GetSupportedGrantTypes(app *object.Application) []string
	GetSupportedResponseTypes(app *object.Application) []string

	// 权限授予
	GetApplicationPermissions(appID string) ([]*object.Permission, error)
	GetApplicationRoles(appID string) ([]*object.Role, error)
	GetApplicationUserCount(appID string) (int64, error)

	// 客户端凭证管理
	RotateClientCredentials(owner, name string) (*object.Application, error)
	GenerateClientCredentials() (string, string)

	// 业务操作
	EnableApplication(owner, name string) (bool, error)
	DisableApplication(owner, name string) (bool, error)
	UpdateApplicationTheme(owner, name, theme string) (bool, error)
	RegenerateClientSecret(owner, name string) (string, error)
}

type applicationService struct {
	appRepo repository.ApplicationRepository
}

// NewApplicationService 创建应用Service实例
func NewApplicationService() ApplicationService {
	return &applicationService{
		appRepo: repository.NewApplicationRepository(),
	}
}

// GetApplication 获取应用
func (s *applicationService) GetApplication(owner, name string) (*object.Application, error) {
	return s.appRepo.GetByID(owner, name)
}

// GetApplicationByClientID 根据ClientID获取应用
func (s *applicationService) GetApplicationByClientID(clientID string) (*object.Application, error) {
	return s.appRepo.GetByClientID(clientID)
}

// GetApplicationByClientSecret 根据ClientSecret获取应用
func (s *applicationService) GetApplicationByClientSecret(clientSecret string) (*object.Application, error) {
	return s.appRepo.GetByClientSecret(clientSecret)
}

// ListApplications 获取应用列表
func (s *applicationService) ListApplications(owner string, page, pageSize int, field, value, sortField, sortOrder string) ([]*object.Application, int64, error) {
	if page < 1 {
		page = common.DefaultPage
	}
	if pageSize < 1 || pageSize > common.MaxPageSize {
		pageSize = common.DefaultPageSize
	}

	offset := (page - 1) * pageSize
	apps, err := s.appRepo.List(owner, offset, pageSize, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.appRepo.Count(owner, field, value)
	if err != nil {
		return nil, 0, err
	}

	return apps, total, nil
}

// CreateApplication 创建应用
func (s *applicationService) CreateApplication(app *object.Application) (bool, error) {
	// 生成客户端凭证
	if app.ClientId == "" {
		app.ClientId, app.ClientSecret = s.GenerateClientCredentials()
	}

	// 设置默认值
	if app.Status == "" {
		app.Status = common.StatusEnabled
	}
	if app.TokenTimeout == 0 {
		app.TokenTimeout = 86400 // 默认24小时
	}
	if app.RefreshTokenTimeout == 0 {
		app.RefreshTokenTimeout = 604800 // 默认7天
	}
	if app.FormOffset == "" {
		app.FormOffset = "+0"
	}
	if app.SamlResponseSigned == "" {
		app.SamlResponseSigned = "true"
	}
	if app.SamlAssertionSigned == "" {
		app.SamlAssertionSigned = "true"
	}
	if app.SamlNameidFormat == "" {
		app.SamlNameidFormat = "unspecified"
	}

	return s.appRepo.Create(app)
}

// UpdateApplication 更新应用
func (s *applicationService) UpdateApplication(app *object.Application, columns ...string) (bool, error) {
	return s.appRepo.Update(app, columns...)
}

// DeleteApplication 删除应用
func (s *applicationService) DeleteApplication(app *object.Application) (bool, error) {
	return s.appRepo.Delete(app)
}

// GetAllApplications 获取所有应用
func (s *applicationService) GetAllApplications(owner string) ([]*object.Application, error) {
	return s.appRepo.GetAll(owner)
}

// SearchApplications 搜索应用
func (s *applicationService) SearchApplications(owner string, cond builder.Cond) ([]*object.Application, error) {
	return s.appRepo.GetWithFilter(owner, cond)
}

// GetApplicationsByOrganization 获取组织下的所有应用
func (s *applicationService) GetApplicationsByOrganization(orgName string) ([]*object.Application, error) {
	return s.appRepo.GetByOrganization(orgName)
}

// ValidateClientCredentials 验证客户端凭证
func (s *applicationService) ValidateClientCredentials(clientID, clientSecret string) (*object.Application, error) {
	app, err := s.appRepo.GetByClientID(clientID)
	if err != nil {
		return nil, err
	}

	if app.ClientSecret != clientSecret {
		return nil, common.NewBusinessError(common.ErrCodeUnauthorized, "Invalid client secret")
	}

	if app.Status != common.StatusEnabled {
		return nil, common.NewBusinessError(common.ErrCodeForbidden, "Application is disabled")
	}

	return app, nil
}

// ValidateRedirectURI 验证重定向URI
func (s *applicationService) ValidateRedirectURI(app *object.Application, redirectURI string) bool {
	if app.RedirectUris == "" {
		return true
	}

	allowedURIs := strings.Split(app.RedirectUris, ",")
	for _, uri := range allowedURIs {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}
		// 支持通配符匹配
		if strings.HasSuffix(uri, "*") {
			prefix := strings.TrimSuffix(uri, "*")
			if strings.HasPrefix(redirectURI, prefix) {
				return true
			}
		} else if uri == redirectURI {
			return true
		}
	}

	return false
}

// GetSupportedGrantTypes 获取支持的授权类型
func (s *applicationService) GetSupportedGrantTypes(app *object.Application) []string {
	grantTypes := make([]string, 0)
	if app.GrantTypePassword {
		grantTypes = append(grantTypes, "password")
	}
	if app.GrantTypeCode {
		grantTypes = append(grantTypes, "authorization_code")
	}
	if app.GrantTypeClientCredentials {
		grantTypes = append(grantTypes, "client_credentials")
	}
	if app.GrantTypeRefreshToken {
		grantTypes = append(grantTypes, "refresh_token")
	}
	if app.GrantTypeOtp {
		grantTypes = append(grantTypes, "otp")
	}
	return grantTypes
}

// GetSupportedResponseTypes 获取支持的响应类型
func (s *applicationService) GetSupportedResponseTypes(app *object.Application) []string {
	responseTypes := make([]string, 0)
	if app.ResponseTypeCode {
		responseTypes = append(responseTypes, "code")
	}
	if app.ResponseTypeIdToken {
		responseTypes = append(responseTypes, "id_token")
	}
	if app.ResponseTypeToken {
		responseTypes = append(responseTypes, "token")
	}
	return responseTypes
}

// GetApplicationPermissions 获取应用权限列表
func (s *applicationService) GetApplicationPermissions(appID string) ([]*object.Permission, error) {
	return s.appRepo.GetPermissions(appID)
}

// GetApplicationRoles 获取应用角色列表
func (s *applicationService) GetApplicationRoles(appID string) ([]*object.Role, error) {
	return s.appRepo.GetRoles(appID)
}

// GetApplicationUserCount 获取应用用户数量
func (s *applicationService) GetApplicationUserCount(appID string) (int64, error) {
	return s.appRepo.GetUserCountByApp(appID)
}

// RotateClientCredentials 轮换客户端凭证
func (s *applicationService) RotateClientCredentials(owner, name string) (*object.Application, error) {
	app, err := s.appRepo.GetByID(owner, name)
	if err != nil {
		return nil, err
	}

	app.ClientId, app.ClientSecret = s.GenerateClientCredentials()
	_, err = s.appRepo.Update(app, "client_id", "client_secret")
	if err != nil {
		return nil, err
	}

	return app, nil
}

// GenerateClientCredentials 生成客户端凭证
func (s *applicationService) GenerateClientCredentials() (string, string) {
	clientID := util.GenerateClientId()
	clientSecret := util.GenerateClientSecret()
	return clientID, clientSecret
}

// EnableApplication 启用应用
func (s *applicationService) EnableApplication(owner, name string) (bool, error) {
	app, err := s.appRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	app.Status = common.StatusEnabled
	return s.appRepo.Update(app, "status")
}

// DisableApplication 禁用应用
func (s *applicationService) DisableApplication(owner, name string) (bool, error) {
	app, err := s.appRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	app.Status = common.StatusDisabled
	return s.appRepo.Update(app, "status")
}

// UpdateApplicationTheme 更新应用主题
func (s *applicationService) UpdateApplicationTheme(owner, name, theme string) (bool, error) {
	app, err := s.appRepo.GetByID(owner, name)
	if err != nil {
		return false, err
	}

	app.Theme = theme
	return s.appRepo.Update(app, "theme")
}

// RegenerateClientSecret 重新生成客户端密钥
func (s *applicationService) RegenerateClientSecret(owner, name string) (string, error) {
	app, err := s.appRepo.GetByID(owner, name)
	if err != nil {
		return "", err
	}

	newSecret := util.GenerateClientSecret()
	app.ClientSecret = newSecret
	_, err = s.appRepo.Update(app, "client_secret")
	if err != nil {
		return "", err
	}

	return newSecret, nil
}
