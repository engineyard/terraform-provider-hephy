package hephy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hephy "github.com/teamhephy/controller-sdk-go"
	"github.com/teamhephy/controller-sdk-go/api"
	"github.com/teamhephy/controller-sdk-go/keys"
	//"github.com/teamhephy/controller-sdk-go/users"
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyCreate,
		ReadContext:   resourceKeyRead,
		DeleteContext: resourceKeyDelete,
		Schema: map[string]*schema.Schema{
			"owner_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hephy.Client)

	var diags diag.Diagnostics

	ownerToken := d.Get("owner_token").(string)
	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	client.Token = ownerToken

	err := CheckConnection(client)
	if err != nil {
		return append(diags, diagFromError("Error creating hephy client", err))
	}

	key, err := keys.New(client, name, publicKey)
	if err != nil {
		return append(diags, diagFromError("Error creating key", err))
	}

	d.SetId(key.ID)

	return diags
}

func resourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hephy.Client)

	ownerToken := d.Get("owner_token").(string)

	client.Token = ownerToken

	var diags diag.Diagnostics

	err := CheckConnection(client)
	if err != nil {
		return append(diags, diagFromError("Error creating hephy client", err))
	}

	keyID := d.Id()

	var key api.Key

	collection, _, err := keys.List(client, -1)
	if err != nil {
		return append(diags, diagFromError("Error getting key list", err))
	}

	for _, candidate := range collection {
		if candidate.ID == keyID {
			key = candidate
			break
		}
	}

	if key.ID == "" {
		return append(diags, diagFromError("Error getting key info", fmt.Errorf("no such key")))
	}

	d.Set("name", key.ID)
	d.SetId(key.ID)

	return diags
}

func resourceKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hephy.Client)

	ownerToken := d.Get("owner_token").(string)

	client.Token = ownerToken

	var diags diag.Diagnostics
	err := CheckConnection(client)
	if err != nil {
		return append(diags, diagFromError("Error creating hephy client", err))
	}

	keyID := d.Id()

	err = keys.Delete(client, keyID)
	if err != nil {
		return append(diags, diagFromError("Error deleting key", err))
	}

	d.SetId("")
	return diags
}
