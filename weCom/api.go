package weCom

import (
	"encoding/json"
	"errors"
	"github.com/ayu-666/gopkg/weCom/errorGroup"
	"log"
	"strings"
)

//https://developer.work.weixin.qq.com/document/path/90236

type AppMsgSendTarget struct {
	//成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向关注该企业应用的全部成员发送
	Touser []string
	//部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	Toparty []string
	//标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	Totag []string
}
type AppMsgParams struct {
	//成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向关注该企业应用的全部成员发送
	Touser []string
	//部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	Toparty []string
	//标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
	Totag []string

	//表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，默认为0
	Safe uint8

	/*以下为消息内容*/
	//image,voice,video,file
	MediaId string
	//textcard
	Title string
	//textcard
	Description string
	//text
	Content string
	//textcard
	Url string
	//textcard
	BtnText string
}

// AppMsgText 发送应用消息-text
func (_this WeCom) AppMsgText(params AppMsgSendTarget, text string) error {
	return _this.AppMsg(AppMsgParams{
		Touser:  params.Touser,
		Toparty: params.Toparty,
		Totag:   params.Totag,
		Content: text,
	}, "text")
}

// AppMsg 发送应用消息
func (_this WeCom) AppMsg(params AppMsgParams, msgType string) error {
	toUser := "@all"
	if len(params.Totag)+len(params.Toparty)+len(params.Touser) > 0 {
		atAll := false
		for _, s := range params.Touser {
			if s == "@all" {
				atAll = true
			}
		}
		if atAll == false {
			toUser = strings.Join(params.Touser, "|")
		}
	}

	token, err := _this.GetAccessToken()
	if errors.Is(err, errorGroup.AccessTokenInvalidError{}) {
		log.Println("access_token过期重取")
		_this.Cache.Del("access_token")
		token, err = _this.GetAccessToken()
	}
	if err != nil {
		return err
	}
	url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + token
	body := map[string]interface{}{
		//表示是否开启重复消息检查，0表示否，1表示是，默认0
		"enable_duplicate_check": 1,
		//表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时
		"duplicate_check_interval": 3600,
		"msgtype":                  msgType,
		"touser":                   toUser,
		"toparty":                  strings.Join(params.Toparty, "|"),
		"totag":                    strings.Join(params.Totag, "|"),
		"agentid":                  _this.AgentId,
	}
	if params.Safe != 1 {
		body["safe"] = 0
	}
	if msgType != "" {
		mediaMap := map[string]string{}
		if params.MediaId != "" {
			mediaMap["media_id"] = params.MediaId
		}
		if params.Title != "" {
			mediaMap["title"] = params.Title
		}
		if params.Description != "" {
			mediaMap["description"] = params.Description
		}
		if params.Content != "" {
			mediaMap["content"] = params.Content
		}
		if params.Url != "" {
			mediaMap["url"] = params.Url
		}
		body[msgType] = mediaMap
	}
	jsonParam, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = Post(url, jsonParam)
	if err != nil {
		return err
	}
	return nil
}

// AppMsgMarkdown 发送应用消息-markdown
func (_this *WeCom) AppMsgMarkdown(params AppMsgSendTarget, content string) error {
	return _this.AppMsg(AppMsgParams{
		Touser:  params.Touser,
		Toparty: params.Toparty,
		Totag:   params.Totag,
		Content: content,
	}, "markdown")
}

// GetAccessToken 获取access_token
func (_this *WeCom) GetAccessToken() (string, error) {
	var err error
	token, ok := _this.Cache.Get("access_token")
	if ok {
		log.Println("缓存获取access_token")
		return token, nil
	}
	log.Println("获取新access_token")
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + _this.CorpId + "&corpsecret=" + _this.AppSecret
	result, err := Get(url)
	if err != nil {
		return "", err
	}
	token, ok = result["access_token"].(string)
	if ok {
		_this.Cache.Put("access_token", token, 1000*3600)
		return token, nil
	} else {
		return token, errors.New("响应结果中不存在access_token")
	}
}
