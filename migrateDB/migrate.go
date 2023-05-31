package migratedb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	path = "./migrateDB/migrations"
)

func CreateTable(db *sql.DB) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range dir {
		info, err := file.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", path, info.Name()))
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(data)); err != nil {
			log.Println("error razaq", err)
		}
	}
	return nil
}

func DropAllDB(db *sql.DB) error {
	records := `DROP TABLE IF EXISTS`

	tabls, err := SelectAllTable(db)
	if err != nil {
		return err
	}
	for _, table := range tabls {
		_, err := db.Exec(fmt.Sprintf("%s %s", records, table))
		if err != nil {
			return err
		}
	}
	return nil
}

func SelectAllTable(db *sql.DB) ([]string, error) {
	records := `SELECT name FROM sqlite_master WHERE type='table';`

	stmt, err := db.Prepare(records)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var tabls []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, err
		} else if table == "sqlite_sequence" {
			continue
		}

		tabls = append(tabls, table)
	}
	return tabls, nil
}
