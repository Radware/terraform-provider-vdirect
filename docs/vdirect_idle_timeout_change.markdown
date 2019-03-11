# vdirect_idle_timeout_change

`vdirect_idle_timeout_change` allow change alteon idle timeout through vdirect

## Example Usage


```hcl

resource "vdirect_idle_timeout_change" "idle-timeout" {
	vdirect_ip = "10.210.71.71"
	vdirect_username = "root"
	vdirect_password = "radware"
	adc_name = "tf-test"
  	idle_timeout = 50

}


```       

## Argument Reference

* `vdirect_ip`
* `vdirect_username`
* `vdirect_password`
* `adc_name`
* `idle_timeout` - timeout to set
