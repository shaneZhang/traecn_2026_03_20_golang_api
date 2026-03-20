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

	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/model"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/xorm"
)

// OrganizationRepository defines organization repository interface
type OrganizationRepository interface {
	// Basic CRUD
	GetByID(ctx context.Context, id string) (*model.Organization, error)
	GetByOwnerAndName(ctx context.Context, owner, name string) (*model.Organization, error)
	Create(ctx context.Context, org *model.Organization) error
	Update(ctx context.Context, org *model.Organization, columns []string) error
	Delete(ctx context.Context, id string) error

	// List operations
	List(ctx context.Context, filter OrganizationFilter) ([]*model.Organization, error)
	ListByOwner(ctx context.Context, owner string, offset, limit int) ([]*model.Organization, error)
	ListWithPagination(ctx context.Context, filter OrganizationFilter, offset, limit int) ([]*model.Organization, int64, error)
	Count(ctx context.Context, filter OrganizationFilter) (int64, error)

	// Batch operations
	BatchCreate(ctx context.Context, orgs []*model.Organization) error
	BatchUpdate(ctx context.Context, orgs []*model.Organization, columns []string) error
	BatchDelete(ctx context.Context, ids []string) error

	// Hierarchy operations
	GetChildren(ctx context.Context, parentID string) ([]*model.Organization, error)
	GetAncestors(ctx context.Context, orgID string) ([]*model.Organization, error)
	GetDescendants(ctx context.Context, orgID string, maxDepth int) ([]*model.Organization, error)

	// Statistics
	GetStatistics(ctx context.Context, orgID string) (*OrganizationStatistics, error)

	// Search
	Search(ctx context.Context, keyword string, fields []string) ([]*model.Organization, error)

	// Get by fields
	GetByFields(ctx context.Context, owner string, fields []string) ([]*model.Organization, error)
}

// OrganizationFilter represents organization filter criteria
type OrganizationFilter struct {
	Owner     string
	Name      string
	Field     string
	Value     string
	SortField string
	SortOrder string
}

// OrganizationStatistics represents organization statistics
type OrganizationStatistics struct {
	UserCount        int64
	ApplicationCount int64
	RoleCount        int64
	PermissionCount  int64
	TotalBalance     float64
}

// organizationRepository implements OrganizationRepository
type organizationRepository struct {
	db *xorm.Engine
}

// NewOrganizationRepository creates new organization repository
func NewOrganizationRepository(db *xorm.Engine) OrganizationRepository {
	return &organizationRepository{db: db}
}

// GetByID gets organization by ID
func (r *organizationRepository) GetByID(ctx context.Context, id string) (*model.Organization, error) {
	owner, name, err := parseOrgID(id)
	if err != nil {
		return nil, common.WrapError(err, "invalid organization ID")
	}
	return r.GetByOwnerAndName(ctx, owner, name)
}

// GetByOwnerAndName gets organization by owner and name
func (r *organizationRepository) GetByOwnerAndName(ctx context.Context, owner, name string) (*model.Organization, error) {
	org := &model.Organization{}
	exists, err := r.db.Context(ctx).Where("owner = ? AND name = ?", owner, name).Get(org)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrOrganizationNotFound
	}
	return org, nil
}

// Create creates a new organization
func (r *organizationRepository) Create(ctx context.Context, org *model.Organization) error {
	_, err := r.db.Context(ctx).Insert(org)
	if err != nil {
		return common.WrapError(err, "failed to create organization")
	}
	return nil
}

// Update updates organization
func (r *organizationRepository) Update(ctx context.Context, org *model.Organization, columns []string) error {
	session := r.db.Context(ctx).ID([]interface{}{org.Owner, org.Name})
	if len(columns) > 0 {
		session = session.Cols(columns...)
	}
	_, err := session.Update(org)
	if err != nil {
		return common.WrapError(err, "failed to update organization")
	}
	return nil
}

// Delete deletes organization
func (r *organizationRepository) Delete(ctx context.Context, id string) error {
	owner, name, err := parseOrgID(id)
	if err != nil {
		return err
	}
	_, err = r.db.Context(ctx).Where("owner = ? AND name = ?", owner, name).Delete(&model.Organization{})
	if err != nil {
		return common.WrapError(err, "failed to delete organization")
	}
	return nil
}

// List lists organizations with filter
func (r *organizationRepository) List(ctx context.Context, filter OrganizationFilter) ([]*model.Organization, error) {
	session := r.buildFilterSession(ctx, filter)

	var orgs []*model.Organization
	err := session.Find(&orgs)
	if err != nil {
		return nil, common.WrapError(err, "failed to list organizations")
	}
	return orgs, nil
}

// ListByOwner gets organizations by owner
func (r *organizationRepository) ListByOwner(ctx context.Context, owner string, offset, limit int) ([]*model.Organization, error) {
	var orgs []*model.Organization
	session := r.db.Context(ctx).Where("owner = ?", owner)
	if limit > 0 {
		session = session.Limit(limit, offset)
	}
	err := session.Find(&orgs)
	return orgs, err
}

