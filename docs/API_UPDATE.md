# Casdoor API 更新文档

## 1. 概述

本文档描述了 Casdoor API 模块重构后的接口变更和新增功能。

## 2. API 接口列表

### 2.1 用户管理 API

#### 2.1.1 获取用户列表

**接口地址**：`GET /api/get-users`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 否 | 组织所有者 |
| pageSize | int | 否 | 每页数量，默认 10 |
| p | int | 否 | 页码，默认 1 |
| field | string | 否 | 搜索字段 |
| value | string | 否 | 搜索值 |
| sortField | string | 否 | 排序字段 |
| sortOrder | string | 否 | 排序方向（asc/desc） |
| groupName | string | 否 | 组名过滤 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": [
        {
            "owner": "built-in",
            "name": "admin",
            "createdTime": "2024-01-01T00:00:00+08:00",
            "id": "built-in/admin",
            "type": "normal-user",
            "displayName": "Admin",
            "email": "admin@example.com",
            "phone": "",
            "isAdmin": true,
            "groups": [],
            "roles": ["admin"],
            "permissions": []
        }
    ],
    "data2": 100
}
```

#### 2.1.2 获取单个用户

**接口地址**：`GET /api/get-user`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 否 | 用户 ID（owner/name 格式） |
| owner | string | 否 | 组织所有者 |
| email | string | 否 | 邮箱 |
| phone | string | 否 | 手机号 |
| userId | string | 否 | 用户 ID |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": {
        "owner": "built-in",
        "name": "admin",
        "createdTime": "2024-01-01T00:00:00+08:00",
        "id": "built-in/admin",
        "type": "normal-user",
        "displayName": "Admin",
        "email": "admin@example.com",
        "phone": "",
        "isAdmin": true,
        "groups": [],
        "roles": ["admin"],
        "permissions": [],
        "multiFactorAuths": []
    }
}
```

#### 2.1.3 创建用户

**接口地址**：`POST /api/add-user`

**请求体**：

```json
{
    "owner": "built-in",
    "name": "newuser",
    "type": "normal-user",
    "displayName": "New User",
    "email": "newuser@example.com",
    "phone": "",
    "password": "password123",
    "groups": [],
    "properties": {}
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.1.4 更新用户

**接口地址**：`POST /api/update-user`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 用户 ID |
| columns | string | 否 | 更新字段列表（逗号分隔） |

**请求体**：

```json
{
    "displayName": "Updated Name",
    "email": "updated@example.com"
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.1.5 删除用户

**接口地址**：`POST /api/delete-user`

**请求体**：

```json
{
    "owner": "built-in",
    "name": "user-to-delete"
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.1.6 设置密码

**接口地址**：`POST /api/set-password`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userOwner | string | 是 | 用户所属组织 |
| userName | string | 是 | 用户名 |
| oldPassword | string | 是 | 旧密码 |
| newPassword | string | 是 | 新密码 |
| code | string | 否 | 验证码 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "OK"
}
```

#### 2.1.7 批量导入用户

**接口地址**：`POST /api/upload-users`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 组织所有者 |
| file | file | 是 | CSV/XLSX 文件 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.1.8 导出用户

**接口地址**：`GET /api/export-users`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 组织所有者 |
| field | string | 否 | 过滤字段 |
| value | string | 否 | 过滤值 |
| sortField | string | 否 | 排序字段 |
| sortOrder | string | 否 | 排序方向 |

**响应**：CSV 文件下载

### 2.2 组织架构 API

#### 2.2.1 获取组织列表

**接口地址**：`GET /api/get-organizations`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 否 | 所有者 |
| pageSize | int | 否 | 每页数量 |
| p | int | 否 | 页码 |
| field | string | 否 | 搜索字段 |
| value | string | 否 | 搜索值 |
| sortField | string | 否 | 排序字段 |
| sortOrder | string | 否 | 排序方向 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": [
        {
            "owner": "admin",
            "name": "built-in",
            "createdTime": "2024-01-01T00:00:00+08:00",
            "displayName": "Built-in Organization",
            "websiteUrl": "",
            "passwordType": "plain",
            "passwordSalt": "",
            "mfaRememberInHours": 12
        }
    ],
    "data2": 10
}
```

#### 2.2.2 获取单个组织

**接口地址**：`GET /api/get-organization`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 组织 ID |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": {
        "owner": "admin",
        "name": "built-in",
        "createdTime": "2024-01-01T00:00:00+08:00",
        "displayName": "Built-in Organization",
        "websiteUrl": "",
        "passwordType": "plain",
        "passwordSalt": "",
        "mfaRememberInHours": 12
    }
}
```

#### 2.2.3 创建组织

**接口地址**：`POST /api/add-organization`

**请求体**：

