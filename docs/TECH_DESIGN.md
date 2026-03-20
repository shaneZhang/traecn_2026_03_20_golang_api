# 技术设计文档

## 1. 项目概述

### 1.1 项目背景
本项目是一个企业级身份访问管理(IAM)系统，采用前后端分离架构，后端提供RESTful API服务。本次重构主要针对用户管理、组织架构、应用授权三个核心模块，以及批量导入导出、MFA多因素认证等功能。

### 1.2 重构目标
1. **架构优化**: 采用分层架构，提高代码的可维护性和可扩展性
2. **性能优化**: 优化数据库查询和操作，提升系统性能
3. **代码质量**: 遵循Go语言最佳实践，统一代码风格
4. **功能完善**: 增强批量操作和MFA功能
5. **兼容性**: 保持与现有代码风格和功能的兼容性

### 1.3 范围说明
- **用户管理模块**: 用户CRUD、批量导入导出、MFA认证
- **组织架构模块**: 组织结构管理、层级关系
- **应用授权模块**: OAuth应用管理、权限授予

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway / Beego                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Handler 层 (API接口层)                      │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐       │
│  │   User   │   Org    │   App    │  Batch   │   MFA    │       │
│  │ Handler  │ Handler  │ Handler  │ Handler  │ Handler  │       │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Service 层 (业务逻辑层)                      │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐       │
│  │   User   │   Org    │   App    │  Batch   │   MFA    │       │
│  │ Service  │ Service  │ Service  │ Service  │ Service  │       │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Repository 层 (数据访问层)                      │
│  ┌──────────┬──────────┬──────────┐                              │
│  │   User   │   Org    │   App    │                              │
│  │   Repo   │   Repo   │   Repo   │                              │
│  └──────────┴──────────┴──────────┘                              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                Common 层 (公共组件层)                            │
│  ┌──────────┬──────────┬──────────┬──────────┐                  │
│  │  Errors  │ Response │ Constants│ Session  │                  │
│  └──────────┴──────────┴──────────┴──────────┘                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Database (MySQL / MariaDB)                    │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 层次职责

#### 2.2.1 Handler层 (API接口层)
- **职责**: 处理HTTP请求，参数解析，响应格式化
- **设计原则**: 
  - 瘦Handler，胖Service
  - 统一参数校验
  - 统一响应格式
  - 错误处理和国际化

#### 2.2.2 Service层 (业务逻辑层)
- **职责**: 核心业务逻辑实现，事务管理，数据校验
- **设计原则**:
  - 业务逻辑封装
  - 事务管理
  - 跨模块业务编排
  - 可测试性

#### 2.2.3 Repository层 (数据访问层)
- **职责**: 数据库CRUD操作，查询优化
- **设计原则**:
  - 单一数据源操作
  - 查询性能优化
  - 与ORM框架解耦
  - 批量操作支持

#### 2.2.4 Common层 (公共组件层)
- **职责**: 通用工具，常量定义，错误码，响应格式
- **设计原则**:
  - 无状态
  - 可复用
  - 与业务逻辑无关

### 2.3 关键设计决策

#### 2.3.1 数据库会话管理
```go
// SessionBuilder 模式，支持链式调用
type SessionBuilder struct {
    engine     *xorm.Engine
    owner      string
    offset     int
    limit      int
    field      string
    value      string
    sortField  string
    sortOrder  string
    conditions []builder.Cond
}

// 使用示例
sb := common.NewSessionBuilder(owner).
    SetPagination(offset, limit).
    SetFilter(field, value).
    SetSort(sortField, sortOrder).
    AddCondition(cond)

session := sb.Build()
```

**优点**:
- 链式调用，代码简洁
- 统一的查询构建逻辑
- 易于扩展新的查询条件
- 自动处理字段映射（驼峰转下划线）

#### 2.3.2 错误处理机制

```go
// BusinessError 业务错误结构
type BusinessError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Err     error     `json:"-"`
}

// 预定义错误
var (
    ErrUserNotFound = NewBusinessError(ErrCodeUserNotFound, "User not found")
    ErrBadRequest   = NewBusinessError(ErrCodeBadRequest, "Invalid request parameters")
)
```

