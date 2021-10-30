package database

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

type virtualMachineMapper struct {
	log        *logrus.Entry
	findByID   *sql.Stmt
	findByName *sql.Stmt
	insert     *sql.Stmt
	update     *sql.Stmt
	delete     *sql.Stmt
}

func getVirtualMachineMapper(db *sql.DB, log *logrus.Entry) (*virtualMachineMapper, error) {
	ret := &virtualMachineMapper{log: log}
	var err error

	ret.findByID, err = db.Prepare(`SELECT id, name, display_name, ram_size, network_id, network_ip, public_ip
		FROM virtual_machines WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'findByID' statement: ", err)
		return nil, err
	}

	ret.findByName, err = db.Prepare(`SELECT id, name, display_name, ram_size, network_id, network_ip, public_ip
		FROM virtual_machines WHERE name = ?`)
	if err != nil {
		log.Error("Cannot prepare 'findByName' statement: ", err)
		return nil, err
	}

	ret.insert, err = db.Prepare(`
		INSERT INTO virtual_machines(name, display_name, ram_size, network_id, network_ip, public_ip)
		VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Error("Cannot prepare 'insert' statement: ", err)
		return nil, err
	}

	ret.update, err = db.Prepare(`UPDATE virtual_machines
		SET display_name = ?, ram_size = ?, network_id = ?, network_ip = ?, public_ip = ?
		WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'update' statement: ", err)
		return nil, err
	}

	ret.delete, err = db.Prepare(`DELETE FROM virtual_machines WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'delete' statement: ", err)
		return nil, err
	}

	return ret, nil
}

func (m *virtualMachineMapper) FindByID(id int64) (*entities.VirtualMachine, error) {
	ret := &entities.VirtualMachine{}
	m.log.WithField("id", id).Trace("Finding by ID")
	err := m.findByID.
		QueryRow(id).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.RAMSizeMB, &ret.NetworkID, &ret.NetworkIP, &ret.PublicIP)

	return ret, err
}

func (m *virtualMachineMapper) FindByName(name string) (*entities.VirtualMachine, error) {
	ret := &entities.VirtualMachine{}
	m.log.WithField("name", name).Trace("Finding by name")
	err := m.findByName.
		QueryRow(name).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.RAMSizeMB, &ret.NetworkID, &ret.NetworkIP, &ret.PublicIP)

	return ret, err
}

func (m *virtualMachineMapper) Save(evm *entities.VirtualMachine) (*entities.VirtualMachine, error) {
	var err error
	if evm.ID > 0 {
		m.log.WithField("id", evm.ID).Trace("Updating")
		err = handleChangeResult(
			m.update.Exec(evm.DisplayName, evm.RAMSizeMB, evm.NetworkID, evm.NetworkIP, evm.PublicIP, evm.ID),
		)
	} else {
		m.log.WithField("name", evm.Name).Trace("Inserting")
		evm.ID, err = handleInsertResult(
			m.insert.Exec(evm.Name, evm.DisplayName, evm.RAMSizeMB, evm.NetworkID, evm.NetworkIP, evm.PublicIP),
		)
	}

	return evm, err
}

func (m *virtualMachineMapper) Delete(id int64) error {
	m.log.WithField("id", id).Trace("Deleting")
	return handleChangeResult(m.delete.Exec(id))
}
