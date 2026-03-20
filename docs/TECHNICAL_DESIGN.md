# Casdoor API 模块重构技术设计文档

## 1. 项目概述

### 1.1 重构目标
本次重构旨在对 Casdoor 项目的 API 模块进行架构升级，主要包括：
- 用户管理模块（用户 CRUD、批量导入导出、MFA 多因素认证）
- 组织架构模块（组织结构管理、层级关系）
- 应用授权模块（OAuth 应用管理、权限授予）

### 1.2 重构原则
- **分层架构**：采用清晰的分层架构，实现关注点分离
- **性能优化**：优化数据库查询，引入缓存机制
- **代码质量**：遵循 Go 语言最佳实践
- **兼容性**：保持与现有代码风格的兼容

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                             │
│  (Controllers/Handlers - HTTP Request/Response Handling)    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│     (Business Logic - Validation, Transactions)             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Repository Layer                          │
│       (Data Access - Database Operations)                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Cache Layer                              │
│           (Memory Cache - Performance Optimization)          │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 目录结构

```
casdoor/
├── api/                    # API 层（Controller）
│   ├── base.go            # 基础控制器
│   ├── user_api.go        # 用户 API
│   ├── organization_api.go # 组织架构 API
│   └── application_api.go # 应用授权 API
├── service/               # 服务层（业务逻辑）
│   ├── user_service.go    # 用户服务
│   ├── organization_service.go # 组织服务
│   ├── application_service.go  # 应用服务
│   └── mfa_service.go     # MFA 服务
├── repository/            # 数据访问层
│   ├── base_repository.go # 基础仓储
│   ├── user_repository.go # 用户仓储
│   ├── organization_repository.go # 组织仓储
│   ├── application_repository.go  # 应用仓储
│   └── mfa_repository.go  # MFA 仓储
├── cache/                 # 缓存层
│   ├── memory_cache.go    # 内存缓存实现
│   └── entity_cache.go    # 实体缓存
├── dto/                   # 数据传输对象
│   ├── common.go          # 通用 DTO
│   ├── user_dto.go        # 用户 DTO
│   ├── organization_dto.go # 组织 DTO
│   ├── application_dto.go # 应用 DTO
│   └── mfa_dto.go         # MFA DTO
├── container/             # 依赖注入容器
│   └── container.go       # IoC 容器
└── object/                # 领域模型（保持原有）
```

## 3. 各层详细设计

### 3.1 DTO 层（Data Transfer Object）

#### 设计原则
- 定义清晰的请求和响应结构
- 与 API 接口契约保持一致
- 包含必要的验证标签

#### 主要结构

```go
// 通用分页请求
type PaginationRequest struct {
    PageSize   int    `json:"pageSize" form:"pageSize"`
    Page       int    `json:"p" form:"p"`
    Field      string `json:"field" form:"field"`
    Value      string `json:"value" form:"value"`
    SortField  string `json:"sortField" form:"sortField"`
    SortOrder  string `json:"sortOrder" form:"sortOrder"`
}

// 标准响应
type Response struct {
    Status string      `json:"status"`
    Msg    string      `json:"msg"`
    Data   interface{} `json:"data"`
    Data2  interface{} `json:"data2"`
}
```

### 3.2 Repository 层

#### 设计原则
- 抽象数据访问逻辑
- 定义清晰的接口契约
- 支持事务操作
- 优化查询性能

#### 接口定义

```go
type UserRepository interface {
    GetById(id string) (*object.User, error)
    GetByEmail(owner, email string) (*object.User, error)
    List(owner string, offset, limit int, ...) ([]*object.User, error)
    Count(owner, field, value, groupName string) (int64, error)
    Create(user *object.User) (int64, error)
    Update(id string, user *object.User, columns []string, isAdmin bool) (int64, error)
    Delete(user *object.User) (int64, error)
}
```

#### 性能优化策略

1. **批量查询优化**
   - 使用 `IN` 查询替代循环单条查询
   - 实现预加载关联数据

2. **索引优化**
   - 为常用查询字段添加复合索引
   - 优化排序查询

3. **查询缓存**
   - 实现二级缓存机制
   - 缓存热点数据

### 3.3 Service 层

#### 设计原则
- 封装业务逻辑
- 实现事务管理
- 数据验证和转换
- 调用 Repository 完成数据操作

#### 核心服务

