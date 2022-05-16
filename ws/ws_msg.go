/**
 * @Author: wenliangzhang
 * @Description:
 * @File: ws
 * @Version: 1.0.0
 * @Date: 2022/5/12 6:15 PM
 */
package ws

const (
	SystemDefault   = 0
	TypeOnlineHeart = 1 //app->server, 客户端心跳1(min)

	SystemSignaling   = 1000
	TypeApp2ServerRtc = 1001 //app->server,客户端上传的RTC信令消息
	TypeServer2AppRtc = 1002 //server->app,RTC调度服务回传给app的信令消息

	SystemMood           = 2000
	TypeUpdateMoodStatus = 2001 //server->app,更新心情状态
	TypeUpdateMoodConfig = 2002 //server->app,心情配置更新

	SystemVideo = 3000
)

type Msg2Server struct {
	CurrentVersion     string `json:"current_version"`
	Platform           string `json:"platform"`            // 端 android,ios,web
	System             int    `json:"system"`              //目标系统，例如10086聊天,10087商品,10088心情
	Type               int    `json:"type"`                //消息类型
	Timestamp          int64  `json:"timestamp"`           //时间戳,用于有序消息传递
	ReferenceSignature string `json:"reference_signature"` //md5串
	Data               struct {
		Content interface{} `json:"content"`
	} `json:"data"`
}

type ServerMsg struct {
	CurrentVersion     string `json:"current_version"`
	UserID             string `json:"user_id"`
	Platform           string `json:"platform"`            // 端 android,ios,web
	System             int    `json:"system"`              //目标系统，例如10086聊天,10087商品,10088心情
	Type               int    `json:"type"`                //消息类型
	Timestamp          int64  `json:"timestamp"`           //时间戳,用于有序消息传递
	ReferenceSignature string `json:"reference_signature"` //md5串
	Data               struct {
		Content interface{} `json:"content"`
	} `json:"data"`
}
