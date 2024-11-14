package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"hw36a.4.1/internal/conf"
	"hw36a.4.1/internal/rss"
)

// Store - хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// New - конструктор объекта хранилища.
func New(cBD *conf.BDConfig) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:%d/%s", cBD.User, cBD.Password, cBD.Port, cBD.Name))
	if err != nil {
		return nil, err
	}

	// Пересоздаем таблицу
	sql := fmt.Sprintf(`
DROP TABLE IF EXISTS %s;

CREATE TABLE %s (
	id TEXT NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	pub_time BIGINT NOT NULL,
	link TEXT NOT NULL,
CHECK((id !='') AND (title !='') AND (content !='') AND (pub_time !=0) AND (link !=''))

);
`, cBD.Table, cBD.Table)
	_, err = db.Exec(context.Background(), sql)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// AddPost - добавить новость в БД
func (s *Store) AddPost(p rss.Post) error {
	sql := `INSERT INTO posts (id, title, content, pub_time, link) VALUES ($1, $2, $3, $4, $5);`
	_, err := s.db.Exec(context.Background(), sql, &p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)

	return fmt.Errorf("ошибка добавления новости в БД: %s", err)
}

// Posts - вернуть n последних новостей
func (s *Store) Posts(n int) ([]rss.Post, error) {
	sql := `SELECT * FROM posts ORDER BY pub_time DESC LIMIT $1;`
	rows, err := s.db.Query(context.Background(), sql, n)

	if err != nil {
		return nil, err
	}

	var posts []rss.Post
	// Проходим по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var post rss.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, post)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}
