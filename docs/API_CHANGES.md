# API 更新文档

## 概述
本文档描述了项目重构后的API接口变化，包括新增的接口、修改的接口以及废弃的接口。

## 架构变化

### 新架构层次
```
internal/
├── common/         # 公共组件
│   ├── errors.go       # 统一错误处理
│   ├── response.go     # 统一响应格式
│   ├── constants.go    # 常量定义
│   └── session.go      # 数据库会话管理
├── repository/     # 数据访问层
│   ├── user_repository.go
│   ├── organization_repository.go
│   └── application_repository.go
├── service/        # 业务逻辑层
│   ├── user_service.go
│   ├── organization_service.go
│   ├── application_service.go
│   ├── batch_service.go      # 批量操作
│   └── mfa_service.go        # MFA认证
└── handler/        # API接口层
    ├── user_handler.go
    ├── organization_handler.go
    ├── application_handler.go
    ├── batch_handler.go
    └── mfa_handler.go
```

## 用户管理模块 API

### 新增接口

#### 获取用户列表（分页）
- **接口**: `GET /api/get-users`
- **参数**:
  - `owner`: 组织名称（必填）
  - `page`: 页码（可选，默认：1）
  - `pageSize`: 每页大小（可选，默认：10）
  - `field`: 搜索字段（可选）
  - `value`: 搜索值（可选）
  - `sortField`: 排序字段（可选）
  - `sortOrder`: 排序方向（可选，asc/desc）
- **响应格式**:
  ```json
  {
      "status": "ok",
      "data": {
          "list": [...],
          "total": 100,
          "page": 1,
          "size": 10
      }
  }
  ```

#### 获取全局用户列表
- **接口**: `GET /api/get-global-users`
- **说明**: 跨组织获取所有用户（管理员权限）

#### 获取单个用户详情
- **接口**: `GET /api/get-user`
- **参数**:
  - `id`: 用户ID，格式：`owner/name`

#### 更新用户信息
- **接口**: `POST /api/update-user`
- **请求体**: User对象
- **说明**: 智能处理密码更新，如果密码为空则不更新密码

#### 添加用户
- **接口**: `POST /api/add-user`
- **请求体**: User对象
- **说明**: 自动处理密码加密、用户名小写、邮箱格式化等

#### 删除用户
- **接口**: `POST /api/delete-user`
- **请求体**: `{"owner": "org", "name": "user"}`

#### 更新用户密码
- **接口**: `POST /api/update-password`
- **参数**:
  - `userOwner`: 用户所属组织
  - `userName`: 用户名
  - `oldPassword`: 旧密码
  - `newPassword`: 新密码

#### 重置用户密码（管理员）
- **接口**: `POST /api/reset-password`
- **参数**:
  - `userOwner`: 用户所属组织
  - `userName`: 用户名
  - `newPassword`: 新密码

#### 启用/禁用用户
- **启用**: `POST /api/enable-user`
- **禁用**: `POST /api/disable-user`
- **参数**: `owner`, `name`

#### 获取排序用户列表
- **接口**: `GET /api/get-sorted-users`
- **参数**: `owner`, `sorter`（排序字段）, `limit`

#### 获取在线用户数量
- **接口**: `GET /api/get-online-user-count`
- **参数**: `owner`

## 组织架构模块 API

### 新增接口

#### 获取组织列表
- **接口**: `GET /api/get-organizations`
- **参数**: 分页、搜索、排序参数

#### 获取单个组织详情
- **接口**: `GET /api/get-organization`
- **参数**: `id`（格式：`admin/name`）

#### 根据用户ID获取组织
- **接口**: `GET /api/get-organization-by-userid`
- **参数**: `userID`

#### 更新组织信息
- **接口**: `POST /api/update-organization`
- **请求体**: Organization对象

#### 添加组织
- **接口**: `POST /api/add-organization`
- **请求体**: Organization对象

#### 删除组织
- **接口**: `POST /api/delete-organization`
- **请求体**: `{"name": "org"}`

#### 获取所有组织
- **接口**: `GET /api/get-all-organizations`

