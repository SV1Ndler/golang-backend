package sql

import (
	"context"
	"fmt"
	"time"
	"url-shortener/pkg/models"
)

// CreateTask создаёт новый задачу в хранилище.
func (s *Storage) CreateImage(image []byte, created time.Time) (int, error) {
	const op = "storage.sql.CreateImage"

	info, st, err := s.imgur_cli.UploadImage(image, "", "base64", "", "")
	if st != 200 || err != nil {
		//TODO
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO image(imgur_id, link, created) VALUES($1, $2, $3) RETURNING id")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(info.ID, info.Link, created).Scan(&id)
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

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) DeleteImage(id int) error {
	const op = "storage.sql.DeleteImage"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("DELETE FROM image WHERE id = $1")
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

func (s *Storage) GetImage(id int) (models.Image, error) {
	const op = "storage.sql.GetImage"
	resErr := models.Image{}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM image WHERE id = $1")
	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	var res models.Image

	err = stmt.QueryRow(id).Scan(&res.ID, &res.Imgur_ID, &res.Link, &res.Created)
	if err != nil {
		return resErr, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

// GetAllTasks возвращает из хранилища все задачи в произвольном порядке.
func (s *Storage) GetAllImages() ([]models.Image, error) {
	const op = "storage.sql.GetAllImages"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM image")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	allImages := make([]models.Image, 0, 1)
	rows, err := stmt.Query()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for id := 0; rows.Next(); id++ {
		var img models.Image

		if err = rows.Scan(&img.ID, &img.Imgur_ID, &img.Link, &img.Created); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		allImages = append(allImages, img)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allImages, nil
}
