# 技术设计文档

## 1. 概述

### 1.1 项目背景

本项目是对 Casdoor 身份认证系统的 API 模块进行重构，目标是提升代码质量、可维护性和系统性能。

### 1.2 设计目标

- **高内聚低耦合**：通过分层架构实现职责分离
- **可测试性**：各层独立，便于单元测试
- **可扩展性**：支持新功能快速迭代
- **性能优化**：数据库查询优化和缓存策略

### 1.3 技术栈

- **语言**: Go 1.21+
- **Web 框架**: Beego v2
- **ORM**: Xorm
- **数据库**: MySQL/PostgreSQL/SQLite
- **缓存**: Redis（可选）
- **认证**: JWT, OAuth 2.0

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      API Gateway                        │
│              (Rate Limit, Auth, Logging)                │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│                    Handler Layer                        │
│         (HTTP Request/Response, Validation)             │
├─────────────────────────────────────────────────────────┤
│                    Service Layer                        │
│      (Business Logic, Transaction Management)           │
├─────────────────────────────────────────────────────────┤
│                  Repository Layer                       │
│         (Data Access, Query Optimization)               │
├─────────────────────────────────────────────────────────┤
│                     Model Layer                         │
│              (Entity Definitions)                       │
├─────────────────────────────────────────────────────────┤
│                      DTO Layer                          │
│         (Data Transfer Objects, Validation)             │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│                   Data Sources                          │
│         (Database, Cache, External APIs)                │
└─────────────────────────────────────────────────────────┘
```

### 2.2 分层职责

#### Handler Layer
- 处理 HTTP 请求和响应
- 参数解析和验证
- 调用 Service 层
- 返回统一格式的响应

#### Service Layer
- 实现业务逻辑
- 管理事务
- 协调多个 Repository
- 数据转换（DTO ↔ Model）

#### Repository Layer
- 数据库访问抽象
- 查询优化
- 批量操作
- 缓存集成点

#### Model Layer
- 定义实体结构
- 数据库映射（Xorm tags）
- 业务方法

#### DTO Layer
- 定义请求/响应结构
- 输入验证 tags
- 数据序列化

## 3. 模块设计

### 3.1 用户管理模块

#### 3.1.1 类图

```
┌─────────────────────┐
│   UserHandler       │
├─────────────────────┤
│ + CreateUser()      │
│ + GetUser()         │
│ + UpdateUser()      │
│ + DeleteUser()      │
│ + ListUsers()       │
│ + ImportUsers()     │
│ + ExportUsers()     │
│ + SetupMFA()        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   UserService       │
├─────────────────────┤
│ + CreateUser()      │
│ + GetUser()         │
│ + UpdateUser()      │
│ + DeleteUser()      │
│ + ListUsers()       │
│ + BatchCreate()     │
│ + ImportUsers()     │
│ + ExportUsers()     │
│ + SetupMFA()        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  UserRepository     │
├─────────────────────┤
│ + Create()          │
│ + GetByID()         │
│ + GetByEmail()      │
│ + Update()          │
│ + Delete()          │
│ + List()            │
│ + BatchCreate()     │
│ + Search()          │
└─────────────────────┘
```

#### 3.1.2 核心流程

**用户创建流程：**
```
1. Handler 接收请求并验证参数
2. Service 检查用户是否存在
3. Service 加密密码
4. Repository 保存到数据库
5. Service 返回响应
6. Handler 返回 JSON 响应
```

**用户导入流程：**
```
1. Handler 接收上传的文件
2. Service 解析文件（CSV/Excel/JSON）
3. Service 验证每行数据
4. Repository 批量插入
5. Service 返回导入结果
```

### 3.2 组织架构模块

#### 3.2.1 类图

```
┌─────────────────────────┐
│   OrganizationHandler   │
├─────────────────────────┤
│ + CreateOrganization()  │
│ + GetOrganization()     │
│ + UpdateOrganization()  │
│ + DeleteOrganization()  │
│ + GetHierarchy()        │
│ + MoveOrganization()    │
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│   OrganizationService   │
├─────────────────────────┤
│ + CreateOrganization()  │
│ + GetOrganization()     │
│ + UpdateOrganization()  │
│ + DeleteOrganization()  │
│ + GetHierarchy()        │
│ + MoveOrganization()    │
│ + GetChildren()         │
│ + GetAncestors()        │
└──────────┬──────────────┘
           │
           ▼
