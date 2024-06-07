package sql

import (
	"context"
	"fmt"
	"time"
	"url-shortener/pkg/models"
)

// CreateTask создаёт новый задачу в хранилище.
func (s *Storage) CreatePost(title string, content string, created time.Time) (int, error) {
	const op = "storage.sql.CreatePost"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO post(title, content, created) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(title, content, created).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

// GetPost получает пост из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) GetPost(id int) (models.Post, error) {
	const op = "storage.sql.GetPost"
	resErr := models.Post{}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM post WHERE id = $1")
	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	var res models.Post

	err = stmt.QueryRow(id).Scan(&res.ID, &res.Title, &res.Content, &res.Created)
	if err != nil {
		return resErr, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

// DeletePost удаляет пост с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) DeletePost(id int) error {
	const op = "storage.sql.DeletePost"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("DELETE FROM post WHERE id = $1")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetAllPosts возвращает из хранилища все посты в произвольном порядке.
func (s *Storage) GetAllPosts() ([]models.Post, error) {
	const op = "storage.sql.GetAllPosts"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM post")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	allPosts := make([]models.Post, 0, 1)
	rows, err := stmt.Query()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for id := 0; rows.Next(); id++ {
		var post models.Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Created); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		allPosts = append(allPosts, post)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allPosts, nil
}

func (s *Storage) UpdatePost(id int, title string, content string) (int, error) {
	const op = "storage.sql.UpdatePost"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare(`
	UPDATE post SET title = $1, content = $2
		WHERE id = $3
		RETURNING id`)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRow(title, content, id).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}
