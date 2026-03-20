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
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/i18n"
	"github.com/casdoor/casdoor/internal/repository"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/casdoor/casdoor/xlsx"
)

// BatchService 批量操作服务接口
type BatchService interface {
	// 用户导入导出
	ImportUsers(owner string, filePath string, operator *object.User, lang string) (*ImportResult, error)
	ImportUsersFromCSV(owner string, reader io.Reader, operator *object.User, lang string) (*ImportResult, error)
	ExportUsers(owner string, field, value string) ([][]string, error)
	ExportUsersToCSV(owner string, field, value string, writer io.Writer) error

	// 批量操作
	BatchDeleteUsers(owner string, userNames []string) (*BatchResult, error)
	BatchUpdateUserStatus(owner string, userNames []string, isDisabled bool) (*BatchResult, error)
}

// ImportResult 导入结果
type ImportResult struct {
	SuccessCount int      `json:"successCount"`
	FailCount    int      `json:"failCount"`
	SkipCount    int      `json:"skipCount"`
	TotalCount   int      `json:"totalCount"`
	Errors       []string `json:"errors,omitempty"`
}

// BatchResult 批量操作结果
type BatchResult struct {
	SuccessCount int      `json:"successCount"`
	FailCount    int      `json:"failCount"`
	TotalCount   int      `json:"totalCount"`
	Errors       []string `json:"errors,omitempty"`
}

type batchService struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
}

// NewBatchService 创建批量服务实例
func NewBatchService() BatchService {
	return &batchService{
		userRepo: repository.NewUserRepository(),
		orgRepo:  repository.NewOrganizationRepository(),
	}
}

// ImportUsers 从Excel文件导入用户
func (s *batchService) ImportUsers(owner string, filePath string, operator *object.User, lang string) (*ImportResult, error) {
	table := xlsx.ReadXlsxFile(filePath)

	if len(table) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}

	// 清理表头（移除#前缀）
	for idx, row := range table[0] {
		splitRow := strings.Split(row, "#")
		if len(splitRow) > 1 {
			table[0][idx] = splitRow[1]
		}
	}

	return s.importUsersFromTable(owner, table, operator, lang)
}

// ImportUsersFromCSV 从CSV文件导入用户
func (s *batchService) ImportUsersFromCSV(owner string, reader io.Reader, operator *object.User, lang string) (*ImportResult, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	table, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %v", err)
	}

	if len(table) == 0 {
		return nil, fmt.Errorf("empty CSV file")
	}

	return s.importUsersFromTable(owner, table, operator, lang)
}

// importUsersFromTable 从表格数据导入用户
func (s *batchService) importUsersFromTable(owner string, table [][]string, operator *object.User, lang string) (*ImportResult, error) {
	result := &ImportResult{
		TotalCount: len(table) - 1, // 减去表头
		Errors:     make([]string, 0),
	}

	uploadedUsers, err := s.stringArrayToUsers(table)
	if err != nil {
		return nil, err
	}

	if len(uploadedUsers) == 0 {
		return nil, fmt.Errorf("no valid user data found")
	}

	// 确定组织名称
	organizationName := uploadedUsers[0].Owner
	if organizationName == "" || !operator.IsGlobalAdmin() {
		organizationName = owner
	}

	// 验证组织是否存在
	org, err := s.orgRepo.GetByID(organizationName)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf(i18n.Translate(lang, "auth:The organization: %s does not exist"), organizationName)
	}

	// 获取现有用户Map用于查重
	existingUsers, err := s.userRepo.List(organizationName, -1, -1, "", "", "", "")
	if err != nil {
		return nil, err
	}

	existingUserMap := make(map[string]*object.User)
	for _, user := range existingUsers {
		existingUserMap[util.GetId(user.Owner, user.Name)] = user
	}

	// 预处理新用户
	newUsers := make([]*object.User, 0)
	for i, user := range uploadedUsers {
		userID := util.GetId(organizationName, user.Name)

		// 跳过已存在的用户
		if _, exists := existingUserMap[userID]; exists {
			result.SkipCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: User '%s' already exists, skipped", i+2, user.Name))
			continue
		}

		// 设置用户属性默认值
		s.preprocessNewUser(user, organizationName, org, operator)

		// 验证必填字段
		if err := s.validateUser(user); err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %s", i+2, err.Error()))
			continue
		}

		newUsers = append(newUsers, user)
	}

	// 批量插入用户
	if len(newUsers) > 0 {
		success, err := s.userRepo.CreateBatch(newUsers)
		if err != nil {
			return nil, err
		}
		if success {
			result.SuccessCount = len(newUsers)
		}
	}

	result.FailCount = result.TotalCount - result.SuccessCount - result.SkipCount
	return result, nil
}

