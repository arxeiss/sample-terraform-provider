package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Must have
)

type DB struct {
	db              *sql.DB // Main SQLite object
	VirtualMachines *virtualMachineMapper
	Storages        *storageMapper
}

// Open dbFile with SQLite3 driver
// Schema is used only, if dbFile is not exists
func Open(dbFile string, schemaFile string) (*DB, error) {
	createSchema := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		createSchema = true
	}
	db, err := sql.Open("sqlite3", "file:"+dbFile+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}
	if createSchema {
		file, err := ioutil.ReadFile(filepath.Clean(schemaFile))
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(string(file))
		if err != nil {
			return nil, err
		}
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	vms, err := getVirtualMachineMapper(db)
	if err != nil {
		return nil, err
	}

	storages, err := getStorageMapper(db)
	if err != nil {
		return nil, err
	}

	return &DB{
		db:              db,
		VirtualMachines: vms,
		Storages:        storages,
	}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
