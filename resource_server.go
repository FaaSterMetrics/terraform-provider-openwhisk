package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/apache/openwhisk-wskdeploy/cmd"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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

type action struct {
	Function string            `yaml:"function"`
	Runtime  string            `yaml:"runtime"`
	Web      string            `yaml:"web"`
	Inputs   map[string]string `yaml:"inputs"`
}

type manifestYaml struct {
	Packages struct {
		Faastermetrics struct {
			Actions map[string]action `yaml:"actions"`
		} `yaml:"faastermetrics"`
	} `yaml:"packages"`
}

var mux sync.Mutex

func destroyFunction(name string) {
	mux.Lock()
	defer mux.Unlock()
	env := make(map[string]string)
	smallestZip := []byte{0x50, 0x4b, 0x05, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	ioutil.WriteFile("smallest.zip", smallestZip, 0644)
	defer os.Remove("smallest.zip")
	act := action{Function: "smallest.zip",
		Runtime: "nodejs:10",
		Web:     "yes",
		Inputs:  env,
	}
	manifest := manifestYaml{}

	manifest.Packages.Faastermetrics.Actions = make(map[string]action)
	manifest.Packages.Faastermetrics.Actions[name] = act
	data, _ := yaml.Marshal(manifest)
	ioutil.WriteFile("manifest.yml", data, 0644)
	defer os.Remove("manifest.yml")
	cmd.Undeploy(&cobra.Command{})
}

func deployFunction(zipPath string, name string, environment map[string]interface{}) {
	mux.Lock()
	defer mux.Unlock()

	prefixedEnv := make(map[string]string)
	for k, v := range environment {
		prefixedEnv["__env_"+k] = v.(string)
	}

	act := action{Function: zipPath,
		Runtime: "nodejs:10",
		Web:     "yes",
		Inputs:  prefixedEnv,
	}
	manifest := manifestYaml{}

	manifest.Packages.Faastermetrics.Actions = make(map[string]action)
	manifest.Packages.Faastermetrics.Actions[name] = act
	data, _ := yaml.Marshal(manifest)
	ioutil.WriteFile("manifest.yml", data, 0644)
	defer os.Remove("manifest.yml")
	cmd.Deploy(&cobra.Command{})
}

//TODO
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	zipPath := d.Get("zip_path").(string)
	environment := d.Get("environment").(map[string]interface{})
	deployFunction(zipPath, name, environment)
	d.SetId(name + ":" + hashFile(zipPath))

	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	zipPath := d.Get("zip_path").(string)
	hash := hashFile(zipPath)
	if name+":"+hash != d.Id() {
		d.SetId("")
	} else {
		d.SetId(name + ":" + hash)
	}

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	zipPath := d.Get("zip_path").(string)
	environment := d.Get("environment").(map[string]interface{})
	id := d.Id()
	destroyFunction(strings.Split(id, ":")[0])

	deployFunction(zipPath, name, environment)
	d.SetId(name + ":" + hashFile(zipPath))

	return resourceServerRead(d, m)
}

//TODO
func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	destroyFunction(name)
	d.SetId("")
	return nil
}

//TODO
func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zip_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},
		},
	}
}
