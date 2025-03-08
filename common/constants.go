package common

import "time"

var StartTime = time.Now().Unix() // unit: second
var Version = "v1.0.0"            // this hard coding will be replaced automatically when building, no need to manually change

type HixModelInfo struct {
	Model   string
	Credit  int
	ModelID int
}

// 创建映射表（假设用 model 名称作为 key）
var modelRegistry = map[string]HixModelInfo{
	"deepseek-r1":          {"deepseek-r1", 1, 85426},
	"deepseek-v3":          {"deepseek-v3", 1, 85427},
	"grok-2":               {"grok-2", 100, 85428},
	"claude-3-7-sonnet":    {"claude-3-7-sonnet", 50, 85429},
	"chatgpt":              {"chatgpt", 4, 2},
	"gpt-3.5-turbo":        {"gpt-3.5-turbo", 4, 3},
	"gpt-3.5-turbo-16k":    {"gpt-3.5-turbo-16k", 12, 4},
	"gpt-4o":               {"gpt-4o", 30, 5},
	"gpt-4o-128k":          {"gpt-4o-128k", 125, 6},
	"gpt-4-turbo":          {"gpt-4-turbo", 35, 8},
	"gpt-4-turbo-128k":     {"gpt-4-turbo-128k", 250, 9},
	"gpt-4":                {"gpt-4", 350, 10},
	"claude":               {"claude", 4, 42},
	"claude-instant-100k":  {"claude-instant-100k", 8, 44},
	"claude-2":             {"claude-2", 35, 45},
	"claude-2-100k":        {"claude-2-100k", 75, 47},
	"claude-2-1-200k":      {"claude-2-1-200k", 300, 49},
	"claude-3-5-sonnet-v2": {"claude-3-5-sonnet-v2", 100, 50},
	"claude-3-sonnet":      {"claude-3-sonnet", 20, 51},
	"claude-3-haiku":       {"claude-3-haiku", 4, 52},
	"claude-3-opus":        {"claude-3-opus", 200, 53},
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
