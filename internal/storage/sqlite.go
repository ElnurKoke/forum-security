package storage

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "forum.db?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	migrations := []string{"userTable.sql",
		"postTable.sql",
		"categoriesTable.sql",
		"commentTable.sql",
		"reactionCommentTable.sql",
		"reactionPostTable.sql",
		"notificationTable.sql"}
	for _, migrationFile := range migrations {
		content, err := ioutil.ReadFile(filepath.Join("migrations", migrationFile))
		if err != nil {
			log.Fatal(err)
		}
		if _, err = db.Exec(string(content)); err != nil {
			log.Fatal(err)
		}
	}
	return db
}
