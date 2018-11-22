package aurora

import "github.com/hashicorp/terraform/helper/schema"

func httpHealthCheckConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"healthcheck_shell"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"endpoint": {
					Type:        schema.TypeString,
					Default:     "/health",
					Optional:    true,
					Description: "HTTP endpoint to check (Default: /health)",
				},
				"expected_response": {
					Type:        schema.TypeString,
					Default:     "ok",
					Optional:    true,
					Description: "If not empty, fail the HTTP health check if the response differs. Case insensitive. (Default: ok)",
				},
				"expected_response_code": {
					Type:          schema.TypeInt,
					Optional:      true,
					Default:       0,
					Description:   "If not zero, fail the HTTP health check if the response code differs. (Default: 0)",
					ConflictsWith: []string{"healthcheck_http.expected_response"},
				},
			},
		},
	}
}

func shellHealthCheckConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"healthcheck_http"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"shell_command": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "An alternative to HTTP health checking. Specifies a shell command that will be executed. Any non-zero exit status will be interpreted as a health check failure.",
				},
			},
		},
	}
}