```go
type UserService interface {
    GetUser(id string) (*object.User, error)
    GetUsers(owner string, offset, limit int, ...) ([]*object.User, int64, error)
    CreateUser(user *object.User, lang string) (bool, error)
    UpdateUser(id string, user *object.User, ...) (bool, error)
    DeleteUser(user *object.User) (bool, error)
    SetPassword(userOwner, userName, oldPassword, newPassword, code, lang string) error
    ImportUsers(owner string, path string, userObj *object.User, lang string) (bool, error)
    ExportUsers(owner string, ...) ([]*object.User, error)
}
```

### 3.4 API 层（Controller）

#### 设计原则
- 处理 HTTP 请求/响应
- 参数解析和验证
- 调用 Service 层
- 统一响应格式

#### 基础控制器

```go
type BaseController struct {
    web.Controller
}

func (c *BaseController) ResponseOk(data ...interface{})
func (c *BaseController) ResponseError(msg string)
func (c *BaseController) IsGlobalAdmin() bool
func (c *BaseController) IsAdmin() bool
```

### 3.5 Cache 层

#### 设计原则
- 内存缓存实现
- 支持过期时间
- 自动清理过期数据

#### 缓存策略

```go
const (
    UserCacheTTL = 10 * time.Minute
    OrgCacheTTL  = 15 * time.Minute
    AppCacheTTL  = 10 * time.Minute
)
```

#### 缓存键设计

```
user:{id}          - 用户缓存
org:{id}           - 组织缓存
app:{id}           - 应用缓存
app:client:{id}    - 应用客户端缓存
```

## 4. 模块详细设计

### 4.1 用户管理模块

#### 功能清单
| 功能 | API | 方法 | 描述 |
|------|-----|------|------|
| 获取用户列表 | /api/get-users | GET | 分页获取用户 |
| 获取单个用户 | /api/get-user | GET | 获取用户详情 |
| 创建用户 | /api/add-user | POST | 新建用户 |
| 更新用户 | /api/update-user | POST | 更新用户信息 |
| 删除用户 | /api/delete-user | POST | 删除用户 |
| 设置密码 | /api/set-password | POST | 修改密码 |
| 导入用户 | /api/upload-users | POST | 批量导入 |
| 导出用户 | /api/export-users | GET | 批量导出 |

#### 性能优化点
1. 用户列表查询添加缓存
2. 批量导入使用事务批量插入
3. 用户详情查询预加载关联数据

### 4.2 组织架构模块

#### 功能清单
| 功能 | API | 方法 | 描述 |
|------|-----|------|------|
| 获取组织列表 | /api/get-organizations | GET | 分页获取组织 |
| 获取单个组织 | /api/get-organization | GET | 获取组织详情 |
| 创建组织 | /api/add-organization | POST | 新建组织 |
| 更新组织 | /api/update-organization | POST | 更新组织信息 |
| 删除组织 | /api/delete-organization | POST | 删除组织 |
| 获取组列表 | /api/get-groups | GET | 获取组列表 |
| 创建组 | /api/add-group | POST | 新建组 |
| 更新组 | /api/update-group | POST | 更新组信息 |
| 删除组 | /api/delete-group | POST | 删除组 |

#### 层级关系设计
- 组织（Organization）作为顶层容器
- 组（Group）支持树形结构
- 用户可属于多个组

### 4.3 应用授权模块

#### 功能清单
| 功能 | API | 方法 | 描述 |
|------|-----|------|------|
| 获取应用列表 | /api/get-applications | GET | 分页获取应用 |
| 获取单个应用 | /api/get-application | GET | 获取应用详情 |
| 创建应用 | /api/add-application | POST | 新建应用 |
| 更新应用 | /api/update-application | POST | 更新应用信息 |
| 删除应用 | /api/delete-application | POST | 删除应用 |
| OAuth 授权 | /api/oauth-grant | POST | OAuth 授权 |
| 获取 Token | /api/oauth-token | POST | 获取访问令牌 |

### 4.4 MFA 多因素认证模块

#### 功能清单
| 功能 | API | 方法 | 描述 |
|------|-----|------|------|
| 初始化 MFA | /api/mfa/setup/initiate | POST | 初始化 MFA 设置 |
| 验证 MFA | /api/mfa/setup/verify | POST | 验证 MFA 配置 |
| 启用 MFA | /api/mfa/setup/enable | POST | 启用 MFA |
| 删除 MFA | /api/mfa/delete | POST | 删除 MFA 配置 |
| 设置首选 MFA | /api/mfa/set-preferred | POST | 设置首选认证方式 |

#### 支持的 MFA 类型
- TOTP（基于时间的一次性密码）
- SMS（短信验证码）
- Email（邮件验证码）
- Radius（Radius 认证）
- Push（推送通知）

## 5. 依赖注入设计

### 5.1 容器设计

