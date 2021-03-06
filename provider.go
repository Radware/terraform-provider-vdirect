package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"vdirect_container_register":  registerContainer(),
			"vdirect_runnable_template":   runnableTemplate(),
			"vdirect_idle_timeout_change": idleTimeoutChange(),
		},
	}
}
