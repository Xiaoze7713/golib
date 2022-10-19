/**
 * @Author: wenliangzhang
 * @Description:
 * @File: event
 * @Version: 1.0.0
 * @Date: 2022/5/10 3:06 PM
 * https://wiki.singularity-ai.com/pages/viewpage.action?pageId=192413717
 */
package event

import (
	"time"

	jsoniter "github.com/json-iterator/go"

	"git.singularity-ai.com/backend/library/common"
)

const (
	SystemDefault = 0
	TypeCreateWs  = 1 // 长链建立事件
	TypeLoginIn   = 2 // 用户登录事件
	TypeCloseWs   = 3 // 长链断开事件

	SystemVideo      = 1000
	TypeVideoStart   = 1001 // 用户建立RTC连接,开始视频
	TypeVideoHeart   = 1002 // 视频心跳事件
	TypeVideoEnd     = 1003 // 视频结束事件
	TypeSpeakStart   = 1011 // 用户开始说话
	TypeSpeakEnd     = 1012 // 用户结束说话
	TypeAISpeakStart = 1015 // AI开始说话
	TypeAISpeakEnd   = 1016 // AI结束说话

	SystemTask = 2000

	SystemAppTrans = 3000 // app透传

	SystemAppFeatures      = 4000 // app功能api
	TypeSendSoundMessage   = 4001 // 用户发送音视频消息
	TypeSendImMessage      = 4002 // 用户发送IM消息
	TypeMarkMessage        = 4003 // 用户标记消息(点赞、点踩等)
	TypeReplySoundMessage  = 4101 // 策略回复音视频消息
	TypeReplyImMessage     = 4102 // 策略回复IM消息
	TypeUpdateFeature      = 4300 // 更新特征
	TypeUpdateUserFeature  = 4301 // 更新用户特效
	TypeUpdateRobotFeature = 4302 // 更新robot特效

	SystemDelayPolicy = 6000 // 延迟策略触发
	TypeOffline3D     = 6001 // 用户3天未登录或回复消息
	TypeOffline7D     = 6002 // 用户7天未登录或回复消息
	TypeOffline30D    = 6003 // 用户30天未登录或回复消息
	TypeBirthDay      = 6004 // 用户当天生日
)

const EventKafkaTopic = "event"

type EventMsg struct {
	ID        int64       `json:"id"`
	System    int         `json:"system"`
	Type      int         `json:"type"`
	Data      interface{} `json:"data"`
	Version   string      `json:"version"`
	Timestamp int64       `json:"timestamp"`
	TraceID   string      `json:"trace_id"`
}

// DefaultEventData 默认事件data结构
type DefaultEventData struct {
	UserID string `json:"user_id"`
}

type Event struct {
	SendMQHandle func(msgByte []byte) error
}

func (e *Event) SendEventTrace(data interface{}, systemID, typeID int, traceID string) error {

	eventMsg := EventMsg{
		ID:        int64(common.GetSFInstance().GetUniqueId()),
		System:    systemID,
		Type:      typeID,
		Data:      data,
		Version:   "",
		Timestamp: time.Now().UnixNano() / 1e6,
		TraceID:   traceID,
	}
	msg, _ := jsoniter.Marshal(eventMsg)
	return e.SendMQHandle(msg)
}

func (e *Event) SendEvent(data interface{}, systemID, typeID int) error {

	eventMsg := EventMsg{
		ID:        int64(common.GetSFInstance().GetUniqueId()),
		System:    systemID,
		Type:      typeID,
		Data:      data,
		Version:   "",
		Timestamp: time.Now().UnixNano() / 1e6,
	}
	msg, _ := jsoniter.Marshal(eventMsg)
	return e.SendMQHandle(msg)
}

func (e *Event) SendEventWithVer(data interface{}, systemID, typeID int, version string) error {
	eventMsg := EventMsg{
		ID:        int64(common.GetSFInstance().GetUniqueId()),
		System:    systemID,
		Type:      typeID,
		Data:      data,
		Version:   version,
		Timestamp: time.Now().UnixNano() / 1e6,
	}
	msg, _ := jsoniter.Marshal(eventMsg)
	return e.SendMQHandle(msg)
}
