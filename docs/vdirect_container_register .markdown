# vdirect_container_register

`vdirect_container_register` allow register adc container in vdirect

## Example Usage


```hcl

resource "vdirect_container_register" "my-container" {
    vdirect_ip = "10.210.71.71"
    vdirect_username = "root"
    vdirect_password = "radware"
    adc_name = "tf-test"
    adc_ip = "10.210.71.33"
    https_port = "443"
    ssh_port = "22"
    adc_username = "admin"
    adc_password = "admin"
}


```       

## Argument Reference

* `vdirect_ip`
* `vdirect_username`
* `vdirect_password`
* `adc_name`
* `adc_ip`
* `https_port`
* `ssh_port`
* `adc_username`
* `adc_password`
