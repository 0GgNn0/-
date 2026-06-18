package models

type Category struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	SortOrder float64 `json:"sort_order"`
}

func GetCategoriesByType(categoryType string) ([]Category, error) {
	rows, err := DB.Query("SELECT id, name, type, sort_order FROM categories WHERE type=? ORDER BY sort_order ASC", categoryType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.SortOrder); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if categories == nil {
		categories = []Category{}
	}
	return categories, nil
}

func CreateCategory(name, categoryType string) (Category, error) {
	var maxOrder float64
	DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM categories WHERE type=?", categoryType).Scan(&maxOrder)
	nextOrder := maxOrder + 1.0

	result, err := DB.Exec("INSERT INTO categories (name, type, sort_order) VALUES (?, ?, ?)", name, categoryType, nextOrder)
	if err != nil {
		return Category{}, err
	}
	id, _ := result.LastInsertId()
	return Category{ID: int(id), Name: name, Type: categoryType, SortOrder: nextOrder}, nil
}

func UpdateCategory(id int, name string, sortOrder float64) error {
	_, err := DB.Exec("UPDATE categories SET name=?, sort_order=? WHERE id=?", name, sortOrder, id)
	return err
}

func DeleteCategory(id int) error {
	_, err := DB.Exec("DELETE FROM categories WHERE id=?", id)
	return err
}
