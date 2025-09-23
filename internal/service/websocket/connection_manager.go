package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager WebSocket连接管理器
type ConnectionManager struct {
	// 每张图片的编辑状态，key: pictureId, value: 当前正在编辑的用户 ID
	pictureEditingUsers map[int64]int64
	// 保存所有连接的会话，key: pictureId, value: 用户会话集合
	pictureSessions map[int64]map[*websocket.Conn]*UserSession
	// 读写锁
	mutex sync.RWMutex
}

// UserSession 用户会话信息
type UserSession struct {
	UserID    int64           `json:"userId"`
	UserName  string          `json:"userName"`
	PictureID int64           `json:"pictureId"`
	Conn      *websocket.Conn `json:"-"`
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		pictureEditingUsers: make(map[int64]int64),
		pictureSessions:     make(map[int64]map[*websocket.Conn]*UserSession),
	}
}

// AddConnection 添加连接
func (cm *ConnectionManager) AddConnection(pictureID int64, userSession *UserSession) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.pictureSessions[pictureID] == nil {
		cm.pictureSessions[pictureID] = make(map[*websocket.Conn]*UserSession)
	}
	cm.pictureSessions[pictureID][userSession.Conn] = userSession
}

// RemoveConnection 移除连接
func (cm *ConnectionManager) RemoveConnection(pictureID int64, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if sessions, exists := cm.pictureSessions[pictureID]; exists {
		// 检查是否是当前编辑者
		if userSession, sessionExists := sessions[conn]; sessionExists {
			if cm.pictureEditingUsers[pictureID] == userSession.UserID {
				delete(cm.pictureEditingUsers, pictureID)
			}
		}
		delete(sessions, conn)

		// 如果没有会话了，删除整个图片的会话集合
		if len(sessions) == 0 {
			delete(cm.pictureSessions, pictureID)
		}
	}
}

// SetEditingUser 设置正在编辑的用户
func (cm *ConnectionManager) SetEditingUser(pictureID int64, userID int64) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 检查是否已有用户在编辑
	if _, exists := cm.pictureEditingUsers[pictureID]; exists {
		return false
	}
	cm.pictureEditingUsers[pictureID] = userID
	return true
}

// RemoveEditingUser 移除正在编辑的用户
func (cm *ConnectionManager) RemoveEditingUser(pictureID int64, userID int64) bool {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if editingUserID, exists := cm.pictureEditingUsers[pictureID]; exists && editingUserID == userID {
		delete(cm.pictureEditingUsers, pictureID)
		return true
	}
	return false
}

// GetEditingUser 获取正在编辑的用户ID
func (cm *ConnectionManager) GetEditingUser(pictureID int64) (int64, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	userID, exists := cm.pictureEditingUsers[pictureID]
	return userID, exists
}

// IsUserEditing 检查用户是否正在编辑
func (cm *ConnectionManager) IsUserEditing(pictureID int64, userID int64) bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	editingUserID, exists := cm.pictureEditingUsers[pictureID]
	return exists && editingUserID == userID
}

// GetPictureSessions 获取图片的所有会话
func (cm *ConnectionManager) GetPictureSessions(pictureID int64) map[*websocket.Conn]*UserSession {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if sessions, exists := cm.pictureSessions[pictureID]; exists {
		// 返回副本以避免并发问题
		result := make(map[*websocket.Conn]*UserSession)
		for conn, session := range sessions {
			result[conn] = session
		}
		return result
	}
	return nil
}

// BroadcastToPicture 向图片的所有用户广播消息
func (cm *ConnectionManager) BroadcastToPicture(pictureID int64, message []byte, excludeConn *websocket.Conn) {
	sessions := cm.GetPictureSessions(pictureID)
	if sessions == nil {
		return
	}

	for conn, _ := range sessions {
		// 排除指定的连接
		if excludeConn != nil && conn == excludeConn {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// 发送失败，移除连接
			cm.RemoveConnection(pictureID, conn)
			conn.Close()
		}
	}
}
