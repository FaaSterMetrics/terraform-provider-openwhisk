package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func hashFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

const ENV_PREFIX = "__env_"

func environmentToParams(d *schema.ResourceData) whisk.KeyValueArr {
	params := make(whisk.KeyValueArr, 0)
	var environment map[string]interface{}

	if v, ok := d.GetOk("environment"); ok {
		environment = v.(map[string]interface{})
	} else {
		return params
	}

	for k, v := range environment {
		params = params.AddOrReplace(&whisk.KeyValue{
			Key:   ENV_PREFIX + k,
			Value: v.(string),
		})
	}
	return params
}

func paramsToEnvironment(params whisk.KeyValueArr) map[string]string {
	environment := make(map[string]string)
	for _, v := range params {
		if strings.HasPrefix(v.Key, ENV_PREFIX) {
			environment[strings.Replace(v.Key, ENV_PREFIX, "", 1)] = v.Value.(string)
		}
	}
	return environment
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	source := d.Get("source").(string)
	client := m.(*whisk.Client)

	code, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	codeStr := base64.StdEncoding.EncodeToString(code)
	action := &whisk.Action{
		Name: name,
		Annotations: whisk.KeyValueArr{
			{Key: "web-export", Value: true},
			{Key: "raw-http", Value: false},
			{Key: "final", Value: true},
			{Key: "provide-api-key", Value: false},
		},
		Exec: &whisk.Exec{
			Kind: "nodejs:10",
			Code: &codeStr,
		},
		Parameters: environmentToParams(d),
	}

	resAction, _, err := client.Actions.Insert(action, false)
	if err != nil {
		return err
	}

	d.SetId(resAction.Name)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	client := m.(*whisk.Client)
	action, _, err := client.Actions.Get(id, false)
	if err != nil && strings.Contains(err.Error(), "The requested resource does not exist") {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	d.Set("environment", paramsToEnvironment(action.Parameters))
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	client := m.(*whisk.Client)
	if d.HasChange("environment") {
		_, _, err := client.Actions.Insert(&whisk.Action{
			Name:       id,
			Parameters: environmentToParams(d),
		}, true)
		if err != nil {
			return err
		}
	}
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	client := m.(*whisk.Client)
	_, err := client.Actions.Delete(id)
	return err
}

func customDiff(d *schema.ResourceDiff, m interface{}) error {
	d.SetNew("source_hash", hashFile(d.Get("source").(string)))
	return nil
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceServerCreate,
		Read:          resourceServerRead,
		Update:        resourceServerUpdate,
		Delete:        resourceServerDelete,
		CustomizeDiff: customDiff,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source_hash": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"environment": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
