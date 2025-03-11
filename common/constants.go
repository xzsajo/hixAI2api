package common

import "time"

var StartTime = time.Now().Unix() // unit: second
var Version = "v1.1.0"            // this hard coding will be replaced automatically when building, no need to manually change

type HixModelInfo struct {
	Model   string
	Credit  int
	Type    string
	ModelID int
}

// 创建映射表（假设用 model 名称作为 key）
var modelRegistry = map[string]HixModelInfo{
	"deepseek-r1":            {"deepseek-r1", 1, "STANDARD", 85426},
	"deepseek-v3":            {"deepseek-v3", 1, "STANDARD", 85427},
	"claude-3-7-sonnet":      {"claude-3-7-sonnet", 20, "ADVANCED", 85429},
	"claude-3-5-haiku":       {"claude-3-5-haiku", 10, "STANDARD", 85422},
	"openai-o3-mini":         {"openai-o3-mini", 200, "STANDARD", 85424},
	"openai-o1":              {"openai-o1", 40, "ADVANCED", 1181},
	"openai-o1-mini":         {"openai-o1-mini", 200, "STANDARD", 1182},
	"grok-2":                 {"grok-2", 100, "STANDARD", 85428},
	"gpt-4o":                 {"gpt-4o", 30, "STANDARD", 5},
	"gpt-4o-128k":            {"gpt-4o-128k", 125, "STANDARD", 6},
	"gpt-4o-mini":            {"gpt-4o-mini", 4, "STANDARD", 86},
	"gpt-4-turbo":            {"gpt-4-turbo", 20, "ADVANCED", 8},
	"gpt-4-turbo-128k":       {"gpt-4-turbo-128k", 20, "ADVANCED", 9},
	"gpt4":                   {"gpt4", 45, "ADVANCED", 10},
	"claude":                 {"claude", 4, "STANDARD", 42},
	"claude-3-5-sonnet":      {"claude-3-5-sonnet", 100, "STANDARD", 50},
	"claude-3-haiku":         {"claude-3-haiku", 4, "STANDARD", 52},
	"claude-3-opus":          {"claude-3-opus", 45, "ADVANCED", 53},
	"claude-3-5-haiku-200k":  {"claude-3-5-haiku-200k", 4, "STANDARD", 85423},
	"claude-3-5-sonnet-200k": {"claude-3-5-sonnet-200k", 100, "STANDARD", 54},
	"claude-3-sonnet-200k":   {"claude-3-sonnet-200k", 100, "STANDARD", 55},
	"claude-3-haiku-200k":    {"claude-3-haiku-200k", 20, "STANDARD", 56},
	"claude-3-opus-200k":     {"claude-3-opus-200k", 120, "ADVANCED", 57},
	"gemini-1-5-flash":       {"gemini-1-5-flash", 4, "STANDARD", 59},
	"gemini-1-5-pro":         {"gemini-1-5-pro", 18, "STANDARD", 60},
	"gemini-1-5-flash-128k":  {"gemini-1-5-flash-128k", 30, "STANDARD", 62},
	"gemini-1-5-pro-128k":    {"gemini-1-5-pro-128k", 175, "STANDARD", 63},
	"gemini-1-5-flash-1m":    {"gemini-1-5-flash-1m", 170, "STANDARD", 65},
	"gemini-1-5-pro-1m":      {"gemini-1-5-pro-1m", 2500, "STANDARD", 67},
	"chatgpt":                {"chatgpt", 4, "STANDARD", 2},
	"gpt-3-5-turbo":          {"gpt-3-5-turbo", 4, "STANDARD", 3},
	"gpt-3-5-turbo-16k":      {"gpt-3-5-turbo-16k", 12, "STANDARD", 4},
	"claude-instant-100k":    {"claude-instant-100k", 8, "STANDARD", 44},
	"claude-2":               {"claude-2", 35, "STANDARD", 45},
	"claude-2-100k":          {"claude-2-100k", 75, "STANDARD", 47},
	"claude-2-1-200k":        {"claude-2-1-200k", 300, "STANDARD", 49},
	"claude-3-sonnet":        {"claude-3-sonnet", 20, "STANDARD", 51},
	"gemini":                 {"gemini", 4, "STANDARD", 83},
	"gemini-1-0-pro":         {"gemini-1-0-pro", 4, "STANDARD", 58},
}

// 通过 model 名称查询的方法
func GetHixModelInfo(modelName string) (HixModelInfo, bool) {
	info, exists := modelRegistry[modelName]
	return info, exists
}

func GetHixModelList() []string {
	var modelList []string
	for k := range modelRegistry {
		modelList = append(modelList, k)
	}
	return modelList
}
