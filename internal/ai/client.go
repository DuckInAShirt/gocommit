package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewClient(apiKey, baseURL, model string) *Client {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) GenerateCommitMessage(diff string) (string, error) {
	prompt := buildPrompt(diff)

	reqBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
		MaxTokens:   500,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("parse response failed: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	message := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	message = cleanMessage(message)

	return message, nil
}

const systemPrompt = `你是一位专业的Git提交信息撰写专家。你的任务是根据代码变更内容生成简洁、准确的中文提交信息。

规则：
1. 必须使用 Conventional Commits 格式：类型(范围): 描述
2. 类型必须是以下之一：feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
3. 描述用中文，简洁明了，不超过50个字符
4. 范围是可选的，用英文小写
5. 如果变更涉及多个类型，选择最主要的那个
6. 不要添加多余的解释、注释或前缀
7. 只输出一行提交信息，不要输出其他内容

示例：
feat: 添加用户登录功能
fix(auth): 修复token过期未刷新的问题
refactor: 重构配置加载逻辑
docs: 更新API文档
chore(deps): 升级依赖版本`

func buildPrompt(diff string) string {
	truncated := diff
	if len(diff) > 8000 {
		truncated = diff[:8000]
		if idx := strings.LastIndex(truncated, "\n"); idx > 0 {
			truncated = truncated[:idx]
		}
		truncated += "\n... (diff truncated)"
	}

	return fmt.Sprintf(`请根据以下 git diff 生成一条中文提交信息：

%s`, truncated)
}

func cleanMessage(msg string) string {
	msg = strings.TrimPrefix(msg, "```")
	msg = strings.TrimSuffix(msg, "```")
	msg = strings.TrimSpace(msg)

	for _, prefix := range []string{"提交信息：", "提交信息:", "commit:", "Commit:"} {
		msg = strings.TrimPrefix(msg, prefix)
	}

	return strings.TrimSpace(msg)
}
