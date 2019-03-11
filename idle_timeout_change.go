package main

import (
	"fmt"
	"vdirect"

	"github.com/hashicorp/terraform/helper/schema"
)

func idleTimeoutChange() *schema.Resource {
	return &schema.Resource{
		Create: idleTimeoutChangeCreate,
		Read:   idleTimeoutChangeRead,
		Update: idleTimeoutChangeUpdate,
		Delete: idleTimeoutChangeDelete,

		Schema: map[string]*schema.Schema{
			"vdirect_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vdirect_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vdirect_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"adc_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"idle_timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func idleTimeoutChangeCreate(d *schema.ResourceData, m interface{}) error {
	vdirectIP := d.Get("vdirect_ip").(string)
	username := d.Get("vdirect_username").(string)
	password := d.Get("vdirect_password").(string)
	adcName := d.Get("adc_name").(string)
	idleTimeout := d.Get("idle_timeout").(int)
	message := map[string]interface{}{
		"AgSystem.agNewCfgIdleCLITimeout": idleTimeout,
	}

	client := vdirect.NewClient(vdirectIP, username, password, vdirect.NewClientConfig(true, 120, false, 120))
	resp := client.ADC.UpdateConfiguration(message, adcName, "AgSystem")

	// check response code
	if resp.StatusCode != 204 {
		return fmt.Errorf("%s", resp.ToString())
	}

	resp = client.ADC.Control1(adcName, "apply")

	// check response code
	if resp.StatusCode != 204 {
		return fmt.Errorf("%s", resp.ToString())
	}

	resp = client.ADC.Control1(adcName, "save")

	// check response code
	if resp.StatusCode != 204 {
		return fmt.Errorf("%s", resp.ToString())
	}

	d.SetId(adcName)
	return idleTimeoutChangeRead(d, m)
}

func idleTimeoutChangeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func idleTimeoutChangeUpdate(d *schema.ResourceData, m interface{}) error {
	return registerContainerRead(d, m)
}

func idleTimeoutChangeDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
