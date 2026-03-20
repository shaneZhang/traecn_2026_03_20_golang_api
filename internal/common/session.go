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

package common

import (
	"fmt"
	"strings"

	"github.com/casdoor/casdoor/object"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/core"
	"xorm.io/xorm"
)

// SessionBuilder 会话构建器
type SessionBuilder struct {
	engine     *xorm.Engine
	tableName  string
	owner      string
	offset     int
	limit      int
	field      string
	value      string
	sortField  string
	sortOrder  string
	conditions []builder.Cond
	columns    []string
}

// NewSessionBuilder 创建会话构建器
func NewSessionBuilder(owner string) *SessionBuilder {
	return &SessionBuilder{
		engine:     object.GetEngine(),
		owner:      owner,
		offset:     -1,
		limit:      -1,
		conditions: make([]builder.Cond, 0),
	}
}

// SetTableName 设置表名
func (sb *SessionBuilder) SetTableName(tableName string) *SessionBuilder {
	sb.tableName = tableName
	return sb
}

// SetPagination 设置分页
func (sb *SessionBuilder) SetPagination(offset, limit int) *SessionBuilder {
	sb.offset = offset
	sb.limit = limit
	return sb
}

// SetFilter 设置过滤条件
func (sb *SessionBuilder) SetFilter(field, value string) *SessionBuilder {
	sb.field = field
	sb.value = value
	return sb
}

// SetSort 设置排序
func (sb *SessionBuilder) SetSort(sortField, sortOrder string) *SessionBuilder {
	sb.sortField = sortField
	sb.sortOrder = sortOrder
	return sb
}

// AddCondition 添加查询条件
func (sb *SessionBuilder) AddCondition(cond builder.Cond) *SessionBuilder {
	if cond != nil {
		sb.conditions = append(sb.conditions, cond)
	}
	return sb
}

// SetColumns 设置查询列
func (sb *SessionBuilder) SetColumns(columns ...string) *SessionBuilder {
	sb.columns = columns
	return sb
}

// Build 构建会话
func (sb *SessionBuilder) Build() *xorm.Session {
	session := sb.engine.NewSession()

	if sb.owner != "" {
		session = session.Where("owner = ?", sb.owner)
	}

	// 添加过滤条件
	if sb.field != "" && sb.value != "" {
		fieldName := camelToSnake(sb.field)
		session = session.Where(fmt.Sprintf("%s like ?", fieldName), "%"+sb.value+"%")
	}

	// 添加自定义条件
	for _, cond := range sb.conditions {
		session = session.Where(cond)
	}

	// 排序
	if sb.sortField != "" {
		sortField := camelToSnake(sb.sortField)
		if sb.sortOrder == "asc" {
			session = session.Asc(sortField)
		} else {
			session = session.Desc(sortField)
		}
	} else {
		session = session.Desc("created_time")
	}

	// 分页
	if sb.limit != -1 && sb.offset != -1 {
		session = session.Limit(sb.limit, sb.offset)
	}

	// 指定列
	if len(sb.columns) > 0 {
		session = session.Cols(sb.columns...)
	}

	return session
}

// BuildForUser 构建用户专属会话
func (sb *SessionBuilder) BuildForUser() *xorm.Session {
	session := sb.engine.NewSession()

	if sb.owner != "" {
		session = session.Where("owner = ?", sb.owner)
	}

	// 添加过滤条件
	if sb.field != "" && sb.value != "" {
		fieldName := camelToSnake(sb.field)
		session = session.Where(fmt.Sprintf("%s like ?", fieldName), "%"+sb.value+"%")
	}

	// 添加自定义条件
	for _, cond := range sb.conditions {
		session = session.Where(cond)
	}

	// 排序
	if sb.sortField != "" {
		sortField := camelToSnake(sb.sortField)
		if sb.sortOrder == "asc" {
			session = session.Asc(sortField)
		} else {
			session = session.Desc(sortField)
		}
	} else {
		session = session.Desc("created_time")
	}

	// 分页
	if sb.limit != -1 && sb.offset != -1 {
		session = session.Limit(sb.limit, sb.offset)
	}

	return session
}

// Transaction 执行事务
func Transaction(f func(*xorm.Session) error) error {
	session := object.GetEngine().NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	if err := f(session); err != nil {
		if err := session.Rollback(); err != nil {
			return err
		}
		return err
	}

	return session.Commit()
}

// GetByID 根据ID获取记录
func GetByID(owner, name string, bean interface{}) (bool, error) {
	return object.GetEngine().ID(core.PK{owner, name}).Get(bean)
}

// camelToSnake 驼峰转下划线
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
