package models

type ChatSession struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Model        string `json:"model"`
	MessageCount int    `json:"message_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ChatMessage struct {
	ID        int    `json:"id"`
	SessionID int    `json:"session_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func GetAllChatSessions() ([]ChatSession, error) {
	rows, err := DB.Query(`
		SELECT s.id, s.title, s.model, s.created_at, s.updated_at,
			COALESCE((SELECT COUNT(*) FROM chat_messages WHERE session_id = s.id), 0) as msg_count
		FROM chat_sessions s ORDER BY s.updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []ChatSession
	for rows.Next() {
		var s ChatSession
		if err := rows.Scan(&s.ID, &s.Title, &s.Model, &s.CreatedAt, &s.UpdatedAt, &s.MessageCount); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if sessions == nil {
		sessions = []ChatSession{}
	}
	return sessions, nil
}

func CreateChatSession(title, model string) (ChatSession, error) {
	if title == "" {
		title = "新对话"
	}
	result, err := DB.Exec("INSERT INTO chat_sessions (title, model) VALUES (?, ?)", title, model)
	if err != nil {
		return ChatSession{}, err
	}
	id, _ := result.LastInsertId()
	return ChatSession{ID: int(id), Title: title, Model: model}, nil
}

func UpdateChatSession(id int, title, model string) error {
	_, err := DB.Exec("UPDATE chat_sessions SET title=?, model=?, updated_at=datetime('now','localtime') WHERE id=?", title, model, id)
	return err
}

func UpdateChatSessionTime(id int) error {
	_, err := DB.Exec("UPDATE chat_sessions SET updated_at=datetime('now','localtime') WHERE id=?", id)
	return err
}

func DeleteChatSession(id int) error {
	_, err := DB.Exec("DELETE FROM chat_sessions WHERE id=?", id)
	return err
}

func GetChatSessionByID(id int) (*ChatSession, error) {
	var s ChatSession
	err := DB.QueryRow("SELECT id, title, model, created_at, updated_at FROM chat_sessions WHERE id=?", id).Scan(
		&s.ID, &s.Title, &s.Model, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetChatMessages(sessionID int) ([]ChatMessage, error) {
	rows, err := DB.Query("SELECT id, session_id, role, content, created_at FROM chat_messages WHERE session_id=? ORDER BY id ASC", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if messages == nil {
		messages = []ChatMessage{}
	}
	return messages, nil
}

func GetRecentChatMessages(sessionID int, limit int) ([]ChatMessage, error) {
	rows, err := DB.Query("SELECT id, session_id, role, content, created_at FROM chat_messages WHERE session_id=? ORDER BY id DESC LIMIT ?", sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	// Reverse to chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	if messages == nil {
		messages = []ChatMessage{}
	}
	return messages, nil
}

func CreateChatMessage(sessionID int, role, content string) (ChatMessage, error) {
	result, err := DB.Exec("INSERT INTO chat_messages (session_id, role, content) VALUES (?, ?, ?)", sessionID, role, content)
	if err != nil {
		return ChatMessage{}, err
	}
	id, _ := result.LastInsertId()
	// Update session timestamp
	UpdateChatSessionTime(sessionID)
	return ChatMessage{ID: int(id), SessionID: sessionID, Role: role, Content: content}, nil
}

func DeleteChatMessages(sessionID int) error {
	_, err := DB.Exec("DELETE FROM chat_messages WHERE session_id=?", sessionID)
	return err
}

func GetChatMessageCount(sessionID int) int {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM chat_messages WHERE session_id=?", sessionID).Scan(&count)
	return count
}

func TruncateChatMessages(sessionID int, keep int) error {
	_, err := DB.Exec(`
		DELETE FROM chat_messages WHERE session_id=? AND id NOT IN (
			SELECT id FROM chat_messages WHERE session_id=? ORDER BY id DESC LIMIT ?
		)`, sessionID, sessionID, keep)
	return err
}