**设计要点**:
- 统一错误码定义
- 支持错误嵌套
- 国际化支持
- 与HTTP状态码映射

#### 2.3.3 响应格式统一

```go
// ApiResponse 统一响应格式
type ApiResponse struct {
    Status string      `json:"status"` // "ok" or "error"
    Msg    string      `json:"msg,omitempty"`
    Data   interface{} `json:"data,omitempty"`
    Data2  interface{} `json:"data2,omitempty"`
}

// PageResponse 分页响应格式
type PageResponse struct {
    List  interface{} `json:"list"`
    Total int64       `json:"total"`
    Page  int         `json:"page"`
    Size  int         `json:"size"`
}
```

**设计要点**:
- 前端友好的格式
- 支持分页数据
- 支持多数据返回（Data, Data2）
- 错误信息国际化

## 3. 核心模块设计

### 3.1 用户管理模块

#### 3.1.1 模块结构

```
user/
├── handler/
│   └── user_handler.go      # API接口处理
├── service/
│   └── user_service.go      # 业务逻辑实现
└── repository/
    └── user_repository.go   # 数据访问实现
```

#### 3.1.2 核心功能时序图

```
用户创建时序:

┌─────────┐     ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│ Client  │────▶│  UserHandler    │────▶│  UserService    │────▶│ UserRepository  │
└─────────┘     └─────────────────┘     └─────────────────┘     └─────────────────┘
                      │                       │                       │
                      │ 1. 参数校验            │ 2. 业务校验          │ 3. 唯一性检查
                      │                       │    - 密码加密         │    - 用户名唯一
                      │                       │    - 默认值填充       │    - 邮箱唯一
                      │                       │                       │    - 手机号唯一
                      │                       │                       │
                      │                       │ 4. 事务提交          │ 5. 插入数据
                      │◀──────────────────────│◀──────────────────────│
                      │ 6. 返回响应            │                       │
                      │                       │                       │
```

#### 3.1.3 数据库优化

**查询优化**:
```go
// 1. 分页查询使用Limit和Offset，避免全表扫描
func (r *userRepository) List(owner string, offset, limit int, ...) {
    sb := common.NewSessionBuilder(owner).
        SetPagination(offset, limit).  // LIMIT ?, ?
        SetFilter(field, value).        // WHERE field LIKE ?
        SetSort(sortField, sortOrder)   // ORDER BY ...
}

// 2. 条件查询使用索引
// 已建索引: owner, name, email, phone, id

// 3. 批量插入优化
func (r *userRepository) CreateBatch(users []*User) {
    batchSize := conf.GetConfigBatchSize()  // 可配置批量大小
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        batch := users[i:end]
        engine.Insert(batch)  // 批量插入
    }
}
```

### 3.2 组织架构模块

#### 3.2.1 模块结构

```
organization/
├── handler/
│   └── organization_handler.go  # API接口处理
├── service/
│   └── organization_service.go  # 业务逻辑实现
│       └── 组织层级计算
│       └── 统计信息聚合
└── repository/
    └── organization_repository.go  # 数据访问实现
```

#### 3.2.2 组织层级设计

**数据模型**:
```go
type Organization struct {
    Owner       string `xorm:"varchar(100) notnull pk"`
    Name        string `xorm:"varchar(100) notnull pk"`
    ParentId    string `xorm:"varchar(100)"`  // 父组织ID
    // ...
}
```

**层级查询优化**:
```go
// BFS遍历获取所有子组织（避免递归查询）
func (r *orgRepository) GetAllChildOrganizations(orgName string) {
    queue := []string{orgName}
    visited := make(map[string]bool)
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        if visited[current] {
            continue
        }
        visited[current] = true
        
        // 批量查询当前层的所有子组织
        children := r.GetChildOrganizations(current)  // IN 查询
        for _, child := range children {
            if !visited[child.Name] {
                result = append(result, child)
                queue = append(queue, child.Name)
            }
        }
    }
}
```

### 3.3 应用授权模块

#### 3.3.1 模块结构

```
application/
├── handler/
│   └── application_handler.go  # API接口处理
├── service/
│   └── application_service.go  # 业务逻辑实现
│       └── Client凭证管理
│       └── OAuth配置验证
└── repository/
    └── application_repository.go  # 数据访问实现
```

