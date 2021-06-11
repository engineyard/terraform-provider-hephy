package hephy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hephy "github.com/teamhephy/controller-sdk-go"
	"github.com/teamhephy/controller-sdk-go/api"
	"github.com/teamhephy/controller-sdk-go/auth"
	"github.com/teamhephy/controller-sdk-go/users"
	//"github.com/teamhephy/controller-sdk-go/users"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"creator_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hephy.Client)

	var diags diag.Diagnostics

	username := d.Get("username").(string)
	email := d.Get("email").(string)
	password := d.Get("password").(string)
	creatorToken := d.Get("creator_token").(string)

	// Use the specified credentials if provided
	if len(creatorToken) > 0 {
		client.Token = creatorToken
	}

	err := CheckConnection(client)
	if err != nil {
		return append(diags, diagFromError("Error creating hephy client", err))
	}

	err = auth.Register(client, username, password, email)
	if err != nil {
		return append(diags, diagFromError("Error creating user", err))
	}

	token, err := auth.Login(client, username, password)
	if err != nil {
		return append(diags, diagFromError("Error retrieving admin token", err))
	}

	d.Set("username", username)
	d.Set("password", password)
	d.Set("email", email)
	d.Set("token", token)
	d.SetId(username)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hephy.Client)

	creatorToken := d.Get("creator_token").(string)

	// Use specified credentials if provided
	if len(creatorToken) > 0 {
		client.Token = creatorToken
	}

	var diags diag.Diagnostics

	err := CheckConnection(client)
	if err != nil {
		return append(diags, diagFromError("Error creating hephy client", err))
	}

	userID := d.Id()

	var user api.User

	collection, _, err := users.List(client, -1)
	if err != nil {
		return append(diags, diagFromError("Error getting user list", err))
	}

	for _, candidate := range collection {
		if candidate.Username == userID {
			user = candidate
			break
		}
	}

	if user.Username == "" {
		return append(diags, diagFromError("Error getting user info", fmt.Errorf("no such user")))
	}

	d.Set("username", userID)
	d.Set("email", user.Email)
	d.SetId(userID)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*hephy.Client)

	creatorToken := d.Get("creator_token").(string)

	// Use specified credentials if provided
	if len(creatorToken) > 0 {
		client.Token = creatorToken
	}

	var diags diag.Diagnostics
	err := CheckConnection(client)
	if err != nil {
		return append(diags, diagFromError("Error creating hephy client", err))
	}

	userID := d.Id()

	err = auth.Delete(client, userID)
	if err != nil {
		return append(diags, diagFromError("Error deleting user", err))
	}

	d.SetId("")
	return diags
}
