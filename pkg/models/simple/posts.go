package simple

import (
	"fmt"
	"sync"
	"time"
	"url-shortener/pkg/models"
)

type PostModel struct {
	sync.Mutex

	posts  map[int]models.Post
	nextID int
}

func NewPost() *PostModel {
	pm := &PostModel{}
	pm.posts = make(map[int]models.Post)
	pm.nextID = 1

	return pm
}

// CreateTask создаёт новый задачу в хранилище.
func (pm *PostModel) CreatePost(title string, content string, created time.Time) (int, error) {
	pm.Lock()
	defer pm.Unlock()

	post := models.Post{
		ID:      pm.nextID,
		Title:   title,
		Content: content,
		Created: created}

	pm.posts[pm.nextID] = post
	pm.nextID++
	return post.ID, nil
}

// GetTask получает задачу из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (pm *PostModel) GetPost(id int) (models.Post, error) {
	pm.Lock()
	defer pm.Unlock()

	p, ok := pm.posts[id]
	if ok {
		return p, nil
	} else {
		return models.Post{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (pm *PostModel) DeletePost(id int) error {
	pm.Lock()
	defer pm.Unlock()

	if _, ok := pm.posts[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(pm.posts, id)
	return nil
}

// DeleteAllTasks удаляет из хранилища все задачи.
func (pm *PostModel) DeleteAllPosts() error {
	pm.Lock()
	defer pm.Unlock()

	pm.posts = make(map[int]models.Post)
	return nil
}

// GetAllTasks возвращает из хранилища все задачи в произвольном порядке.
func (pm *PostModel) GetAllPosts() ([]models.Post, error) {
	pm.Lock()
	defer pm.Unlock()

	allPosts := make([]models.Post, 0, len(pm.posts))
	for _, post := range pm.posts {
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}
