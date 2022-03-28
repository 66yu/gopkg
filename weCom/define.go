package weCom

import "github.com/ayu-666/gopkg/fileCache"

type Params struct {
	//https://work.weixin.qq.com/wework_admin/frame#profile
	//网页版-我的企业-企业信息-企业ID
	CorpId string

	//https://work.weixin.qq.com/wework_admin/frame#apps
	//网页版-应用管理-点击自己创建的应用(没有则创建)-AgentId
	AgentId string

	//https://work.weixin.qq.com/wework_admin/frame#apps
	//网页版-应用管理-点击自己创建的应用(没有则创建)-Secret
	AppSecret string

	//缓存文件路径
	CacheFilePath string
	//缓存文件名
	CacheFilename string
}
type WeCom struct {
	CorpId    string
	AgentId   string
	AppSecret string
	Cache     *fileCache.FcDb
}
