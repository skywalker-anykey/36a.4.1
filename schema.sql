/*
    Схема БД

// Post Публикация, получаемая из RSS.
type Post struct {
	ID      string `json:"guid,omitempty"`        // номер записи
	Title   string `json:"title,omitempty"`       // заголовок публикации
	Content string `json:"description,omitempty"` // содержание публикации
	PubTime int64  `json:"pubDate,omitempty"`     // время публикации
	Link    string `json:"link,omitempty"`        // ссылка на источник
}
*/

DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
    id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time BIGINT NOT NULL,
    link TEXT NOT NULL
);

-- Тестовые данные
--INSERT INTO logs (id, title, content, pub_time, link) VALUES
--();



