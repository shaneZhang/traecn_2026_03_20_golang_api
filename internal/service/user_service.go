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

package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/casdoor/casdoor/internal/common"
	"github.com/casdoor/casdoor/internal/dto"
	"github.com/casdoor/casdoor/internal/model"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/util"
	"github.com/casdoor/casdoor/xlsx"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

// UserService defines user service interface
type UserService interface {
	// CRUD operations
	GetUser(ctx context.Context, id string) (*dto.UserResponse, error)
	GetUserByEmail(ctx context.Context, owner, email string) (*dto.UserResponse, error)
	GetUserByPhone(ctx context.Context, owner, phone string) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error

	// List operations
	ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	ListGlobalUsers(ctx context.Context, page, pageSize int, field, value, sortField, sortOrder string) (*dto.ListUsersResponse, error)

	// Batch operations
	BatchCreateUsers(ctx context.Context, req *dto.ImportUsersRequest) (*dto.ImportUsersResponse, error)
	BatchUpdateUsers(ctx context.Context, operation *dto.BatchUserOperation) error
	BatchDeleteUsers(ctx context.Context, ids []string) error

	// Import/Export
	ImportUsers(ctx context.Context, owner string, file io.Reader, fileType string) (*dto.ImportUsersResponse, error)
	ExportUsers(ctx context.Context, req *dto.ExportUsersRequest) ([]byte, string, error)

	// MFA operations
	SetupMFA(ctx context.Context, userID string, req *dto.MFASetupRequest) (*dto.MFASetupResponse, error)
	VerifyMFASetup(ctx context.Context, userID string, req *dto.MFASetupRequest) error
	EnableMFA(ctx context.Context, userID string, req *dto.MFASetupRequest) error
	DisableMFA(ctx context.Context, userID string) error
	VerifyMFACode(ctx context.Context, userID string, req *dto.MFAVerifyRequest) error
	GetMFAStatus(ctx context.Context, userID string) ([]*dto.MFASetupResponse, error)
	RecoverMFA(ctx context.Context, userID, recoveryCode string) error

	// Group operations
	AddUserToGroup(ctx context.Context, userID, groupID string) error
	RemoveUserFromGroup(ctx context.Context, userID, groupID string) error

	// Statistics
	GetUserStatistics(ctx context.Context, owner string) (*repository.UserStatistics, error)
}

// userService implements UserService
type userService struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
}

// NewUserService creates new user service
func NewUserService(userRepo repository.UserRepository, orgRepo repository.OrganizationRepository) UserService {
	return &userService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
	}
}

