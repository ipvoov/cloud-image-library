package wsmodel

import "cloud/internal/model/entity"

// PictureEditMessageType 图片编辑消息类型枚举
type PictureEditMessageType string

const (
	MessageTypeInfo       PictureEditMessageType = "INFO"
	MessageTypeError      PictureEditMessageType = "ERROR"
	MessageTypeEnterEdit  PictureEditMessageType = "ENTER_EDIT"
	MessageTypeExitEdit   PictureEditMessageType = "EXIT_EDIT"
	MessageTypeEditAction PictureEditMessageType = "EDIT_ACTION"
)

// PictureEditAction 图片编辑动作枚举
type PictureEditAction string

const (
	ActionZoomIn      PictureEditAction = "ZOOM_IN"
	ActionZoomOut     PictureEditAction = "ZOOM_OUT"
	ActionRotateLeft  PictureEditAction = "ROTATE_LEFT"
	ActionRotateRight PictureEditAction = "ROTATE_RIGHT"
)

// PictureEditRequestMessage 图片编辑请求消息
type PictureEditRequestMessage struct {
	Type       PictureEditMessageType `json:"type"`       // 消息类型
	EditAction PictureEditAction      `json:"editAction"` // 执行的编辑动作
}

// PictureEditResponseMessage 图片编辑响应消息
type PictureEditResponseMessage struct {
	Type       PictureEditMessageType `json:"type"`       // 消息类型
	Message    string                 `json:"message"`    // 信息
	EditAction PictureEditAction      `json:"editAction"` // 执行的编辑动作
	User       *entity.User           `json:"user"`       // 用户信息
}

// GetActionText 获取动作文本描述
func (action PictureEditAction) GetActionText() string {
	switch action {
	case ActionZoomIn:
		return "放大操作"
	case ActionZoomOut:
		return "缩小操作"
	case ActionRotateLeft:
		return "左旋操作"
	case ActionRotateRight:
		return "右旋操作"
	default:
		return "未知操作"
	}
}
