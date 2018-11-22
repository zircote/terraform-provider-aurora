package aurora

import "github.com/hashicorp/terraform/helper/schema"

func slaCountPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		Optional:      true,
		ConflictsWith: []string{"sla_percentage_policy", "sla_coordinator_policy"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"duration_secs": {
					Type:        schema.TypeString,
					Description: "Minimum time duration a task needs to be RUNNING to be treated as active.",
					Optional:    true,
				},
				"task_count": {
					Type:        schema.TypeString,
					Description: "The number of active instances required every durationSecs.",
					Optional:    true,
				},
			},
		},
	}
}

func slaPercPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		ConflictsWith: []string{"sla_count_policy", "sla_coordinator_policy"},
		Optional:      true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"duration_secs": {
					Type:        schema.TypeString,
					Description: "Minimum time duration a task needs to be RUNNING to be treated as active.",
					Optional:    true,
				},
				"percentage": {
					Type:        schema.TypeFloat,
					Description: "The percentage of active instances required every durationSecs",
					Optional:    true,
				},
			},
		},
	}
}

func slaCoordinatorPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		ConflictsWith: []string{"sla_count_policy", "sla_percentage_policy"},
		Optional:      true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"coordinator_url": {
					Type:        schema.TypeString,
					Description: "The URL to the Coordinator service to be contacted before performing SLA affecting actions (job updates, host drains etc).",
					Optional:    true,
				},
				"status_key": {
					Type:        schema.TypeString,
					Description: "The field in the Coordinator response that indicates the SLA status for working on the task. (Default: drain)",
					Default:     "drain",
					Optional:    true,
				},
			},
		},
	}
}