// ExportUsers 导出用户为表格数据
func (s *batchService) ExportUsers(owner string, field, value string) ([][]string, error) {
	users, err := s.userRepo.List(owner, -1, -1, field, value, "", "")
	if err != nil {
		return nil, err
	}

	return s.usersToStringArray(users), nil
}

// ExportUsersToCSV 导出用户为CSV格式
func (s *batchService) ExportUsersToCSV(owner string, field, value string, writer io.Writer) error {
	table, err := s.ExportUsers(owner, field, value)
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	return csvWriter.WriteAll(table)
}

// BatchDeleteUsers 批量删除用户
func (s *batchService) BatchDeleteUsers(owner string, userNames []string) (*BatchResult, error) {
	result := &BatchResult{
		TotalCount: len(userNames),
		Errors:     make([]string, 0),
	}

	for _, userName := range userNames {
		user := &object.User{Owner: owner, Name: userName}
		success, err := s.userRepo.Delete(user)
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("User '%s': %s", userName, err.Error()))
			continue
		}
		if success {
			result.SuccessCount++
		} else {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("User '%s': not found", userName))
		}
	}

	return result, nil
}

// BatchUpdateUserStatus 批量更新用户状态
func (s *batchService) BatchUpdateUserStatus(owner string, userNames []string, isDisabled bool) (*BatchResult, error) {
	result := &BatchResult{
		TotalCount: len(userNames),
		Errors:     make([]string, 0),
	}

	for _, userName := range userNames {
		user, err := s.userRepo.GetByID(owner, userName)
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("User '%s': %s", userName, err.Error()))
			continue
		}

		user.IsDisabled = isDisabled
		user.UpdateHash()

		success, err := s.userRepo.Update(user, "is_disabled", "hash")
		if err != nil {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("User '%s': %s", userName, err.Error()))
			continue
		}
		if success {
			result.SuccessCount++
		} else {
			result.FailCount++
			result.Errors = append(result.Errors, fmt.Sprintf("User '%s': update failed", userName))
		}
	}

	return result, nil
}

// 辅助方法

// stringArrayToUsers 将字符串数组转换为用户对象
func (s *batchService) stringArrayToUsers(table [][]string) ([]*object.User, error) {
	if len(table) < 2 {
		return nil, fmt.Errorf("table must have at least header and one data row")
	}

	header := table[0]
	users := make([]*object.User, 0, len(table)-1)

	for i := 1; i < len(table); i++ {
		row := table[i]
		user := &object.User{}

		for j, colName := range header {
			if j >= len(row) {
				continue
			}

			value := row[j]
			if err := s.setUserField(user, colName, value); err != nil {
				return nil, fmt.Errorf("row %d, column '%s': %v", i+1, colName, err)
			}
		}

		users = append(users, user)
	}

	return users, nil
}

// usersToStringArray 将用户对象转换为字符串数组
func (s *batchService) usersToStringArray(users []*object.User) [][]string {
	// 定义导出字段顺序
	headers := []string{
		"owner", "name", "id", "type", "password", "displayName", "avatar",
		"email", "phone", "countryCode", "region", "location", "address",
		"affiliation", "title", "tag", "gender", "birthday", "education",
		"score", "ranking", "isAdmin", "isGlobalAdmin", "isForbidden", "isDeleted",
		"signupApplication", "createdTime", "updatedTime",
	}

	table := make([][]string, 0, len(users)+1)
	table = append(table, headers)

	for _, user := range users {
		row := make([]string, len(headers))
		for i, field := range headers {
			row[i] = s.getUserField(user, field)
		}
		table = append(table, row)
	}

	return table
}

