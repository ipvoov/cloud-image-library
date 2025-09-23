package controller

import (
	"cloud/internal/logic/picture"
	"cloud/internal/logic/space"
	"cloud/internal/model/entity"
	wsModel "cloud/internal/model/websocket"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

var (
	WebSocket = cWebSocket{}
	// 存储连接信息的全局map
	pictureConnections  = make(map[int64]map[*ghttp.WebSocket]*WebSocketSession)
	pictureEditingUsers = make(map[int64]int64) // pictureId -> userId
)

type cWebSocket struct{}

type WebSocketSession struct {
	UserID    int64
	UserName  string
	PictureID int64
	Conn      *ghttp.WebSocket
}

// PictureEdit 图片编辑 WebSocket 处理
func (c *cWebSocket) PictureEdit(r *ghttp.Request) {
	ctx := r.GetCtx()

	// 参数验证
	pictureIDStr := r.Get("pictureId").String()
	if pictureIDStr == "" {
		r.Response.WriteJsonExit(g.Map{
			"code":    40000,
			"message": "缺少图片参数",
		})
		return
	}

	pictureID, err := strconv.ParseInt(pictureIDStr, 10, 64)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    40000,
			"message": "图片ID格式错误",
		})
		return
	}

	// 用户身份验证
	loginUser := c.getLoginUserFromRequest(ctx, r)
	if loginUser == nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    40100,
			"message": "用户未登录",
		})
		return
	}

	// 权限验证
	if !c.checkPictureEditPermission(ctx, pictureID, loginUser) {
		r.Response.WriteJsonExit(g.Map{
			"code":    40300,
			"message": "没有编辑权限",
		})
		return
	}

	// 升级为 WebSocket 连接
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(ctx, "WebSocket升级失败:", err)
		return
	}
	defer ws.Close()

	glog.Info(ctx, "WebSocket连接建立成功，用户:", loginUser.UserName, "图片ID:", pictureID)

	// 创建用户会话
	session := &WebSocketSession{
		UserID:    loginUser.Id,
		UserName:  loginUser.UserName,
		PictureID: pictureID,
		Conn:      ws,
	}

	// 添加连接到管理器
	c.addConnection(pictureID, session)
	defer func() {
		glog.Info(ctx, "WebSocket连接关闭，清理资源，用户:", loginUser.UserName, "图片ID:", pictureID)
		c.removeConnection(pictureID, session)
	}()

	// 发送用户加入编辑的通知
	c.broadcastMessage(ctx, pictureID, wsModel.PictureEditResponseMessage{
		Type:    wsModel.MessageTypeInfo,
		Message: fmt.Sprintf("用户 %s 加入编辑", loginUser.UserName),
		User:    loginUser,
	}, session)

	glog.Info(ctx, "开始监听WebSocket消息，用户:", loginUser.UserName)

	// 消息处理循环
	for {
		var requestMsg wsModel.PictureEditRequestMessage
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			glog.Error(ctx, "读取WebSocket消息失败:", err)
			break
		}

		if err := json.Unmarshal(msgBytes, &requestMsg); err != nil {
			glog.Error(ctx, "解析WebSocket消息失败:", err)
			continue
		}

		// 处理消息
		c.handleMessage(ctx, pictureID, session, requestMsg)
	}

	// 连接关闭时的清理工作
	c.handleUserExit(ctx, pictureID, session)
}

// checkPictureEditPermission 检查图片编辑权限
func (c *cWebSocket) checkPictureEditPermission(ctx context.Context, pictureID int64, loginUser *entity.User) bool {
	// 获取图片信息
	pic, err := c.getPictureById(ctx, pictureID)
	if err != nil || pic == nil {
		g.Log().Error(ctx, "图片不存在:", pictureID)
		return false
	}

	// 如果是私有图片，检查是否是图片所有者
	if pic.SpaceId == 0 {
		return pic.UserId == loginUser.Id
	}

	// 如果是团队空间图片，检查空间权限
	spaceInfo, err := c.getSpaceById(ctx, pic.SpaceId)
	if err != nil || spaceInfo == nil {
		g.Log().Error(ctx, "空间不存在:", pic.SpaceId)
		return false
	}

	// 检查用户在空间中是否有编辑权限
	// 这里需要根据你的空间权限逻辑来实现
	// 简化实现：检查用户是否是空间成员
	return c.hasEditPermission(ctx, pic.SpaceId, loginUser.Id)
}

