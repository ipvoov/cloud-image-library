package websocket

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/model/entity"
	wsmodel "cloud/internal/model/websocket"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gorilla/websocket"
)

// PictureEdit 添加用户
func (s *sWebSocket) PictureEdit(ctx context.Context, r *ghttp.Request, req *v1.WebSocketPictureEditReq) {

	// 用户身份验证
	loginUser := s.getLoginUserFromRequest(ctx, r)
	if loginUser == nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    40100,
			"message": "用户未登录",
		})
		return
	}

	// 权限验证
	if !s.checkPictureEditPermission(ctx, req.PictureId, loginUser) {
		r.Response.WriteJsonExit(g.Map{
			"code":    40300,
			"message": "没有编辑权限",
		})
		return
	}

	ws, err := wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    40500,
			"message": "WebSocket升级失败",
			"data":    err.Error(),
		})
		return
	}
	defer ws.Close()
	// 创建会话
	session := &WebSocketSession{
		UserID:    loginUser.Id,
		UserName:  loginUser.UserName,
		PictureID: req.PictureId,
	}
	// 添加连接
	s.addConnection(req.PictureId, ws, session)
	defer func() {
		// 确保资源清理
		s.removeConnection(req.PictureId, ws, session)
		s.handleUserExit(ctx, req.PictureId, ws, session)
	}()

	// 广播进入编辑消息
	s.broadcastMessage(ctx, req.PictureId, wsmodel.PictureEditResponseMessage{
		Type:    wsmodel.MessageTypeInfo,
		Message: fmt.Sprintf("用户 %s 加入编辑", loginUser.UserName),
		User:    loginUser,
	}, session)

	var requestMsg wsmodel.PictureEditRequestMessage
	var msgBytes []byte
	// 消息处理循环
	for {
		_, msgBytes, err = ws.ReadMessage()
		if err != nil {
			g.Log().Error(ctx, "读取WebSocket消息失败:", err)
			break
		}

		if err = gjson.DecodeTo(msgBytes, &requestMsg); err != nil {
			g.Log().Error(ctx, "解析WebSocket消息失败:", err)
			continue
		}

		// 处理消息
		s.handleMessage(ctx, req.PictureId, ws, session, requestMsg)
	}
	// handleUserExit 在 defer 中已经处理
}