#### 3.3.2 Client凭证管理

```go
// 安全生成凭证
func (s *appService) GenerateClientCredentials() {
    clientID := util.GenerateClientId()      // UUID v4
    clientSecret := util.GenerateClientSecret()  // 32字节随机数
    
    // 存储时加密（可选）
    // clientSecret = hash(clientSecret)
}

// 凭证轮换
func (s *appService) RotateClientCredentials(owner, name string) {
    // 1. 生成新凭证
    // 2. 更新数据库
    // 3. 使旧凭证失效（立即或 grace period）
}
```

### 3.4 批量操作模块

#### 3.4.1 导入流程设计

```
导入流程:

┌─────────────────────────────────────────────────────────────────┐
│ 1. 文件上传验证                                                  │
│    - 文件格式检查 (.xlsx, .csv)                                  │
│    - 文件大小限制                                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. 数据解析                                                     │
│    - Excel/CSV解析                                             │
│    - 表头映射（支持#前缀）                                       │
│    - 类型转换                                                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. 数据校验                                                     │
│    - 必填字段检查                                               │
│    - 格式验证（邮箱、手机号）                                     │
│    - 唯一性检查（用户名、邮箱、手机号）                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. 数据预处理                                                   │
│    - 默认值填充                                                 │
│    - 密码加密                                                   │
│    - 格式标准化（用户名小写、邮箱小写、手机号格式化）             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 5. 批量入库                                                     │
│    - 事务包裹                                                   │
│    - 批量插入优化                                               │
│    - 错误处理（部分失败？全部失败？可配置）                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 6. 结果统计                                                     │
│    - 成功/失败/跳过计数                                         │
│    - 错误详情报告                                               │
└─────────────────────────────────────────────────────────────────┘
```

#### 3.4.2 导出流程设计

```
导出流程:

┌─────────────────────────────────────────────────────────────────┐
│ 1. 查询条件解析                                                 │
│    - owner过滤                                                 │
│    - 搜索条件                                                   │
│    - 导出字段选择                                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 2. 数据查询                                                     │
│    - 分页查询（避免内存溢出）                                   │
│    - 只查询需要的列                                             │
│    - 关联数据预加载（如需）                                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 3. 数据转换                                                     │
│    - 敏感字段掩码（密码、密钥）                                 │
│    - 格式标准化                                                 │
│    - 枚举值转义                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 4. 文件生成                                                     │
│    - CSV写入（流式写入，内存友好）                               │
│    - UTF-8 BOM处理（Excel兼容）                                 │
│    - 大文件支持（分块写入）                                     │
└─────────────────────────────────────────────────────────────────┘
```

### 3.5 MFA多因素认证模块

#### 3.5.1 MFA类型支持

| 类型 | 说明 | 安全性 | 用户体验 |
|------|------|--------|----------|
| TOTP (Google Authenticator) | 时间-based一次性密码 | 高 | 中 |
| SMS | 短信验证码 | 中 | 好 |
| Email | 邮件验证码 | 中 | 中 |
| Radius | 半径认证 | 高 | 中 |
| Push | 推送通知 | 高 | 好 |

#### 3.5.2 MFA流程设计

**启用流程**:
```
┌─────────┐     ┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  User   │────▶│  选择MFA类型    │────▶│  初始化配置     │────▶│  验证并启用     │
└─────────┘     └─────────────────┘     └─────────────────┘     └─────────────────┘
                                                          │
                                                          ▼
                                                ┌─────────────────┐
                                                │ 生成恢复码      │
                                                └─────────────────┘
                                                          │
                                                          ▼
                                                ┌─────────────────┐
                                                │ 设置首选MFA     │
                                                └─────────────────┘
```

**登录验证流程**:
```
┌─────────┐     ┌─────────────────┐     ┌─────────────────────────────┐
│  User   │────▶│  用户名密码验证  │────▶│  MFA验证（首选类型）        │
└─────────┘     └─────────────────┘     └─────────────────────────────┘
                                            │ 验证失败？
                                            │ 是  ↓  否
                                            ▼     ───▶ 登录成功
                                  ┌─────────────────┐
                                  │ 恢复码验证？     │
                                  └─────────────────┘
                                            │ 是  ↓  否
                                            ▼     ───▶ 验证失败
                                  ┌─────────────────┐
                                  │ 使用恢复码登录   │
                                  └─────────────────┘
```

