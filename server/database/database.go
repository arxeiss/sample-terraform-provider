package database

import (
	"database/sql"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Must have
	"github.com/sirupsen/logrus"
)

type DB struct {
	db              *sql.DB // Main SQLite object
	VirtualMachines *virtualMachineMapper
	Storages        *storageMapper
	Networks        *networkMapper
}

// Open dbFile with SQLite3 driver
// Schema is used only, if dbFile is not exists
func Open(dbFile string, schemaFile string, log *logrus.Entry) (*DB, error) {
	createSchema := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		log.Info("DB file does not exists, will import schema")
		createSchema = true
	}
	db, err := sql.Open("sqlite3", "file:"+dbFile+"?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}
	log.Trace("Open DB connection succeeded")
	if createSchema {
		log.Trace("Reading schema")
		path, err := filepath.Abs(schemaFile)
		if err != nil {
			return nil, err
		}
		file, err := ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		log.Trace("Importing schema")
		_, err = db.Exec(string(file))
		if err != nil {
			return nil, err
		}
	}

	log.Trace("Turning on foreign keys")
	if _, err = db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	log.Trace("Getting VM mapper")
	vms, err := getVirtualMachineMapper(db, log.WithField("entity", "vm"))
	if err != nil {
		return nil, err
	}

	log.Trace("Getting Storage mapper")
	storages, err := getStorageMapper(db, log.WithField("entity", "storage"))
	if err != nil {
		return nil, err
	}

	log.Trace("Getting Network mapper")
	networks, err := getNetworkMapper(db, log.WithField("entity", "network"))
	if err != nil {
		return nil, err
	}

	return &DB{
		db:              db,
		VirtualMachines: vms,
		Storages:        storages,
		Networks:        networks,
	}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
