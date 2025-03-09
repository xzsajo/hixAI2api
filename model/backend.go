package model

type AuthVerifyReq struct {
	AccessKey string `json:"accessKey"`
}

type ApiKeySaveReq struct {
	ApiKey string `json:"apiKey"`
}

type ApiKeyUpdateReq struct {
	Id     string `json:"id"`
	ApiKey string `json:"apiKey"`
}

type ApiKeyResp struct {
	Id         string `json:"id"`
	ApiKey     string `json:"apiKey"`
	CreateTime string `json:"createTime"`
}
