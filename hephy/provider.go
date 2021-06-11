package hephy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	hephy "github.com/teamhephy/controller-sdk-go"
)

// Provider function which defines the schema for the EYK provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"controller_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"admin_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ssl_verify": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hephy_initial_admin": resourceInitialAdmin(),
			"hephy_user":          resourceUser(),
			"hephy_key":           resourceKey(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	controllerUrl := d.Get("controller_url").(string)
	token := d.Get("admin_token").(string)
	sslVerify := d.Get("ssl_verify").(bool)

	var diags diag.Diagnostics

	// hephy client setup
	c, err := hephy.New(sslVerify, controllerUrl, token)
	if err != nil {
		return nil, append(diags, diagFromError("Error initializing hephy client", err))
	}

	return c, diags
}
