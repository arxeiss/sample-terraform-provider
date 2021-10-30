package database

import "database/sql"

type VirtualMachine struct {
	ID            int64
	Name          string
	DisplayName   *string
	RAMSize       int
	NetworkID     *int
	NetworkIPAddr *string
	PublicIPAddr  *string
}

type virtualMachineMapper struct {
	findByID   *sql.Stmt
	findByName *sql.Stmt
	insert     *sql.Stmt
	update     *sql.Stmt
	delete     *sql.Stmt
}

func getVirtualMachineMapper(db *sql.DB) (*virtualMachineMapper, error) {
	ret := &virtualMachineMapper{}
	var err error

	ret.findByID, err = db.Prepare(`SELECT id, name, display_name, ram_size, network_id, network_ip, public_ip
		FROM virtual_machines WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	ret.findByName, err = db.Prepare(`SELECT id, name, display_name, ram_size, network_id, network_ip, public_ip
		FROM virtual_machines WHERE name = ?`)
	if err != nil {
		return nil, err
	}

	ret.insert, err = db.Prepare(`
		INSERT INTO virtual_machines(name, display_name, ram_size, network_id, network_ip, public_ip)
		VALUES(?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	ret.update, err = db.Prepare(`UPDATE virtual_machines
		SET display_name = ?, ram_size = ?, network_id = ?, network_ip = ?, public_ip = ?
		WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	ret.delete, err = db.Prepare(`DELETE FROM virtual_machines WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *virtualMachineMapper) FindByID(id int) (*VirtualMachine, error) {
	ret := &VirtualMachine{}
	err := m.findByID.
		QueryRow(id).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.RAMSize, &ret.NetworkID, &ret.NetworkIPAddr, &ret.PublicIPAddr)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return ret, err
}

func (m *virtualMachineMapper) FindByName(name string) (*VirtualMachine, error) {
	ret := &VirtualMachine{}
	err := m.findByID.
		QueryRow(name).
		Scan(&ret.ID, &ret.Name, &ret.DisplayName, &ret.RAMSize, &ret.NetworkID, &ret.NetworkIPAddr, &ret.PublicIPAddr)

	return ret, err
}

func (m *virtualMachineMapper) Save(vm *VirtualMachine) (*VirtualMachine, error) {
	var err error
	if vm.ID > 0 {
		err = handleChangeResult(
			m.update.Exec(vm.DisplayName, vm.RAMSize, vm.NetworkID, vm.NetworkIPAddr, vm.PublicIPAddr, vm.ID),
		)
	} else {
		vm.ID, err = handleInsertResult(
			m.insert.Exec(vm.Name, vm.DisplayName, vm.RAMSize, vm.NetworkID, vm.NetworkIPAddr, vm.PublicIPAddr),
		)
	}

	return vm, err
}