┌─────────────────────────┐
│  OrganizationRepository │
├─────────────────────────┤
│ + Create()              │
│ + GetByID()             │
│ + GetByOwnerAndName()   │
│ + Update()              │
│ + Delete()              │
│ + GetChildren()         │
│ + GetAncestors()        │
│ + GetDescendants()      │
└─────────────────────────┘
```

#### 3.2.2 层级关系存储

使用 **Materialized Path** 模式存储层级关系：

```sql
-- 组织表结构
CREATE TABLE organization (
    owner VARCHAR(100),
    name VARCHAR(100),
    parent_id VARCHAR(100),  -- 父组织 ID
    path VARCHAR(500),       -- 路径，如: /admin/root/org1/org2
    level INT,               -- 层级深度
    PRIMARY KEY (owner, name)
);
```

**查询优势：**
- 获取子组织：`WHERE parent_id = ?`
- 获取祖先：`WHERE ? LIKE path || '%'`
- 获取后代：`WHERE path LIKE ? || '%'`

### 3.3 应用授权模块

#### 3.3.1 OAuth 2.0 流程

```
┌─────────┐                                    ┌─────────┐
│  Client │                                    │ Server  │
└────┬────┘                                    └────┬────┘
     │                                              │
     │ 1. Authorization Request                     │
     │ ─────────────────────────────────────────>   │
     │    ?client_id=xxx&redirect_uri=xxx           │
     │                                              │
     │ 2. Authorization Grant (Code)                │
     │ <─────────────────────────────────────────   │
     │    redirect?code=xxx                         │
     │                                              │
     │ 3. Token Request                             │
     │ ─────────────────────────────────────────>   │
     │    grant_type=authorization_code             │
     │    &code=xxx&client_secret=xxx               │
     │                                              │
     │ 4. Access Token                              │
     │ <─────────────────────────────────────────   │
     │    {access_token, refresh_token, expires_in} │
     │                                              │
     │ 5. API Request with Token                    │
     │ ─────────────────────────────────────────>   │
     │    Authorization: Bearer xxx                 │
     │                                              │
     │ 6. Protected Resource                        │
     │ <─────────────────────────────────────────   │
     │                                              │
```

#### 3.3.2 Token 设计

**JWT Token 结构：**
```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT"
  },
  "payload": {
    "iss": "casdoor",
    "sub": "admin/user1",
    "aud": "application-client-id",
    "exp": 1234567890,
    "iat": 1234567890,
    "scope": "read write",
    "owner": "admin",
    "name": "user1"
  }
}
```

## 4. 数据库设计

### 4.1 用户表 (user)

```sql
CREATE TABLE "user" (
    owner VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_time VARCHAR(100),
    display_name VARCHAR(100),
    avatar VARCHAR(500),
    email VARCHAR(100),
    phone VARCHAR(100),
    password VARCHAR(200),
    password_type VARCHAR(100),
    password_salt VARCHAR(100),
    password_change_required BOOLEAN DEFAULT FALSE,
    is_forbidden BOOLEAN DEFAULT FALSE,
    is_deleted BOOLEAN DEFAULT FALSE,
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    signup_application VARCHAR(100),
    mfa_config VARCHAR(500),
    groups TEXT,
    properties TEXT,
    score INT DEFAULT 0,
    PRIMARY KEY (owner, name)
);

-- 索引
CREATE INDEX idx_user_owner_email ON "user" (owner, email);
CREATE INDEX idx_user_owner_phone ON "user" (owner, phone);
CREATE INDEX idx_user_created_time ON "user" (created_time);
```

### 4.2 组织表 (organization)

```sql
CREATE TABLE organization (
    owner VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_time VARCHAR(100),
    display_name VARCHAR(100),
    website_url VARCHAR(100),
    logo VARCHAR(200),
    favicon VARCHAR(200),
    parent_id VARCHAR(100),
    password_type VARCHAR(100),
    password_salt VARCHAR(100),
    password_options TEXT,
    default_avatar VARCHAR(200),
    default_application VARCHAR(100),
    tags TEXT,
    languages VARCHAR(255),
    master_password VARCHAR(200),
    enable_soft_deletion BOOLEAN DEFAULT FALSE,
    is_profile_public BOOLEAN DEFAULT FALSE,
    client_id VARCHAR(100),
    client_secret VARCHAR(100),
    PRIMARY KEY (owner, name)
);

