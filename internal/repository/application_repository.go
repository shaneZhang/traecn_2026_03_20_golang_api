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

package repository

import (
	"context"
	"fmt"

	"github.com/casdoor/casdoor/internal/model"
	"github.com/casdoor/casdoor/internal/common"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/xorm"
)

// ApplicationRepository defines application repository interface
type ApplicationRepository interface {
	// Basic CRUD
	GetByID(ctx context.Context, id string) (*model.Application, error)
	GetByOwnerAndName(ctx context.Context, owner, name string) (*model.Application, error)
	GetByClientID(ctx context.Context, clientID string) (*model.Application, error)
	Create(ctx context.Context, app *model.Application) error
	Update(ctx context.Context, app *model.Application, columns []string) error
	Delete(ctx context.Context, id string) error
	
	// List operations
	List(ctx context.Context, filter ApplicationFilter) ([]*model.Application, error)
	ListWithPagination(ctx context.Context, filter ApplicationFilter, offset, limit int) ([]*model.Application, int64, error)
	Count(ctx context.Context, filter ApplicationFilter) (int64, error)
	
	// Organization operations
	GetByOrganization(ctx context.Context, owner, organization string) ([]*model.Application, error)
	GetByOrganizationWithPagination(ctx context.Context, owner, organization string, offset, limit int) ([]*model.Application, int64, error)
	CountByOrganization(ctx context.Context, owner, organization string) (int64, error)
	
	// Batch operations
	BatchCreate(ctx context.Context, apps []*model.Application) error
	BatchUpdate(ctx context.Context, apps []*model.Application, columns []string) error
	BatchDelete(ctx context.Context, ids []string) error
	
	// OAuth operations
	ValidateClientCredentials(ctx context.Context, clientID, clientSecret string) (*model.Application, error)
	ValidateRedirectURI(ctx context.Context, appID, redirectURI string) error
	GetByDomain(ctx context.Context, domain string) (*model.Application, error)
	
	// Search
	Search(ctx context.Context, owner, keyword string, fields []string) ([]*model.Application, error)
	
	// Statistics
	GetStatistics(ctx context.Context, owner string) (*ApplicationStatistics, error)
}

// ApplicationFilter represents application filter criteria
type ApplicationFilter struct {
	Owner       string
	Field       string
	Value       string
	SortField   string
	SortOrder   string
	IsShared    *bool
}

// ApplicationStatistics represents application statistics
type ApplicationStatistics struct {
	Total          int64
	ByCategory     map[string]int64
	ByType         map[string]int64
	OAuthApps      int64
	SAMLApps       int64
}

// applicationRepository implements ApplicationRepository
type applicationRepository struct {
	db *xorm.Engine
}

// NewApplicationRepository creates new application repository
func NewApplicationRepository(db *xorm.Engine) ApplicationRepository {
	return &applicationRepository{db: db}
}

// GetByID gets application by ID
func (r *applicationRepository) GetByID(ctx context.Context, id string) (*model.Application, error) {
	owner, name, err := parseAppID(id)
	if err != nil {
		return nil, common.WrapError(err, "invalid application ID")
	}
	return r.GetByOwnerAndName(ctx, owner, name)
}

// GetByOwnerAndName gets application by owner and name
func (r *applicationRepository) GetByOwnerAndName(ctx context.Context, owner, name string) (*model.Application, error) {
	app := &model.Application{}
	exists, err := r.db.Context(ctx).Where("owner = ? AND name = ?", owner, name).Get(app)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrApplicationNotFound
	}
	return app, nil
}

// GetByClientID gets application by client ID
func (r *applicationRepository) GetByClientID(ctx context.Context, clientID string) (*model.Application, error) {
	app := &model.Application{}
	exists, err := r.db.Context(ctx).Where("client_id = ?", clientID).Get(app)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrApplicationNotFound
	}
	return app, nil
}

// Create creates a new application
func (r *applicationRepository) Create(ctx context.Context, app *model.Application) error {
	_, err := r.db.Context(ctx).Insert(app)
	if err != nil {
		return common.WrapError(err, "failed to create application")
	}
	return nil
}

