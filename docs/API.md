# API 文档

## 概述

智能云图库提供RESTful API接口，支持图片管理、用户管理、空间管理等功能。

**基础URL**: `http://localhost:8123/api`

## 认证

大部分API需要用户认证，请在请求头中包含认证信息：

```
Authorization: Bearer <token>
```

## 用户管理 API

### 用户注册

**POST** `/user/register`

注册新用户账户。

**请求参数:**
```json
{
  "userAccount": "string",    // 用户账号
  "userPassword": "string",  // 用户密码
  "checkPassword": "string", // 确认密码
  "userName": "string"       // 用户昵称
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 1,
    "userAccount": "testuser",
    "userName": "测试用户"
  }
}
```

### 用户登录

**POST** `/user/login`

用户登录获取访问令牌。

**请求参数:**
```json
{
  "userAccount": "string",   // 用户账号
  "userPassword": "string"  // 用户密码
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "id": 1,
    "userAccount": "testuser",
    "userName": "测试用户",
    "userAvatar": "https://example.com/avatar.jpg",
    "userRole": "user"
  }
}
```

### 获取当前用户信息

**GET** `/user/get/login`

获取当前登录用户的详细信息。

**响应示例:**
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": 1,
    "userAccount": "testuser",
    "userName": "测试用户",
    "userAvatar": "https://example.com/avatar.jpg",
    "userRole": "user",
    "createTime": "2024-01-01T00:00:00Z"
  }
}
```

### 用户登出

**POST** `/user/logout`

用户登出，清除会话信息。

**响应示例:**
```json
{
  "code": 0,
  "message": "登出成功"
}
```

## 图片管理 API

### 上传图片

**POST** `/picture/add`

上传单张或多张图片。

**请求参数:**
```json
{
  "name": "string",           // 图片名称
  "introduction": "string",   // 图片简介
  "category": "string",       // 图片分类
  "tags": ["string"],         // 图片标签
  "spaceId": 1               // 所属空间ID
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "id": 1,
    "name": "示例图片",
    "url": "https://example.com/image.jpg",
    "picSize": 1024000,
    "picWidth": 1920,
    "picHeight": 1080,
    "picFormat": "jpg"
  }
}
```

### 分页查询图片

**POST** `/picture/list/page/vo`

分页查询图片列表，支持多种筛选条件。

**请求参数:**
```json
{
  "current": 1,              // 当前页码
  "pageSize": 10,           // 每页大小
  "category": "string",      // 分类筛选
  "tags": ["string"],        // 标签筛选
  "searchText": "string",    // 搜索关键词
  "spaceId": "string",       // 空间ID筛选
  "sortField": "string",     // 排序字段
  "sortOrder": "string"      // 排序方向
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "records": [
      {
        "id": 1,
        "name": "示例图片",
        "url": "https://example.com/image.jpg",
        "introduction": "这是一张示例图片",
        "category": "风景",
        "tags": ["自然", "风景"],
        "picSize": 1024000,
        "picWidth": 1920,
        "picHeight": 1080,
        "picFormat": "jpg",
        "createTime": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "current": 1,
    "pageSize": 10
  }
}
```

### 编辑图片信息

**POST** `/picture/edit`

编辑图片的基本信息。

**请求参数:**
```json
{
  "id": 1,                   // 图片ID
  "name": "string",          // 新名称
  "introduction": "string",  // 新简介
  "category": "string",      // 新分类
  "tags": ["string"],        // 新标签
  "spaceId": 1              // 新空间ID
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "编辑成功",
  "data": {
    "success": true
  }
}
```

### 删除图片

**POST** `/picture/delete`

删除指定的图片。

**请求参数:**
```json
{
  "id": 1  // 图片ID
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "success": true
  }
}
```

## 空间管理 API

### 创建空间

**POST** `/space/add`

创建新的图片空间。

**请求参数:**
```json
{
  "spaceName": "string",     // 空间名称
  "spaceLevel": 0,           // 空间级别：0-普通版 1-专业版 2-旗舰版
  "maxSize": 1000000000,     // 最大存储大小（字节）
  "maxCount": 1000,          // 最大图片数量
  "spaceType": 0             // 空间类型：0-私有 1-团队
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    "spaceName": "我的空间",
    "spaceLevel": 0,
    "maxSize": 1000000000,
    "maxCount": 1000,
    "spaceType": 0
  }
}
```

### 分页查询空间

**POST** `/space/list/page/vo`

分页查询用户的空间列表。

**请求参数:**
```json
{
  "current": 1,              // 当前页码
  "pageSize": 10,           // 每页大小
  "spaceName": "string",     // 空间名称筛选
  "spaceType": 0             // 空间类型筛选
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "records": [
      {
        "id": 1,
        "spaceName": "我的空间",
        "spaceLevel": 0,
        "maxSize": 1000000000,
        "maxCount": 1000,
        "totalSize": 500000000,
        "totalCount": 50,
        "spaceType": 0,
        "createTime": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5,
    "current": 1,
    "pageSize": 10
  }
}
```

### 空间分析

**GET** `/space/analyze`

获取空间使用情况分析数据。

**请求参数:**
- `spaceId`: 空间ID（查询参数）

**响应示例:**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "spaceUsage": {
      "totalSize": 500000000,
      "maxSize": 1000000000,
      "usagePercent": 50.0
    },
    "categoryStats": [
      {
        "category": "风景",
        "count": 20,
        "size": 200000000
      },
      {
        "category": "人物",
        "count": 15,
        "size": 150000000
      }
    ],
    "tagStats": [
      {
        "tag": "自然",
        "count": 25
      },
      {
        "tag": "城市",
        "count": 10
      }
    ]
  }
}
```

## 空间用户管理 API

### 添加空间成员

**POST** `/spaceUser/add`

向空间添加新成员。

**请求参数:**
```json
{
  "spaceId": 1,              // 空间ID
  "userId": 2,               // 用户ID
  "spaceRole": "editor"      // 角色：admin/editor/viewer
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "添加成功",
  "data": {
    "success": true
  }
}
```

### 查询空间成员

**POST** `/spaceUser/list/page/vo`

分页查询空间的成员列表。

**请求参数:**
```json
{
  "current": 1,              // 当前页码
  "pageSize": 10,           // 每页大小
  "spaceId": 1              // 空间ID
}
```

**响应示例:**
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "records": [
      {
        "id": 1,
        "spaceId": 1,
        "userId": 2,
        "spaceRole": "editor",
        "userName": "成员用户",
        "userAvatar": "https://example.com/avatar.jpg",
        "joinTime": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5,
    "current": 1,
    "pageSize": 10
  }
}
```

## WebSocket 协同编辑 API

### 连接WebSocket

**WebSocket** `/ws/picture/edit`

建立图片协同编辑的WebSocket连接。

**连接参数:**
- `pictureId`: 图片ID
- `token`: 用户认证令牌

**消息格式:**

**进入编辑状态:**
```json
{
  "type": "enter_edit",
  "pictureId": 1,
  "userId": 1,
  "userName": "用户名"
}
```

**退出编辑状态:**
```json
{
  "type": "exit_edit",
  "pictureId": 1,
  "userId": 1
}
```

**编辑操作同步:**
```json
{
  "type": "edit_operation",
  "pictureId": 1,
  "userId": 1,
  "operation": {
    "type": "crop",
    "data": {
      "x": 100,
      "y": 100,
      "width": 200,
      "height": 200
    }
  }
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40000 | 请求参数错误 |
| 40001 | 未登录 |
| 40003 | 无权限 |
| 50000 | 系统内部错误 |

## 响应格式

所有API响应都遵循统一格式：

```json
{
  "code": 0,           // 状态码
  "message": "string",  // 消息
  "data": {}          // 数据
}
```

## 分页格式

分页查询的响应格式：

```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "records": [],     // 数据列表
    "total": 100,      // 总记录数
    "current": 1,      // 当前页码
    "pageSize": 10    // 每页大小
  }
}
```
