/**
 * @Author: wenliangzhang
 * @Description:
 * @File: push
 * @Version: 1.0.0
 * @Date: 2022/10/8 10:54 AM
 */
package push

import (
	"github.com/ylywyn/jpush-api-go-client"

	"git.singularity-ai.com/backend/library/arch/web"
)

const (
	appKey = "56740c7590d2900b76a86913"
	secret = "3f85cae4e5d51ebf1353a7ca"
)

type IosAlert struct {
	Body  string `json:"body"`
	Title string `json:"title"`
}

func PushMsg(ctx *web.WebContext, userID []string, title, content string) error {
	// Platform
	var pf jpushclient.Platform
	pf.All()
	// Audience
	var ad jpushclient.Audience
	ad.SetAlias(userID)
	// Notice
	var notice Notice
	notice.SetAlert(content)
	notice.SetAndroidNotice(&AndroidNotice{Alert: content, Title: title, Intent: AndroidIntent{"intent:#Intent;action=android.intent.action.MAIN;end"}})
	notice.SetIOSNotice(&IOSNotice{Alert: IosAlert{content, title}})
	notice.SetWinPhoneNotice(&WinPhoneNotice{Alert: content, Title: title})
	// payload
	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.Notification = notice
	bytes, _ := payload.ToBytes()
	// push
	c := jpushclient.NewPushClient(secret, appKey)
	req, err := c.Send(bytes)
	if err != nil {
		ctx.Warnf("push err = %s", err.Error())
		return err
	}
	ctx.Infof("push req = %s", req)
	return nil
}

func PushMsgForAll(ctx *web.WebContext, title, content string) error {
	// Platform
	var pf jpushclient.Platform
	pf.All()
	// Audience
	var ad jpushclient.Audience
	ad.All()
	// Notice
	var notice Notice
	notice.SetAlert(content)
	notice.SetAndroidNotice(&AndroidNotice{Alert: content, Title: title, Intent: AndroidIntent{"intent:#Intent;action=android.intent.action.MAIN;end"}})
	notice.SetIOSNotice(&IOSNotice{Alert: IosAlert{content, title}})
	notice.SetWinPhoneNotice(&WinPhoneNotice{Alert: content, Title: title})
	// payload
	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.Notification = notice
	bytes, _ := payload.ToBytes()
	// push
	c := jpushclient.NewPushClient(secret, appKey)
	req, err := c.Send(bytes)
	if err != nil {
		ctx.Warnf("push err = %s", err.Error())
		return err
	}
	ctx.Infof("push req = %s", req)
	return nil
}