#### 获取组织层级关系
- **父组织**: `GET /api/get-parent-organizations`（参数：`name`）
- **直接子组织**: `GET /api/get-child-organizations`（参数：`name`）
- **所有子组织**: `GET /api/get-all-child-organizations`（参数：`name`）
- **完整层级结构**: `GET /api/get-organization-hierarchy`（参数：`name`）

#### 获取组织统计信息
- **接口**: `GET /api/get-organization-stats`
- **返回**: 用户数、应用数等统计数据

#### 启用/禁用组织
- **启用**: `POST /api/enable-organization`
- **禁用**: `POST /api/disable-organization`

#### 更新组织主题
- **接口**: `POST /api/update-organization-theme`
- **参数**: `name`, `theme`

## 应用授权模块 API

### 新增接口

#### 获取应用列表
- **接口**: `GET /api/get-applications`
- **参数**: `owner`（组织名称）、分页、搜索、排序参数

#### 获取单个应用详情
- **接口**: `GET /api/get-application`
- **参数**: `id`（格式：`owner/name`）

#### 根据ClientID获取应用
- **接口**: `GET /api/get-application-by-clientid`
- **参数**: `clientId`

#### 更新应用信息
- **接口**: `POST /api/update-application`
- **请求体**: Application对象

#### 添加应用
- **接口**: `POST /api/add-application`
- **请求体**: Application对象
- **说明**: 自动生成ClientID和ClientSecret

#### 删除应用
- **接口**: `POST /api/delete-application`
- **请求体**: `{"owner": "org", "name": "app"}`

#### 获取所有应用
- **接口**: `GET /api/get-all-applications`
- **参数**: `owner`

#### 获取组织下的应用
- **接口**: `GET /api/get-applications-by-organization`
- **参数**: `orgName`

#### 获取应用权限和角色
- **权限列表**: `GET /api/get-application-permissions`（参数：`appID`）
- **角色列表**: `GET /api/get-application-roles`（参数：`appID`）

#### 获取应用用户数量
- **接口**: `GET /api/get-application-user-count`
- **参数**: `appID`

#### 启用/禁用应用
- **启用**: `POST /api/enable-application`
- **禁用**: `POST /api/disable-application`

#### 更新应用主题
- **接口**: `POST /api/update-application-theme`
- **参数**: `owner`, `name`, `theme`

#### 客户端凭证管理
- **轮换凭证**: `POST /api/rotate-client-credentials`（生成新的ClientID和ClientSecret）
- **重新生成密钥**: `POST /api/regenerate-client-secret`（仅生成新的ClientSecret）

#### 验证重定向URI
- **接口**: `GET /api/validate-redirect-uri`
- **参数**: `clientId`, `redirectUri`

## 批量操作模块 API

### 新增接口

#### 导入用户
- **接口**: `POST /api/import-users`
- **参数**: `owner`, `file`（Excel/CSV文件）
- **支持格式**: .xlsx, .csv
- **响应**: 导入结果统计（成功、失败、跳过数量及错误详情）

#### 导出用户
- **接口**: `GET /api/export-users`
- **参数**: `owner`, `field`, `value`, `format`（csv/xlsx）
- **输出**: 文件下载

#### 批量删除用户
- **接口**: `POST /api/batch-delete-users`
- **参数**: `owner`, `userNames`（JSON数组）

#### 批量启用/禁用用户
- **批量禁用**: `POST /api/batch-disable-users`
- **批量启用**: `POST /api/batch-enable-users`
- **参数**: `owner`, `userNames`（JSON数组）

## MFA多因素认证模块 API

### 新增接口

#### MFA设置流程
- **初始化设置**: `POST /api/mfa/setup/initiate`（参数：`mfaType`）
- **验证设置**: `POST /api/mfa/setup/verify`（参数：`mfaType`, `passcode`）
- **启用MFA**: `POST /api/mfa/setup/enable`（参数：`mfaType`, `passcode`）