-- 索引
CREATE INDEX idx_organization_parent_id ON organization (parent_id);
CREATE INDEX idx_organization_client_id ON organization (client_id);
```

### 4.3 应用表 (application)

```sql
CREATE TABLE application (
    owner VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_time VARCHAR(100),
    display_name VARCHAR(100),
    logo VARCHAR(200),
    homepage_url VARCHAR(100),
    description VARCHAR(100),
    organization VARCHAR(100),
    client_id VARCHAR(100),
    client_secret VARCHAR(100),
    redirect_uris TEXT,
    token_format VARCHAR(100),
    expire_in_hours FLOAT DEFAULT 168,
    refresh_expire_in_hours FLOAT DEFAULT 720,
    enable_password BOOLEAN DEFAULT TRUE,
    enable_sign_up BOOLEAN DEFAULT TRUE,
    enable_signin_session BOOLEAN DEFAULT FALSE,
    grant_types TEXT,
    tags TEXT,
    is_shared BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (owner, name)
);

-- 索引
CREATE INDEX idx_application_client_id ON application (client_id);
CREATE INDEX idx_application_organization ON application (organization);
```

## 5. 性能优化策略

### 5.1 数据库优化

#### 5.1.1 索引策略

1. **覆盖索引**: 包含查询所需的所有字段
2. **复合索引**: 多字段联合查询优化
3. **前缀索引**: 长文本字段的前缀索引

#### 5.1.2 查询优化

```go
// 优化前：查询所有字段
session.Find(&users)

// 优化后：只查询需要的字段
session.Cols("owner", "name", "email", "created_time").Find(&users)
```

#### 5.1.3 批量操作

```go
// 使用批量插入替代单条插入
func (r *userRepository) BatchCreate(ctx context.Context, users []*model.User) error {
    session := r.db.Context(ctx).NewSession()
    defer session.Close()
    
    err := session.Begin()
    if err != nil {
        return err
    }
    
    for _, user := range users {
        _, err := session.Insert(user)
        if err != nil {
            session.Rollback()
            return err
        }
    }
    
    return session.Commit()
}
```

### 5.2 缓存策略

#### 5.2.1 多级缓存

```
┌─────────────┐
│   Client    │ ← 浏览器缓存
└──────┬──────┘
       │
┌──────▼──────┐
│    CDN      │ ← 静态资源缓存
└──────┬──────┘
       │
┌──────▼──────┐
│   Redis     │ ← 分布式缓存
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │ ← 持久化存储
└─────────────┘
```

#### 5.2.2 缓存模式

**Cache-Aside 模式：**
```go
func (s *userService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%s", id)
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil {
        return cached.(*dto.UserResponse), nil
    }
    
    // 2. 从数据库获取
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    resp := s.toUserResponse(user)
    s.cache.Set(ctx, cacheKey, resp, time.Hour)
    
    return resp, nil
}
```

### 5.3 连接池优化

```go
// Xorm 连接池配置
engine, err := xorm.NewEngine(driverName, dataSourceName)
if err != nil {
    return nil, err
}

