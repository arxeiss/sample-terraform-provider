package database

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

type storageMapper struct {
	log        *logrus.Entry
	findByID   *sql.Stmt
	findByName *sql.Stmt
	insert     *sql.Stmt
	update     *sql.Stmt
	delete     *sql.Stmt
}

func getStorageMapper(db *sql.DB, log *logrus.Entry) (*storageMapper, error) {
	ret := &storageMapper{log: log}
	var err error

	ret.findByID, err = db.Prepare(`
		SELECT id, name, display_name, size, network_id, network_ip, virtual_machine_id, mount_path
		FROM storages WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'findByID' statement: ", err)
		return nil, err
	}

	ret.findByName, err = db.Prepare(`
		SELECT id, name, display_name, size, network_id, network_ip, virtual_machine_id, mount_path
		FROM storages WHERE name = ?`)
	if err != nil {
		log.Error("Cannot prepare 'findByName' statement: ", err)
		return nil, err
	}

	ret.insert, err = db.Prepare(`
		INSERT INTO storages(name, display_name, size, network_id, network_ip, virtual_machine_id, mount_path)
		VALUES(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Error("Cannot prepare 'insert' statement: ", err)
		return nil, err
	}

	ret.update, err = db.Prepare(`UPDATE storages
		SET display_name = ?, size = ?, network_id = ?, network_ip = ?, virtual_machine_id = ?, mount_path = ?
		WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'update' statement: ", err)
		return nil, err
	}

	ret.delete, err = db.Prepare(`DELETE FROM storages WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'delete' statement: ", err)
		return nil, err
	}

	return ret, nil
}

func (m *storageMapper) FindByID(id int64) (*entities.Storage, error) {
	ret := &entities.Storage{}
	m.log.WithField("id", id).Trace("Finding by ID")
	err := m.findByID.
		QueryRow(id).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.SizeMB,
			&ret.NetworkID, &ret.NetworkID, &ret.VirtualMachineID, &ret.MountPath)

	return ret, err
}

func (m *storageMapper) FindByName(name string) (*entities.Storage, error) {
	ret := &entities.Storage{}
	m.log.WithField("name", name).Trace("Finding by name")
	err := m.findByName.
		QueryRow(name).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.SizeMB,
			&ret.NetworkID, &ret.NetworkID, &ret.VirtualMachineID, &ret.MountPath)

	return ret, err
}

func (m *storageMapper) Save(es *entities.Storage) (*entities.Storage, error) {
	var err error
	if es.ID > 0 {
		m.log.WithField("id", es.ID).Trace("Updating")
		err = handleChangeResult(m.update.Exec(
			es.DisplayName, es.SizeMB, es.NetworkID, es.NetworkIP, es.VirtualMachineID, es.MountPath, es.ID,
		))
	} else {
		m.log.WithField("name", es.Name).Trace("Inserting")
		es.ID, err = handleInsertResult(m.insert.Exec(
			es.Name, es.DisplayName, es.SizeMB, es.NetworkID, es.NetworkIP, es.VirtualMachineID, es.MountPath,
		))
	}

	return es, err
}

func (m *storageMapper) Delete(id int64) error {
	m.log.WithField("id", id).Trace("Deleting")
	return handleChangeResult(m.delete.Exec(id))
}
