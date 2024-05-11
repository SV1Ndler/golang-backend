package simple

import (
	"fmt"
	"sync"
	"url-shortener/pkg/models"
)

type PostImageMappingModel struct {
	sync.Mutex

	imageModel       *ImageModel
	postImageMapping map[int]models.PostImageMapping
	nextID           int
}

func NewPostImageMapping(imageModel *ImageModel) *PostImageMappingModel {
	m := &PostImageMappingModel{}
	m.postImageMapping = make(map[int]models.PostImageMapping)
	m.imageModel = imageModel
	m.nextID = 1

	return m
}

// CreateTask создаёт новый задачу в хранилище.
func (m *PostImageMappingModel) CreateMapping(imageID int, postID int) (int, error) {
	m.Lock()
	defer m.Unlock()

	mapping := models.PostImageMapping{
		ID:      m.nextID,
		ImageID: imageID,
		PostID:  postID}

	m.postImageMapping[m.nextID] = mapping
	m.nextID++
	return mapping.ID, nil
}

// GetTask получает задачу из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (m *PostImageMappingModel) GetMapping(id int) (models.PostImageMapping, error) {
	m.Lock()
	defer m.Unlock()

	item, ok := m.postImageMapping[id]
	if ok {
		return item, nil
	} else {
		return models.PostImageMapping{}, fmt.Errorf("task with id=%d not found", id)
	}
}

func (m *PostImageMappingModel) GetMappingWithImage(id int) (models.PostImageMapping, models.Image, error) {
	m.Lock()
	defer m.Unlock()

	item, ok := m.postImageMapping[id]
	img := m.imageModel.images[item.ImageID]
	if ok {
		return item, img, nil
	} else {
		return models.PostImageMapping{}, models.Image{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (m *PostImageMappingModel) DeleteMapping(id int) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.postImageMapping[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(m.postImageMapping, id)
	return nil
}

func (m *PostImageMappingModel) DeleteMappingByParams(postID int, imageID int) error {
	m.Lock()
	defer m.Unlock()

	for id := range m.postImageMapping {
		if m.postImageMapping[id].ImageID == imageID && m.postImageMapping[id].PostID == postID {
			delete(m.postImageMapping, id)
		}
	}

	// if _, ok := m.postImageMapping[id]; !ok {
	// 	return fmt.Errorf("task with id=%d not found", id)
	// }

	return nil
}

// DeleteAllTasks удаляет из хранилища все задачи.
// func (pm *PostImageMapping) DeleteAllPosts() error {
// 	pm.Lock()
// 	defer pm.Unlock()

// 	pm.posts = make(map[int]models.Post)
// 	return nil
// }

// GetAllTasks возвращает из хранилища все задачи в произвольном порядке.
func (m *PostImageMappingModel) GetAllMappings() ([]models.PostImageMapping, error) {
	m.Lock()
	defer m.Unlock()

	allMappings := make([]models.PostImageMapping, 0, len(m.postImageMapping))
	for _, item := range m.postImageMapping {
		allMappings = append(allMappings, item)
	}
	return allMappings, nil
}

func (m *PostImageMappingModel) GetAllMappingsWithLink() ([]models.PostImageMappingWithLink, error) {
	m.Lock()
	defer m.Unlock()

	allMappings := make([]models.PostImageMappingWithLink, 0, len(m.postImageMapping))
	for _, item := range m.postImageMapping {
		itemWithLink := models.PostImageMappingWithLink{
			ID: item.ID,
			PostID: item.PostID,
			ImageID: item.ImageID,
			Link: m.imageModel.images[item.ImageID].Link,
		}
		
		allMappings = append(allMappings, itemWithLink)
	}
	return allMappings, nil
}