// 连接池设置
engine.SetMaxOpenConns(100)        // 最大连接数
engine.SetMaxIdleConns(10)         // 最大空闲连接数
engine.SetConnMaxLifetime(time.Hour) // 连接最大生命周期
```

## 6. 安全设计

### 6.1 认证机制

#### 6.1.1 JWT 认证

```go
// Token 生成
func GenerateToken(user *model.User, expireTime time.Duration) (string, error) {
    claims := jwt.MapClaims{
        "sub":   fmt.Sprintf("%s/%s", user.Owner, user.Name),
        "iss":   "casdoor",
        "iat":   time.Now().Unix(),
        "exp":   time.Now().Add(expireTime).Unix(),
        "owner": user.Owner,
        "name":  user.Name,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(privateKey)
}
```

#### 6.1.2 OAuth 2.0 安全

- **PKCE**: 支持授权码流程的 PKCE 扩展
- **State 参数**: 防止 CSRF 攻击
- **Scope 限制**: 最小权限原则
- **Token 过期**: 短期访问令牌 + 长期刷新令牌

### 6.2 密码安全

```go
// 密码加密
func HashPassword(password, salt string) string {
    // 使用 PBKDF2 或 bcrypt
    hash := pbkdf2.Key([]byte(password), []byte(salt), 10000, 64, sha256.New)
    return base64.StdEncoding.EncodeToString(hash)
}

// 密码验证
func VerifyPassword(password, salt, hash string) bool {
    return HashPassword(password, salt) == hash
}
```

### 6.3 输入验证

```go
// DTO 验证
type CreateUserRequest struct {
    Owner       string `json:"owner" binding:"required"`
    Name        string `json:"name" binding:"required,min=2,max=100"`
    Email       string `json:"email" binding:"omitempty,email"`
    Phone       string `json:"phone" binding:"omitempty,e164"`
    Password    string `json:"password" binding:"required,min=6"`
}
```

## 7. 测试策略

### 7.1 单元测试

```go
// Repository 单元测试
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB()
    repo := NewUserRepository(db)
    
    user := &model.User{
        Owner: "admin",
        Name:  "testuser",
        Email: "test@example.com",
    }
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    
    // 验证
    found, err := repo.GetByID(context.Background(), "admin/testuser")
    assert.NoError(t, err)
    assert.Equal(t, user.Email, found.Email)
}
```

### 7.2 集成测试

```go
// API 集成测试
func TestUserHandler_CreateUser(t *testing.T) {
    router := setupTestRouter()
    
    body := `{"owner":"admin","name":"test","email":"test@example.com","password":"123456"}`
    req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### 7.3 性能测试

```bash
# 使用 wrk 进行压力测试
wrk -t12 -c400 -d30s http://localhost:8000/api/users

# 使用 go test 进行基准测试
go test -bench=. -benchmem ./...
```

## 8. 部署架构

### 8.1 单机部署

```
┌─────────────────────────────────────┐
│           Load Balancer             │
│         (Nginx/HAProxy)             │
└─────────────────────────────────────┘
                   │
        ┌─────────┴─────────┐
        │                   │
┌───────▼───────┐   ┌───────▼───────┐
│   Casdoor     │   │   Casdoor     │
│   Instance 1  │   │   Instance 2  │
└───────┬───────┘   └───────┬───────┘
        │                   │
        └─────────┬─────────┘
                  │
        ┌─────────▼─────────┐
        │     Database      │
        │   (MySQL/Postgre) │
        └───────────────────┘
```

### 8.2 容器化部署

```yaml
# docker-compose.yml
version: '3.8'
services:
  casdoor:
    image: casdoor/casdoor:latest
    ports:
      - "8000:8000"
    environment:
      - DB_DRIVER=mysql
      - DB_HOST=db
      - DB_PORT=3306
    depends_on:
      - db
      - redis
    deploy:
      replicas: 3
      
  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=casdoor
    volumes:
      - db_data:/var/lib/mysql
      
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
      
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - casdoor
```

## 9. 监控与日志

### 9.1 监控指标

| 指标 | 类型 | 描述 |
|------|------|------|
| http_requests_total | Counter | HTTP 请求总数 |
| http_request_duration_seconds | Histogram | HTTP 请求处理时间 |
| db_query_duration_seconds | Histogram | 数据库查询时间 |
| active_users | Gauge | 当前活跃用户 |
| cache_hit_ratio | Gauge | 缓存命中率 |

### 9.2 日志规范

```go
// 结构化日志
log.Info().
    Str("method", "POST").
    Str("path", "/api/users").
    Str("user", "admin").
    Int("status", 201).
    Dur("latency", time.Since(start)).
    Msg("User created")
```

## 10. 扩展性设计

### 10.1 插件机制

```go
// 插件接口
type Plugin interface {
    Name() string
    Initialize(config map[string]interface{}) error
    BeforeRequest(ctx context.Context, req *http.Request) error
    AfterResponse(ctx context.Context, resp *http.Response) error
}

// 插件管理器
type PluginManager struct {
    plugins []Plugin
}

func (pm *PluginManager) Register(plugin Plugin) {
    pm.plugins = append(pm.plugins, plugin)
}
```

### 10.2 多租户支持

```go
// 租户上下文
type TenantContext struct {
    TenantID   string
    TenantName string
    Plan       string
    Quotas     map[string]int64
}

// 中间件注入租户信息
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")
        tenant := loadTenant(tenantID)
        c.Set("tenant", tenant)
        c.Next()
    }
}
```

## 11. 版本管理

### 11.1 API 版本策略

- **URL 版本**: `/api/v1/users`
- **Header 版本**: `Accept: application/vnd.casdoor.v1+json`

### 11.2 向后兼容

- 保留旧版本 API 至少 6 个月
- 使用 deprecation headers
- 提供迁移指南

## 12. 总结

本技术设计文档详细描述了重构后的系统架构、模块设计、数据库设计、性能优化、安全设计和部署方案。通过分层架构和清晰的职责分离，系统具备了良好的可维护性、可扩展性和高性能特性。
