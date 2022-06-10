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
	"git.singularity-ai.com/backend/library/common"
	jsoniter "github.com/json-iterator/go"
	"time"
)

const (
	SystemDefault = 0
	TypeCreateWs  = 1 //长链建立事件
	TypeLoginIn   = 2 //用户登录事件
	TypeCloseWs   = 3 //长链断开事件

	SystemVideo    = 1000
	TypeVideoStart = 1001 //用户建立RTC连接,开始视频
	TypeVideoHeart = 1002 //视频心跳事件
	TypeVideoEnd   = 1003 //视频结束事件

	SystemTask = 2000

	SystemAppTrans = 3000 //app透传

	SystemAppFeatures = 4000 //app功能api
	TypeSendImMessage = 4002
	TypeMarkMessage   = 4003
)

const EventKafkaTopic = "event"

type EventMsg struct {
	ID        int64       `json:"id"`
	System    int         `json:"system"`
	Type      int         `json:"type"`
	Data      interface{} `json:"data"`
	Version   string      `json:"version"`
	Timestamp int64       `json:"timestamp"`
}

// DefaultEventData 默认事件data结构
type DefaultEventData struct {
	UserID string `json:"user_id"`
}

type Event struct {
	SendMQHandle func(msgByte []byte) error
}

func (e *Event) SendEvent(data interface{}, systemID, typeID int) error {

	eventMsg := EventMsg{
		int64(common.GetSFInstance().GetUniqueId()),
		systemID,
		typeID,
		data,
		"",
		time.Now().UnixNano() / 1e6,
	}
	msg, _ := jsoniter.Marshal(eventMsg)
	return e.SendMQHandle(msg)
}

func (e *Event) SendEventWithVer(data interface{}, systemID, typeID int, version string) error {
	eventMsg := EventMsg{
		int64(common.GetSFInstance().GetUniqueId()),
		systemID,
		typeID,
		data,
		version,
		time.Now().UnixNano() / 1e6,
	}
	msg, _ := jsoniter.Marshal(eventMsg)
	return e.SendMQHandle(msg)
}
