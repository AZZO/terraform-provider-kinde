package kinde

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	validApplicationTypes = []string{"reg", "spa", "m2m"}
)

func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationCreate,
		Read:   resourceApplicationRead,
		Update: resourceApplicationUpdate,
		Delete: resourceApplicationDelete,

		Schema: map[string]*schema.Schema{
			// required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(validApplicationTypes, false),
				ForceNew:     true,
			},

			// optional
			"redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"logout_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// computed
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	req := &ApiRequestApplicationCreate{}
	req.Name = d.Get("name").(string)
	req.Type = d.Get("type").(string)

	resp := &ApiResponseApplicationCreate{}

	err := PostJson(config.client, EndpointApplicationCreate(config), req, resp)

	if err != nil {
		return err
	}

	// urisLogout := d.Get("logout_uris").([]string)
	urisLogoutRaw := d.Get("logout_uris").([]interface{})
	urisLogout := make([]string, len(urisLogoutRaw))
	for i, raw := range urisLogoutRaw {
		urisLogout[i] = raw.(string)
	}

	// urisRedirect := d.Get("redirect_uris").([]string)
	urisRedirectRaw := d.Get("redirect_uris").([]interface{})
	urisRedirect := make([]string, len(urisRedirectRaw))
	for i, raw := range urisRedirectRaw {
		urisRedirect[i] = raw.(string)
	}

	if len(urisLogout) > 0 || len(urisRedirect) > 0 {
		uris := &ApiRequestApplicationUpdate{}
		uris.RedirectUris = urisRedirect
		uris.LogoutUris = urisRedirect

		err = PostJson(config.client, EndpointApplicationUpdate(config, resp.Application.Id), uris, resp)
		if err != nil {
			return err
		}
	}

	d.SetId(resp.Application.Id)
	d.Set("application_id", resp.Application.Id)
	d.Set("client_id", resp.Application.ClientId)
	d.Set("client_secret", resp.Application.ClientSecret)

	return nil
}

func resourceApplicationRead(d *schema.ResourceData, meta interface{}) error {

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

	return nil
}

func resourceApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	req := &ApiRequestApplicationUpdate{}
	req.Name = d.Get("name").(string)

	// urisLogout := d.Get("logout_uris").([]string)
	urisLogoutRaw := d.Get("logout_uris").([]interface{})
	req.LogoutUris = make([]string, len(urisLogoutRaw))
	for i, raw := range urisLogoutRaw {
		req.LogoutUris[i] = raw.(string)
	}

	// urisRedirect := d.Get("redirect_uris").([]string)
	urisRedirectRaw := d.Get("redirect_uris").([]interface{})
	req.RedirectUris = make([]string, len(urisRedirectRaw))
	for i, raw := range urisRedirectRaw {
		req.RedirectUris[i] = raw.(string)
	}

	err := PatchJson(config.client, EndpointApplicationUpdate(config, d.Get("application_id").(string)), req)

	if err != nil {
		return err
	}

	return nil
}

func resourceApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	err := Delete(config.client, EndpointApplicationDelete(config, d.Get("application_id").(string)))

	if err != nil {
		return err
	}

	return nil
}
