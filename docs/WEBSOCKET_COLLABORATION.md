# 图片协同编辑功能文档

## 功能概述

本系统实现了基于 WebSocket 的图片协同编辑功能，允许多个用户在团队空间中实时协同编辑图片。该功能参考了 Java 实现，使用 Go 语言重新开发。

## 架构设计

### 后端架构 (Go)

```
internal/
├── controller/
│   └── websocket.go           # WebSocket 控制器
├── service/websocket/
│   └── connection_manager.go  # 连接管理器
├── model/websocket/
│   └── picture_edit_message.go # 消息模型
└── logic/
    ├── picture/
    │   └── websocket_helper.go # 图片相关辅助方法
    └── space/
        └── websocket_helper.go # 空间相关辅助方法
```

### 前端架构 (Vue3 + TypeScript)

```
src/
├── utils/
│   └── pictureEditWebSocket.ts  # WebSocket 客户端
├── components/
│   └── ImageCropper.vue         # 图片编辑组件
└── constants/
    └── picture.ts               # 消息类型常量
```

## 核心功能

### 1. 连接管理
- **建立连接**: 用户进入图片编辑页面时自动建立 WebSocket 连接
- **权限验证**: 握手时验证用户身份和编辑权限
- **连接管理**: 自动管理连接的生命周期

### 2. 编辑状态管理
- **进入编辑**: 用户可以申请进入编辑状态
- **编辑锁定**: 同一时间只允许一个用户编辑
- **退出编辑**: 用户可以主动退出编辑状态
- **自动退出**: 连接断开时自动退出编辑状态

### 3. 实时同步
- **操作同步**: 编辑操作实时同步给其他用户
- **状态通知**: 用户进入/退出编辑状态的通知
- **消息广播**: 向同一图片的所有用户广播消息

## 消息类型

### 请求消息 (`PictureEditRequestMessage`)
```go
type PictureEditRequestMessage struct {
    Type       PictureEditMessageType `json:"type"`       // 消息类型
    EditAction PictureEditAction      `json:"editAction"` // 编辑动作
}
```

### 响应消息 (`PictureEditResponseMessage`)
```go
type PictureEditResponseMessage struct {
    Type       PictureEditMessageType `json:"type"`       // 消息类型
    Message    string                 `json:"message"`    // 消息内容
    EditAction PictureEditAction      `json:"editAction"` // 编辑动作
    User       *entity.User           `json:"user"`       // 用户信息
}
```

### 消息类型枚举
- `INFO`: 通知消息
- `ERROR`: 错误消息
- `ENTER_EDIT`: 进入编辑状态
- `EXIT_EDIT`: 退出编辑状态
- `EDIT_ACTION`: 编辑操作

### 编辑动作枚举
- `ZOOM_IN`: 放大操作
- `ZOOM_OUT`: 缩小操作
- `ROTATE_LEFT`: 左旋操作
- `ROTATE_RIGHT`: 右旋操作

## API 接口

### WebSocket 端点
```
GET /api/ws/picture/edit?pictureId={pictureId}
```

**参数说明**:
- `pictureId`: 图片ID

**认证要求**:
- 用户必须已登录
- 用户必须有该图片的编辑权限

## 使用流程

### 1. 前端连接流程
```typescript
// 1. 创建 WebSocket 实例
const websocket = new PictureEditWebSocket(pictureId)

// 2. 建立连接
websocket.connect()

// 3. 监听事件
websocket.on('INFO', (msg) => {
  console.log('收到通知:', msg.message)
})

// 4. 发送消息
websocket.sendMessage({
  type: 'ENTER_EDIT'
})
```

### 2. 协同编辑流程
1. **用户A** 进入图片编辑页面，WebSocket 连接建立
2. **用户B** 也进入同一图片编辑页面，收到用户A加入的通知
3. **用户A** 点击"进入编辑"，获得编辑权限
4. **用户A** 执行旋转操作，操作实时同步给用户B
5. **用户A** 退出编辑或断开连接，释放编辑权限

## 权限机制

### 1. 连接权限
- 用户必须已登录
- 对于私有图片: 必须是图片所有者
- 对于团队空间图片: 必须是空间成员且有编辑权限

### 2. 编辑权限
- 同一时间只允许一个用户编辑
- 编辑者断开连接时自动释放权限
- 非编辑者只能观看编辑过程

## 部署配置

### 1. 后端配置
确保 Go 应用包含以下依赖：
```go
github.com/gorilla/websocket v1.5.3
github.com/gogf/gf/v2 v2.9.3
```

### 2. 前端配置
WebSocket 连接地址配置：
```typescript
const DEV_BASE_URL = "ws://localhost:8123"
const url = `${DEV_BASE_URL}/api/ws/picture/edit?pictureId=${pictureId}`
```

## 错误处理

### 常见错误
1. **连接被拒绝**: 检查用户登录状态和权限
2. **消息发送失败**: 检查 WebSocket 连接状态
3. **权限不足**: 确认用户有编辑权限

### 错误代码
- `400`: 参数错误
- `401`: 用户未登录
- `403`: 权限不足
- `404`: 图片不存在

## 性能优化

### 1. 连接管理
- 使用连接池管理 WebSocket 连接
- 自动清理无效连接
- 限制单个图片的最大连接数

### 2. 消息处理
- 使用 JSON 格式减少传输大小
- 批量处理多个操作
- 防抖动避免频繁操作

## 安全考虑

### 1. 认证授权
- 严格的权限检查
- Session 验证
- CSRF 防护

### 2. 输入验证
- 参数类型检查
- 消息格式验证
- 操作权限验证

## 测试

### 运行测试
```bash
# 运行 WebSocket 相关测试
go test ./internal/controller -v

# 运行所有测试
go test ./... -v
```

### 测试用例
- 消息类型常量测试
- 编辑动作枚举测试
- 连接管理功能测试
- 权限验证测试

## 故障排除

### 1. 连接问题
- 检查网络连接
- 确认服务端口开放
- 验证 WebSocket 协议支持

### 2. 权限问题
- 检查用户登录状态
- 验证图片访问权限
- 确认空间成员身份

### 3. 同步问题
- 检查消息格式
- 验证事件监听器
- 确认连接状态

## 扩展功能

### 未来可扩展的功能
1. **历史记录**: 保存编辑历史
2. **撤销重做**: 支持操作撤销
3. **批量操作**: 支持批量编辑
4. **实时预览**: 更细粒度的实时同步
5. **协作指针**: 显示其他用户的操作位置

---

## 开发团队

本功能由原 Java 代码参考实现，使用 Go 语言重新开发，保持了原有的设计思想和功能特性。
