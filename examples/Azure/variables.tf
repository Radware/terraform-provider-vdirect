variable "resource_group" {
  default = "nati"
  description = "The name of the resource group in which to create the virtual network."
}
variable "rg_prefix" {
  default = "nati"
  description = "The shortened abbreviation to represent your resource group that will go on the front of some resources."
}
variable "hostname" {
  default = "nati-terraform"
  description = "VM name referenced also in storage-related names."
}
variable "dns_name" {
  default = "nati-terraform"
  description = " Label for the Domain Name. Will be used to make up the FQDN. If a domain name label is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system."
}
variable "location" {
  default = "West Europe"
  description = "The location/region where the virtual network is created. Changing this forces a new resource to be created."
}
variable "virtual_network_name" {
  default = "nati"
  description = "The name for the virtual network."
}
variable "address_space" {
  default = "10.1.1.0/24"
  description = "The address space that is used by the virtual network. You can supply more than one address space. Changing this forces a new resource to be created."
}
variable "subnet_prefix" {
  default = "10.1.1.0/24"
  description = "The address prefix to use for the subnet."
}
variable "vm_size" {
  default = "Standard_DS1_v2"
  description = "Specifies the size of the virtual machine."
}
variable "admin_username" {
  default = "nati"
  description = "administrator user name"
}
variable "admin_password" {
  default = "Aa123456123456"
  description = "administrator password (recommended to disable password auth)"
}