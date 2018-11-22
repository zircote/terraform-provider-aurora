package aurora

import "github.com/hashicorp/terraform/helper/schema"

func containerSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"mesos":  mesosContainerSchema(),
				"docker": dockerContainerSchema(),
			},
		},
	}
}

func dockerContainerSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		Description:   "A Docker container to use (via Docker engine)",
		ConflictsWith: []string{"container.mesos"},
		Optional:      true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"image": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the docker image to execute. If the image does not exist locally it will be pulled with docker pull",
				},
				"parameters": dockerContainerParameterSchema(),
			},
		},
	}
}

func dockerContainerParameterSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Additional parameters to pass to the Docker engine",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the docker parameter. E.g. volume",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "The value of the parameter. E.g. /usr/local/bin:/usr/bin:rw",
					Required:    true,
				},
			},
		},
	}
}

func mesosContainerSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		MaxItems:      1,
		Description:   "A native Mesos container to use.",
		ConflictsWith: []string{"container.docker"},
		Optional:      true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"image":   mesosContainerImageSchema(),
				"volumes": mesosContainerVolumeSchema(),
			},
		},
	}
}

func mesosContainerImageSchema() *schema.Schema {
	return &schema.Schema{ // container.mesos.image.
		Type:        schema.TypeList,
		Description: "An optional filesystem image to use within this container.",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "The name of the appc or docker image.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"image_id": {
					Type:          schema.TypeString,
					Description:   "The image id of the appc image.",
					ConflictsWith: []string{"container.mesos.image.tag"},
					Optional:      true,
				},
				"tag": {
					Type:          schema.TypeString,
					Description:   "The tag that identifies the docker image.",
					ConflictsWith: []string{"container.mesos.image.image_id"},
					Optional:      true,
				},
			},
		},
	}
}

func mesosContainerVolumeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"container_path": {
					Type:        schema.TypeString,
					Description: "Path on the host to mount.",
					Required:    true,
				},
				"host_path": {
					Type:        schema.TypeString,
					Description: "Mount point in the container.",
					Required:    true,
				},
				"mode": {
					Type:        schema.TypeString,
					Description: "Mode of the mount, can be ‘RW’ or 'RO’",
					Required:    true,
				},
			},
		},
	}
}
