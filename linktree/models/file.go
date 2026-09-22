package models

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type File struct {
	ID            int     `json:"id"`
	CategoryID    *int    `json:"category_id"`
	OriginalName  string  `json:"original_name"`
	StoredName    string  `json:"stored_name"`
	MimeType      string  `json:"mime_type"`
	Size          int64   `json:"size"`
	FormattedSize string  `json:"formatted_size"`
	IsPublic      bool    `json:"is_public"`
	DownloadCount int     `json:"download_count"`
	SortOrder     float64 `json:"sort_order"`
	CreatedAt     string  `json:"created_at"`
}

// allowedMimeTypes defines which file types are permitted for upload
var allowedMimeTypes = map[string]bool{
	"image/jpeg":                                       true,
	"image/png":                                        true,
	"image/gif":                                        true,
	"image/webp":                                       true,
	"image/svg+xml":                                    true,
	"application/pdf":                                  true,
	"text/plain":                                       true,
	"application/msword":                               true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel":                         true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"application/zip":                                  true,
	"application/x-rar-compressed":                     true,
	"application/gzip":                                 true,
	"video/mp4":                                        true,
	"video/webm":                                       true,
	"video/ogg":                                        true,
	"audio/mpeg":                                       true,
	"audio/mp3":                                        true,
	"audio/wav":                                        true,
	"audio/ogg":                                        true,
	"audio/flac":                                       true,
	"audio/aac":                                        true,
	"audio/mp4":                                        true,
	"audio/x-m4a":                                      true,
	"audio/webm":                                       true,
	// 注意：不允许 application/octet-stream —— DetectContentType 对未知二进制一律返回该类型，
	// 允许它等于允许任意可执行文件上传
}

func IsAllowedMimeType(mimeType string) bool {
	return allowedMimeTypes[mimeType]
}

func GetFilesByCategory(categoryID *int) ([]File, error) {
	var rows *sql.Rows
	var err error
	if categoryID != nil {
		rows, err = DB.Query("SELECT id, category_id, original_name, stored_name, mime_type, size, is_public, download_count, sort_order, created_at FROM files WHERE category_id=? ORDER BY sort_order ASC", *categoryID)
	} else {
		rows, err = DB.Query("SELECT id, category_id, original_name, stored_name, mime_type, size, is_public, download_count, sort_order, created_at FROM files ORDER BY sort_order ASC")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.CategoryID, &f.OriginalName, &f.StoredName, &f.MimeType, &f.Size, &f.IsPublic, &f.DownloadCount, &f.SortOrder, &f.CreatedAt); err != nil {
			return nil, err
		}
		f.FormattedSize = FormatFileSize(f.Size)
		files = append(files, f)
	}
	if files == nil {
		files = []File{}
	}
	return files, nil
}

func GetPublicFiles() ([]File, error) {
	rows, err := DB.Query("SELECT id, category_id, original_name, stored_name, mime_type, size, is_public, download_count, sort_order, created_at FROM files WHERE is_public=1 ORDER BY sort_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.CategoryID, &f.OriginalName, &f.StoredName, &f.MimeType, &f.Size, &f.IsPublic, &f.DownloadCount, &f.SortOrder, &f.CreatedAt); err != nil {
			return nil, err
		}
		f.FormattedSize = FormatFileSize(f.Size)
		files = append(files, f)
	}
	if files == nil {
		files = []File{}
	}
	return files, nil
}

func GetFileByID(id int) (*File, error) {
	var f File
	err := DB.QueryRow("SELECT id, category_id, original_name, stored_name, mime_type, size, is_public, download_count, sort_order, created_at FROM files WHERE id=?", id).Scan(
		&f.ID, &f.CategoryID, &f.OriginalName, &f.StoredName, &f.MimeType, &f.Size, &f.IsPublic, &f.DownloadCount, &f.SortOrder, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	f.FormattedSize = FormatFileSize(f.Size)
	return &f, nil
}

func CreateFile(categoryID *int, originalName, storedName, mimeType string, size int64, isPublic bool) (File, error) {
	var maxOrder float64
	DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM files").Scan(&maxOrder)
	nextOrder := maxOrder + 1.0

	publicInt := 0
	if isPublic {
		publicInt = 1
	}

	result, err := DB.Exec("INSERT INTO files (category_id, original_name, stored_name, mime_type, size, is_public, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)",
		categoryID, originalName, storedName, mimeType, size, publicInt, nextOrder)
	if err != nil {
		return File{}, err
	}
	id, _ := result.LastInsertId()
	return File{ID: int(id), CategoryID: categoryID, OriginalName: originalName, StoredName: storedName, MimeType: mimeType, Size: size, IsPublic: isPublic, SortOrder: nextOrder}, nil
}

func IncrementDownloadCount(id int) error {
	_, err := DB.Exec("UPDATE files SET download_count = download_count + 1 WHERE id=?", id)
	return err
}

func DeleteFile(id int) error {
	_, err := DB.Exec("DELETE FROM files WHERE id=?", id)
	return err
}

func GenerateStoredName(originalName string) string {
	ext := filepath.Ext(originalName)
	return strings.ToLower(uuid.New().String()) + ext
}

func GetUploadDir() string {
	dir := os.Getenv("DATA_DIR")
	if dir == "" {
		dir = "data"
	}
	return filepath.Join(dir, "uploads")
}

func GetMaxFileSize() int64 {
	return 200 * 1024 * 1024 // 200MB
}

func FormatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
	}
	return fmt.Sprintf("%.1fGB", float64(size)/(1024*1024*1024))
}