## 4. 性能优化设计

### 4.1 数据库层面优化

#### 4.1.1 索引优化

**现有索引分析**:
```sql
-- 用户表索引
PRIMARY KEY (owner, name),
KEY idx_owner (owner),
KEY idx_email (email),
KEY idx_phone (phone),
KEY idx_id (id),
KEY idx_owner_id (owner, id)

-- 优化建议: 添加复合索引
KEY idx_owner_created (owner, created_time),
KEY idx_owner_is_online (owner, is_online)
```

#### 4.1.2 查询优化

**N+1查询问题解决**:
```go
// 避免循环查询（反模式）
for _, user := range users {
    org := GetOrganization(user.Owner)  // N次查询
}

// 优化: 批量IN查询
orgNames := collectOwnerNames(users)       // 收集所有owner名称
orgs := GetOrganizationsByNames(orgNames)  // 一次IN查询
orgMap := toMap(orgs)                      // 建立映射

for _, user := range users {
    user.Org = orgMap[user.Owner]          // 内存中关联
}
```

#### 4.1.3 批量操作优化

**批量插入**:
```go
// 单条插入（反模式）
for _, user := range users {
    engine.Insert(user)  // N次网络往返
}

// 批量插入（优化）
batchSize := 1000
for i := 0; i < len(users); i += batchSize {
    batch := users[i:min(i+batchSize, len(users))]
    engine.Insert(batch)  // 每1000条一次网络往返
}
```

**性能提升**: ~10x - 100x（取决于批量大小）

### 4.2 应用层面优化

#### 4.2.1 连接池配置

```go
// XORM连接池配置
engine.SetMaxOpenConns(100)    // 最大打开连接数
engine.SetMaxIdleConns(20)     // 最大空闲连接数
engine.SetConnMaxLifetime(300) // 连接最大生命周期（秒）

// 连接池状态监控
// SHOW STATUS LIKE 'Threads_%';
// Threads_connected: 当前连接数
// Threads_running: 活跃连接数
```

#### 4.2.2 缓存策略

**可缓存数据识别**:
| 数据类型 | 缓存建议 | TTL |
|---------|---------|-----|
| 应用配置（ClientID, ClientSecret） | ✅ 高优先级 | 1h |
| 组织配置 | ✅ 高优先级 | 30m |
| 用户权限 | ✅ 中优先级 | 15m |
| 用户基础信息 | ✅ 中优先级 | 5m |
| MFA配置 | ✅ 中优先级 | 5m |
| 在线状态 | ❌ 不建议 | - |
| 动态统计数据 | ❌ 不建议 | - |

**缓存实现建议**:
```go
// 使用Redis作为分布式缓存
// 注意: 缓存更新策略（写穿透、写回）
// 注意: 缓存击穿、穿透、雪崩防护

type CacheLayer struct {
    redis *redis.Client
    db    *xorm.Engine
}

func (c *CacheLayer) GetAppByClientID(clientID string) {
    // 1. 查缓存
    // 2. 缓存命中 -> 返回
    // 3. 缓存未命中 -> 查DB -> 写入缓存 -> 返回
}
```

#### 4.2.3 异步处理

**可异步场景**:
1. **批量导入** - 大文件导入可异步处理，返回jobID轮询结果
2. **邮件/短信发送** - MFA验证码、通知等
3. **日志审计** - 操作日志写入
4. **统计报表** - 非实时统计计算

## 5. 安全设计

### 5.1 密码安全

```go
// 密码加密策略（已实现）
func GetEncryptedPassword(password string, salt ...string) string {
    // 1. 盐值处理
    // 2. PBKDF2密钥派生
    // 3. SHA256哈希
    // 4. Base64编码
}

// 密码复杂度校验（建议添加）
func ValidatePassword(password string) error {
    // - 最小长度（建议8+）
    // - 包含大小写
    // - 包含数字
    // - 包含特殊字符
    // - 排除常见密码
}
```

