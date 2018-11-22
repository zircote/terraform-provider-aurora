package aurora

import "github.com/hashicorp/terraform/helper/schema"

func taskResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Task name",
				Required:    true,
			},
			"max_failures": {
				Type:        schema.TypeInt,
				Description: "Maximum process failures before being considered failed (Default: 1)",
				Default:     1,
				Optional:    true,
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of concurrent processes (Default: 0, unlimited concurrency.)",
				Default:     0,
			},
			"finalization_wait": {
				Type:        schema.TypeInt,
				Description: "Amount of time allocated for finalizing processes, in seconds. (Default: 30)",
				Default:     30,
				Optional:    true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources": resourceSchema(),
			"process":   processSchema(),
		},
	}
}

func taskSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem:     taskResource(),
	}
}