// setUserField 设置用户字段值
func (s *batchService) setUserField(user *object.User, fieldName string, value string) error {
	fieldName = strings.ToLower(fieldName)
	value = strings.TrimSpace(value)

	switch fieldName {
	case "owner", "organization":
		user.Owner = value
	case "name", "username":
		user.Name = value
	case "id", "userid":
		user.Id = value
	case "type":
		user.Type = value
	case "password":
		user.Password = value
	case "displayname", "display_name", "display":
		user.DisplayName = value
	case "avatar", "icon", "photo":
		user.Avatar = value
	case "email", "mail":
		user.Email = value
	case "phone", "mobile", "telephone":
		user.Phone = value
	case "countrycode", "country_code":
		user.CountryCode = value
	case "region", "timezone":
		user.Region = value
	case "location", "city":
		user.Location = value
	case "address":
		if value != "" {
			user.Address = strings.Split(value, ";")
		}
	case "affiliation", "company", "org":
		user.Affiliation = value
	case "title", "jobtitle", "position":
		user.Title = value
	case "tag", "role", "group":
		user.Tag = value
	case "gender", "sex":
		user.Gender = value
	case "birthday", "dob":
		user.Birthday = value
	case "education", "degree", "school":
		user.Education = value
	case "score", "points":
		user.Score, _ = strconv.Atoi(value)
	case "ranking", "rank":
		user.Ranking, _ = strconv.Atoi(value)
	case "isadmin", "admin", "is_admin":
		user.IsAdmin = parseBool(value)
	case "isglobaladmin", "globaladmin", "is_global_admin":
		user.IsGlobalAdmin = parseBool(value)
	case "isforbidden", "forbidden", "is_forbidden":
		user.IsForbidden = parseBool(value)
	case "isdeleted", "deleted", "is_deleted":
		user.IsDeleted = parseBool(value)
	case "signupapplication", "signup_application", "app":
		user.SignupApplication = value
	case "createdtime", "created_time", "createdat":
		user.CreatedTime = value
	case "updatedtime", "updated_time", "updatedat":
		user.UpdatedTime = value
	}

	return nil
}

// getUserField 获取用户字段值
func (s *batchService) getUserField(user *object.User, fieldName string) string {
	fieldName = strings.ToLower(fieldName)

	switch fieldName {
	case "owner":
		return user.Owner
	case "name":
		return user.Name
	case "id":
		return user.Id
	case "type":
		return user.Type
	case "password":
		return "" // 不导出密码
	case "displayname":
		return user.DisplayName
	case "avatar":
		return user.Avatar
	case "email":
		return user.Email
	case "phone":
		return user.Phone
	case "countrycode":
		return user.CountryCode
	case "region":
		return user.Region
	case "location":
		return user.Location
	case "address":
		return strings.Join(user.Address, ";")
	case "affiliation":
		return user.Affiliation
	case "title":
		return user.Title
	case "tag":
		return user.Tag
	case "gender":
		return user.Gender
	case "birthday":
		return user.Birthday
	case "education":
		return user.Education
	case "score":
		return strconv.Itoa(user.Score)
	case "ranking":
		return strconv.Itoa(user.Ranking)
	case "isadmin":
		return strconv.FormatBool(user.IsAdmin)
	case "isglobaladmin":
		return strconv.FormatBool(user.IsGlobalAdmin)
	case "isforbidden":
		return strconv.FormatBool(user.IsForbidden)
	case "isdeleted":
		return strconv.FormatBool(user.IsDeleted)
	case "signupapplication":
		return user.SignupApplication
	case "createdtime":
		return user.CreatedTime
	case "updatedtime":
		return user.UpdatedTime
	default:
		return ""
	}
}

// preprocessNewUser 预处理新用户
func (s *batchService) preprocessNewUser(user *object.User, orgName string, org *object.Organization, operator *object.User) {
	user.Owner = orgName

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
	if user.Avatar == "" {
		user.Avatar = org.DefaultAvatar
	}
	if user.Region == "" {
		user.Region = operator.Region
	}
	if user.Address == nil {
		user.Address = []string{}
	}
	if user.CountryCode == "" {
		user.CountryCode = operator.CountryCode
	}
	if user.SignupApplication == "" {
		user.SignupApplication = org.DefaultApplication
	}
	if user.RegisterType == "" {
		user.RegisterType = "Upload Users"
	}
	if user.RegisterSource == "" {
		user.RegisterSource = util.GetId(operator.Owner, operator.Name)
	}

	// 密码加密
	if user.Password != "" {
		user.Password = conf.GetEncryptedPassword(user.Password, user.Owner, user.Name)
	}

	// 用户名转小写
	isUsernameLowered := conf.GetConfigBool("isUsernameLowered")
	if isUsernameLowered {
		user.Name = strings.ToLower(user.Name)
	}

	// 邮箱转小写
	if user.Email != "" {
		user.Email = strings.ToLower(user.Email)
	}

	// 格式化手机号
	if user.Phone != "" {
		user.Phone = util.GetSeperatedPhone(user.Phone)
	}

	user.UpdateHash()
}

// validateUser 验证用户必填字段
func (s *batchService) validateUser(user *object.User) error {
	if user.Name == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}

// parseBool 解析布尔值
func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "yes" || s == "1" || s == "y"
}
