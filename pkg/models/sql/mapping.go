package sql

import (
	"context"
	"fmt"
	"url-shortener/pkg/models"
)

// type PostImageMappingModel struct {
// 	sync.Mutex

// 	imageModel       *ImageModel
// 	postImageMapping map[int]models.PostImageMapping
// 	nextID           int
// }

// func NewPostImageMapping(imageModel *ImageModel) *PostImageMappingModel {
// 	m := &PostImageMappingModel{}
// 	m.postImageMapping = make(map[int]models.PostImageMapping)
// 	m.imageModel = imageModel
// 	m.nextID = 1

// 	return m
// }

// CreateTask создаёт новый задачу в хранилище.
func (s *Storage) CreateMapping(imageID int, postID int) (int, error) {
	const op = "storage.sql.CreateMapping"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO post_image_mapping(image_id, post_id) VALUES($1, $2) RETURNING id")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int
	err = stmt.QueryRow(imageID, postID).Scan(&id)
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

// GetTask получает задачу из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) GetMapping(id int) (models.PostImageMapping, error) {
	const op = "storage.sql.GetMapping"
	resErr := models.PostImageMapping{}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM post_image_mapping WHERE id = $1")
	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	var res models.PostImageMapping

	err = stmt.QueryRow(id).Scan(&res.ID, &res.ImageID, &res.PostID)
	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	return resErr, nil
}

func (s *Storage) GetMappingWithLink(id int) (models.PostImageMappingWithLink, error) {
	const op = "storage.sql.GetMappingWithLink"
	resErr := models.PostImageMappingWithLink{}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare(`
	SELECT m.*, image.Link FROM post_image_mapping AS m 
	JOIN image ON m.image_id = image.id
	WHERE m.id = $1;
	`)
	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	var res models.PostImageMappingWithLink

	err = stmt.QueryRow(id).Scan(&res.ID, &res.ImageID, &res.PostID, &res.Link)
	if err != nil {
		return resErr, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err != nil {
		tx.Rollback()
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return resErr, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (s *Storage) DeleteMapping(id int) error {
	const op = "storage.sql.DeleteMapping"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("DELETE FROM post_image_mapping WHERE id = $1")
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

func (s *Storage) DeleteMappingByParams(postID int, imageID int) error {
	const op = "storage.sql.DeleteMappingByParams"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("DELETE FROM post_image_mapping WHERE post_id = $1 AND image_id = $2;")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(postID, imageID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetAllTasks возвращает из хранилища все задачи в произвольном порядке.
func (s *Storage) GetAllMappings() ([]models.PostImageMapping, error) {
	const op = "storage.sql.GetAllMappings"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare(`
	SELECT m.*, image.Link FROM post_image_mapping AS m 
	JOIN image ON m.image_id = image.id;
	`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	allMappings := make([]models.PostImageMapping, 0, 1)
	rows, err := stmt.Query()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for id := 0; rows.Next(); id++ {
		var mapping models.PostImageMapping

		if err = rows.Scan(&mapping.ID, &mapping.ImageID, &mapping.PostID); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		allMappings = append(allMappings, mapping)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allMappings, nil
}

func (s *Storage) GetAllMappingsWithLink() ([]models.PostImageMappingWithLink, error) {
	const op = "storage.sql.GetAllMappingsWithLink"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare("SELECT * FROM post_image_mapping AS m JOIN image ON m.image_id = image.id")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	allMappings := make([]models.PostImageMappingWithLink, 0, 1)
	rows, err := stmt.Query()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for id := 0; rows.Next(); id++ {
		var mapping models.PostImageMappingWithLink
		if err = rows.Scan(&mapping.ID, &mapping.ImageID, &mapping.PostID, &mapping.Link); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		allMappings = append(allMappings, mapping)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return allMappings, nil
}

func (s *Storage) GetPostImages(postID int) ([]models.Image, error) {
	const op = "storage.sql.GetPostImages"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := s.db.Prepare(`
	SELECT image.* FROM post_image_mapping AS m 
	JOIN image ON m.image_id = image.id
	WHERE m.post_id = $1
	`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	allImages := make([]models.Image, 0, 1)
	rows, err := stmt.Query(postID)
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
