package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/dto"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/repository"
	"github.com/casdoor/casdoor/util"
)

type UserService interface {
	GetUser(id string) (*object.User, error)
	GetUserByEmail(owner, email string) (*object.User, error)
	GetUserByPhone(owner, phone string) (*object.User, error)
	GetUserByUserId(owner, userId string) (*object.User, error)
	GetUserByFields(organization, username string) (*object.User, error)

	GetUsers(owner string, offset, limit int, field, value, sortField, sortOrder, groupName string) ([]*object.User, int64, error)
	GetGlobalUsers(offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, int64, error)
	GetSortedUsers(owner, sorter string, limit int) ([]*object.User, error)
	GetUserCount(owner, field, value, groupName string) (int64, error)
	GetOnlineUserCount(owner string, isOnline int) (int64, error)

	CreateUser(user *object.User, lang string) (bool, error)
	UpdateUser(id string, user *object.User, columns []string, isAdmin bool, lang string) (bool, error)
	DeleteUser(user *object.User) (bool, error)

	SetPassword(userOwner, userName, oldPassword, newPassword, code string, lang string) error
	CheckUserPassword(owner, name, password, lang string, enableCaptcha, ignorePassword, isPasswordWithLdapEnabled bool) (*object.User, error)

	ImportUsers(owner string, path string, userObj *object.User, lang string) (bool, error)
	ExportUsers(owner string, field, value, sortField, sortOrder string) ([]*object.User, error)

	GetMaskedUsers(users []*object.User) ([]*object.User, error)
	GetMaskedUser(user *object.User, isAdminOrSelf bool) (*object.User, error)
	ExtendUserWithRolesAndPermissions(user *object.User) error
	CheckUserPermission(requestUserId, userId string, strict bool, lang string) (bool, error)
}

type userService struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
	appRepo  repository.ApplicationRepository
}

func NewUserService(userRepo repository.UserRepository, orgRepo repository.OrganizationRepository, appRepo repository.ApplicationRepository) UserService {
	return &userService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
		appRepo:  appRepo,
	}
}

func (s *userService) GetUser(id string) (*object.User, error) {
	return s.userRepo.GetById(id)
}

func (s *userService) GetUserByEmail(owner, email string) (*object.User, error) {
	return s.userRepo.GetByEmail(owner, email)
}

func (s *userService) GetUserByPhone(owner, phone string) (*object.User, error) {
	return s.userRepo.GetByPhone(owner, phone)
}

func (s *userService) GetUserByUserId(owner, userId string) (*object.User, error) {
	return s.userRepo.GetByUserId(owner, userId)
}

func (s *userService) GetUserByFields(organization, username string) (*object.User, error) {
	return s.userRepo.GetByFields(organization, username)
}

