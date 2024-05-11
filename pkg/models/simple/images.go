package simple

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"url-shortener/pkg/models"

	"github.com/koffeinsource/go-imgur"
)

type ImageModel struct {
	sync.Mutex
	imgur_cli *imgur.Client

	images map[int]models.Image
	nextID int
}

func NewImage(clientID string) (*ImageModel, error) {
	var err error
	im := &ImageModel{}

	im.images = make(map[int]models.Image)
	im.nextID = 1
	im.imgur_cli, err = imgur.NewClient(new(http.Client), clientID, "")

	return im, err
}

// CreateTask создаёт новый задачу в хранилище.
func (im *ImageModel) CreateImage(image []byte, created time.Time) (int, error) {
	im.Lock()
	defer im.Unlock()

	info, st, err := im.imgur_cli.UploadImage(image, "", "base64", "", "")
	if st != 200 || err != nil {
		//TODO
		return 0, err
	}

	item := models.Image{
		ID:       im.nextID,
		Imgur_ID: info.ID,
		Link:     info.Link,
		Created:  created}

	im.images[im.nextID] = item
	im.nextID++
	return item.ID, nil
}

// GetTask получает задачу из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (im *ImageModel) GetImage(id int) (models.Image, error) {
	im.Lock()
	defer im.Unlock()

	i, ok := im.images[id]
	if ok {
		return i, nil
	} else {
		return models.Image{}, fmt.Errorf("image with id=%d not found", id)
	}
}

// DeleteTask удаляет задачу с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (im *ImageModel) DeleteImage(id int) error {
	im.Lock()
	defer im.Unlock()

	if _, ok := im.images[id]; !ok {
		return fmt.Errorf("image with id=%d not found", id)
	}

	delete(im.images, id)
	return nil
}

// // DeleteAllTasks удаляет из хранилища все задачи.
// func (pm *PostModel) DeleteAllPosts() error {
// 	pm.Lock()
// 	defer pm.Unlock()

// 	pm.posts = make(map[int]models.Post)
// 	return nil
// }

// GetAllTasks возвращает из хранилища все задачи в произвольном порядке.
func (im *ImageModel) GetAllImages() ([]models.Image, error) {
	im.Lock()
	defer im.Unlock()

	allImages := make([]models.Image, 0, len(im.images))
	for _, image := range im.images {
		allImages = append(allImages, image)
	}
	return allImages, nil
}
