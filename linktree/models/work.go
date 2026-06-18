package models

import "database/sql"

type Work struct {
	ID          int         `json:"id"`
	CategoryID  *int        `json:"category_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CoverURL    string      `json:"cover_url"`
	ContentType string      `json:"content_type"`
	ExternalURL string      `json:"external_url"`
	SortOrder   float64     `json:"sort_order"`
	CreatedAt   string      `json:"created_at"`
	Media       []WorkMedia `json:"media,omitempty"`
}

type WorkMedia struct {
	ID        int     `json:"id"`
	WorkID    int     `json:"work_id"`
	FileID    *int    `json:"file_id"`
	MediaURL  string  `json:"media_url"`
	MediaType string  `json:"media_type"`
	SortOrder float64 `json:"sort_order"`
}

func GetWorksByCategory(categoryID *int) ([]Work, error) {
	var rows *sql.Rows
	var err error
	if categoryID != nil {
		rows, err = DB.Query("SELECT id, category_id, title, description, cover_url, content_type, external_url, sort_order, created_at FROM works WHERE category_id=? ORDER BY sort_order ASC", *categoryID)
	} else {
		rows, err = DB.Query("SELECT id, category_id, title, description, cover_url, content_type, external_url, sort_order, created_at FROM works ORDER BY sort_order ASC")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []Work
	for rows.Next() {
		var w Work
		if err := rows.Scan(&w.ID, &w.CategoryID, &w.Title, &w.Description, &w.CoverURL, &w.ContentType, &w.ExternalURL, &w.SortOrder, &w.CreatedAt); err != nil {
			return nil, err
		}
		works = append(works, w)
	}
	if works == nil {
		works = []Work{}
	}
	return works, nil
}

func GetPublicWorks() ([]Work, error) {
	rows, err := DB.Query("SELECT id, category_id, title, description, cover_url, content_type, external_url, sort_order, created_at FROM works ORDER BY sort_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var works []Work
	for rows.Next() {
		var w Work
		if err := rows.Scan(&w.ID, &w.CategoryID, &w.Title, &w.Description, &w.CoverURL, &w.ContentType, &w.ExternalURL, &w.SortOrder, &w.CreatedAt); err != nil {
			return nil, err
		}
		works = append(works, w)
	}
	if works == nil {
		works = []Work{}
	}
	return works, nil
}

func GetWorkByID(id int) (*Work, error) {
	var w Work
	err := DB.QueryRow("SELECT id, category_id, title, description, cover_url, content_type, external_url, sort_order, created_at FROM works WHERE id=?", id).Scan(
		&w.ID, &w.CategoryID, &w.Title, &w.Description, &w.CoverURL, &w.ContentType, &w.ExternalURL, &w.SortOrder, &w.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func CreateWork(categoryID *int, title, description, coverURL, contentType, externalURL string) (Work, error) {
	var maxOrder float64
	DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM works").Scan(&maxOrder)
	nextOrder := maxOrder + 1.0

	result, err := DB.Exec("INSERT INTO works (category_id, title, description, cover_url, content_type, external_url, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)",
		categoryID, title, description, coverURL, contentType, externalURL, nextOrder)
	if err != nil {
		return Work{}, err
	}
	id, _ := result.LastInsertId()
	return Work{ID: int(id), CategoryID: categoryID, Title: title, Description: description, CoverURL: coverURL, ContentType: contentType, ExternalURL: externalURL, SortOrder: nextOrder}, nil
}

func UpdateWork(id int, categoryID *int, title, description, coverURL, contentType, externalURL string) error {
	_, err := DB.Exec("UPDATE works SET category_id=?, title=?, description=?, cover_url=?, content_type=?, external_url=? WHERE id=?",
		categoryID, title, description, coverURL, contentType, externalURL, id)
	return err
}

func DeleteWork(id int) error {
	_, err := DB.Exec("DELETE FROM works WHERE id=?", id)
	return err
}

func GetWorkMedia(workID int) ([]WorkMedia, error) {
	rows, err := DB.Query("SELECT id, work_id, file_id, media_url, media_type, sort_order FROM work_media WHERE work_id=? ORDER BY sort_order ASC", workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []WorkMedia
	for rows.Next() {
		var m WorkMedia
		if err := rows.Scan(&m.ID, &m.WorkID, &m.FileID, &m.MediaURL, &m.MediaType, &m.SortOrder); err != nil {
			return nil, err
		}
		media = append(media, m)
	}
	if media == nil {
		media = []WorkMedia{}
	}
	return media, nil
}

func AddWorkMedia(workID int, fileID *int, mediaURL, mediaType string) (WorkMedia, error) {
	var maxOrder float64
	DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM work_media WHERE work_id=?", workID).Scan(&maxOrder)
	nextOrder := maxOrder + 1.0

	result, err := DB.Exec("INSERT INTO work_media (work_id, file_id, media_url, media_type, sort_order) VALUES (?, ?, ?, ?, ?)",
		workID, fileID, mediaURL, mediaType, nextOrder)
	if err != nil {
		return WorkMedia{}, err
	}
	id, _ := result.LastInsertId()
	return WorkMedia{ID: int(id), WorkID: workID, FileID: fileID, MediaURL: mediaURL, MediaType: mediaType, SortOrder: nextOrder}, nil
}

func DeleteWorkMedia(id int) error {
	_, err := DB.Exec("DELETE FROM work_media WHERE id=?", id)
	return err
}

func LoadWorkMedia(works []Work) error {
	for i := range works {
		media, err := GetWorkMedia(works[i].ID)
		if err != nil {
			return err
		}
		works[i].Media = media
	}
	return nil
}
