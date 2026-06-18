package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"linktree/models"
)

func ListChatSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := models.GetAllChatSessions()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, sessions)
}

func CreateChatSession(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
		Model string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	session, err := models.CreateChatSession(body.Title, body.Model)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "session": session})
}

func DeleteChatSession(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	if err := models.DeleteChatSession(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func UpdateChatSession(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	var body struct {
		Title string `json:"title"`
		Model string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}

	if err := models.UpdateChatSession(id, body.Title, body.Model); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true})
}

func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	messages, err := models.GetChatMessages(id)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, messages)
}

func SendChatMessage(w http.ResponseWriter, r *http.Request) {
	// Check AI config
	apiKey := models.GetAIAPIKey()
	if apiKey == "" {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"未配置 API Key，请先在管理后台 > 设置 中配置 AI API Key\"}\n\n")
		w.(http.Flusher).Flush()
		return
	}

	var body struct {
		SessionID int    `json:"session_id"`
		Content   string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请求格式错误"})
		return
	}
	if body.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "消息不能为空"})
		return
	}

	// Auto-create session if needed
	if body.SessionID == 0 {
		title := body.Content
		if len([]rune(title)) > 20 {
			title = string([]rune(title)[:20]) + "..."
		}
		session, err := models.CreateChatSession(title, models.GetAIModel())
		if err != nil {
			writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
			return
		}
		body.SessionID = session.ID
		// Return session info as first event
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		sessionJSON, _ := json.Marshal(map[string]interface{}{
			"type":    "session_created",
			"session": session,
		})
		fmt.Fprintf(w, "data: %s\n\n", sessionJSON)
		w.(http.Flusher).Flush()
	}

	// Save user message
	_, err := models.CreateChatMessage(body.SessionID, "user", body.Content)
	if err != nil {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"保存消息失败\"}\n\n")
		w.(http.Flusher).Flush()
		return
	}

	// Get recent messages for context (limit to 20)
	messages, err := models.GetRecentChatMessages(body.SessionID, 20)
	if err != nil {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"获取历史消息失败\"}\n\n")
		w.(http.Flusher).Flush()
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	// Build messages array for API
	apiMessages := make([]map[string]string, 0, len(messages)+1)
	apiMessages = append(apiMessages, map[string]string{
		"role":    "system",
		"content": "你是一个有帮助的 AI 助手。请用中文回复。",
	})
	for _, m := range messages {
		apiMessages = append(apiMessages, map[string]string{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	// Call AI API
	baseURL := models.GetAIBaseURL()
	model := body.Content
	session, _ := models.GetChatSessionByID(body.SessionID)
	if session != nil && session.Model != "" {
		model = session.Model
	} else {
		model = models.GetAIModel()
	}

	apiBody := map[string]interface{}{
		"model":    model,
		"messages": apiMessages,
		"stream":   true,
	}
	apiJSON, _ := json.Marshal(apiBody)

	url := strings.TrimRight(baseURL, "/") + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewReader(apiJSON))
	if err != nil {
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"创建请求失败: %s\"}\n\n", err.Error())
		flusher.Flush()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"AI 服务请求失败: %s\"}\n\n", err.Error())
		flusher.Flush()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("AI 服务返回错误 (HTTP %d): %s", resp.StatusCode, string(respBody))
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"%s\"}\n\n", errMsg)
		flusher.Flush()
		return
	}

	// Send start event
	fmt.Fprintf(w, "data: {\"type\":\"start\"}\n\n")
	flusher.Flush()

	// Read SSE stream
	var fullContent strings.Builder
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
				FinishReason *string `json:"finish_reason"`
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
			fullContent.WriteString(delta)
			deltaJSON, _ := json.Marshal(map[string]interface{}{
				"type":    "delta",
				"content": delta,
			})
			fmt.Fprintf(w, "data: %s\n\n", deltaJSON)
			flusher.Flush()
		}
	}

	// Save assistant message
	assistantContent := fullContent.String()
	if assistantContent != "" {
		models.CreateChatMessage(body.SessionID, "assistant", assistantContent)
	}

	// Auto-generate title from first user message
	msgCount := models.GetChatMessageCount(body.SessionID)
	if msgCount == 1 && session != nil && session.Title == "新对话" {
		title := body.Content
		if len([]rune(title)) > 20 {
			title = string([]rune(title)[:20]) + "..."
		}
		models.UpdateChatSession(body.SessionID, title, session.Model)
	}

	// Send done event
	fmt.Fprintf(w, "data: {\"type\":\"done\"}\n\n")
	flusher.Flush()
}
