package weCom

import (
	"github.com/ayu-666/gopkg/fileCache"
)

// Init 登录https://work.weixin.qq.com/
func Init(params Params) (*WeCom, error) {
	instance := &WeCom{
		CorpId:    params.CorpId,
		AgentId:   params.AgentId,
		AppSecret: params.AppSecret,
	}
	cachePath := params.CacheFilePath
	cacheFilename := params.CacheFilename
	if cachePath == "" {
		cachePath = "./_cache"
	}
	if cacheFilename == "" {
		cacheFilename = "wecom"
	}
	cache, err := fileCache.Init(cachePath, cacheFilename)
	instance.Cache = cache
	return instance, err
}
