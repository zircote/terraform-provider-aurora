package aurora

import "github.com/hashicorp/terraform/helper/schema"

func announceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"primary_port": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Which named port to register as the primary endpoint in the ServerSet (Default: http)",
				},
				"portmap": {
					Type:        schema.TypeMap,
					Description: "A mapping of additional endpoints to be announced in the ServerSet (Default: { 'aurora': '{{primary_port}}' })",
					Optional:    true,
				},
				"zk_path": {
					Type:        schema.TypeString,
					Description: "Zookeeper serverset path override (executor must be started with the --announcer-allow-custom-serverset-path parameter)",
					Optional:    true,
				},
			},
		},
	}
}

func constraintsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "",
				},
				"value": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "",
				},
			},
		},
	}
}

func executorConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the executor to use for this task. Must match the name of an executor in custom_executor_config or Thermos (AuroraExecutor). (Default: AuroraExecutor)",
					Default:     "AuroraExecutor",
				},
				"data": {
					Type:        schema.TypeString,
					Description: "Data blob to pass on to the executor. (Default: “”)",
					Default:     "",
					Optional:    true,
				},
			},
		},
	}
}

func metadataSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
		},
	}
}

func resourceSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		MaxItems:    1,
		Description: "Resource footprint.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cpu": {
					Type:        schema.TypeFloat,
					Description: "Fractional number of cores required by the task.",
					Optional:    true,
				},
				"ram_mb": {
					Type:        schema.TypeInt,
					Description: "MBytes of RAM required by the task.",
					Optional:    true,
				},
				"disk_mb": {
					Type:        schema.TypeInt,
					Description: "MBytes of disk required by the task.",
					Optional:    true,
				},
			},
		},
	}
}

func loggerSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "defining the log behavior for the process.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"destination": {
					Type:        schema.TypeString,
					Description: "Destination of logs.",
					Default:     "file",
					Optional:    true,
				},
				"mode": {
					Type:        schema.TypeString,
					Description: "Mode of the logger.",
					Default:     "standard",
					Optional:    true,
				},
				"rotate": {
					Type:        schema.TypeList,
					Description: "An optional rotation policy.",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"log_size": {
								Type:        schema.TypeInt,
								Description: "Maximum size (in bytes) of an individual log file.",
								Default:     100,
								Optional:    true,
							},
							"backups": {
								Type:        schema.TypeInt,
								Description: "The maximum number of backups to retain.",
								Default:     5,
								Optional:    true,
							},
						},
					},
				},
			},
		},
	}
}

func mesosFetcherUri() *schema.Schema {

	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of URIs to fetch",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"uri": {
					Type:        schema.TypeString,
					Description: "the URI to fetch",
					Required:    true,
				},
				"extract": {
					Type:        schema.TypeBool,
					Description: "Extract the URI target",
					Optional:    true,
					Default:     false,
				},
				"cache": {
					Type:        schema.TypeBool,
					Description: "Cache the URI target",
					Default:     false,
					Optional:    true,
				},
			},
		},
	}
}