### 5.2 Client凭证安全

```go
// ClientID: UUID v4，全局唯一，不可预测
// ClientSecret: 32字节加密随机数（256位熵）

// 存储安全建议:
// 1. 数据库中加密存储（AES-256）
// 2. 内存中使用后立即清零
// 3. 日志中掩码显示（只显示前4后4位）
// 4. 定期轮换机制
```

### 5.3 MFA安全

```go
// TOTP验证安全窗口
// - 默认可接受当前±1个时间窗口（共3个窗口）
// - 每个窗口30秒
// - 可接受窗口数可配置（1-10）

// 防暴力破解:
// - 连续失败锁定（如: 5次失败锁定15分钟）
// - 验证码尝试次数限制
// - 恢复码尝试次数限制

// 恢复码安全:
// - 每个恢复码只能使用一次
// - 生成后立即显示，服务端只存哈希值
// - 10个恢复码，建议用户安全保存
```

## 6. 测试策略

### 6.1 单元测试

**测试覆盖范围**:
- Repository层: CRUD操作、查询正确性
- Service层: 业务逻辑、边界条件、错误场景
- Common层: 工具函数、错误处理、响应格式

**测试覆盖率目标**:
- 核心模块: >80%
- 工具类: >90%
- 复杂业务逻辑: >85%

### 6.2 集成测试

**测试场景**:
1. **用户生命周期**: 创建→查询→更新→删除
2. **组织层级**: 创建层级结构→查询父/子组织→删除
3. **应用OAuth**: 创建应用→验证凭证→授权流程
4. **批量导入导出**: 各种格式、边界情况
5. **MFA完整流程**: 启用→验证→禁用→恢复

### 6.3 性能测试

**关键指标**:
| 操作 | 单用户响应时间 | 并发100用户响应时间 | TPS |
|------|---------------|-------------------|-----|
| 用户查询 | <100ms | <500ms | >500 |
| 用户创建 | <200ms | <1000ms | >200 |
| 批量导入(1万条) | <30s | - | - |
| MFA验证 | <200ms | <1000ms | >300 |

**测试工具**:
- Go标准库`testing.B`（单元性能）
- `wrk` / `hey`（HTTP接口）
- `locust` / `JMeter`（复杂场景）

## 7. 部署与运维

### 7.1 配置建议

**数据库配置**:
```ini
[database]
max_open_conns = 100
max_idle_conns = 20
conn_max_lifetime = 300
slow_query_threshold = 500  # 慢查询阈值(ms)
```

**批量操作配置**:
```ini
[batch]
max_batch_size = 1000       # 单批最大记录数
max_upload_size = 10485760  # 上传文件限制(10MB)
timeout = 300               # 批量操作超时
```

### 7.2 监控指标

**业务指标**:
- API请求量/错误率
- 各模块响应时间分布
- 批量操作成功率/耗时
- MFA启用率/验证成功率

**系统指标**:
- 数据库连接池状态
- 慢查询日志告警
- 内存/CPU使用率
- Goroutine数量

### 7.3 日志规范

```go
// 结构化日志建议
logger.WithFields(Fields{
    "module":   "user",
    "action":   "create",
    "user_id":  userID,
    "duration": durationMs,
    "status":   "success",  // or "error"
}).Info("User created")

// 错误日志包含堆栈
if err != nil {
    logger.WithError(err).
           WithField("stack", debug.Stack()).
           Error("Failed to create user")
}
```

## 8. 演进路线

### 8.1 短期优化(0-3个月)
1. 完成现有接口迁移到新架构
2. 性能基准测试和瓶颈修复
3. 监控告警配置
4. 缓存层引入（Redis）

### 8.2 中期优化(3-6个月)
1. 读写分离架构
2. 数据库分库分表准备
3. 微服务拆分可行性评估
4. GraphQL API支持

### 8.3 长期演进(6个月+)
1. 服务网格(Service Mesh)
2. 多活数据中心
3. 可观测性完善(OpenTelemetry)
4. AI驱动的安全分析

---

**文档版本**: v1.0  
**最后更新**: 2024年  
**作者**: Architecture Team
