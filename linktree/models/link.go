package models

import "database/sql"

type Link struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	SortOrder float64 `json:"sort_order"`
	Icon      string  `json:"icon"`
}

type OrderEntry struct {
	ID        int     `json:"id"`
	SortOrder float64 `json:"sort_order"`
}

func GetAllLinks() ([]Link, error) {
	rows, err := DB.Query("SELECT id, title, url, sort_order, icon FROM links ORDER BY sort_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.Title, &l.URL, &l.SortOrder, &l.Icon); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	if links == nil {
		links = []Link{}
	}
	return links, nil
}

func CreateLink(title, url string) (Link, error) {
	var maxOrder sql.NullFloat64
	DB.QueryRow("SELECT MAX(sort_order) FROM links").Scan(&maxOrder)
	nextOrder := 1.0
	if maxOrder.Valid {
		nextOrder = maxOrder.Float64 + 1.0
	}

	result, err := DB.Exec("INSERT INTO links (title, url, sort_order) VALUES (?, ?, ?)", title, url, nextOrder)
	if err != nil {
		return Link{}, err
	}
	id, _ := result.LastInsertId()
	return Link{ID: int(id), Title: title, URL: url, SortOrder: nextOrder}, nil
}

func UpdateLink(id int, title, url string, sortOrder float64) error {
	_, err := DB.Exec("UPDATE links SET title=?, url=?, sort_order=? WHERE id=?", title, url, sortOrder, id)
	return err
}

func DeleteLink(id int) error {
	_, err := DB.Exec("DELETE FROM links WHERE id=?", id)
	return err
}

func ReorderLinks(orders []OrderEntry) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, o := range orders {
		if _, err := tx.Exec("UPDATE links SET sort_order=? WHERE id=?", o.SortOrder, o.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func SeedDefaultLinks() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM links").Scan(&count)
	if count == 0 {
		DB.Exec("INSERT INTO links (title, url, sort_order) VALUES (?, ?, 1.0)", "CSDN", "https://blog.csdn.net/")
		DB.Exec("INSERT INTO links (title, url, sort_order) VALUES (?, ?, 2.0)", "GitHub", "https://github.com/")
		DB.Exec("INSERT INTO links (title, url, sort_order) VALUES (?, ?, 3.0)", "中转站", "http://62.234.92.232:3000")
	}
}
