package main

import (
	"fmt"
	"vdirect"

	"github.com/hashicorp/terraform/helper/schema"
)

func registerContainer() *schema.Resource {
	return &schema.Resource{
		Create: registerContainerCreate,
		Read:   registerContainerRead,
		Update: registerContainerUpdate,
		Delete: registerContainerDelete,

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
			"adc_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"https_port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"ssh_port": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"adc_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"adc_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func registerContainerCreate(d *schema.ResourceData, m interface{}) error {
	vdirectIP := d.Get("vdirect_ip").(string)
	username := d.Get("vdirect_username").(string)
	password := d.Get("vdirect_password").(string)
	adcName := d.Get("adc_name").(string)
	adcIp := d.Get("adc_ip").(string)
	adcUsername := d.Get("adc_username").(string)
	adcPassword := d.Get("adc_password").(string)
	httpsPort := d.Get("https_port").(string)
	sshPort := d.Get("ssh_port").(string)
	message := map[string]interface{}{
		"name": adcName,
		"type": "AlteonDedicated",
		"configuration": map[string]string{
			"configProtocol": "HTTPS",
			"host":           adcIp,
			"cli.user":       adcUsername,
			"cli.password":   adcPassword,
			"cli.ssh":        "true",
			"cli.port":       sshPort,
			"https.port":     httpsPort,
			"https.user":     adcUsername,
			"https.password": adcPassword,
		},
	}

	client := vdirect.NewClient(vdirectIP, username, password, vdirect.NewClientConfig(true, 120, false, 120))
	resp := client.Container.Create0(message, "false")

	// check response code
	if resp.StatusCode != 201 {
		return fmt.Errorf("%s", resp.ToString())
	}

	d.SetId(adcName)
	return registerContainerRead(d, m)
}

func registerContainerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func registerContainerUpdate(d *schema.ResourceData, m interface{}) error {
	return registerContainerRead(d, m)
}

func registerContainerDelete(d *schema.ResourceData, m interface{}) error {
	vdirectIP := d.Get("vdirect_ip").(string)
	username := d.Get("vdirect_username").(string)
	password := d.Get("vdirect_password").(string)
	adcName := d.Get("adc_name").(string)
	client := vdirect.NewClient(vdirectIP, username, password, vdirect.NewClientConfig(true, 120, false, 120))
	resp := client.Container.Delete(adcName)

	// check response code
	if resp.StatusCode != 204 {
		return fmt.Errorf("%s", resp.ToString())
	}
	d.SetId("")
	return nil
}
