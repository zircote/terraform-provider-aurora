package aurora

import "github.com/hashicorp/terraform/helper/schema"

func processSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "List of Process objects bound to this task. (Required)",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "Process name",
					Required:    true,
				},
				"cmdline": {
					Type:        schema.TypeString,
					Description: "Command line",
					Required:    true,
				},
				"max_failures": {
					Type:        schema.TypeInt,
					Description: "Maximum process failures",
					Default:     1,
					Optional:    true,
				},
				"daemon": {
					Type:        schema.TypeBool,
					Description: "When True, this is a daemon process.",
					Default:     false,
					Optional:    true,
				},
				"ephemeral": {
					Type:        schema.TypeBool,
					Description: "When True, this is an ephemeral process.",
					Default:     false,
					Optional:    true,
				},
				"min_duration": {
					Type:        schema.TypeInt,
					Description: "Minimum duration between process restarts in seconds.",
					Default:     5,
					Optional:    true,
				},
				"final": {
					Type:        schema.TypeBool,
					Description: "When True, this process is a finalizing one that should run last. (Default: False)",
					Default:     false,
					Optional:    true,
				},
				"logger": loggerSchema(),
			},
		},
	}
}
