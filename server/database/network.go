package database

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/arxeiss/sample-terraform-provider/entities"
)

type networkMapper struct {
	log        *logrus.Entry
	findByID   *sql.Stmt
	findByName *sql.Stmt
	insert     *sql.Stmt
	update     *sql.Stmt
	delete     *sql.Stmt
}

func getNetworkMapper(db *sql.DB, log *logrus.Entry) (*networkMapper, error) {
	ret := &networkMapper{log: log}
	var err error

	ret.findByID, err = db.Prepare(`SELECT id, name, display_name, ip_range, use_dhcp
		FROM networks WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'findByID' statement: ", err)
		return nil, err
	}

	ret.findByName, err = db.Prepare(`SELECT id, name, display_name, ip_range, use_dhcp
		FROM networks WHERE name = ?`)
	if err != nil {
		log.Error("Cannot prepare 'findByName' statement: ", err)
		return nil, err
	}

	ret.insert, err = db.Prepare(`INSERT INTO networks(name, display_name, ip_range, use_dhcp)
		VALUES(?, ?, ?, ?)`)
	if err != nil {
		log.Error("Cannot prepare 'insert' statement: ", err)
		return nil, err
	}

	ret.update, err = db.Prepare(`UPDATE networks
		SET display_name = ?, ip_range = ?, use_dhcp = ?
		WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'update' statement: ", err)
		return nil, err
	}

	ret.delete, err = db.Prepare(`DELETE FROM networks WHERE id = ?`)
	if err != nil {
		log.Error("Cannot prepare 'delete' statement: ", err)
		return nil, err
	}

	return ret, nil
}

func (m *networkMapper) FindByID(id int64) (*entities.Network, error) {
	ret := &entities.Network{}
	m.log.WithField("id", id).Trace("Finding by ID")
	err := m.findByID.
		QueryRow(id).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.IPRange, &ret.UseDHCP)

	return ret, err
}

func (m *networkMapper) FindByName(name string) (*entities.Network, error) {
	ret := &entities.Network{}
	m.log.WithField("name", name).Trace("Finding by name")
	err := m.findByName.
		QueryRow(name).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.IPRange, &ret.UseDHCP)

	return ret, err
}

func (m *networkMapper) Save(en *entities.Network) (*entities.Network, error) {
	var err error
	if en.ID > 0 {
		m.log.WithField("id", en.ID).Trace("Updating")
		err = handleChangeResult(m.update.Exec(en.DisplayName, en.IPRange, en.UseDHCP, en.ID))
	} else {
		m.log.WithField("name", en.Name).Trace("Inserting")
		en.ID, err = handleInsertResult(m.insert.Exec(en.Name, en.DisplayName, en.IPRange, en.UseDHCP))
	}

	return en, err
}

func (m *networkMapper) Delete(id int64) error {
	m.log.WithField("id", id).Trace("Deleting")
	return handleChangeResult(m.delete.Exec(id))
}
