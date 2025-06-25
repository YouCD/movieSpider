package nameParser

import (
	"encoding/json"
	"fmt"
	"io"
	"movieSpider/internal/config"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/youcd/toolkit/log"
)

type nameParserReq struct {
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream           bool    `json:"stream"`
	CachePrompt      bool    `json:"cache_prompt,omitempty"`
	Samplers         string  `json:"samplers,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	DynatempRange    int     `json:"dynatemp_range,omitempty"`
	DynatempExponent int     `json:"dynatemp_exponent,omitempty"`
	TopK             int     `json:"top_k,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	MinP             float64 `json:"min_p,omitempty"`
	TypicalP         int     `json:"typical_p,omitempty"`
	XtcProbability   int     `json:"xtc_probability,omitempty"`
	XtcThreshold     float64 `json:"xtc_threshold,omitempty"`
	RepeatLastN      int     `json:"repeat_last_n,omitempty"`
	RepeatPenalty    int     `json:"repeat_penalty,omitempty"`
	PresencePenalty  int     `json:"presence_penalty,omitempty"`
	FrequencyPenalty int     `json:"frequency_penalty,omitempty"`
	DryMultiplier    int     `json:"dry_multiplier,omitempty"`
	DryBase          float64 `json:"dry_base,omitempty"`
	DryAllowedLength int     `json:"dry_allowed_length,omitempty"`
	DryPenaltyLastN  int     `json:"dry_penalty_last_n,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	TimingsPerToken  bool    `json:"timings_per_token,omitempty"`
}

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10, // 复用连接
			MaxIdleConnsPerHost: 10,
		},
		//Timeout: 30 * time.Second, // 超时时间
	}
)

func (n *nameParserReq) String() string {
	marshal, err := json.Marshal(n)
	if err != nil {
		return ""
	}
	return string(marshal)
}
func NameParserModelHandler(name string) (string, string, string, string, error) {
	str := newNameParserReq(name).String()
	log.Debug("ai req", str)

	req, err := http.NewRequest("POST", config.Config.Global.NameParserModel+"/v1/chat/completions", strings.NewReader(str))
	if err != nil {
		return "", "", "", "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", "", "", fmt.Errorf("请求失败: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Errorf("关闭请求失败: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("HTTP 错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("读取响应失败: %w", err)
	}
	log.Debug("ai resp", string(body))

	content := jsoniter.Get(body, "choices", 0, "message", "content").ToString()
	split := strings.Split(content, ",")
	if len(split) > 3 {
		typeStr := split[0]
		year := split[1]
		newName := split[2]
		resolution := split[3]
		return typeStr, newName, year, resolution, nil
	} else {
		log.Warnw("content", "parser content", content, "name", name)
	}

	return "", "", "", "", fmt.Errorf("解析失败: %s", name)
}
func newNameParserReq(name string) *nameParserReq {
	return &nameParserReq{
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: name,
			},
		},
		Stream: false,
	}
}
