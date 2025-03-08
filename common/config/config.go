package config

import (
	"hixai2api/common/env"
	"os"
	"strings"
	"time"
)

var ApiSecret = os.Getenv("API_SECRET")
var ApiSecrets = strings.Split(os.Getenv("API_SECRET"), ",")

var MysqlDsn = os.Getenv("MYSQL_DSN")

// var IpBlackList = os.Getenv("IP_BLACK_LIST")
var IpBlackList = strings.Split(os.Getenv("IP_BLACK_LIST"), ",")

var ProxyUrl = env.String("PROXY_URL", "")

// 隐藏思考过程
var ReasoningHide = env.Int("REASONING_HIDE", 0)

// 前置message
var PRE_MESSAGES_JSON = env.String("PRE_MESSAGES_JSON", "")

// 路由前缀
var RoutePrefix = env.String("ROUTE_PREFIX", "")
var SwaggerEnable = os.Getenv("SWAGGER_ENABLE")
var OnlyOpenaiApi = os.Getenv("ONLY_OPENAI_API")

var DebugEnabled = os.Getenv("DEBUG") == "true"

var RateLimitKeyExpirationDuration = 20 * time.Minute

var RequestOutTimeDuration = 5 * time.Minute

var (
	RequestRateLimitNum            = env.Int("REQUEST_RATE_LIMIT", 60)
	RequestRateLimitDuration int64 = 1 * 60
)