// Update updates application
func (r *applicationRepository) Update(ctx context.Context, app *model.Application, columns []string) error {
	session := r.db.Context(ctx).ID([]interface{}{app.Owner, app.Name})
	if len(columns) > 0 {
		session = session.Cols(columns...)
	}
	_, err := session.Update(app)
	if err != nil {
		return common.WrapError(err, "failed to update application")
	}
	return nil
}

// Delete deletes application
func (r *applicationRepository) Delete(ctx context.Context, id string) error {
	owner, name, err := parseAppID(id)
	if err != nil {
		return err
	}
	_, err = r.db.Context(ctx).Where("owner = ? AND name = ?", owner, name).Delete(&model.Application{})
	if err != nil {
		return common.WrapError(err, "failed to delete application")
	}
	return nil
}

// List lists applications with filter
func (r *applicationRepository) List(ctx context.Context, filter ApplicationFilter) ([]*model.Application, error) {
	session := r.buildFilterSession(ctx, filter)
	
	var apps []*model.Application
	err := session.Find(&apps)
	if err != nil {
		return nil, common.WrapError(err, "failed to list applications")
	}
	return apps, nil
}

// ListWithPagination lists applications with pagination
func (r *applicationRepository) ListWithPagination(ctx context.Context, filter ApplicationFilter, offset, limit int) ([]*model.Application, int64, error) {
	session := r.buildFilterSession(ctx, filter)
	
	total, err := session.Count(&model.Application{})
	if err != nil {
		return nil, 0, common.WrapError(err, "failed to count applications")
	}
	
	var apps []*model.Application
	err = session.Limit(limit, offset).Find(&apps)
	if err != nil {
		return nil, 0, common.WrapError(err, "failed to list applications")
	}
	
	return apps, total, nil
}

// Count counts applications
func (r *applicationRepository) Count(ctx context.Context, filter ApplicationFilter) (int64, error) {
	session := r.buildFilterSession(ctx, filter)
	return session.Count(&model.Application{})
}

// GetByOrganization gets applications by organization
func (r *applicationRepository) GetByOrganization(ctx context.Context, owner, organization string) ([]*model.Application, error) {
	var apps []*model.Application
	err := r.db.Context(ctx).
		Where("owner = ? AND (organization = ? OR is_shared = ?)", owner, organization, true).
		Find(&apps)
	if err != nil {
		return nil, common.WrapError(err, "failed to get applications by organization")
	}
	return apps, nil
}

// GetByOrganizationWithPagination gets applications by organization with pagination
func (r *applicationRepository) GetByOrganizationWithPagination(ctx context.Context, owner, organization string, offset, limit int) ([]*model.Application, int64, error) {
	session := r.db.Context(ctx).
		Where("owner = ? AND (organization = ? OR is_shared = ?)", owner, organization, true)
	
	total, err := session.Count(&model.Application{})
	if err != nil {
		return nil, 0, err
	}
	
	var apps []*model.Application
	err = session.Limit(limit, offset).Find(&apps)
	if err != nil {
		return nil, 0, err
	}
	
	return apps, total, nil
}

// CountByOrganization counts applications by organization
func (r *applicationRepository) CountByOrganization(ctx context.Context, owner, organization string) (int64, error) {
	return r.db.Context(ctx).
		Where("owner = ? AND (organization = ? OR is_shared = ?)", owner, organization, true).
		Count(&model.Application{})
}

// BatchCreate creates multiple applications
func (r *applicationRepository) BatchCreate(ctx context.Context, apps []*model.Application) error {
	session := r.db.NewSession()
	defer session.Close()
	
	err := session.Begin()
	if err != nil {
		return err
	}
	
	for _, app := range apps {
		_, err = session.Context(ctx).Insert(app)
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch create applications")
		}
	}
	
	return session.Commit()
}

// BatchUpdate updates multiple applications
func (r *applicationRepository) BatchUpdate(ctx context.Context, apps []*model.Application, columns []string) error {
	session := r.db.NewSession()
	defer session.Close()
	
	err := session.Begin()
	if err != nil {
		return err
	}
	
	for _, app := range apps {
		s := session.Context(ctx).ID([]interface{}{app.Owner, app.Name})
		if len(columns) > 0 {
			s = s.Cols(columns...)
		}
		_, err = s.Update(app)
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch update applications")
		}
	}
	
	return session.Commit()
}

