package model

type AuthVerifyReq struct {
	AccessKey string `json:"accessKey"`
}

type ApiKeySaveReq struct {
	ApiKey string `json:"apiKey"`
	Remark string `json:"remark"`
}

type ApiKeyUpdateReq struct {
	Id     string `json:"id"`
	ApiKey string `json:"apiKey"`
	Remark string `json:"remark"`
}

type ApiKeyResp struct {
	Id         string `json:"id"`
	ApiKey     string `json:"apiKey"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}

type CookieSaveReq struct {
	Cookie string `json:"cookie"`
	Remark string `json:"remark"`
}

type CookieUpdateReq struct {
	Id     string `json:"id"`
	Cookie string `json:"cookie"`
	Remark string `json:"remark"`
}

type CookieResp struct {
	Id         string `json:"id"`
	Cookie     string `json:"cookie"`
	Credit     int    `json:"credit"`
	Remark     string `json:"remark"`
	CreateTime string `json:"createTime"`
}
