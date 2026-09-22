package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"linktree/models"
)

// ==================== 公开 AI 问答（访客用） ====================
// 设计：无状态问答，不保存会话；固定 system prompt（资料库）；
// IP 维度限流防滥用；问题长度限制；SSE 流式输出。

// 问答限流：每分钟 3 次 + 每天 30 次
var (
	askMinute   = map[string][]time.Time{}
	askDaily    = map[string][]time.Time{}
	askLimiterM sync.Mutex
)

const (
	askPerMinute = 3
	askPerDay    = 30
	askMaxLen    = 200
)

func isAskAllowed(ip string) bool {
	askLimiterM.Lock()
	defer askLimiterM.Unlock()
	now := time.Now()

	// 清理过期记录
	var recentMin []time.Time
	for _, t := range askMinute[ip] {
		if now.Sub(t) < time.Minute {
			recentMin = append(recentMin, t)
		}
	}
	askMinute[ip] = recentMin
	if len(recentMin) >= askPerMinute {
		return false
	}

	var recentDay []time.Time
	for _, t := range askDaily[ip] {
		if now.Sub(t) < 24*time.Hour {
			recentDay = append(recentDay, t)
		}
	}
	askDaily[ip] = recentDay
	if len(recentDay) >= askPerDay {
		return false
	}
	return true
}

func recordAsk(ip string) {
	askLimiterM.Lock()
	defer askLimiterM.Unlock()
	now := time.Now()
	askMinute[ip] = append(askMinute[ip], now)
	askDaily[ip] = append(askDaily[ip], now)
}

// 资料库 system prompt（AI 据此回答关于"啊芃"的问题）
const askSystemPrompt = `你是"啊芃"个人主页上的 AI 问答助手。你的职责是回答访客关于啊芃的问题。

# 啊芃的资料
- 名字：啊芃
- 职位：Full-Stack / AI Developer（全栈 / AI 开发者）
- 简介：一个热爱技术的开发者
- 理念：用代码把想法变成产品 · 全栈开发 · AI 应用实践者
- 技术栈：Go、Python、TypeScript、React、Vue、LLM/AI 应用、Docker、SQL
- 经验：3+ 年开发经验、20+ 项目作品、5+ AI 应用实践

# 工作经历
- 2024 至今：全栈开发者 · 独立开发者——主导开发多款 AI 应用与自动化工具，负责前后端与部署运维
- 2022-2024：后端开发者——负责 API 设计与服务端开发，熟悉高并发与数据存储方案
- 2018-2022：计算机科学与技术 · 本科

# 代表项目
1. LinkTree 个人主页系统：基于 Go + SQLite 开发的个人主页系统，Bento 网格布局、时钟/日历/天气小组件、作品与文件管理、AI 对话、后台管理，Docker 一键部署。
2. 中转站：自建工具中转服务，集中管理多类在线工具入口。
3. AI 对话助手：基于大模型 API 的对话应用，支持多轮会话、Markdown 渲染与流式输出。

# 在线链接
- CSDN 技术博客：https://blog.csdn.net/
- GitHub：https://github.com/
- 中转站：http://62.234.92.232:3000

# 行为规则
- 只回答与啊芃相关的问题（项目、技能、经历、求职意向等）
- 与啊芃无关的问题，礼貌拒绝："这个问题超出了我的范围，我专门用来介绍啊芃。你可以问问他的项目经验或技术栈～"
- 不要编造资料中没有的信息；不确定就说"这个信息我这边没有，建议直接联系他"
- 用中文回答，语气友好专业，回复简洁（200 字以内）
- 被问到联系方式时：建议通过主页上的链接（CSDN/GitHub）或邮箱联系`

// AskPage 公开问答页
func AskPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "ask.html", nil)
}

// AskAI 公开问答 API（SSE 流式）
func AskAI(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)

	var body struct {
		Question string `json:"question"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	question := strings.TrimSpace(body.Question)
	if question == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请输入问题"})
		return
	}
	if len([]rune(question)) > askMaxLen {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": fmt.Sprintf("问题过长，请控制在 %d 字以内", askMaxLen)})
		return
	}

	if !isAskAllowed(ip) {
		w.WriteHeader(http.StatusTooManyRequests)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "提问太频繁啦，请稍后再试（每分钟最多 3 次）"})
		return
	}

	apiKey := models.GetAIAPIKey()
	if apiKey == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "AI 服务未配置"})
		return
	}

	recordAsk(ip)

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// 构造 AI 请求（无历史，单轮问答）
	apiMessages := []map[string]string{
		{"role": "system", "content": askSystemPrompt},
		{"role": "user", "content": question},
	}
	apiBody := map[string]interface{}{
		"model":    models.GetAIModel(),
		"messages": apiMessages,
		"stream":   true,
	}
	apiJSON, _ := json.Marshal(apiBody)

	baseURL := models.GetAIBaseURL()
	url := strings.TrimRight(baseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(r.Context(), "POST", url, bytes.NewReader(apiJSON))
	if err != nil {
		sendAskError(w, flusher, "创建请求失败")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		sendAskError(w, flusher, "AI 服务请求失败")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		io.ReadAll(resp.Body) // 丢弃错误体，不回显给访客
		sendAskError(w, flusher, "AI 服务暂时不可用")
		return
	}

	fmt.Fprintf(w, "data: {\"type\":\"start\"}\n\n")
	flusher.Flush()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		delta := chunk.Choices[0].Delta.Content
		if delta != "" {
			deltaJSON, _ := json.Marshal(map[string]interface{}{"type": "delta", "content": delta})
			fmt.Fprintf(w, "data: %s\n\n", deltaJSON)
			flusher.Flush()
		}
	}
	fmt.Fprintf(w, "data: {\"type\":\"done\"}\n\n")
	flusher.Flush()
}

func sendAskError(w http.ResponseWriter, flusher http.Flusher, msg string) {
	errJSON, _ := json.Marshal(map[string]interface{}{"type": "error", "message": msg})
	fmt.Fprintf(w, "data: %s\n\n", errJSON)
	flusher.Flush()
}