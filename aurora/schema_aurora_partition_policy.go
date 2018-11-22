package aurora

import "github.com/hashicorp/terraform/helper/schema"

func partitionPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"reschedule": {
					Type:        schema.TypeBool,
					Description: "Whether or not to reschedule when running tasks become partitioned (Default: True)",
					Default:     true,
					Optional:    true,
				},
				"delay_secs": {
					Type:        schema.TypeInt,
					Default:     0,
					Description: "How long to delay transitioning to LOST when running tasks are partitioned. (Default: 0)",
					Optional:    true,
				},
			},
		},
	}
}