#### MFA状态管理
- **获取状态**: `GET /api/mfa/status`
- **禁用MFA**: `POST /api/mfa/delete`（参数：`mfaType`）
- **禁用所有MFA**: `POST /api/mfa/delete-all`
- **设置首选MFA**: `POST /api/mfa/set-preferred`（参数：`mfaType`）

#### MFA验证
- **验证**: `POST /api/mfa/verify`（参数：`mfaType`, `passcode`, `userId`）
- **恢复码验证**: `POST /api/mfa/recover`（参数：`recoveryCode`, `userId`）

#### 恢复码管理
- **获取恢复码**: `GET /api/mfa/recovery-codes`
- **重新生成**: `POST /api/mfa/regenerate-recovery-codes`

## 响应格式统一

### 成功响应
```json
{
    "status": "ok",
    "data": {...},
    "data2": {...}  // 可选，用于返回第二组数据
}
```

### 分页响应
```json
{
    "status": "ok",
    "data": {
        "list": [...],
        "total": 100,
        "page": 1,
        "size": 10
    }
}
```

### 错误响应
```json
{
    "status": "error",
    "msg": "错误信息（已国际化）",
    "data": {...}  // 可选，错误详情
}
```

## 错误码统一

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| `BAD_REQUEST` | 请求参数错误 | 400 |
| `UNAUTHORIZED` | 未授权访问 | 401 |
| `FORBIDDEN` | 权限不足 | 403 |
| `NOT_FOUND` | 资源不存在 | 404 |
| `INTERNAL_ERROR` | 服务器内部错误 | 500 |
| `VALIDATION_ERROR` | 验证失败 | 400 |
| `DUPLICATE_ENTRY` | 重复条目 | 409 |
| `USER_NOT_FOUND` | 用户不存在 | 404 |
| `USER_ALREADY_EXISTS` | 用户已存在 | 409 |
| `ORG_NOT_FOUND` | 组织不存在 | 404 |
| `APP_NOT_FOUND` | 应用不存在 | 404 |
| `MFA_NOT_ENABLED` | MFA未启用 | 400 |
| `MFA_INVALID_CODE` | MFA验证码无效 | 400 |

## 迁移指南

### 1. 更新路由配置
将原有的路由指向新的Handler实现：

```go
// 原路由
beego.Router("/get-users", &controllers.UserController{}, "GET:GetUsers")

// 新路由
userHandler := handler.NewUserHandler()
beego.Router("/api/get-users", userHandler, "GET:GetUsers")
```

### 2. 更新响应处理
新的API统一了响应格式，前端需要相应调整：

```javascript
// 原响应格式
{
    "code": 200,
    "message": "success",
    "result": [...]
}

// 新响应格式
{
    "status": "ok",
    "data": {
        "list": [...],
        "total": 100,
        "page": 1,
        "size": 10
    }
}
```

### 3. 错误处理
新的API返回统一的错误格式，包含国际化的错误信息：

```javascript
// 错误响应
{
    "status": "error",
    "msg": "用户名或密码错误"  // 根据Accept-Language自动国际化
}
```

### 4. 分页参数变化
原有的`p`参数改为`page`，`pageSize`保持一致：

```
// 原请求
GET /get-users?p=2&pageSize=20

// 新请求
GET /api/get-users?page=2&pageSize=20
```

## 向后兼容性

为了保持向后兼容性，建议：

1. **保留原有接口**: 原有接口继续可用，新接口使用`/api/`前缀
2. **逐步迁移**: 前端可以逐步迁移到新的API接口
3. **文档更新**: 更新API文档，推荐使用新的接口

## 废弃接口列表

以下原有接口建议在新版本中废弃：

- `/get-users` → 建议使用 `/api/get-users`
- `/get-user` → 建议使用 `/api/get-user`
- `/update-user` → 建议使用 `/api/update-user`
- `/add-user` → 建议使用 `/api/add-user`
- `/delete-user` → 建议使用 `/api/delete-user`

（其他模块以此类推...）

---

**文档版本**: v1.0  
**最后更新**: 2024年  
**适用版本**: v2.x
