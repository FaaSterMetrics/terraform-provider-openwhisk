package main

import (
	"net/http"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//TODO
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"openwhisk_function": resourceServer(),
		},
		ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
			return whisk.NewClient(http.DefaultClient, nil)
		},
	}
}
