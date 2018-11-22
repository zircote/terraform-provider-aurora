package aurora

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_SERVER_URL", ""),
				Description: "The hostname (in form of URI) of Aurora master.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_USERNAME", ""),
				Description: "Username to use for authorization",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_PASSWORD", ""),
				Description: "Password to use for authorization",
			},
			"zk_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_ZK_URL", ""),
				Description: "zookeeper url",
			},
			"ca_certs_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_CA_PATH", ""),
				Description: "Path to CA certs on local machine.",
			},
			"client_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_CLIENT_CERT", ""),
				Description: "Client certificate to use to connect to Aurora.",
			},
			"client_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_CLIENT_KEY", ""),
				Description: "Client private key to use to connect to Aurora.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AURORA_CLIENT_TIMEOUT", 20000),
				Description: "connection timeout in ms",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"aurora_job":  resourceAuroraJob(),
			"aurora_role": resourceAuroraRole(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		SchedulerUrl: d.Get("server_url").(string),
		Timeout:      d.Get("timeout").(int),
		CaCertsPath:  d.Get("ca_certs_path").(string),
		ClientKey:    d.Get("client_key").(string),
		ClientCert:   d.Get("client_cert").(string),
		ZkUrl:        d.Get("zk_url").(string),
	}
	if err := config.CreateAuroraClient(); err != nil {
		return nil, err
	}
	return config, nil
}
