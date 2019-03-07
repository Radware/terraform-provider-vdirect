# Configure the Microsoft Azure Provider
provider "azurerm"{
  # if you're using a Service Principal (shared account) then either set the environment variables, or fill these in: 
  #subscription_id = "${var.subscription_id}"
  #client_id       = "${var.client_id}"
  # client_secret   = "${var.client_secret}"
  # tenant_id       = "${var.tenant_id}"
}
resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}
resource "azurerm_virtual_network" "vnet" {
  name                = "${var.virtual_network_name}"
  location            = "${var.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  address_space       = ["${var.address_space}"]
}
resource "azurerm_subnet" "subnet" {
  name                 = "${var.rg_prefix}subnet"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  address_prefix       = "${var.subnet_prefix}"
}
resource "azurerm_network_interface" "nic" {
  name                = "${var.rg_prefix}nic"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  ip_configuration {
    name                          = "${var.rg_prefix}ipconfig"
    subnet_id                     = "${azurerm_subnet.subnet.id}"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = "${azurerm_public_ip.pip.id}"
  }
}
resource "azurerm_public_ip" "pip" {
  name                         = "${var.rg_prefix}-ip"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "Dynamic"
  domain_name_label            = "${var.dns_name}"
}
resource "azurerm_virtual_machine" "vm" {
  name                  = "${var.rg_prefix}vm"
  resource_group_name   = "${azurerm_resource_group.rg.name}"
  location            = "${azurerm_resource_group.rg.location}"
  network_interface_ids = ["${azurerm_network_interface.nic.id}"]
  vm_size               = "${var.vm_size}"

  storage_image_reference {
    publisher = "radware"
    offer     = "radware-alteon-va"
    sku       = "radware-alteon-ng-va-adc"
    version   = "latest"
  }

  plan {
        name = "radware-alteon-ng-va-adc"
        publisher = "radware"
        product = "radware-alteon-va"
        }

  storage_os_disk {
    name              = "${var.hostname}-osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "${var.hostname}"
    admin_username = "${var.admin_username}"
    admin_password = "${var.admin_password}"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

provider "radware" {
  
}

resource "vdirect_container_register" "my-container" {
	address = "10.210.71.71"
	username = "root"
	password = "radware"
  adc_name = "tf-test-azure"
  adc_ip = "${azurerm_public_ip.pip.ip_address}"
  https_port = "8443"
  ssh_port = "22"
	adc_username = "admin"
	adc_password = "admin"

  depends_on = ["azurerm_public_ip.pip"]
}

resource "vdirect_runnable_template" "my-template" {
	address = "10.210.71.71"
	username = "root"
	password = "radware"
	runnable_type = "ConfigurationTemplate"
	runnable_name = "real_servers_crud.vm"
	action = "run"
	paramaters = "{\"alteon\":{\"type\":\"Container\",\"name\":\"tf-test-azure\"},\"real_servers\":[{\"name\":\"test\",\"address\":\"1.1.1.1\",\"weight\":1,\"state\":\"disable\",\"action\":\"create\"}],\"__dryRun\":false}"

  depends_on = ["radware_container_register.my-container"]
}
