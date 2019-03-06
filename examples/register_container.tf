resource "vdirect_container_register" "my-container" {
	vdirect_ip = "10.210.71.71"
	vdirect_username = "root"
	vdirect_password = "radware"
	adc_name = "tf-test"
  	adc_ip = "10.210.71.33"
	https_port = "443"
}

resource "vdirect_runnable_template" "my-template" {
	vdirect_ip = "10.210.71.71"
	vdirect_username = "root"
	vdirect_password = "radware"
	runnable_type = "ConfigurationTemplate"
	runnable_name = "real_servers_crud.vm"
	action = "run"
	paramaters = "{\"alteon\":{\"type\":\"Container\",\"name\":\"tf-test\"},\"real_servers\":[{\"name\":\"test\",\"address\":\"1.1.1.1\",\"weight\":1,\"state\":\"disable\",\"action\":\"create\"}],\"__dryRun\":false}"

	depends_on = ["radware_container_register.my-container"]
}