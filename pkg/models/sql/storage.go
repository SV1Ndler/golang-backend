package sql

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/koffeinsource/go-imgur"
)

type Storage struct {
	db        *sql.DB
	imgur_cli *imgur.Client
}

func dbURL() string {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
}

func New() (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("pgx", dbURL())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS post(
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  		title TEXT NOT NULL,
  		content TEXT NOT NULL,
  		created TIMESTAMP WITHOUT TIME ZONE NOT NULL);

	CREATE TABLE IF NOT EXISTS image (
  		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  		imgur_id TEXT NOT NULL,
  		link TEXT NOT NULL,
  		created TIMESTAMP WITHOUT TIME ZONE NOT NULL);

	CREATE TABLE IF NOT EXISTS  post_image_mapping (
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		image_id INTEGER NOT NULL REFERENCES Image(id) ON DELETE CASCADE,
		post_id INTEGER NOT NULL  REFERENCES Post(id) ON DELETE CASCADE,
		UNIQUE (image_id, post_id));

	CREATE TABLE IF NOT EXISTS  user_ (
		id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		email TEXT NOT NULL,
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE);
	`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	imgur_cli, err := imgur.NewClient(new(http.Client), "05b624614c5f987", "") // TODO
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db, imgur_cli: imgur_cli}, nil
}
