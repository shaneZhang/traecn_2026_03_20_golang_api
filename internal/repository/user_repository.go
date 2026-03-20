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

// UserRepository defines user repository interface
type UserRepository interface {
	// Basic CRUD
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByOwnerAndName(ctx context.Context, owner, name string) (*model.User, error)
	GetByEmail(ctx context.Context, owner, email string) (*model.User, error)
	GetByPhone(ctx context.Context, owner, phone string) (*model.User, error)
	GetByUserID(ctx context.Context, owner, userID string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User, columns []string) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error

	// List operations
	List(ctx context.Context, filter UserFilter) ([]*model.User, error)
	ListWithPagination(ctx context.Context, filter UserFilter, offset, limit int) ([]*model.User, int64, error)
	Count(ctx context.Context, filter UserFilter) (int64, error)

	// Batch operations
	BatchCreate(ctx context.Context, users []*model.User) error
	BatchUpdate(ctx context.Context, users []*model.User, columns []string) error
	BatchDelete(ctx context.Context, ids []string) error

	// Global operations
	GetGlobalUsers(ctx context.Context, filter UserFilter) ([]*model.User, error)
	GetGlobalUsersWithPagination(ctx context.Context, filter UserFilter, offset, limit int) ([]*model.User, int64, error)
	CountGlobal(ctx context.Context, field, value string) (int64, error)

	// Group operations
	GetUsersByGroup(ctx context.Context, groupID string) ([]*model.User, error)
	AddUserToGroup(ctx context.Context, userID, groupID string) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID string) error

	// MFA operations
	UpdateMFA(ctx context.Context, userID string, mfaConfig *model.MfaConfig) error
	GetMFAConfig(ctx context.Context, userID string) (*model.MfaConfig, error)

	// Search operations
	Search(ctx context.Context, owner, keyword string, fields []string) ([]*model.User, error)

	// Statistics
	GetStatistics(ctx context.Context, owner string) (*UserStatistics, error)
}

// UserFilter represents user filter criteria
type UserFilter struct {
	Owner       string
	GroupName   string
	Field       string
	Value       string
	SortField   string
	SortOrder   string
	IsAdmin     *bool
	IsForbidden *bool
	IsDeleted   *bool
}

// UserStatistics represents user statistics
type UserStatistics struct {
	Total    int64
	Active   int64
	Inactive int64
	Admin    int64
	Online   int64
	ByType   map[string]int64
	ByRegion map[string]int64
}

// MfaConfig represents MFA configuration
type MfaConfig struct {
	PreferredMfaType  string
	RecoveryCodes     []string
	TotpSecret        string
	MfaPhoneEnabled   bool
	MfaEmailEnabled   bool
	MfaRadiusEnabled  bool
	MfaRadiusUsername string
	MfaRadiusProvider string
	MfaPushEnabled    bool
	MfaPushReceiver   string
	MfaPushProvider   string
}

// userRepository implements UserRepository
type userRepository struct {
	db *xorm.Engine
}

// NewUserRepository creates new user repository
func NewUserRepository(db *xorm.Engine) UserRepository {
	return &userRepository{db: db}
}

// GetByID gets user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	owner, name, err := parseID(id)
	if err != nil {
		return nil, common.WrapError(err, "invalid user ID")
	}
	return r.GetByOwnerAndName(ctx, owner, name)
}

