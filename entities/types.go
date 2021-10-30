package entities

import "errors"

type VirtualMachine struct {
	ID          int64   `json:"id,omitempty"`
	Name        string  `json:"name"`
	DisplayName *string `json:"display_name,omitempty"`
	RAMSizeMB   int     `json:"ram_size_mb"`
	NetworkID   *int    `json:"network_id,omitempty"`
	NetworkIP   *string `json:"network_ip,omitempty"`
	PublicIP    *string `json:"public_ip,omitempty"`
}

func (vm VirtualMachine) Validate() error {
	if vm.ID < 0 {
		return errors.New("id cannot be negative number")
	}
	if err := ValidateName(vm.Name); err != nil {
		return err
	}
	if err := ValidateDisplayName(vm.DisplayName); err != nil {
		return err
	}

	if vm.RAMSizeMB < 100 || vm.RAMSizeMB > 100*1024 {
		return errors.New("ram_size_mb must be between 100 MB and 100 GB")
	}
	if vm.NetworkID != nil && *vm.NetworkID < 0 {
		return errors.New("network_id cannot be negative number")
	}
	if vm.NetworkIP != nil {
		if !ValidIP(*vm.NetworkIP) {
			return errors.New("network_ip is not valid")
		}
		if vm.NetworkID == nil {
			return errors.New("cannot set network_ip when network_id is null")
		}
	}
	if vm.PublicIP != nil && !ValidIP(*vm.PublicIP) {
		return errors.New("public_ip is not valid")
	}
	return nil
}

type Storage struct {
	ID               int64   `json:"id,omitempty"`
	Name             string  `json:"name"`
	DisplayName      *string `json:"display_name,omitempty"`
	SizeMB           int     `json:"size_mb"`
	NetworkID        *int    `json:"network_id,omitempty"`
	NetworkIP        *string `json:"network_ip,omitempty"`
	VirtualMachineID *int    `json:"virtual_machine_id,omitempty"`
	MountPath        *string `json:"mount_path,omitempty"`
}

func (n Storage) Validate() error { //nolint:cyclop
	if n.ID < 0 {
		return errors.New("id cannot be negative number")
	}
	if err := ValidateName(n.Name); err != nil {
		return err
	}
	if err := ValidateDisplayName(n.DisplayName); err != nil {
		return err
	}

	if n.SizeMB < 1024 || n.SizeMB > 200*1024*1024 {
		return errors.New("size_mb must be between 1 GB and 200 TB")
	}

	if n.NetworkID != nil && n.VirtualMachineID != nil || n.NetworkID == nil && n.VirtualMachineID == nil {
		return errors.New("exactly one of network_id or virtual_machine_id must be set")
	}

	if n.NetworkID != nil && *n.NetworkID < 0 {
		return errors.New("network_id cannot be negative number")
	}
	if n.NetworkIP != nil {
		if !ValidIP(*n.NetworkIP) {
			return errors.New("network_ip is not valid")
		}
		if n.NetworkID == nil {
			return errors.New("cannot set network_ip when network_id is null")
		}
	}

	if n.VirtualMachineID != nil && *n.VirtualMachineID < 0 {
		return errors.New("virtual_machine_id cannot be negative number")
	}
	if n.MountPath != nil && n.VirtualMachineID == nil {
		return errors.New("cannot set mount_path when virtual_machine_id is null")
	}
	if n.MountPath == nil && n.VirtualMachineID != nil {
		return errors.New("mount_path is required when virtual_machine_id is set")
	}

	return nil
}

type Network struct {
	ID          int64   `json:"id,omitempty"`
	Name        string  `json:"name"`
	DisplayName *string `json:"display_name,omitempty"`
	IPRange     string  `json:"ip_range"`
	UseDHCP     bool    `json:"use_dhcp"`
}

func (n Network) Validate() error {
	if n.ID < 0 {
		return errors.New("id cannot be negative number")
	}
	if err := ValidateName(n.Name); err != nil {
		return err
	}
	if err := ValidateDisplayName(n.DisplayName); err != nil {
		return err
	}

	if !ValidIPRange(n.IPRange) {
		return errors.New("IP range is not valid, use 8.8.8.8/24 format")
	}
	return nil
}