func (s *userService) GetUsers(owner string, offset, limit int, field, value, sortField, sortOrder, groupName string) ([]*object.User, int64, error) {
	if offset == -1 || limit == -1 {
		users, err := s.userRepo.List(owner, -1, -1, field, value, sortField, sortOrder, groupName)
		if err != nil {
			return nil, 0, err
		}
		return users, int64(len(users)), nil
	}

	count, err := s.userRepo.Count(owner, field, value, groupName)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.userRepo.List(owner, offset, limit, field, value, sortField, sortOrder, groupName)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *userService) GetGlobalUsers(offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, int64, error) {
	if offset == -1 || limit == -1 {
		users, err := s.userRepo.ListGlobal(-1, -1, field, value, sortField, sortOrder)
		if err != nil {
			return nil, 0, err
		}
		return users, int64(len(users)), nil
	}

	count, err := s.userRepo.CountGlobal(field, value)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.userRepo.ListGlobal(offset, limit, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *userService) GetSortedUsers(owner, sorter string, limit int) ([]*object.User, error) {
	return s.userRepo.ListSorted(owner, sorter, limit)
}

func (s *userService) GetUserCount(owner, field, value, groupName string) (int64, error) {
	return s.userRepo.Count(owner, field, value, groupName)
}

func (s *userService) GetOnlineUserCount(owner string, isOnline int) (int64, error) {
	return s.userRepo.CountOnline(owner, isOnline)
}

func (s *userService) CreateUser(user *object.User, lang string) (bool, error) {
	emptyUser := &object.User{}
	msg := object.CheckUpdateUser(emptyUser, user, lang)
	if msg != "" {
		return false, errors.New(msg)
	}

	if user.RegisterType == "" {
		user.RegisterType = "Add User"
	}

	if user.CreatedTime == "" {
		user.CreatedTime = util.GetCurrentTime()
	}

	if user.Id == "" {
		user.Id = util.GenerateId()
	}

	if user.Type == "" {
		user.Type = "normal-user"
	}

	if user.DisplayName == "" {
		user.DisplayName = user.Name
	}

	affected, err := s.userRepo.Create(user)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *userService) UpdateUser(id string, user *object.User, columns []string, isAdmin bool, lang string) (bool, error) {
	oldUser, err := s.userRepo.GetById(id)
	if err != nil {
		return false, err
	}
	if oldUser == nil {
		return false, fmt.Errorf(i18n.Translate(lang, "general:The user: %s doesn't exist"), id)
	}

	if oldUser.Owner == "built-in" && oldUser.Name == "admin" && (user.Owner != "built-in" || user.Name != "admin") {
		return false, errors.New(i18n.Translate(lang, "auth:Unauthorized operation"))
	}

	if user.MfaEmailEnabled && user.Email == "" {
		return false, errors.New(i18n.Translate(lang, "user:MFA email is enabled but email is empty"))
	}

	if user.MfaPhoneEnabled && user.Phone == "" {
		return false, errors.New(i18n.Translate(lang, "user:MFA phone is enabled but phone number is empty"))
	}

	msg := object.CheckUpdateUser(oldUser, user, lang)
	if msg != "" {
		return false, errors.New(msg)
	}

	isUsernameLowered := conf.GetConfigBool("isUsernameLowered")
	if isUsernameLowered {
		user.Name = strings.ToLower(user.Name)
	}

	affected, err := s.userRepo.Update(id, user, columns, isAdmin)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *userService) DeleteUser(user *object.User) (bool, error) {
	if user.Owner == "built-in" && user.Name == "admin" {
		return false, errors.New("unauthorized operation")
	}

	affected, err := s.userRepo.Delete(user)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *userService) SetPassword(userOwner, userName, oldPassword, newPassword, code string, lang string) error {
	userId := util.GetId(userOwner, userName)

	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf(i18n.Translate(lang, "general:The user: %s doesn't exist"), userId)
	}

	org, err := s.orgRepo.GetByUser(user)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf(i18n.Translate(lang, "auth:the organization: %s is not found"), user.Owner)
	}

	if strings.Contains(newPassword, " ") {
		return errors.New(i18n.Translate(lang, "user:New password cannot contain blank space."))
	}

	msg := object.CheckPasswordComplexity(user, newPassword, lang)
	if msg != "" {
		return errors.New(msg)
	}

	if !object.CheckPasswordNotSameAsCurrent(user, newPassword, org) {
		return errors.New(i18n.Translate(lang, "user:The new password must be different from your current password"))
	}

	user.Password = newPassword
	user.UpdateUserPassword(org)
	user.NeedUpdatePassword = false
	user.LastChangePasswordTime = util.GetCurrentTime()

	_, err = s.userRepo.Update(userId, user, []string{"password", "password_salt", "need_update_password", "password_type", "last_change_password_time"}, false)
	return err
}

func (s *userService) CheckUserPassword(owner, name, password, lang string, enableCaptcha, ignorePassword, isPasswordWithLdapEnabled bool) (*object.User, error) {
	return object.CheckUserPassword(owner, name, password, lang, enableCaptcha, ignorePassword, isPasswordWithLdapEnabled)
}

func (s *userService) ImportUsers(owner string, path string, userObj *object.User, lang string) (bool, error) {
	return object.UploadUsers(owner, path, userObj, lang)
}

func (s *userService) ExportUsers(owner string, field, value, sortField, sortOrder string) ([]*object.User, error) {
	return s.userRepo.List(owner, -1, -1, field, value, sortField, sortOrder, "")
}

func (s *userService) GetMaskedUsers(users []*object.User) ([]*object.User, error) {
	return object.GetMaskedUsers(users)
}

func (s *userService) GetMaskedUser(user *object.User, isAdminOrSelf bool) (*object.User, error) {
	return object.GetMaskedUser(user, isAdminOrSelf)
}

func (s *userService) ExtendUserWithRolesAndPermissions(user *object.User) error {
	return object.ExtendUserWithRolesAndPermissions(user)
}

func (s *userService) CheckUserPermission(requestUserId, userId string, strict bool, lang string) (bool, error) {
	return object.CheckUserPermission(requestUserId, userId, strict, lang)
}

func ParsePaginationRequest(req *dto.PaginationRequest) (offset, limit int) {
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	return (req.Page - 1) * req.PageSize, req.PageSize
}
