package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//TODO
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"openwhisk_function": resourceServer(),
		},
	}
}
