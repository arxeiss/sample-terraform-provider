package database

import "database/sql"

type Storage struct {
	ID          int
	Name        string
	DisplayName *string
	Size        int
	Type        string
}

type storageMapper struct {
	findByID *sql.Stmt
}

func getStorageMapper(db *sql.DB) (*storageMapper, error) {
	ret := &storageMapper{}
	var err error

	ret.findByID, err = db.Prepare(`SELECT id, name, display_name, size, type
		FROM virtual_machines WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (vmm *storageMapper) FindByID(id int) (*Storage, error) {
	ret := &Storage{}
	err := vmm.findByID.
		QueryRow(id).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.Size, &ret.Type)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ret, err
}