// getLoginUserFromRequest 从请求中获取登录用户
func (s *sWebSocket) getLoginUserFromRequest(ctx context.Context, r *ghttp.Request) *entity.User {
	// 从session中获取用户信息
	userObj, _ := r.Session.Get(consts.LoginState) // 使用正确的session key
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

// checkPictureEditPermission 检查图片编辑权限
func (s *sWebSocket) checkPictureEditPermission(ctx context.Context, pictureID int64, loginUser *entity.User) bool {
	// 获取图片信息
	pic, err := GetPictureById(ctx, pictureID)
	if err != nil || pic == nil {
		g.Log().Error(ctx, "图片不存在:", pictureID)
		return false
	}

	// 如果是私有图片，检查是否是图片所有者
	if pic.SpaceId == 0 {
		return pic.UserId == loginUser.Id
	}

	// 如果是团队空间图片，检查空间权限
	spaceInfo, err := GetSpaceById(ctx, pic.SpaceId)
	if err != nil || spaceInfo == nil {
		g.Log().Error(ctx, "空间不存在:", pic.SpaceId)
		return false
	}

	// 检查用户在空间中是否有编辑权限
	return CheckEditPermission(ctx, pic.SpaceId, loginUser.Id)
}

// addConnection 添加连接
func (s *sWebSocket) addConnection(pictureID int64, conn *websocket.Conn, session *WebSocketSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pictureConnections[pictureID] == nil {
		s.pictureConnections[pictureID] = make(map[*websocket.Conn]*WebSocketSession)
	}
	s.pictureConnections[pictureID][conn] = session
}

// removeConnection 移除连接
func (s *sWebSocket) removeConnection(pictureID int64, conn *websocket.Conn, session *WebSocketSession) {
	//todo
	s.mu.Lock()
	defer s.mu.Unlock()
	if connMap, exists := s.pictureConnections[pictureID]; exists {
		// 如果是当前编辑者，清除编辑状态
		if s.pictureEditingUsers[pictureID] == session.UserID {
			delete(s.pictureEditingUsers, pictureID)
		}
		delete(connMap, conn)

		// 如果没有连接了，删除整个图片的连接映射
		if len(connMap) == 0 {
			delete(s.pictureConnections, pictureID)
		}
	}
}

// broadcastMessage 广播消息
func (s *sWebSocket) broadcastMessage(ctx context.Context, pictureID int64, message wsmodel.PictureEditResponseMessage, excludeSession *WebSocketSession) {
	messageBytes, err := gjson.Encode(message)
	if err != nil {
		g.Log().Error(ctx, "序列化消息失败:", err)
		return
	}

	// 获取读锁并复制连接映射，避免在网络IO时持有锁
	s.mu.RLock()
	connMap, exists := s.pictureConnections[pictureID]
	if !exists {
		s.mu.RUnlock()
		return
	}

	// 复制连接映射到临时变量，避免长时间持锁
	connectionsCopy := make(map[*websocket.Conn]*WebSocketSession)
	for conn, session := range connMap {
		connectionsCopy[conn] = session
	}
	s.mu.RUnlock()

	// 在锁外进行网络IO操作
	var failedConnections []struct {
		conn    *websocket.Conn
		session *WebSocketSession
	}

	for conn, session := range connectionsCopy {
		// 排除指定的连接
		if excludeSession != nil && session.UserID == excludeSession.UserID {
			continue
		}

		if err = conn.WriteMessage(ghttp.WsMsgText, messageBytes); err != nil {
			g.Log().Warning(ctx, "发送WebSocket消息失败:", err)
			// 记录失败的连接，稍后清理
			failedConnections = append(failedConnections, struct {
				conn    *websocket.Conn
				session *WebSocketSession
			}{conn, session})
		}
	}

	// 清理失败的连接
	for _, failed := range failedConnections {
		s.removeConnection(pictureID, failed.conn, failed.session)
	}
}

// handleMessage 处理WebSocket消息
func (s *sWebSocket) handleMessage(ctx context.Context, pictureID int64, conn *websocket.Conn, session *WebSocketSession, requestMsg wsmodel.PictureEditRequestMessage) {
	switch requestMsg.Type {
	case wsmodel.MessageTypeEnterEdit:
		s.handleEnterEdit(ctx, pictureID, conn, session)
	case wsmodel.MessageTypeExitEdit:
		s.handleExitEdit(ctx, pictureID, conn, session)
	case wsmodel.MessageTypeEditAction:
		s.handleEditAction(ctx, pictureID, conn, session, requestMsg)
	default:
		s.sendErrorMessage(conn, "未知的消息类型")
	}
}

// handleEnterEdit 处理进入编辑状态
func (s *sWebSocket) handleEnterEdit(ctx context.Context, pictureID int64, conn *websocket.Conn, session *WebSocketSession) {
	// 先尝试获取编辑权限
	s.mu.Lock()
	if _, exists := s.pictureEditingUsers[pictureID]; !exists {
		// 设置当前用户为编辑者
		s.pictureEditingUsers[pictureID] = session.UserID
		s.mu.Unlock()

		// 先给发送者发送成功消息
		successMsg := wsmodel.PictureEditResponseMessage{
			Type:    wsmodel.MessageTypeEnterEdit,
			Message: "您已成功进入编辑模式",
			User:    &entity.User{Id: session.UserID, UserName: session.UserName},
		}
		messageBytes, _ := gjson.Encode(successMsg)
		conn.WriteMessage(ghttp.WsMsgText, messageBytes)

		// 然后广播给其他用户
		s.broadcastMessage(ctx, pictureID, wsmodel.PictureEditResponseMessage{
			Type:    wsmodel.MessageTypeEnterEdit,
			Message: fmt.Sprintf("用户 %s 开始编辑图片", session.UserName),
			User:    &entity.User{Id: session.UserID, UserName: session.UserName},
		}, session)
	} else {
		s.mu.Unlock()
		// 已有用户在编辑，发送错误消息
		s.sendErrorMessage(conn, "已有用户正在编辑该图片")
	}
}

// handleExitEdit 处理退出编辑状态
func (s *sWebSocket) handleExitEdit(ctx context.Context, pictureID int64, conn *websocket.Conn, session *WebSocketSession) {
	// 检查并移除编辑状态
	s.mu.Lock()
	editingUserID, exists := s.pictureEditingUsers[pictureID]
	if exists && editingUserID == session.UserID {
		delete(s.pictureEditingUsers, pictureID)
		s.mu.Unlock()

		// 先给发送者发送成功消息
		successMsg := wsmodel.PictureEditResponseMessage{
			Type:    wsmodel.MessageTypeExitEdit,
			Message: "您已成功退出编辑模式",
			User:    &entity.User{Id: session.UserID, UserName: session.UserName},
		}
		messageBytes, _ := gjson.Encode(successMsg)
		conn.WriteMessage(ghttp.WsMsgText, messageBytes)

		// 然后广播给其他用户
		s.broadcastMessage(ctx, pictureID, wsmodel.PictureEditResponseMessage{
			Type:    wsmodel.MessageTypeExitEdit,
			Message: fmt.Sprintf("用户 %s 退出编辑图片", session.UserName),
			User:    &entity.User{Id: session.UserID, UserName: session.UserName},
		}, session)
	} else {
		s.mu.Unlock()
		// 发送错误消息：用户不是当前编辑者
		s.sendErrorMessage(conn, "您当前不是编辑者")
	}
}

// handleEditAction 处理编辑操作
func (s *sWebSocket) handleEditAction(ctx context.Context, pictureID int64, conn *websocket.Conn, session *WebSocketSession, requestMsg wsmodel.PictureEditRequestMessage) {
	// 检查是否是当前编辑者
	s.mu.RLock()
	editingUserID, exists := s.pictureEditingUsers[pictureID]
	s.mu.RUnlock()

	if !exists || editingUserID != session.UserID {
		s.sendErrorMessage(conn, "您不是当前编辑者")
		return
	}

	// 在锁外广播编辑操作给其他用户（排除发送者）
	s.broadcastMessage(ctx, pictureID, wsmodel.PictureEditResponseMessage{
		Type:       wsmodel.MessageTypeEditAction,
		Message:    fmt.Sprintf("%s 执行 %s", session.UserName, requestMsg.EditAction.GetActionText()),
		EditAction: requestMsg.EditAction,
		User:       &entity.User{Id: session.UserID, UserName: session.UserName},
	}, session)
}

// handleUserExit 处理用户退出
func (s *sWebSocket) handleUserExit(ctx context.Context, pictureID int64, conn *websocket.Conn, session *WebSocketSession) {
	// 如果是编辑者，自动退出编辑状态
	s.mu.Lock()
	if editingUserID, exists := s.pictureEditingUsers[pictureID]; exists && editingUserID == session.UserID {
		delete(s.pictureEditingUsers, pictureID)
	}
	s.mu.Unlock()

	// 在锁外通知其他用户该用户已离开
	s.broadcastMessage(ctx, pictureID, wsmodel.PictureEditResponseMessage{
		Type:    wsmodel.MessageTypeInfo,
		Message: fmt.Sprintf("用户 %s 离开编辑", session.UserName),
		User:    &entity.User{Id: session.UserID, UserName: session.UserName},
	}, session)
}

// sendErrorMessage 发送错误消息
func (s *sWebSocket) sendErrorMessage(conn *websocket.Conn, errorMsg string) {
	errorResponse := wsmodel.PictureEditResponseMessage{
		Type:    wsmodel.MessageTypeError,
		Message: errorMsg,
	}
	responseBytes, _ := json.Marshal(errorResponse)
	conn.WriteMessage(ghttp.WsMsgText, responseBytes)
}
