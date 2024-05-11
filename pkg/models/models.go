package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Post struct {
	ID      int
	Title   string
	Content string
	Created time.Time
}

type Image struct {
	ID       int
	Imgur_ID string
	Link     string
	Created  time.Time
}

type PostImageMapping struct {
	ID      int
	ImageID int
	PostID  int
}

type PostImageMappingWithLink struct {
	ID      int
	ImageID int
	PostID  int
	Link    string
}
