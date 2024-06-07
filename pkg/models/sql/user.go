package sql

import (
	"context"
	"fmt"
	"url-shortener/pkg/models"
)

// CreateTask создаёт новый задачу в хранилище.
func (s *Storage) CreateUser(email string, name string, login string, password string) (int, error) {
	const op = "storage.sql.CreateUser"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO user_(email, name, login, password) VALUES($1, $2, $3, $4) RETURNING id")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(email, name, login, password).Scan(&id)
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

// // GetTask получает задачу из хранилища по ID. Если ID не существует -
// // будет возвращена ошибка.
// func (s *Storage) GetUser(id int) (models.Post, error) {
// 	const op = "storage.sql.GetPost"
// 	resErr := models.Post{}

// 	ctx := context.Background()
// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return resErr, fmt.Errorf("%s: %w", op, err)
// 	}

// 	stmt, err := s.db.Prepare("SELECT * FROM post WHERE id = $1")
// 	if err != nil {
// 		tx.Rollback()
// 		return resErr, fmt.Errorf("%s: %w", op, err)
// 	}

// 	var res models.Post

// 	err = stmt.QueryRow(id).Scan(&res.ID, &res.Title, &res.Content, &res.Created)
// 	if err != nil {
// 		return resErr, fmt.Errorf("%s: execute statement: %w", op, err)
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return resErr, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return res, nil
// }

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) DeleteUser(id int) error {
	const op = "storage.sql.DeleteUser"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("DELETE FROM user_ WHERE id = $1")
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

// GetTask получает задачу из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) GetUserByLoginAndPassword(login string, password string) (models.User, error) {
	const op = "storage.sql.GetUserByLoginAndPassword"
	resErr := models.User{}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM user_ WHERE login = $1 AND password = $2")
	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	var res models.User

	err = stmt.QueryRow(login, password).Scan(&res.ID, &res.Email, &res.Login, &res.Password, &res.Name)
	if err != nil {
		return resErr, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