// GetUser gets user by ID
func (s *userService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

// GetUserByEmail gets user by email
func (s *userService) GetUserByEmail(ctx context.Context, owner, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, owner, email)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

// GetUserByPhone gets user by phone
func (s *userService) GetUserByPhone(ctx context.Context, owner, phone string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByPhone(ctx, owner, phone)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByOwnerAndName(ctx, req.Owner, req.Name)
	if existingUser != nil {
		return nil, common.ErrUserAlreadyExists
	}

	// Get organization for defaults
	org, err := s.orgRepo.GetByOwnerAndName(ctx, "admin", req.Owner)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Owner:             req.Owner,
		Name:              req.Name,
		CreatedTime:       util.GetCurrentTime(),
		Id:                util.GenerateId(),
		Type:              req.Type,
		DisplayName:       req.DisplayName,
		Email:             req.Email,
		Phone:             req.Phone,
		CountryCode:       req.CountryCode,
		Password:          req.Password,
		PasswordType:      "plain",
		Avatar:            org.DefaultAvatar,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Gender:            req.Gender,
		Birthday:          req.Birthday,
		Location:          req.Location,
		Address:           req.Address,
		Affiliation:       req.Affiliation,
		Title:             req.Title,
		Homepage:          req.Homepage,
		Bio:               req.Bio,
		Tag:               req.Tag,
		Region:            req.Region,
		Language:          req.Language,
		Score:             req.Score,
		SignupApplication: req.SignupApplication,
		Properties:        req.Properties,
		Groups:            req.Groups,
	}

	// Set defaults
	if user.Type == "" {
		user.Type = "normal-user"
	}
	if user.DisplayName == "" {
		user.DisplayName = user.Name
	}
	if user.SignupApplication == "" {
		user.SignupApplication = org.DefaultApplication
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// UpdateUser updates user
func (s *userService) UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.CountryCode != "" {
		user.CountryCode = req.CountryCode
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Birthday != "" {
		user.Birthday = req.Birthday
	}
	if req.Location != "" {
		user.Location = req.Location
	}
	if len(req.Address) > 0 {
		user.Address = req.Address
	}
	if req.Affiliation != "" {
		user.Affiliation = req.Affiliation
	}
	if req.Title != "" {
		user.Title = req.Title
	}
	if req.Homepage != "" {
		user.Homepage = req.Homepage
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Tag != "" {
		user.Tag = req.Tag
	}
	if req.Region != "" {
		user.Region = req.Region
	}
	if req.Language != "" {
		user.Language = req.Language
	}
	if req.Score != 0 {
		user.Score = req.Score
	}

	user.IsAdmin = req.IsAdmin
	user.IsForbidden = req.IsForbidden
	user.Properties = req.Properties
	user.Groups = req.Groups

	user.UpdatedTime = util.GetCurrentTime()

	err = s.userRepo.Update(ctx, user, nil)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

// DeleteUser deletes user
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers lists users
func (s *userService) ListUsers(ctx context.Context, req *dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	filter := repository.UserFilter{
		Owner:     req.Owner,
		GroupName: req.GroupName,
		Field:     req.Field,
		Value:     req.Value,
		SortField: req.SortField,
		SortOrder: req.SortOrder,
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.PageSize

	users, total, err := s.userRepo.ListWithPagination(ctx, filter, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.toUserResponse(user)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.ListUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// ListGlobalUsers lists global users
func (s *userService) ListGlobalUsers(ctx context.Context, page, pageSize int, field, value, sortField, sortOrder string) (*dto.ListUsersResponse, error) {
	filter := repository.UserFilter{
		Field:     field,
		Value:     value,
		SortField: sortField,
		SortOrder: sortOrder,
	}

	if pageSize == 0 {
		pageSize = 10
	}
	if page == 0 {
		page = 1
	}

	offset := (page - 1) * pageSize

	users, total, err := s.userRepo.GetGlobalUsersWithPagination(ctx, filter, offset, pageSize)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = s.toUserResponse(user)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dto.ListUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// BatchCreateUsers creates multiple users
func (s *userService) BatchCreateUsers(ctx context.Context, req *dto.ImportUsersRequest) (*dto.ImportUsersResponse, error) {
	resp := &dto.ImportUsersResponse{
		Total: len(req.Users),
	}

	users := make([]*model.User, len(req.Users))
	for i, reqUser := range req.Users {
		org, _ := s.orgRepo.GetByOwnerAndName(ctx, "admin", reqUser.Owner)
		if org == nil {
			resp.Failed++
			resp.Errors = append(resp.Errors, fmt.Sprintf("organization not found for user %s", reqUser.Name))
			continue
		}

		users[i] = &model.User{
			Owner:             reqUser.Owner,
			Name:              reqUser.Name,
			CreatedTime:       util.GetCurrentTime(),
			Id:                util.GenerateId(),
			Type:              reqUser.Type,
			DisplayName:       reqUser.DisplayName,
			Email:             reqUser.Email,
			Phone:             reqUser.Phone,
			CountryCode:       reqUser.CountryCode,
			Password:          reqUser.Password,
			PasswordType:      "plain",
			Avatar:            org.DefaultAvatar,
			FirstName:         reqUser.FirstName,
			LastName:          reqUser.LastName,
			Gender:            reqUser.Gender,
			Birthday:          reqUser.Birthday,
			Location:          reqUser.Location,
			Address:           reqUser.Address,
			Affiliation:       reqUser.Affiliation,
			Title:             reqUser.Title,
			Homepage:          reqUser.Homepage,
			Bio:               reqUser.Bio,
			Tag:               reqUser.Tag,
			Region:            reqUser.Region,
			Language:          reqUser.Language,
			Score:             reqUser.Score,
			SignupApplication: reqUser.SignupApplication,
			Properties:        reqUser.Properties,
			Groups:            reqUser.Groups,
		}

		if users[i].Type == "" {
			users[i].Type = "normal-user"
		}
		if users[i].DisplayName == "" {
			users[i].DisplayName = users[i].Name
		}
		if users[i].SignupApplication == "" {
			users[i].SignupApplication = org.DefaultApplication
		}
	}

	err := s.userRepo.BatchCreate(ctx, users)
	if err != nil {
		resp.Failed = len(req.Users)
		resp.Errors = append(resp.Errors, err.Error())
		return resp, nil
	}

	resp.Success = len(req.Users)
	return resp, nil
}

// BatchUpdateUsers updates multiple users
func (s *userService) BatchUpdateUsers(ctx context.Context, operation *dto.BatchUserOperation) error {
	users := make([]*model.User, len(operation.UserIds))
	for i, id := range operation.UserIds {
		user, err := s.userRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}

		switch operation.Operation {
		case "enable":
			user.IsForbidden = false
		case "disable":
			user.IsForbidden = true
		case "make_admin":
			user.IsAdmin = true
		case "remove_admin":
			user.IsAdmin = false
		}

		users[i] = user
	}

	return s.userRepo.BatchUpdate(ctx, users, nil)
}

// BatchDeleteUsers deletes multiple users
func (s *userService) BatchDeleteUsers(ctx context.Context, ids []string) error {
	return s.userRepo.BatchDelete(ctx, ids)
}

// ImportUsers imports users from file
func (s *userService) ImportUsers(ctx context.Context, owner string, file io.Reader, fileType string) (*dto.ImportUsersResponse, error) {
	var users []*dto.CreateUserRequest

	switch fileType {
	case "xlsx":
		// Read all data from file
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		// Write to temp file and read using xlsx package
		rows, err := xlsx.ReadXlsxFileBytes(data)
		if err != nil {
			return nil, err
		}

		// Parse header
		if len(rows) < 2 {
			return nil, fmt.Errorf("empty file")
		}

		headers := rows[0]
		for i, row := range rows[1:] {
			user := s.parseRowToUser(headers, row)
			user.Owner = owner
			users = append(users, user)
			_ = i
		}

	case "csv":
		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}

		if len(rows) < 2 {
			return nil, fmt.Errorf("empty file")
		}

		headers := rows[0]
		for _, row := range rows[1:] {
			user := s.parseRowToUser(headers, row)
			user.Owner = owner
			users = append(users, user)
		}

	case "json":
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &users)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	return s.BatchCreateUsers(ctx, &dto.ImportUsersRequest{Users: users})
}

// ExportUsers exports users to file
func (s *userService) ExportUsers(ctx context.Context, req *dto.ExportUsersRequest) ([]byte, string, error) {
	filter := repository.UserFilter{
		Owner: req.Owner,
	}

	users, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, "", err
	}

	switch req.Format {
	case "xlsx":
		return s.exportToExcel(users, req.Fields)
	case "csv":
		return s.exportToCSV(users, req.Fields)
	case "json":
		return s.exportToJSON(users, req.Fields)
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", req.Format)
	}
}

// SetupMFA sets up MFA
func (s *userService) SetupMFA(ctx context.Context, userID string, req *dto.MFASetupRequest) (*dto.MFASetupResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	org, err := s.orgRepo.GetByOwnerAndName(ctx, "admin", user.Owner)
	if err != nil {
		return nil, err
	}

	resp := &dto.MFASetupResponse{
		MfaType:            req.MfaType,
		MfaRememberInHours: org.MfaRememberInHours,
	}

	switch req.MfaType {
	case "totp":
		// Generate TOTP secret
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      org.DisplayName,
			AccountName: user.Email,
		})
		if err != nil {
			return nil, err
		}
		resp.Secret = key.Secret()
		resp.URL = key.URL()

	case "sms":
		if user.Phone == "" && req.Dest == "" {
			return nil, fmt.Errorf("phone number required")
		}
		resp.Secret = user.Phone
		if req.Dest != "" {
			resp.Secret = req.Dest
		}

	case "email":
		if user.Email == "" && req.Dest == "" {
			return nil, fmt.Errorf("email required")
		}
		resp.Secret = user.Email
		if req.Dest != "" {
			resp.Secret = req.Dest
		}
	}

	// Generate recovery code
	resp.RecoveryCodes = []string{uuid.NewString()}

	return resp, nil
}

// VerifyMFASetup verifies MFA setup
func (s *userService) VerifyMFASetup(ctx context.Context, userID string, req *dto.MFASetupRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	switch req.MfaType {
	case "totp":
		valid := totp.Validate(req.Passcode, req.Secret)
		if !valid {
			return common.ErrInvalidMFACode
		}
		user.TotpSecret = req.Secret

	case "sms", "email":
		// In real implementation, verify the code sent to phone/email
		if req.Passcode == "" {
			return common.ErrInvalidMFACode
		}
		if req.MfaType == "sms" {
			user.MfaPhoneEnabled = true
		} else {
			user.MfaEmailEnabled = true
		}
	}

	return s.userRepo.Update(ctx, user, nil)
}

// EnableMFA enables MFA
func (s *userService) EnableMFA(ctx context.Context, userID string, req *dto.MFASetupRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify setup first
	err = s.VerifyMFASetup(ctx, userID, req)
	if err != nil {
		return err
	}

	switch req.MfaType {
	case "totp":
		user.TotpSecret = req.Secret
	case "sms":
		user.MfaPhoneEnabled = true
		if req.Dest != "" {
			user.Phone = req.Dest
			user.CountryCode = req.CountryCode
		}
	case "email":
		user.MfaEmailEnabled = true
		if req.Dest != "" {
			user.Email = req.Dest
		}
	}

	user.PreferredMfaType = req.MfaType
	user.RecoveryCodes = []string{uuid.NewString()}

	return s.userRepo.Update(ctx, user, nil)
}

// DisableMFA disables MFA
func (s *userService) DisableMFA(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.PreferredMfaType = ""
	user.RecoveryCodes = nil
	user.TotpSecret = ""
	user.MfaPhoneEnabled = false
	user.MfaEmailEnabled = false
	user.MfaRadiusEnabled = false
	user.MfaPushEnabled = false

	return s.userRepo.Update(ctx, user, nil)
}

// VerifyMFACode verifies MFA code
func (s *userService) VerifyMFACode(ctx context.Context, userID string, req *dto.MFAVerifyRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	switch req.MfaType {
	case "totp":
		valid := totp.Validate(req.Passcode, user.TotpSecret)
		if !valid {
			return common.ErrInvalidMFACode
		}
	case "sms", "email":
		// In real implementation, verify the code
		if req.Passcode == "" {
			return common.ErrInvalidMFACode
		}
	}

	return nil
}

// GetMFAStatus gets MFA status
func (s *userService) GetMFAStatus(ctx context.Context, userID string) ([]*dto.MFASetupResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var statuses []*dto.MFASetupResponse

	// TOTP
	if user.TotpSecret != "" {
		statuses = append(statuses, &dto.MFASetupResponse{
			Enabled:     true,
			MfaType:     "totp",
			IsPreferred: user.PreferredMfaType == "totp",
		})
	}

	// SMS
	if user.MfaPhoneEnabled {
		statuses = append(statuses, &dto.MFASetupResponse{
			Enabled:     true,
			MfaType:     "sms",
			Secret:      util.GetMaskedPhone(user.Phone),
			IsPreferred: user.PreferredMfaType == "sms",
		})
	}

	// Email
	if user.MfaEmailEnabled {
		statuses = append(statuses, &dto.MFASetupResponse{
			Enabled:     true,
			MfaType:     "email",
			Secret:      util.GetMaskedEmail(user.Email),
			IsPreferred: user.PreferredMfaType == "email",
		})
	}

	return statuses, nil
}

// RecoverMFA recovers MFA using recovery code
func (s *userService) RecoverMFA(ctx context.Context, userID, recoveryCode string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	found := false
	newCodes := make([]string, 0, len(user.RecoveryCodes))
	for _, code := range user.RecoveryCodes {
		if code == recoveryCode {
			found = true
			continue
		}
		newCodes = append(newCodes, code)
	}

	if !found {
		return common.ErrRecoveryCodeInvalid
	}

	user.RecoveryCodes = newCodes
	return s.userRepo.Update(ctx, user, []string{"recovery_codes"})
}

// AddUserToGroup adds user to group
func (s *userService) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	return s.userRepo.AddUserToGroup(ctx, userID, groupID)
}

// RemoveUserFromGroup removes user from group
func (s *userService) RemoveUserFromGroup(ctx context.Context, userID, groupID string) error {
	return s.userRepo.RemoveUserFromGroup(ctx, userID, groupID)
}

// GetUserStatistics gets user statistics
func (s *userService) GetUserStatistics(ctx context.Context, owner string) (*repository.UserStatistics, error) {
	return s.userRepo.GetStatistics(ctx, owner)
}

// Helper functions

func (s *userService) toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		Owner:             user.Owner,
		Name:              user.Name,
		CreatedTime:       user.CreatedTime,
		Id:                user.Id,
		Type:              user.Type,
		DisplayName:       user.DisplayName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Avatar:            user.Avatar,
		PermanentAvatar:   user.PermanentAvatar,
		Email:             user.Email,
		EmailVerified:     user.EmailVerified,
		Phone:             user.Phone,
		CountryCode:       user.CountryCode,
		Region:            user.Region,
		Location:          user.Location,
		Affiliation:       user.Affiliation,
		Title:             user.Title,
		Homepage:          user.Homepage,
		Bio:               user.Bio,
		Tag:               user.Tag,
		Language:          user.Language,
		Gender:            user.Gender,
		Birthday:          user.Birthday,
		Education:         user.Education,
		Score:             user.Score,
		Karma:             user.Karma,
		Ranking:           user.Ranking,
		Balance:           user.Balance,
		IsAdmin:           user.IsAdmin,
		IsForbidden:       user.IsForbidden,
		SignupApplication: user.SignupApplication,
		Groups:            user.Groups,
		Properties:        user.Properties,
	}
}

