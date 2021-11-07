provider "sdc" {
  endpoint = var.endpoint
  api_token = var.api_token
}

data sdc_network existing-network {
  id = 1
}
data sdc_vm existing-vm {
  id = 1
}
data sdc_storage existing-storage {
  id = 1
}

output network {
  value       = data.sdc_network.existing-network
}
output vm {
  value       = data.sdc_vm.existing-vm
}
output storage {
  value       = data.sdc_storage.existing-storage
}

resource sdc_network main_network {
  name = "backbone"
  display_name = "Main backbone network"
  ip_range = "192.168.0.0/16"
  use_dhcp = true
}

resource sdc_vm mainframe {
  name = "mainframe"
  display_name = "Mainframe"
  ram_size_mb = 4096
  network_id = sdc_network.main_network.id
  public_ip = "123.123.123.143"
}

resource sdc_storage vmhdd {
  name = "${sdc_vm.mainframe.name}-storage"
  size_mb = 180 * 1024
  virtual_machine_id = sdc_vm.mainframe.id
  mount_path = "/mtn/os"
}

resource sdc_storage network-disk {
  name = "${sdc_network.main_network.name}-storage"
  size_mb = 80 * 1024 * 1024
  network_id = sdc_network.main_network.id
  network_ip = "192.161.100.20"
}
