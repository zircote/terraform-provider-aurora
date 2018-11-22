package aurora

import "github.com/hashicorp/terraform/helper/schema"

func updateConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"batch_size": {
					Type:        schema.TypeInt,
					Description: "Maximum number of shards to be updated in one iteration (Default: 1)",
					Default:     1,
					Optional:    true,
				},
				"watch_secs": {
					Type:        schema.TypeInt,
					Description: "Minimum number of seconds a shard must remain in RUNNING state before considered a success (Default: 45)",
					Default:     45,
					Optional:    true,
				},
				"max_per_shard_failures": {
					Type:        schema.TypeInt,
					Description: "Maximum number of restarts per shard during update. Increments total failure count when this limit is exceeded. (Default: 0)",
					Default:     0,
					Optional:    true,
				},
				"max_total_failures": {
					Type:        schema.TypeInt,
					Description: "Maximum number of shard failures to be tolerated in total during an update. Cannot be greater than or equal to the total number of tasks in a job. (Default: 0)",
					Default:     0,
					Optional:    true,
				},
				"rollback_on_failure": {
					Type:        schema.TypeBool,
					Description: "When False, prevents auto rollback of a failed update (Default: True)",
					Default:     true,
					Optional:    true,
				},
				"wait_for_batch_completion": {
					Type:        schema.TypeBool,
					Description: "When True, all threads from a given batch will be blocked from picking up new instances until the entire batch is updated. This essentially simulates the legacy sequential updater algorithm. (Default: False)",
					Default:     false,
					Optional:    true,
				},
				"pulse_interval_secs": {
					Type:        schema.TypeInt,
					Description: "Indicates a coordinated update. If no pulses are received within the provided interval the update will be blocked. Beta-updater only. Will fail on submission when used with client updater. (Default: None)",
					Default:     nil,
					Optional:    true,
				},
				"sla_aware": {
					Type:        schema.TypeBool,
					Description: "When True, updates will only update an instance if it does not break the task’s specified SLA Requirements. (Default: None)",
					Default:     nil,
					Optional:    true,
				},
			},
		},
	}
}
