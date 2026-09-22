package models

type Experience struct {
	ID          int     `json:"id"`
	Period      string  `json:"period"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	SortOrder   float64 `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
}

func GetAllExperiences() ([]Experience, error) {
	rows, err := DB.Query("SELECT id, period, title, description, sort_order, created_at FROM experiences ORDER BY sort_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Experience
	for rows.Next() {
		var e Experience
		if err := rows.Scan(&e.ID, &e.Period, &e.Title, &e.Description, &e.SortOrder, &e.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	if list == nil {
		list = []Experience{}
	}
	return list, nil
}

func CreateExperience(period, title, description string) (Experience, error) {
	var maxOrder float64
	DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM experiences").Scan(&maxOrder)
	nextOrder := maxOrder + 1.0

	result, err := DB.Exec("INSERT INTO experiences (period, title, description, sort_order) VALUES (?, ?, ?, ?)",
		period, title, description, nextOrder)
	if err != nil {
		return Experience{}, err
	}
	id, _ := result.LastInsertId()
	return Experience{ID: int(id), Period: period, Title: title, Description: description, SortOrder: nextOrder}, nil
}

func UpdateExperience(id int, period, title, description string) error {
	_, err := DB.Exec("UPDATE experiences SET period=?, title=?, description=? WHERE id=?",
		period, title, description, id)
	return err
}

func DeleteExperience(id int) error {
	_, err := DB.Exec("DELETE FROM experiences WHERE id=?", id)
	return err
}

func SeedDefaultExperiences() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM experiences").Scan(&count)
	if count > 0 {
		return
	}
	DB.Exec("INSERT INTO experiences (period, title, description, sort_order) VALUES ('2024 — 至今', '全栈开发者 · 独立开发者', '主导开发多款 AI 应用与自动化工具，负责前后端与部署运维。', 1)")
	DB.Exec("INSERT INTO experiences (period, title, description, sort_order) VALUES ('2022 — 2024', '后端开发者', '负责 API 设计与服务端开发，熟悉高并发与数据存储方案。', 2)")
	DB.Exec("INSERT INTO experiences (period, title, description, sort_order) VALUES ('2018 — 2022', '计算机科学与技术 · 本科', '计算机专业学习，打下扎实的编程与算法基础。', 3)")
}