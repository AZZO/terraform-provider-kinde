package kinde

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApplicationRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"logout_uris": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApplicationRead(d *schema.ResourceData, meta interface{}) error {

	// reader := ResourceApplication()
	// reader.Read(d, meta)

	config := meta.(*Config)
	result := &ApiResponseApplicationGet{}
	err := GetJson(config.client, EndpointApplicationGet(config, d.Get("application_id").(string)), result)

	if err != nil {
		return err
	}

	d.SetId(result.Application.Id)
	d.Set("application_id", result.Application.Id)
	d.Set("name", result.Application.Name)
	d.Set("type", result.Application.Type)
	d.Set("client_id", result.Application.ClientId)
	d.Set("client_secret", result.Application.ClientSecret)

	// TODO not currently returned by API
	d.Set("redirect_uris", []string{})
	d.Set("logout_uris", []string{})

	return nil
}