// GetByOwnerAndName gets user by owner and name
func (r *userRepository) GetByOwnerAndName(ctx context.Context, owner, name string) (*model.User, error) {
	user := &model.User{}
	exists, err := r.db.Context(ctx).Where("owner = ? AND name = ?", owner, name).Get(user)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// GetByEmail gets user by email
func (r *userRepository) GetByEmail(ctx context.Context, owner, email string) (*model.User, error) {
	user := &model.User{}
	exists, err := r.db.Context(ctx).Where("owner = ? AND email = ?", owner, email).Get(user)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// GetByPhone gets user by phone
func (r *userRepository) GetByPhone(ctx context.Context, owner, phone string) (*model.User, error) {
	user := &model.User{}
	exists, err := r.db.Context(ctx).Where("owner = ? AND phone = ?", owner, phone).Get(user)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// GetByUserID gets user by user ID
func (r *userRepository) GetByUserID(ctx context.Context, owner, userID string) (*model.User, error) {
	user := &model.User{}
	exists, err := r.db.Context(ctx).Where("owner = ? AND id = ?", owner, userID).Get(user)
	if err != nil {
		return nil, common.WrapError(err, "database error")
	}
	if !exists {
		return nil, common.ErrUserNotFound
	}
	return user, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Context(ctx).Insert(user)
	if err != nil {
		return common.WrapError(err, "failed to create user")
	}
	return nil
}

// Update updates user
func (r *userRepository) Update(ctx context.Context, user *model.User, columns []string) error {
	session := r.db.Context(ctx).ID([]interface{}{user.Owner, user.Name})
	if len(columns) > 0 {
		session = session.Cols(columns...)
	}
	_, err := session.Update(user)
	if err != nil {
		return common.WrapError(err, "failed to update user")
	}
	return nil
}

// Delete deletes user permanently
func (r *userRepository) Delete(ctx context.Context, id string) error {
	owner, name, err := parseID(id)
	if err != nil {
		return err
	}
	_, err = r.db.Context(ctx).Where("owner = ? AND name = ?", owner, name).Delete(&model.User{})
	if err != nil {
		return common.WrapError(err, "failed to delete user")
	}
	return nil
}

// SoftDelete soft deletes user
func (r *userRepository) SoftDelete(ctx context.Context, id string) error {
	owner, name, err := parseID(id)
	if err != nil {
		return err
	}
	user := &model.User{
		IsDeleted: true,
	}
	_, err = r.db.Context(ctx).Cols("is_deleted", "deleted_time").
		Where("owner = ? AND name = ?", owner, name).Update(user)
	if err != nil {
		return common.WrapError(err, "failed to soft delete user")
	}
	return nil
}

// List lists users with filter
func (r *userRepository) List(ctx context.Context, filter UserFilter) ([]*model.User, error) {
	session := r.buildFilterSession(ctx, filter)

	var users []*model.User
	err := session.Find(&users)
	if err != nil {
		return nil, common.WrapError(err, "failed to list users")
	}
	return users, nil
}

// ListWithPagination lists users with pagination
func (r *userRepository) ListWithPagination(ctx context.Context, filter UserFilter, offset, limit int) ([]*model.User, int64, error) {
	session := r.buildFilterSession(ctx, filter)

	total, err := session.Count(&model.User{})
	if err != nil {
		return nil, 0, common.WrapError(err, "failed to count users")
	}

	var users []*model.User
	err = session.Limit(limit, offset).Find(&users)
	if err != nil {
		return nil, 0, common.WrapError(err, "failed to list users")
	}

	return users, total, nil
}

// Count counts users
func (r *userRepository) Count(ctx context.Context, filter UserFilter) (int64, error) {
	session := r.buildFilterSession(ctx, filter)
	return session.Count(&model.User{})
}

// BatchCreate creates multiple users
func (r *userRepository) BatchCreate(ctx context.Context, users []*model.User) error {
	session := r.db.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	for _, user := range users {
		_, err = session.Context(ctx).Insert(user)
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch create users")
		}
	}

	return session.Commit()
}

// BatchUpdate updates multiple users
func (r *userRepository) BatchUpdate(ctx context.Context, users []*model.User, columns []string) error {
	session := r.db.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	for _, user := range users {
		s := session.Context(ctx).ID([]interface{}{user.Owner, user.Name})
		if len(columns) > 0 {
			s = s.Cols(columns...)
		}
		_, err = s.Update(user)
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch update users")
		}
	}

	return session.Commit()
}

// BatchDelete deletes multiple users
func (r *userRepository) BatchDelete(ctx context.Context, ids []string) error {
	session := r.db.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	for _, id := range ids {
		owner, name, err := parseID(id)
		if err != nil {
			session.Rollback()
			return err
		}
		_, err = session.Context(ctx).Where("owner = ? AND name = ?", owner, name).Delete(&model.User{})
		if err != nil {
			session.Rollback()
			return common.WrapError(err, "failed to batch delete users")
		}
	}

	return session.Commit()
}

// GetGlobalUsers gets all users globally
func (r *userRepository) GetGlobalUsers(ctx context.Context, filter UserFilter) ([]*model.User, error) {
	session := r.db.Context(ctx)
	if filter.Field != "" && filter.Value != "" {
		session = session.Where(builder.Like{filter.Field, filter.Value})
	}

	var users []*model.User
	err := session.Find(&users)
	if err != nil {
		return nil, common.WrapError(err, "failed to get global users")
	}
	return users, nil
}

// GetGlobalUsersWithPagination gets global users with pagination
func (r *userRepository) GetGlobalUsersWithPagination(ctx context.Context, filter UserFilter, offset, limit int) ([]*model.User, int64, error) {
	session := r.db.Context(ctx)
	if filter.Field != "" && filter.Value != "" {
		session = session.Where(builder.Like{filter.Field, filter.Value})
	}

	total, err := session.Count(&model.User{})
	if err != nil {
		return nil, 0, err
	}

	var users []*model.User
	err = session.Limit(limit, offset).Find(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CountGlobal counts global users
func (r *userRepository) CountGlobal(ctx context.Context, field, value string) (int64, error) {
	session := r.db.Context(ctx)
	if field != "" && value != "" {
		session = session.Where(builder.Like{field, value})
	}
	return session.Count(&model.User{})
}

// GetUsersByGroup gets users by group
func (r *userRepository) GetUsersByGroup(ctx context.Context, groupID string) ([]*model.User, error) {
	var users []*model.User
	err := r.db.Context(ctx).Where("groups LIKE ?", "%"+groupID+"%").Find(&users)
	if err != nil {
		return nil, common.WrapError(err, "failed to get users by group")
	}
	return users, nil
}

// AddUserToGroup adds user to group
func (r *userRepository) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	for _, g := range user.Groups {
		if g == groupID {
			return nil // Already in group
		}
	}

	user.Groups = append(user.Groups, groupID)
	return r.Update(ctx, user, []string{"groups"})
}

// RemoveUserFromGroup removes user from group
func (r *userRepository) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	newGroups := make([]string, 0, len(user.Groups))
	for _, g := range user.Groups {
		if g != groupID {
			newGroups = append(newGroups, g)
		}
	}

	user.Groups = newGroups
	return r.Update(ctx, user, []string{"groups"})
}

// UpdateMFA updates MFA configuration
func (r *userRepository) UpdateMFA(ctx context.Context, userID string, mfaConfig *model.MfaConfig) error {
	owner, name, err := parseID(userID)
	if err != nil {
		return err
	}

	user := &model.User{
		PreferredMfaType:  mfaConfig.PreferredMfaType,
		RecoveryCodes:     mfaConfig.RecoveryCodes,
		TotpSecret:        mfaConfig.TotpSecret,
		MfaPhoneEnabled:   mfaConfig.MfaPhoneEnabled,
		MfaEmailEnabled:   mfaConfig.MfaEmailEnabled,
		MfaRadiusEnabled:  mfaConfig.MfaRadiusEnabled,
		MfaRadiusUsername: mfaConfig.MfaRadiusUsername,
		MfaRadiusProvider: mfaConfig.MfaRadiusProvider,
		MfaPushEnabled:    mfaConfig.MfaPushEnabled,
		MfaPushReceiver:   mfaConfig.MfaPushReceiver,
		MfaPushProvider:   mfaConfig.MfaPushProvider,
	}

	_, err = r.db.Context(ctx).Cols(
		"preferred_mfa_type", "recovery_codes", "totp_secret",
		"mfa_phone_enabled", "mfa_email_enabled", "mfa_radius_enabled",
		"mfa_radius_username", "mfa_radius_provider", "mfa_push_enabled",
		"mfa_push_receiver", "mfa_push_provider",
	).Where("owner = ? AND name = ?", owner, name).Update(user)

	if err != nil {
		return common.WrapError(err, "failed to update MFA")
	}
	return nil
}

// GetMFAConfig gets MFA configuration
func (r *userRepository) GetMFAConfig(ctx context.Context, userID string) (*model.MfaConfig, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &model.MfaConfig{
		PreferredMfaType:  user.PreferredMfaType,
		RecoveryCodes:     user.RecoveryCodes,
		TotpSecret:        user.TotpSecret,
		MfaPhoneEnabled:   user.MfaPhoneEnabled,
		MfaEmailEnabled:   user.MfaEmailEnabled,
		MfaRadiusEnabled:  user.MfaRadiusEnabled,
		MfaRadiusUsername: user.MfaRadiusUsername,
		MfaRadiusProvider: user.MfaRadiusProvider,
		MfaPushEnabled:    user.MfaPushEnabled,
		MfaPushReceiver:   user.MfaPushReceiver,
		MfaPushProvider:   user.MfaPushProvider,
	}, nil
}

// Search searches users
func (r *userRepository) Search(ctx context.Context, owner, keyword string, fields []string) ([]*model.User, error) {
	session := r.db.Context(ctx).Where("owner = ?", owner)

	if len(fields) > 0 && keyword != "" {
		cond := builder.NewCond()
		for _, field := range fields {
			cond = cond.Or(builder.Like{field, "%" + keyword + "%"})
		}
		session = session.Where(cond)
	}

	var users []*model.User
	err := session.Find(&users)
	if err != nil {
		return nil, common.WrapError(err, "failed to search users")
	}
	return users, nil
}

// GetStatistics gets user statistics
func (r *userRepository) GetStatistics(ctx context.Context, owner string) (*UserStatistics, error) {
	stats := &UserStatistics{
		ByType:   make(map[string]int64),
		ByRegion: make(map[string]int64),
	}

	session := r.db.Context(ctx).Where("owner = ?", owner)

	total, err := session.Count(&model.User{})
	if err != nil {
		return nil, err
	}
	stats.Total = total

	// Count active users (not forbidden)
	active, err := session.Where("is_forbidden = ?", false).Count(&model.User{})
	if err != nil {
		return nil, err
	}
	stats.Active = active
	stats.Inactive = total - active

	// Count admins
	admin, err := session.Where("is_admin = ?", true).Count(&model.User{})
	if err != nil {
		return nil, err
	}
	stats.Admin = admin

	// Count online users
	online, err := session.Where("is_online = ?", true).Count(&model.User{})
	if err != nil {
		return nil, err
	}
	stats.Online = online

	return stats, nil
}

// buildFilterSession builds filter session
func (r *userRepository) buildFilterSession(ctx context.Context, filter UserFilter) *xorm.Session {
	session := r.db.Context(ctx)

	if filter.Owner != "" {
		session = session.Where("owner = ?", filter.Owner)
	}

	if filter.Field != "" && filter.Value != "" {
		session = session.Where(builder.Like{filter.Field, "%" + filter.Value + "%"})
	}

	if filter.IsAdmin != nil {
		session = session.Where("is_admin = ?", *filter.IsAdmin)
	}

	if filter.IsForbidden != nil {
		session = session.Where("is_forbidden = ?", *filter.IsForbidden)
	}

	if filter.IsDeleted != nil {
		session = session.Where("is_deleted = ?", *filter.IsDeleted)
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

// parseID parses user ID
func parseID(id string) (owner, name string, err error) {
	parts := splitID(id)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid user ID format: %s", id)
	}
	return parts[0], parts[1], nil
}

// splitID splits ID into owner and name
func splitID(id string) []string {
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
