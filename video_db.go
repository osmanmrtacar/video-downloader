package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type VideoDB struct {
	DB *sql.DB
}

func NewVideoDB(path string) (*VideoDB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT UNIQUE,
		filename TEXT,
		description TEXT
	)`); err != nil {
		return nil, err
	}
	return &VideoDB{DB: db}, nil
}

func (vdb *VideoDB) GetVideoMetadataByUrl(url string) (filename, description string, err error) {
	row := vdb.DB.QueryRow("SELECT filename, description FROM videos WHERE url = ? LIMIT 1", url)
	err = row.Scan(&filename, &description)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	return filename, description, err
}

func (vdb *VideoDB) Save(url, filename, description string) error {
	_, err := vdb.DB.Exec("INSERT OR IGNORE INTO videos (url, filename, description) VALUES (?, ?, ?)", url, filename, description)
	return err
}