```go
type Container struct {
    engine *xorm.Engine
    
    // Repositories
    userRepo    repository.UserRepository
    orgRepo     repository.OrganizationRepository
    appRepo     repository.ApplicationRepository
    groupRepo   repository.GroupRepository
    mfaRepo     repository.MfaRepository
    
    // Services
    userService    service.UserService
    orgService     service.OrganizationService
    groupService   service.GroupService
    appService     service.ApplicationService
    oauthService   service.OAuthService
    mfaService     service.MfaService
    
    // Caches
    userCache   cache.UserCache
    orgCache    cache.OrganizationCache
    appCache    cache.ApplicationCache
}
```

### 5.2 初始化流程

```go
func GetContainer(engine *xorm.Engine) *Container {
    once.Do(func() {
        instance = &Container{engine: engine}
        instance.initRepositories()
        instance.initCaches()
        instance.initServices()
    })
    return instance
}
```

## 6. 数据库优化设计

### 6.1 索引优化

```sql
-- 用户表索引
CREATE INDEX idx_user_owner_name ON user (owner, name);
CREATE INDEX idx_user_email ON user (owner, email);
CREATE INDEX idx_user_phone ON user (owner, phone);
CREATE INDEX idx_user_created_time ON user (created_time);

-- 组织表索引
CREATE INDEX idx_org_owner ON organization (owner);
CREATE INDEX idx_org_name ON organization (name);

-- 应用表索引
CREATE INDEX idx_app_owner ON application (owner);
CREATE INDEX idx_app_client_id ON application (client_id);
CREATE INDEX idx_app_organization ON application (organization);

-- 组表索引
CREATE INDEX idx_group_owner ON `group` (owner);
CREATE INDEX idx_group_parent ON `group` (parent_id);
```

### 6.2 查询优化

1. **分页查询优化**
   - 使用 LIMIT + OFFSET 实现分页
   - 对于大数据量使用游标分页

2. **关联查询优化**
   - 使用预加载减少 N+1 查询
   - 合理使用 JOIN

3. **批量操作优化**
   - 使用批量插入替代循环插入
   - 使用事务保证数据一致性

## 7. 错误处理设计

### 7.1 错误类型

```go
type ErrorCode string

const (
    ErrBadRequest    ErrorCode = "BAD_REQUEST"
    ErrUnauthorized  ErrorCode = "UNAUTHORIZED"
    ErrForbidden     ErrorCode = "FORBIDDEN"
    ErrNotFound      ErrorCode = "NOT_FOUND"
    ErrInternal      ErrorCode = "INTERNAL_ERROR"
)
```

### 7.2 统一响应格式

```json
{
    "status": "error",
    "msg": "错误信息",
    "data": null
}
```

## 8. 安全设计

### 8.1 权限控制
- 基于角色的访问控制（RBAC）
- 管理员权限验证
- 资源所有权验证

### 8.2 数据脱敏
- 敏感字段自动脱敏
- API 响应数据过滤

### 8.3 输入验证
- 参数类型验证
- 业务规则验证
- SQL 注入防护

## 9. 测试策略

### 9.1 单元测试
- Repository 层测试
- Service 层测试
- 工具函数测试

### 9.2 集成测试
- API 接口测试
- 数据库操作测试

### 9.3 性能测试
- 并发请求测试
- 数据库查询性能测试

## 10. 部署方案

### 10.1 部署架构
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Nginx     │────▶│  Casdoor    │────▶│  Database   │
│  (反向代理)  │     │  (应用服务)  │     │  (MySQL)    │
└─────────────┘     └─────────────┘     └─────────────┘
```

### 10.2 配置管理
- 环境变量配置
- 配置文件热加载
- 日志配置

### 10.3 监控告警
- API 响应时间监控
- 错误率监控
- 数据库连接池监控

## 11. 迁移计划

### 11.1 迁移步骤
1. 创建新的分层目录结构
2. 实现 Repository 层
3. 实现 Service 层
4. 实现 API 层
5. 更新路由配置
6. 进行测试验证
7. 灰度发布

### 11.2 兼容性保证
- 保持原有 API 接口不变
- 保持响应格式兼容
- 支持渐进式迁移

## 12. 总结

本次重构通过引入分层架构、缓存机制和性能优化策略，显著提升了代码的可维护性和系统性能。主要收益包括：

1. **代码质量提升**：清晰的分层架构，职责分明
2. **性能优化**：缓存机制和数据库查询优化
3. **可测试性增强**：依赖注入和接口抽象
4. **可扩展性提升**：模块化设计，易于扩展新功能