func (s *userService) parseRowToUser(headers, row []string) *dto.CreateUserRequest {
	user := &dto.CreateUserRequest{}
	for i, header := range headers {
		if i >= len(row) {
			continue
		}
		value := row[i]
		switch header {
		case "name":
			user.Name = value
		case "displayName":
			user.DisplayName = value
		case "email":
			user.Email = value
		case "phone":
			user.Phone = value
		case "countryCode":
			user.CountryCode = value
		case "password":
			user.Password = value
		case "type":
			user.Type = value
		case "firstName":
			user.FirstName = value
		case "lastName":
			user.LastName = value
		case "gender":
			user.Gender = value
		case "birthday":
			user.Birthday = value
		case "location":
			user.Location = value
		case "affiliation":
			user.Affiliation = value
		case "title":
			user.Title = value
		case "homepage":
			user.Homepage = value
		case "bio":
			user.Bio = value
		case "tag":
			user.Tag = value
		case "region":
			user.Region = value
		case "language":
			user.Language = value
		case "groups":
			user.Groups = strings.Split(value, ";")
		}
	}
	return user
}

func (s *userService) exportToExcel(users []*model.User, fields []string) ([]byte, string, error) {
	// Prepare data
	data := [][]string{fields}

	// Write data
	for _, user := range users {
		row := make([]string, len(fields))
		for i, field := range fields {
			row[i] = fmt.Sprintf("%v", s.getFieldValue(user, field))
		}
		data = append(data, row)
	}

	buf, err := xlsx.WriteXlsxFileBytes(data)
	if err != nil {
		return nil, "", err
	}

	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

func (s *userService) exportToCSV(users []*model.User, fields []string) ([]byte, string, error) {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write headers
	writer.Write(fields)

	// Write data
	for _, user := range users {
		row := make([]string, len(fields))
		for i, field := range fields {
			row[i] = fmt.Sprintf("%v", s.getFieldValue(user, field))
		}
		writer.Write(row)
	}

	writer.Flush()
	return []byte(buf.String()), "text/csv", nil
}

func (s *userService) exportToJSON(users []*model.User, fields []string) ([]byte, string, error) {
	data := make([]map[string]interface{}, len(users))
	for i, user := range users {
		item := make(map[string]interface{})
		for _, field := range fields {
			item[field] = s.getFieldValue(user, field)
		}
		data[i] = item
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, "", err
	}

	return jsonData, "application/json", nil
}

func (s *userService) getFieldValue(user *model.User, field string) interface{} {
	switch field {
	case "owner":
		return user.Owner
	case "name":
		return user.Name
	case "displayName":
		return user.DisplayName
	case "email":
		return user.Email
	case "phone":
		return user.Phone
	case "type":
		return user.Type
	case "createdTime":
		return user.CreatedTime
	case "isAdmin":
		return user.IsAdmin
	case "isForbidden":
		return user.IsForbidden
	case "groups":
		return strings.Join(user.Groups, ";")
	default:
		return ""
	}
}
