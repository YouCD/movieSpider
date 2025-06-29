package nameParser

import (
	"context"
	"fmt"
	"movieSpider/internal/config"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/youcd/toolkit/log"
)

const prompt = `规范化BT种子名称，移除广告和无关信息，按以下格式输出：类型,年份,规范名称,分辨率`

func NameParserModelHandler(name string) (string, string, string, string, error) {
	// +"/v1/chat/completions
	conf := openai.DefaultAnthropicConfig("", config.Config.Global.NameParserModel+"/v1")
	conf.APIType = openai.APITypeOpenAI

	client := openai.NewClientWithConfig(conf)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: prompt},
				{Role: "user", Content: name},
			},
		},
	)
	if err != nil {
		log.Error("解码响应时发生错误:", err)
		return "", "", "", "", fmt.Errorf("解析失败: %s", name)
	}

	Content := resp.Choices[0].Message.Content

	split := strings.Split(Content, ",")
	if len(split) > 3 {
		typeStr := split[0]
		year := split[1]
		newName := split[2]
		resolution := split[3]
		return typeStr, newName, year, resolution, nil
	} else {
		log.Warnw("content", "parser content", Content, "name", name)
	}
	return "", "", "", "", fmt.Errorf("解析失败: %s", name)
}