// ListWithPagination lists organizations with pagination
func (r *organizationRepository) ListWithPagination(ctx context.Context, filter OrganizationFilter, offset, limit int) ([]*model.Organization, int64, error) {
	session := r.buildFilterSession(ctx, filter)

	total, err := session.Count(&model.Organization{})
	if err != nil {
		return nil, 0, common.WrapError(err, "failed to count organizations")
	}

	var orgs []*model.Organization
	err = session.Limit(limit, offset).Find(&orgs)
	if err != nil {
		return nil, 0, common.WrapError(err, "failed to list organizations")
	}

	return orgs, total, nil
}

// Count counts organizations
func (r *organizationRepository) Count(ctx context.Context, filter OrganizationFilter) (int64, error) {
	session := r.buildFilterSession(ctx, filter)
	return session.Count(&model.Organization{})
}

// BatchCreate creates multiple organizations
func (r *organizationRepository) BatchCreate(ctx context.Context, orgs []*model.Organization) error {
	session := r.db.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	for _, org := range orgs {
		_, err = session.Context(ctx).Insert(org)
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch create organizations")
		}
	}

	return session.Commit()
}

// BatchUpdate updates multiple organizations
func (r *organizationRepository) BatchUpdate(ctx context.Context, orgs []*model.Organization, columns []string) error {
	session := r.db.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	for _, org := range orgs {
		s := session.Context(ctx).ID([]interface{}{org.Owner, org.Name})
		if len(columns) > 0 {
			s = s.Cols(columns...)
		}
		_, err = s.Update(org)
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch update organizations")
		}
	}

	return session.Commit()
}

// BatchDelete deletes multiple organizations
func (r *organizationRepository) BatchDelete(ctx context.Context, ids []string) error {
	session := r.db.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	for _, id := range ids {
		owner, name, err := parseOrgID(id)
		if err != nil {
			session.Rollback()
			return err
		}
		_, err = session.Context(ctx).Where("owner = ? AND name = ?", owner, name).Delete(&model.Organization{})
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch delete organizations")
		}
	}

	return session.Commit()
}

// GetChildren gets child organizations
func (r *organizationRepository) GetChildren(ctx context.Context, parentID string) ([]*model.Organization, error) {
	// This is a placeholder - actual implementation depends on how hierarchy is stored
	// For now, return empty list
	return []*model.Organization{}, nil
}

// GetAncestors gets ancestor organizations
func (r *organizationRepository) GetAncestors(ctx context.Context, orgID string) ([]*model.Organization, error) {
	// This is a placeholder - actual implementation depends on how hierarchy is stored
	return []*model.Organization{}, nil
}

// GetDescendants gets descendant organizations
func (r *organizationRepository) GetDescendants(ctx context.Context, orgID string, maxDepth int) ([]*model.Organization, error) {
	// This is a placeholder - actual implementation depends on how hierarchy is stored
	return []*model.Organization{}, nil
}

// GetStatistics gets organization statistics
func (r *organizationRepository) GetStatistics(ctx context.Context, orgID string) (*OrganizationStatistics, error) {
	// This is a placeholder - actual implementation would query related tables
	return &OrganizationStatistics{}, nil
}

// Search searches organizations
func (r *organizationRepository) Search(ctx context.Context, keyword string, fields []string) ([]*model.Organization, error) {
	session := r.db.Context(ctx)

	if len(fields) > 0 && keyword != "" {
		cond := builder.NewCond()
		for _, field := range fields {
			cond = cond.Or(builder.Like{field, "%" + keyword + "%"})
		}
		session = session.Where(cond)
	}

	var orgs []*model.Organization
	err := session.Find(&orgs)
	if err != nil {
		return nil, common.WrapError(err, "failed to search organizations")
	}
	return orgs, nil
}

// GetByFields gets organizations by specific fields
func (r *organizationRepository) GetByFields(ctx context.Context, owner string, fields []string) ([]*model.Organization, error) {
	var orgs []*model.Organization
	session := r.db.Context(ctx).Where("owner = ?", owner)
	if len(fields) > 0 {
		session = session.Cols(fields...)
	}
	err := session.Find(&orgs)
	if err != nil {
		return nil, common.WrapError(err, "failed to get organizations by fields")
	}
	return orgs, nil
}

// buildFilterSession builds filter session
func (r *organizationRepository) buildFilterSession(ctx context.Context, filter OrganizationFilter) *xorm.Session {
	session := r.db.Context(ctx)

	if filter.Owner != "" {
		session = session.Where("owner = ?", filter.Owner)
	}

	if filter.Name != "" {
		session = session.Where("name = ?", filter.Name)
	}

	if filter.Field != "" && filter.Value != "" {
		session = session.Where(builder.Like{filter.Field, "%" + filter.Value + "%"})
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

// parseOrgID parses organization ID
func parseOrgID(id string) (owner, name string, err error) {
	parts := splitOrgID(id)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid organization ID format: %s", id)
	}
	return parts[0], parts[1], nil
}

// splitOrgID splits ID into owner and name
func splitOrgID(id string) []string {
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