```json
{
    "owner": "admin",
    "name": "new-org",
    "displayName": "New Organization",
    "websiteUrl": "https://example.com",
    "passwordType": "plain",
    "mfaRememberInHours": 12
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.2.4 更新组织

**接口地址**：`POST /api/update-organization`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 组织 ID |

**请求体**：

```json
{
    "displayName": "Updated Organization Name"
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.2.5 删除组织

**接口地址**：`POST /api/delete-organization`

**请求体**：

```json
{
    "owner": "admin",
    "name": "org-to-delete"
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.2.6 获取组列表

**接口地址**：`GET /api/get-groups`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 组织所有者 |
| pageSize | int | 否 | 每页数量 |
| p | int | 否 | 页码 |
| withTree | string | 否 | 是否返回树形结构（true/false） |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": [
        {
            "owner": "built-in",
            "name": "admin-group",
            "createdTime": "2024-01-01T00:00:00+08:00",
            "displayName": "Admin Group",
            "parentId": "",
            "users": []
        }
    ]
}
```

#### 2.2.7 创建组

**接口地址**：`POST /api/add-group`

**请求体**：

```json
{
    "owner": "built-in",
    "name": "new-group",
    "displayName": "New Group",
    "parentId": "",
    "users": []
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

### 2.3 应用授权 API

#### 2.3.1 获取应用列表

**接口地址**：`GET /api/get-applications`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 否 | 所有者 |
| organization | string | 否 | 组织名称 |
| pageSize | int | 否 | 每页数量 |
| p | int | 否 | 页码 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": [
        {
            "owner": "admin",
            "name": "app-built-in",
            "createdTime": "2024-01-01T00:00:00+08:00",
            "displayName": "Built-in App",
            "organization": "built-in",
            "clientId": "client-id",
            "redirectUris": ["http://localhost:8080/callback"],
            "tokenFormat": "JWT"
        }
    ]
}
```

#### 2.3.2 获取单个应用

**接口地址**：`GET /api/get-application`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 应用 ID |
| withKey | string | 否 | 是否包含证书公钥 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": {
        "owner": "admin",
        "name": "app-built-in",
        "createdTime": "2024-01-01T00:00:00+08:00",
        "displayName": "Built-in App",
        "organization": "built-in",
        "clientId": "client-id",
        "clientSecret": "***",
        "redirectUris": ["http://localhost:8080/callback"],
        "tokenFormat": "JWT"
    }
}
```

#### 2.3.3 创建应用

**接口地址**：`POST /api/add-application`

**请求体**：

```json
{
    "owner": "admin",
    "name": "new-app",
    "displayName": "New Application",
    "organization": "built-in",
    "redirectUris": ["http://localhost:8080/callback"],
    "tokenFormat": "JWT"
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.3.4 更新应用

**接口地址**：`POST /api/update-application`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 应用 ID |

**请求体**：

```json
{
    "displayName": "Updated App Name",
    "redirectUris": ["http://localhost:8080/new-callback"]
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

#### 2.3.5 删除应用

**接口地址**：`POST /api/delete-application`

**请求体**：

```json
{
    "owner": "admin",
    "name": "app-to-delete"
}
```

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "Affected"
}
```

### 2.4 MFA 多因素认证 API

#### 2.4.1 初始化 MFA 设置

**接口地址**：`POST /api/mfa/setup/initiate`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 用户所属组织 |
| name | string | 是 | 用户名 |
| mfaType | string | 是 | MFA 类型（totp/sms/email/radius/push） |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": {
        "mfaType": "totp",
        "secret": "JBSWY3DPEHPK3PXP",
        "url": "otpauth://totp/Casdoor:user?secret=JBSWY3DPEHPK3PXP",
        "recoveryCodes": ["code1", "code2"],
        "mfaRememberInHours": 12
    }
}
```

#### 2.4.2 验证 MFA 配置

**接口地址**：`POST /api/mfa/setup/verify`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mfaType | string | 是 | MFA 类型 |
| passcode | string | 是 | 验证码 |
| secret | string | 否 | 密钥（TOTP 必填） |
| dest | string | 否 | 目标地址（SMS/Email 必填） |
| countryCode | string | 否 | 国家代码（SMS 必填） |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "OK"
}
```

#### 2.4.3 启用 MFA

**接口地址**：`POST /api/mfa/setup/enable`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 用户所属组织 |
| name | string | 是 | 用户名 |
| mfaType | string | 是 | MFA 类型 |
| secret | string | 否 | 密钥 |
| dest | string | 否 | 目标地址 |
| countryCode | string | 否 | 国家代码 |
| recoveryCodes | string | 是 | 恢复代码 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": "OK"
}
```

#### 2.4.4 删除 MFA

**接口地址**：`POST /api/mfa/delete`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 用户所属组织 |
| name | string | 是 | 用户名 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": [
        {
            "mfaType": "totp",
            "enabled": false
        }
    ]
}
```

#### 2.4.5 设置首选 MFA

**接口地址**：`POST /api/mfa/set-preferred`

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| owner | string | 是 | 用户所属组织 |
| name | string | 是 | 用户名 |
| mfaType | string | 是 | MFA 类型 |

**响应示例**：

```json
{
    "status": "ok",
    "msg": "",
    "data": [
        {
            "mfaType": "totp",
            "enabled": true,
            "isPreferred": true
        }
    ]
}
```

## 3. 错误响应

所有错误响应遵循统一格式：

```json
{
    "status": "error",
    "msg": "错误描述信息",
    "data": null
}
```

### 常见错误码

| 错误信息 | 说明 |
|----------|------|
| Missing parameter | 缺少必要参数 |
| Unauthorized operation | 无权限操作 |
| The user doesn't exist | 用户不存在 |
| The organization doesn't exist | 组织不存在 |
| The application doesn't exist | 应用不存在 |
| Invalid MFA type | 无效的 MFA 类型 |

## 4. 变更日志

### v2.0.0 (2024-01-01)

#### 新增功能
- 引入分层架构（DTO/Repository/Service/API）
- 添加内存缓存支持
- 优化数据库查询性能
- 支持批量导入导出优化

#### 接口变更
- 所有接口响应格式保持兼容
- 新增 MFA 相关 API
- 优化分页查询性能

#### 性能优化
- 用户查询添加缓存
- 批量操作使用事务
- 数据库索引优化