// addConnection 添加连接
func (c *cWebSocket) addConnection(pictureID int64, session *WebSocketSession) {
	if pictureConnections[pictureID] == nil {
		pictureConnections[pictureID] = make(map[*ghttp.WebSocket]*WebSocketSession)
	}
	pictureConnections[pictureID][session.Conn] = session
}

// removeConnection 移除连接
func (c *cWebSocket) removeConnection(pictureID int64, session *WebSocketSession) {
	if sessions, exists := pictureConnections[pictureID]; exists {
		// 如果是当前编辑者，清除编辑状态
		if pictureEditingUsers[pictureID] == session.UserID {
			delete(pictureEditingUsers, pictureID)
		}
		delete(sessions, session.Conn)

		// 如果没有连接了，删除整个图片的连接映射
		if len(sessions) == 0 {
			delete(pictureConnections, pictureID)
		}
	}
}

// handleMessage 处理WebSocket消息
func (c *cWebSocket) handleMessage(ctx context.Context, pictureID int64, session *WebSocketSession, requestMsg wsModel.PictureEditRequestMessage) {
	switch requestMsg.Type {
	case wsModel.MessageTypeEnterEdit:
		c.handleEnterEdit(ctx, pictureID, session)
	case wsModel.MessageTypeExitEdit:
		c.handleExitEdit(ctx, pictureID, session)
	case wsModel.MessageTypeEditAction:
		c.handleEditAction(ctx, pictureID, session, requestMsg)
	default:
		c.sendErrorMessage(session.Conn, "未知的消息类型")
	}
}

// handleEnterEdit 处理进入编辑状态
func (c *cWebSocket) handleEnterEdit(ctx context.Context, pictureID int64, session *WebSocketSession) {
	// 检查是否已有用户在编辑
	if _, exists := pictureEditingUsers[pictureID]; !exists {
		// 设置当前用户为编辑者
		pictureEditingUsers[pictureID] = session.UserID

		// 广播进入编辑状态的消息
		c.broadcastMessage(ctx, pictureID, wsModel.PictureEditResponseMessage{
			Type:    wsModel.MessageTypeEnterEdit,
			Message: fmt.Sprintf("用户 %s 开始编辑图片", session.UserName),
			User:    c.getUserEntity(session),
		}, nil)
	} else {
		// 已有用户在编辑，发送错误消息
		c.sendErrorMessage(session.Conn, "已有用户正在编辑该图片")
	}
}

// handleExitEdit 处理退出编辑状态
func (c *cWebSocket) handleExitEdit(ctx context.Context, pictureID int64, session *WebSocketSession) {
	// 检查是否是当前编辑者
	if editingUserID, exists := pictureEditingUsers[pictureID]; exists && editingUserID == session.UserID {
		// 移除编辑状态
		delete(pictureEditingUsers, pictureID)

		// 广播退出编辑状态的消息
		c.broadcastMessage(ctx, pictureID, wsModel.PictureEditResponseMessage{
			Type:    wsModel.MessageTypeExitEdit,
			Message: fmt.Sprintf("用户 %s 退出编辑图片", session.UserName),
			User:    c.getUserEntity(session),
		}, nil)
	}
}

// handleEditAction 处理编辑操作
func (c *cWebSocket) handleEditAction(ctx context.Context, pictureID int64, session *WebSocketSession, requestMsg wsModel.PictureEditRequestMessage) {
	// 检查是否是当前编辑者
	editingUserID, exists := pictureEditingUsers[pictureID]
	if !exists || editingUserID != session.UserID {
		c.sendErrorMessage(session.Conn, "您不是当前编辑者")
		return
	}

	// 广播编辑操作给其他用户（排除发送者）
	c.broadcastMessage(ctx, pictureID, wsModel.PictureEditResponseMessage{
		Type:       wsModel.MessageTypeEditAction,
		Message:    fmt.Sprintf("%s 执行 %s", session.UserName, requestMsg.EditAction.GetActionText()),
		EditAction: requestMsg.EditAction,
		User:       c.getUserEntity(session),
	}, session)
}

// handleUserExit 处理用户退出
func (c *cWebSocket) handleUserExit(ctx context.Context, pictureID int64, session *WebSocketSession) {
	// 如果是编辑者，自动退出编辑状态
	if editingUserID, exists := pictureEditingUsers[pictureID]; exists && editingUserID == session.UserID {
		delete(pictureEditingUsers, pictureID)
	}

	// 通知其他用户该用户已离开
	c.broadcastMessage(ctx, pictureID, wsModel.PictureEditResponseMessage{
		Type:    wsModel.MessageTypeInfo,
		Message: fmt.Sprintf("用户 %s 离开编辑", session.UserName),
		User:    c.getUserEntity(session),
	}, session)
}

