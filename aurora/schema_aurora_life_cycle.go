package aurora

import "github.com/hashicorp/terraform/helper/schema"

func lifeCycleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"http": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"port": {
								Type:        schema.TypeString,
								Description: "The named port to send POST commands. (Default: health)",
								Default:     "health",
								Optional:    true,
							},
							"graceful_shutdown_endpoint": {
								Type:        schema.TypeString,
								Description: "Endpoint to hit to indicate that a task should gracefully shutdown. (Default: /quitquitquit)",
								Default:     "/quitquitquit",
								Optional:    true,
							},
							"shutdown_endpoint": {
								Type:        schema.TypeString,
								Description: "Endpoint to hit to give a task its final warning before being killed. (Default: /abortabortabort)",
								Default:     "/abortabortabort",
								Optional:    true,
							},
							"graceful_shutdown_wait_secs": {
								Type:        schema.TypeInt,
								Description: "The amount of time (in seconds) to wait after hitting the graceful_shutdown_endpoint before proceeding with the task termination lifecycle. (Default: 5)",
								Default:     5,
								Optional:    true,
							},
							"shutdown_wait_secs": {
								Type:        schema.TypeInt,
								Description: "The amount of time (in seconds) to wait after hitting the shutdown_endpoint before proceeding with the task termination lifecycle. ",
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
