package common

import "time"

var StartTime = time.Now().Unix() // unit: second
var Version = "v1.1.1"            // this hard coding will be replaced automatically when building, no need to manually change

type HixModelInfo struct {
	Model     string
	Credit    int
	Type      string
	ModelID   int
	MaxTokens int
}

// 创建映射表（假设用 model 名称作为 key）
var ModelRegistry = map[string]HixModelInfo{
	"deepseek-r1":            {"deepseek-r1", 1, "STANDARD", 85426, 8000},
	"deepseek-v3":            {"deepseek-v3", 1, "STANDARD", 85427, 8000},
	"gpt-4o-mini":            {"gpt-4o-mini", 4, "STANDARD", 86, 8000},
	"claude":                 {"claude", 4, "STANDARD", 42, 8000},
	"claude-3-haiku":         {"claude-3-haiku", 4, "STANDARD", 52, 8000},
	"claude-3-5-haiku-200k":  {"claude-3-5-haiku-200k", 4, "STANDARD", 85423, 100000},
	"gemini-1-5-flash":       {"gemini-1-5-flash", 4, "STANDARD", 59, 8000},
	"chatgpt":                {"chatgpt", 4, "STANDARD", 2, 8000},
	"gemini-1-0-pro":         {"gemini-1-0-pro", 4, "STANDARD", 58, 8000},
	"gemini":                 {"gemini", 4, "STANDARD", 83, 8000},
	"gpt-3-5-turbo":          {"gpt-3-5-turbo", 4, "STANDARD", 3, 8000},
	"claude-instant-100k":    {"claude-instant-100k", 8, "STANDARD", 44, 50000},
	"claude-3-5-haiku":       {"claude-3-5-haiku", 10, "STANDARD", 85422, 8000},
	"gpt-3-5-turbo-16k":      {"gpt-3-5-turbo-16k", 12, "STANDARD", 4, 8000},
	"gemini-1-5-pro":         {"gemini-1-5-pro", 18, "STANDARD", 60, 8000},
	"claude-3-haiku-200k":    {"claude-3-haiku-200k", 20, "STANDARD", 56, 100000},
	"claude-3-sonnet":        {"claude-3-sonnet", 20, "STANDARD", 51, 8000},
	"gpt-4o":                 {"gpt-4o", 30, "STANDARD", 5, 8000},
	"gemini-1-5-flash-128k":  {"gemini-1-5-flash-128k", 30, "STANDARD", 62, 64000},
	"claude-2":               {"claude-2", 35, "STANDARD", 45, 8000},
	"claude-2-100k":          {"claude-2-100k", 75, "STANDARD", 47, 50000},
	"grok-2":                 {"grok-2", 100, "STANDARD", 85428, 8000},
	"claude-3-5-sonnet":      {"claude-3-5-sonnet", 100, "STANDARD", 50, 8000},
	"claude-3-5-sonnet-200k": {"claude-3-5-sonnet-200k", 100, "STANDARD", 54, 100000},
	"claude-3-sonnet-200k":   {"claude-3-sonnet-200k", 100, "STANDARD", 55, 100000},
	"gpt-4o-128k":            {"gpt-4o-128k", 125, "STANDARD", 6, 64000},
	"gemini-1-5-pro-128k":    {"gemini-1-5-pro-128k", 175, "STANDARD", 63, 64000},
	"gemini-1-5-flash-1m":    {"gemini-1-5-flash-1m", 170, "STANDARD", 65, 100000},
	"openai-o3-mini":         {"openai-o3-mini", 200, "STANDARD", 85424, 8000},
	"openai-o1-mini":         {"openai-o1-mini", 200, "STANDARD", 1182, 8000},
	"claude-2-1-200k":        {"claude-2-1-200k", 300, "STANDARD", 49, 100000},
	"gemini-1-5-pro-1m":      {"gemini-1-5-pro-1m", 2500, "STANDARD", 67, 100000},

	"gpt-4-turbo":        {"gpt-4-turbo", 20, "ADVANCED", 8, 8000},
	"gpt-4-turbo-128k":   {"gpt-4-turbo-128k", 20, "ADVANCED", 9, 64000},
	"claude-3-7-sonnet":  {"claude-3-7-sonnet", 20, "ADVANCED", 85429, 8000},
	"openai-o1":          {"openai-o1", 40, "ADVANCED", 1181, 8000},
	"claude-3-opus":      {"claude-3-opus", 45, "ADVANCED", 53, 8000},
	"gpt4":               {"gpt4", 45, "ADVANCED", 10, 8000},
	"claude-3-opus-200k": {"claude-3-opus-200k", 120, "ADVANCED", 57, 100000},
}

// 通过 model 名称查询的方法
func GetHixModelInfo(modelName string) (HixModelInfo, bool) {
	info, exists := ModelRegistry[modelName]
	return info, exists
}

func GetHixModelList() []string {
	var modelList []string
	for k := range ModelRegistry {
		modelList = append(modelList, k)
	}
	return modelList
}
