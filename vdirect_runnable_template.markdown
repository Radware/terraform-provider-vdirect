# vdirect_runnable_template

`vdirect_runnable_template` allow using vdirect runnable API

## Example Usage


```hcl

resource "vdirect_runnable_template" "my-template" {
    vdirect_ip = "10.210.71.71"
    vdirect_username = "root"
    vdirect_password = "radware"
    runnable_type = "ConfigurationTemplate"
    runnable_name = "real_servers_crud.vm"
    action = "run"
    paramaters = "{\"alteon\":{\"type\":\"Container\",\"name\":\"tf-test\"},\"real_servers\":[{\"name\":\"test\",\"address\":\"1.1.1.1\",\"weight\":1,\"state\":\"disable\",\"action\":\"create\"}],\"__dryRun\":false}"

}

```       

## Argument Reference

* `vdirect_ip`
* `vdirect_username`
* `vdirect_password`
* `runnable_type` – vDirect runnable type (ConfigurationTemplate or WorkflowTemplate)
* `runnable_name`
* `action` – action to do (if runnable_type=ConfigurationTemplate put “run”)
* `parameters` – parameters JSON as string