// BatchDelete deletes multiple applications
func (r *applicationRepository) BatchDelete(ctx context.Context, ids []string) error {
	session := r.db.NewSession()
	defer session.Close()
	
	err := session.Begin()
	if err != nil {
		return err
	}
	
	for _, id := range ids {
		owner, name, err := parseAppID(id)
		if err != nil {
			session.Rollback()
			return err
		}
		_, err = session.Context(ctx).Where("owner = ? AND name = ?", owner, name).Delete(&model.Application{})
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch delete applications")
		}
	}
	
	return session.Commit()
}

// ValidateClientCredentials validates OAuth client credentials
func (r *applicationRepository) ValidateClientCredentials(ctx context.Context, clientID, clientSecret string) (*model.Application, error) {
	app := &model.Application{}
	exists, err := r.db.Context(ctx).
		Where("client_id = ? AND client_secret = ?", clientID, clientSecret).
		Get(app)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrInvalidClientCredentials
	}
	return app, nil
}

// ValidateRedirectURI validates redirect URI for application
func (r *applicationRepository) ValidateRedirectURI(ctx context.Context, appID, redirectURI string) error {
	app, err := r.GetByID(ctx, appID)
	if err != nil {
		return err
	}
	
	for _, uri := range app.RedirectUris {
		if uri == redirectURI {
			return nil
		}
	}
	
	return common.ErrInvalidRedirectURI
}

// GetByDomain gets application by domain
func (r *applicationRepository) GetByDomain(ctx context.Context, domain string) (*model.Application, error) {
	app := &model.Application{}
	exists, err := r.db.Context(ctx).
		Where("domain = ?", domain).
		Get(app)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrApplicationNotFound
	}
	return app, nil
}

// Search searches applications
func (r *applicationRepository) Search(ctx context.Context, owner, keyword string, fields []string) ([]*model.Application, error) {
	session := r.db.Context(ctx).Where("owner = ?", owner)
	
	if len(fields) > 0 && keyword != "" {
		cond := builder.NewCond()
		for _, field := range fields {
			cond = cond.Or(builder.Like{field, "%" + keyword + "%"})
		}
		session = session.Where(cond)
	}
	
	var apps []*model.Application
	err := session.Find(&apps)
	if err != nil {
		return nil, common.WrapError(err, "failed to search applications")
	}
	return apps, nil
}

// GetStatistics gets application statistics
func (r *applicationRepository) GetStatistics(ctx context.Context, owner string) (*ApplicationStatistics, error) {
	stats := &ApplicationStatistics{
		ByCategory: make(map[string]int64),
		ByType:     make(map[string]int64),
	}
	
	session := r.db.Context(ctx).Where("owner = ?", owner)
	
	total, err := session.Count(&model.Application{})
	if err != nil {
		return nil, err
	}
	stats.Total = total
	
	// Count OAuth apps (those with client_id)
	oauth, err := session.Where("client_id != ?", "").Count(&model.Application{})
	if err != nil {
		return nil, err
	}
	stats.OAuthApps = oauth
	
	return stats, nil
}

// buildFilterSession builds filter session
func (r *applicationRepository) buildFilterSession(ctx context.Context, filter ApplicationFilter) *xorm.Session {
	session := r.db.Context(ctx)
	
	if filter.Owner != "" {
		session = session.Where("owner = ?", filter.Owner)
	}
	
	if filter.Field != "" && filter.Value != "" {
		session = session.Where(builder.Like{filter.Field, "%" + filter.Value + "%"})
	}
	
	if filter.IsShared != nil {
		session = session.Where("is_shared = ?", *filter.IsShared)
	}
	
	// Sorting
	sortField := filter.SortField
	if sortField == "" {
		sortField = "created_time"
	}
	
	if filter.SortOrder == "desc" {
		session = session.Desc(sortField)
	} else {
		session = session.Asc(sortField)
	}
	
	return session
}

// parseAppID parses application ID
func parseAppID(id string) (owner, name string, err error) {
	parts := splitAppID(id)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid application ID format: %s", id)
	}
	return parts[0], parts[1], nil
}

// splitAppID splits ID into owner and name
func splitAppID(id string) []string {
	parts := make([]string, 0)
	current := ""
	for _, ch := range id {
		if ch == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	parts = append(parts, current)
	return parts
}
