package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"vdirect"

	"github.com/hashicorp/terraform/helper/schema"
)

func runnableTemplate() *schema.Resource {
	return &schema.Resource{
		Create: runnableTemplateCreate,
		Read:   runnableTemplateRead,
		Update: runnableTemplateUpdate,
		Delete: runnableTemplateDelete,

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
			"runnable_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"runnable_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"paramaters": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func runnableTemplateCreate(d *schema.ResourceData, m interface{}) error {
	vdirectIP := d.Get("vdirect_ip").(string)
	username := d.Get("vdirect_username").(string)
	password := d.Get("vdirect_password").(string)
	runnableType := d.Get("runnable_type").(string)
	runnableName := d.Get("runnable_name").(string)
	action := d.Get("action").(string)
	paramaters := d.Get("paramaters").(string)

	paramaters = strings.Replace(paramaters, "{", "{\n", -1)
	paramaters = strings.Replace(paramaters, ",", ",\n", -1)
	//payload := strings.NewReader(paramaters)
	//path := "http://" + address + ":2188"

	// check if the runnable type accepted
	if runnableType != "WorkflowTemplate" && runnableType != "ConfigurationTemplate" {
		log.Fatalln("wrong runnable type")
	}

	paramatersMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(paramaters), &paramatersMap)

	if err != nil {
		log.Fatalln(err)
	}

	client := vdirect.NewClient(vdirectIP, username, password, vdirect.NewClientConfig(true, 120, false, 120))
	resp := client.Runnable.Run(paramatersMap, runnableName, runnableType, action)

	// check response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", resp.ToString())
	}

	d.SetId(runnableName)
	return runnableTemplateRead(d, m)
}

func runnableTemplateRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func runnableTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	return runnableTemplateRead(d, m)
}

func runnableTemplateDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
