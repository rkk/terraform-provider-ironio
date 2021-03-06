package ironio

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/iron-io/iron_go3/api"
	"github.com/iron-io/iron_go3/config"
)

// ProjectInfo describes a project.
type ProjectInfo struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	TenantID  int       `json:"tenant_id,omitempty"`
	Name      string    `json:"name"`
	Status    string    `json:"status,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
}

// ProjectRequest describes a project request payload.
type ProjectRequest struct {
	Project ProjectInfo `json:"project"`
}

// resourceProject() manages projects.
func resourceProject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the project",
			},
		},

		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
	}
}

// resourceProjectCreate() creates a project.
func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	clientSettings := m.(config.Settings)
	project := ProjectInfo{
		Name: d.Get("name").(string),
	}

	in := ProjectRequest{
		Project: project,
	}

	var out ProjectInfo

	url := resourceProjectGetEndpoint(clientSettings, "")
	err := url.Req("POST", in, &out)

	if err != nil {
		return err
	}

	d.SetId(out.ID)

	return nil
}

// resourceProjectGetEndpoint() returns an endpoint for a project.
func resourceProjectGetEndpoint(cs config.Settings, id string) *api.URL {
	u := &api.URL{Settings: cs, URL: url.URL{Scheme: cs.Scheme}}

	u.URL.Host = fmt.Sprintf("%s:%d", cs.Host, cs.Port)
	u.URL.Path = fmt.Sprintf("/%s/projects", cs.ApiVersion)

	if id != "" {
		u.URL.Path = fmt.Sprintf("%s/%s", u.URL.Path, id)
	}

	return u
}

// resourceProjectRead reads information about an existing project.
func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	clientSettings := m.(config.Settings)

	var out ProjectInfo

	url := resourceProjectGetEndpoint(clientSettings, d.Id())
	err := url.Req("GET", nil, &out)

	if err != nil {
		if strings.Contains(err.Error(), " 404 ") {
			d.SetId("")

			return nil
		}
		return err
	}

	d.Set("name", out.Name)

	return nil
}

// resourceProjectUpdate updates an existing project.
func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	clientSettings := m.(config.Settings)
	project := ProjectInfo{
		Name: d.Get("name").(string),
	}

	in := ProjectRequest{
		Project: project,
	}

	var out ProjectInfo

	url := resourceProjectGetEndpoint(clientSettings, d.Id())
	err := url.Req("PATCH", in, &out)

	if err != nil {
		return err
	}

	return nil
}

// resourceProjectDelete deletes an existing project.
func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	clientSettings := m.(config.Settings)

	var out struct {
		Message string `json:"msg"`
	}

	url := resourceProjectGetEndpoint(clientSettings, d.Id())
	err := url.Req("DELETE", nil, &out)

	if err != nil {
		if !strings.Contains(err.Error(), " 404 ") {
			return err
		}
	}

	if out.Message != "success" {
		return fmt.Errorf("ERROR: Failed to delete the project due to an unknown error")
	}

	d.SetId("")

	return nil
}
