package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"linktree/models"
)

// PublicResume 公开接口：返回"简历"文件的下载信息。
// 判定顺序：公开文件里按 文件名关键词(resume/简历/cv) + PDF MIME 打分，取最高分的一个。
// 库里没有匹配文件时仍返回 200 + available=false，前端据此隐藏下载按钮（不产生 404 噪音）。
func PublicResume(w http.ResponseWriter, r *http.Request) {
	files, err := models.GetPublicFiles()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}

	var best *models.File
	bestScore := 0
	for i := range files {
		f := &files[i]
		if s := resumeScore(*f); s > bestScore {
			bestScore = s
			best = f
		}
	}

	if best == nil {
		writeJSON(w, map[string]interface{}{
			"ok":        true,
			"available": false,
			"message":   "简历尚未上传。请在后台「文件」页上传 PDF（文件名含 resume 或 简历 即可自动识别）",
		})
		return
	}

	writeJSON(w, map[string]interface{}{
		"ok":             true,
		"available":      true,
		"id":             best.ID,
		"name":           displayFileName(best.OriginalName),
		"size":           best.Size,
		"formatted_size": best.FormattedSize,
		"mime_type":      best.MimeType,
		"download_count": best.DownloadCount,
		"url":            "/api/files/download/" + strconv.Itoa(best.ID),
	})
}

func resumeScore(f models.File) int {
	name := strings.ToLower(displayFileName(f.OriginalName))
	score := 0
	if strings.Contains(name, "resume") {
		score += 100
	}
	if strings.Contains(name, "简历") {
		score += 100
	}
	if strings.Contains(name, "cv") {
		score += 40
	}
	if isPDF(f) {
		score += 30
	}
	// 只有 PDF 或命中关键词的文件才算简历候选，避免随便一个公开图片被当成简历
	if !isPDF(f) && score < 40 {
		return 0
	}
	return score
}

func isPDF(f models.File) bool {
	return strings.Contains(strings.ToLower(f.MimeType), "pdf") ||
		strings.HasSuffix(strings.ToLower(f.OriginalName), ".pdf")
}

// displayFileName 还原可读文件名：旧数据里 original_name 存的是 URL 编码后的值
// （如 %E7%AE%80%E5%8E%86.pdf），既用于展示也用于关键词匹配。
func displayFileName(raw string) string {
	if decoded, err := url.QueryUnescape(raw); err == nil && decoded != "" {
		return decoded
	}
	return raw
}
