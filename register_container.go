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
				Type:     schema.TypeInt,
				Required: true,
			},
			"ssh_port": &schema.Schema{
				Type:     schema.TypeInt,
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
	httpsPort := d.Get("https_port").(int)
	sshPort := d.Get("ssh_port").(int)

	message := new(adcConfig)
	message.Name = adcName
	message.Connect.Cli.Port = sshPort
	message.Connect.Cli.User = adcUsername
	message.Connect.Cli.Password = adcPassword
	message.Connect.Cli.SSH = true
	message.Connect.HTTPS.Port = httpsPort
	message.Connect.HTTPS.User = adcUsername
	message.Connect.HTTPS.Password = adcPassword
	message.Connect.IP = adcIp
	message.Connect.ConfigProtocol = "HTTPS"

	client := vdirect.NewClient(vdirectIP, username, password, vdirect.NewClientConfig(true, 120, false, 120))
	resp := client.ADC.Create(message, "false")

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
	resp := client.ADC.Delete(adcName, "DELETE")

	// check response code
	if resp.StatusCode != 204 {
		return fmt.Errorf("%s", resp.ToString())
	}
	d.SetId("")
	return nil
}

type adcConfig struct {
	Name                string        `json:"name"`
	Tenants             []interface{} `json:"tenants"`
	ExtensionProperties struct {
	} `json:"extensionProperties"`
	Connect struct {
		Snmp struct {
		} `json:"snmp"`
		Cli struct {
			Port     int    `json:"port"`
			User     string `json:"user"`
			Password string `json:"password"`
			SSH      bool   `json:"ssh"`
		} `json:"cli"`
		HTTPS struct {
			User     string `json:"user"`
			Password string `json:"password"`
			Port     int    `json:"port"`
		} `json:"https"`
		IP             string `json:"ip"`
		ConfigProtocol string `json:"configProtocol"`
	} `json:"connect"`
}