// broadcastMessage 广播消息
func (c *cWebSocket) broadcastMessage(ctx context.Context, pictureID int64, message wsModel.PictureEditResponseMessage, excludeSession *WebSocketSession) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		glog.Error(ctx, "序列化消息失败:", err)
		return
	}

	if sessions, exists := pictureConnections[pictureID]; exists {
		for conn, session := range sessions {
			// 排除指定的连接
			if excludeSession != nil && conn == excludeSession.Conn {
				continue
			}

			if err := conn.WriteMessage(ghttp.WsMsgText, messageBytes); err != nil {
				glog.Error(ctx, "发送WebSocket消息失败:", err)
				// 发送失败，移除连接
				c.removeConnection(pictureID, session)
			}
		}
	}
}

// sendErrorMessage 发送错误消息
func (c *cWebSocket) sendErrorMessage(conn *ghttp.WebSocket, errorMsg string) {
	errorResponse := wsModel.PictureEditResponseMessage{
		Type:    wsModel.MessageTypeError,
		Message: errorMsg,
	}
	responseBytes, _ := json.Marshal(errorResponse)
	conn.WriteMessage(ghttp.WsMsgText, responseBytes)
}

// getUserEntity 从用户会话创建用户实体
func (c *cWebSocket) getUserEntity(session *WebSocketSession) *entity.User {
	return &entity.User{
		Id:       session.UserID,
		UserName: session.UserName,
	}
}

// getLoginUserFromRequest 从请求中获取登录用户
func (c *cWebSocket) getLoginUserFromRequest(ctx context.Context, r *ghttp.Request) *entity.User {
	// 从session中获取用户信息
	userObj, _ := r.Session.Get("login") // 使用正确的session key
	if userObj == nil {
		return nil
	}

	var user *entity.User
	if err := userObj.Struct(&user); err != nil {
		g.Log().Error(ctx, "解析用户信息失败:", err)
		return nil
	}

	return user
}

// getPictureById 根据ID获取图片
func (c *cWebSocket) getPictureById(ctx context.Context, pictureID int64) (*entity.Picture, error) {
	return picture.GetPictureById(ctx, pictureID)
}

// getSpaceById 根据ID获取空间
func (c *cWebSocket) getSpaceById(ctx context.Context, spaceID int64) (*entity.Space, error) {
	return space.GetSpaceById(ctx, spaceID)
}

// hasEditPermission 检查用户是否有编辑权限
func (c *cWebSocket) hasEditPermission(ctx context.Context, spaceID int64, userID int64) bool {
	return space.CheckEditPermission(ctx, spaceID, userID)
}

// Test WebSocket 测试方法
func (c *cWebSocket) Test(r *ghttp.Request) {
	ctx := r.GetCtx()

	// 简单的WebSocket测试
	ws, err := r.WebSocket()
	if err != nil {
		glog.Error(ctx, "WebSocket升级失败:", err)
		r.Response.WriteJsonExit(g.Map{
			"code":    500,
			"message": "WebSocket升级失败",
		})
		return
	}
	defer ws.Close()

	glog.Info(ctx, "WebSocket测试连接建立成功")

	// 发送欢迎消息
	welcomeMsg := map[string]interface{}{
		"type":      "welcome",
		"message":   "WebSocket连接测试成功",
		"timestamp": time.Now().Unix(),
	}

	msgBytes, _ := json.Marshal(welcomeMsg)
	if err := ws.WriteMessage(ghttp.WsMsgText, msgBytes); err != nil {
		glog.Error(ctx, "发送欢迎消息失败:", err)
		return
	}

	// 简单的echo服务
	for {
		_, msgBytes, err := ws.ReadMessage()
		if err != nil {
			glog.Info(ctx, "WebSocket测试连接关闭:", err)
			break
		}

		// echo回消息
		response := map[string]interface{}{
			"type":      "echo",
			"original":  string(msgBytes),
			"timestamp": time.Now().Unix(),
		}

		responseBytes, _ := json.Marshal(response)
		if err := ws.WriteMessage(ghttp.WsMsgText, responseBytes); err != nil {
			glog.Error(ctx, "echo消息失败:", err)
			break
		}
	}

	glog.Info(ctx, "WebSocket测试连接结束")
}
