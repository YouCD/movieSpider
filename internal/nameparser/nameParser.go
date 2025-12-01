package nameparser

import (
	"context"
	"encoding/json"
	"fmt"
	"movieSpider/internal/config"
	"movieSpider/internal/types"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/youcd/toolkit/log"
)

const prompt = `你是一位专业的媒体文件元数据解析专家，擅长从混乱的种子名称字符串中提取结构化信息。

**输入**：一个或多个种子名称（每行一个）
**输出**：严格的JSON数组，每个元素包含以下字段：
- 'id': 数字类型，从0开始的递增序号
- 'typeStr': 字符串（movie/tv/空字符串）
- 'newName': 字符串，规范化剧名（英文单词首字母大写，用.连接）
- 'year': 数字类型，4位年份（无法确定则返回'0'）
- 'resolution': 数字类型，分辨率数值（无法确定则返回'0'）

---

### **解析规则（按优先级严格执行）**

**1. 类型判定**
- **tv**：字符串包含'SxxExx'、'Season.X'或其变体（如'1x02'、'01E05'等数字+E/x+数字格式）
- **movie**：不满足tv规则
- **""**: 若无法提取任何有效英文字母作为剧名，则留空

**2. 剧名提取（核心逻辑）**
- **电影**：提取'年份'或'分辨率'之前的所有连续英文字母片段
- **剧集**：提取'SxxExx/1x02'或'Season.X'之前的所有连续英文字母片段
- **多语言混合**：只保留英文字母片段（如示例5，忽略俄文'Подозрительные.лицы'）
- **格式要求**：
  - 每个单词首字母大写，用'.'连接
  - 删除所有标点符号（':', '!', '?', '《》'等）
  - '&'全局替换为'And'
  - 末尾不得有'.'或空格
  - **剧名不得包含年份**

**3. 年份识别**
- 提取4位数字（1900-2099）
- 若提取失败则返回'0'

**4. 分辨率提取**
- 匹配模式：'1080p', '720p', '2160p', '4K', 'UHD'
- **别名映射**：'4K'/'UHD' → '2160'
- 最终输出**纯数字**，若失败则返回'0'

**5. 错误处理**
- 若无法提取剧名，所有字段返回空字符串或'0'
- 静默失败，不添加任何解释或错误字段

**6. 无法识别处理**
若整个字符串中**不存在任何连续英文字母片段**（如纯中文、日文、符号），直接返回：
'{"id":0,"typeStr":"","newName":"","year":0,"resolution":0}'

---

### **修正后的正确示例**

**示例 1（电影）：**
<输入>
A-Heavenly-Vintage-2009-1080p-BluRay-x265-RARBG
<输出>
[{"id":0,"typeStr":"movie","newName":"A.Heavenly.Vintage","year":2009,"resolution":1080}]

**示例 2（剧集无年份）：**
<输入>
Big Antique Adventure With Susan Calman S01 1080p HDTV H264-DARKFLiX[rartv]
<输出>
[{"id":0,"typeStr":"tv","newName":"Big.Antique.Adventure.With.Susan.Calman","year":0,"resolution":1080}]

**示例 3（剧集含年份）：**
<输入>
www.Torrenting.com - Tracker.2024.S02E07.1080p.HEVC.x265-MeGusta
<输出>
[{"id":0,"typeStr":"tv","newName":"Tracker","year":2024,"resolution":1080}]

**示例 4（电影带符号）：**
<输入>
《The Book Of Solutions (2023) [1080p] [BluRay] [5.1] [YTS.MX]》
<输出>
[{"id":0,"typeStr":"movie","newName":"The.Book.Of.Solutions","year":2023,"resolution":1080}]

**示例 5（多语言）：**
<输入>
Подозрительные.лицы.The.Usual.Suspects.1995.JPN.Transfer.BDRip-HEVC.1080p.mkv
<输出>
[{"id":0,"typeStr":"movie","newName":"The.Usual.Suspects","year":1995,"resolution":1080}]

**示例 6（剧集忽略副标题）：**
<输入>
Greys Anatomy S21E13 Dont You Forget About Me 1080p AMZN WEB-DL DDP5 1 HEVC
<输出>
[{"id":0,"typeStr":"tv","newName":"Greys.Anatomy","year":0,"resolution":1080}]

**示例 7（格式变体）：**
<输入>
Alien Earth 1x02 Mr October 1080p WEB-DL H265 Ita Eng AC3 5 1 Multisub iDN CreW
<输出>
[{"id":0,"typeStr":"tv","newName":"Alien.Earth","year":0,"resolution":1080}]

**示例 8（无法识别）：**
<输入>
浮世絵ＥＤＯ－ＬＩＦＥ　べらぼうの世界　出産後も座ったまま！驚きの出産事情
<输出>
[{"id":0,"typeStr":"","newName":"","year":0,"resolution":0}]

**示例 9（批量处理）：**
<输入>
A-Heavenly-Vintage-2009-1080p-BluRay-x265-RARBG
Tracker.2024.S02E07.1080p.HEVC.x265-MeGusta
无效种子名称无关键信息
<输出>
[{"id":0,"typeStr":"movie","newName":"A.Heavenly.Vintage","year":2009,"resolution":1080},{"id":1,"typeStr":"tv","newName":"Tracker","year":2024,"resolution":1080},{"id":2,"typeStr":"","newName":"","year":0,"resolution":0}]

---

**只输出解析后的JSON数组，不要包含任何额外文字或解释。**
`

var (
	ErrNamesIsEmpty = fmt.Errorf("names is empty")
)

func ModelHandler(ctx context.Context, names ...string) (map[string]*types.LLMResult, error) {
	if len(names) == 0 {
		return nil, ErrNamesIsEmpty //nolint:err113
	}
	conf := openai.DefaultAnthropicConfig(config.Config.LLM.ApiKey, config.Config.LLM.BaseURL)
	conf.APIType = openai.APITypeOpenAI

	client := openai.NewClientWithConfig(conf)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: config.Config.LLM.Model,
			Messages: []openai.ChatCompletionMessage{
				{Role: "system", Content: prompt},
				{Role: "user", Content: strings.Join(names, "\n")},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)
	if err != nil {
		//nolint:err113
		return nil, fmt.Errorf("解析失败: %s", names)
	}

	Content := resp.Choices[0].Message.Content

	var result []*types.LLMResult
	err = json.Unmarshal([]byte(Content), &result)
	if err != nil {
		return nil, fmt.Errorf("解析失败,content:%s , err: %w", Content, err)
	}
	log.WithCtx(ctx).Warnw("content", "parser content", Content, "name", strings.Join(names, ";"))

	resultMap := make(map[string]*types.LLMResult)
	if len(names) == len(result) {
		for index, llmResult := range result {
			if llmResult.NewName == "" {
				continue
			}
			resultMap[names[index]] = llmResult
		}
	}
	return resultMap, nil
}
