package handlers

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"linktree/models"
)

func ListFiles(w http.ResponseWriter, r *http.Request) {
	var categoryID *int
	if cidStr := r.URL.Query().Get("category_id"); cidStr != "" {
		cid, err := strconv.Atoi(cidStr)
		if err == nil {
			categoryID = &cid
		}
	}

	files, err := models.GetFilesByCategory(categoryID)
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, files)
}

func ListPublicFiles(w http.ResponseWriter, r *http.Request) {
	files, err := models.GetPublicFiles()
	if err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}
	writeJSON(w, files)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(models.GetMaxFileSize()); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "文件过大或请求格式错误"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "请选择文件"})
		return
	}
	defer file.Close()

	// Detect MIME type
	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	file.Seek(0, io.SeekStart)
	mimeType := http.DetectContentType(buf[:n])
	if mimeType == "application/octet-stream" {
		// Fallback to extension-based detection
		ext := strings.ToLower(filepath.Ext(header.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".png":
			mimeType = "image/png"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		case ".svg":
			mimeType = "image/svg+xml"
		case ".pdf":
			mimeType = "application/pdf"
		case ".txt":
			mimeType = "text/plain"
		case ".doc":
			mimeType = "application/msword"
		case ".docx":
			mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case ".xls":
			mimeType = "application/vnd.ms-excel"
		case ".xlsx":
			mimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case ".zip":
			mimeType = "application/zip"
		case ".rar":
			mimeType = "application/x-rar-compressed"
		case ".gz":
			mimeType = "application/gzip"
		case ".mp4":
			mimeType = "video/mp4"
		case ".webm":
			mimeType = "video/webm"
		case ".ogg":
			mimeType = "video/ogg"
		case ".mp3":
			mimeType = "audio/mpeg"
		case ".wav":
			mimeType = "audio/wav"
		case ".flac":
			mimeType = "audio/flac"
		case ".aac":
			mimeType = "audio/aac"
		case ".m4a":
			mimeType = "audio/x-m4a"
		default:
			mimeType = "application/octet-stream"
		}
	}

	// Check MIME type
	// Allow detected MIME type or fall back to extension check for common types
	ext := strings.ToLower(filepath.Ext(header.Filename))
	extMime := mime.TypeByExtension(ext)

	allowed := models.IsAllowedMimeType(mimeType)
	if !allowed && extMime != "" {
		allowed = models.IsAllowedMimeType(extMime)
	}
	if !allowed {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "不支持的文件类型: " + mimeType})
		return
	}

	// Parse optional fields
	var categoryID *int
	if cidStr := r.FormValue("category_id"); cidStr != "" {
		cid, err := strconv.Atoi(cidStr)
		if err == nil {
			categoryID = &cid
		}
	}

	isPublic := true
	if pubStr := r.FormValue("is_public"); pubStr == "0" || pubStr == "false" {
		isPublic = false
	}

	// Generate stored name
	storedName := models.GenerateStoredName(header.Filename)
	uploadDir := models.GetUploadDir()
	os.MkdirAll(uploadDir, 0755)
	savePath := filepath.Join(uploadDir, storedName)

	// Save file
	dst, err := os.Create(savePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "保存文件失败"})
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(savePath)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "写入文件失败"})
		return
	}

	// Save to database
	record, err := models.CreateFile(categoryID, header.Filename, storedName, mimeType, written, isPublic)
	if err != nil {
		os.Remove(savePath)
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"ok": true, "file": record})
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	file, err := models.GetFileByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "文件不存在"})
		return
	}

	uploadDir := models.GetUploadDir()
	filePath := filepath.Join(uploadDir, file.StoredName)

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "文件已丢失"})
		return
	}

	// Increment download count
	models.IncrementDownloadCount(id)

	// Set headers
	w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+urlEncode(file.OriginalName))
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))

	http.ServeFile(w, r, filePath)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "无效ID"})
		return
	}

	file, err := models.GetFileByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJSON(w, map[string]interface{}{"ok": false, "message": "文件不存在"})
		return
	}

	// Delete physical file
	uploadDir := models.GetUploadDir()
	filePath := filepath.Join(uploadDir, file.StoredName)
	os.Remove(filePath)

	// Delete database record
	if err := models.DeleteFile(id); err != nil {
		writeJSON(w, map[string]interface{}{"ok": false, "message": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"ok": true})
}

func urlEncode(s string) string {
	// Simple URL encoding for non-ASCII filenames
	return strings.ReplaceAll(s, " ", "%20")
}
