package aurora

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paypal/gorealis"
	"github.com/paypal/gorealis/gen-go/apache/aurora"
	"github.com/zircote/thermos-payload"
	"log"
)

func resourceAuroraJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAuroraJobCreate,
		Read:   resourceAuroraJobRead,
		//Exists: resourceAuroraJobExists,
		Update: resourceAuroraJobUpdate,
		Delete: resourceAuroraJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Job name.",
				Required:    true,
			},
			"role": {
				Type:        schema.TypeString,
				Description: "Job role account",
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Description: "Job environment, default devel. By default must be one of prod, devel, test or staging<number> but it can be changed by the Cluster operator using the scheduler option allowed_job_environments",
				Required:    true,
			},
			"contact": {
				Type:        schema.TypeString,
				Description: "Best email address to reach the owner of the job. For production jobs, this is usually a team mailing list.",
				Optional:    true,
			},
			"instances": {
				Type:        schema.TypeInt,
				Description: "Number of instances (sometimes referred to as replicas or shards) of the task to create.",
				Required:    true,
			},
			"cron_schedule": {
				Type:        schema.TypeString,
				Description: "Cron schedule in cron format. May only be used with non-service jobs. See Cron Jobs for more information.",
				Optional:    true,
			},
			"cron_collision_policy": {
				Type:        schema.TypeString,
				Description: "Policy to use when a cron job is triggered while a previous run is still active. KILL_EXISTING Kill the previous run, and schedule the new run CANCEL_NEW Let the previous run continue, and cancel the new run. (Default: KILL_EXISTING)",
				Default:     "KILL_EXISTING",
				Optional:    true,
			},
			"service": {
				Type:        schema.TypeBool,
				Description: "If True, restart tasks regardless of success or failure.",
				Default:     false,
				Optional:    true,
			},
			"max_task_failures": {
				Type:        schema.TypeInt,
				Description: "Maximum number of failures after which the task is considered to have failed (Default: 1) Set to -1 to allow for infinite failures",
				Default:     1,
				Optional:    true,
			},
			"tier": {
				Type:        schema.TypeString,
				Description: "Task tier type. The default scheduler tier configuration allows for 3 tiers: revocable, preemptible, and preferred. If a tier is not elected, Aurora assigns the task to a tier based on its choice of production (that is preferred for production and preemptible for non-production jobs). See the section on Configuration Tiers for more information.",
				Optional:    true,
			},
			"update_config":          updateConfigSchema(),
			"task":                   taskSchema(),
			"constraints":            constraintsSchema(),
			"healthcheck_http":       httpHealthCheckConfigSchema(),
			"healthcheck_shell":      shellHealthCheckConfigSchema(),
			"life_cycle":             lifeCycleSchema(),
			"announce":               announceSchema(),
			"partition_policy":       partitionPolicySchema(),
			"metadata":               metadataSchema(),
			"executor_config":        executorConfigSchema(),
			"sla_count_policy":       slaCountPolicySchema(),
			"sla_percentage_policy":  slaPercPolicySchema(),
			"sla_coordinator_policy": slaCoordinatorPolicySchema(),
			"container":              containerSchema(),
			"mesos_fetcher_uri":      mesosFetcherUri(),
		},
	}
}

func resourceAuroraJobUpdate(d *schema.ResourceData, meta interface{}) error {
	var payload *thermos_payload.ThermosTaskConfig
	job := realis.NewJob()
	taskConfig := job.TaskConfig()

	// Required Params
	if name, ok := d.GetOkExists("name"); ok {
		job.Name(name.(string))
		payload = thermos_payload.DefaultThermosTaskConfig(name.(string))
	} else {
		return fmt.Errorf("name is not set for aurora_job")
	}

	if role, ok := d.GetOkExists("role"); ok {
		job.Role(role.(string))
		payload.Role = role.(string)
	} else {
		return fmt.Errorf("role is not set for aurora_job")
	}

	if environment, ok := d.GetOkExists("environment"); ok {
		job.Environment(environment.(string))
		payload.Environment = environment.(string)
	} else {
		return fmt.Errorf("environment is not set for aurora_job")
	}
	//
	if instances, ok := d.GetOkExists("instances"); ok {
		i := instances.(int)
		job.InstanceCount(int32(i))
	}

	// Optional Params
	//
	if contact, ok := d.GetOkExists("contact"); ok {
		taskConfig.ContactEmail = resourceDataToStringPtr(contact)
	}
	//
	if service, ok := d.GetOkExists("service"); ok {
		job.IsService(resourceDataToBool(service))
		payload.Service = resourceDataToBool(service)
	}
	//
	if cronSchedule, ok := d.GetOkExists("cron_schedule"); ok {
		job.CronSchedule(cronSchedule.(string))
		if payload.Service {
			return fmt.Errorf("service must be false for cron jobs, got %v", payload.Service)
		}
		job.IsService(false)
		payload.Service = false
	}
	//
	if cronCollisionPolicy, ok := d.GetOkExists("cron_collision_policy"); ok {
		p, _ := aurora.CronCollisionPolicyFromString(resourceDataToString(cronCollisionPolicy))
		job.CronCollisionPolicy(p)
		payload.CronCollisionPolicy = resourceDataToString(cronCollisionPolicy)
	}
	//
	if maxTaskFailures, ok := d.GetOkExists("max_task_failures"); ok {
		job.MaxFailure(resourceDataToInt32(maxTaskFailures))
	}
	//
	if tier, ok := d.GetOkExists("tier"); ok {
		job.Tier(resourceDataToString(tier))
		//payload.Tier = resourceDataToString(tier)
	}

	//
	// Task Config

	if finalizationWait, ok := d.GetOkExists("task.0.finalization_wait"); ok {
		payload.Task.FinalizationWait = resourceDataToInt(finalizationWait)
	}
	if enableHooks, ok := d.GetOkExists("enable_hooks"); ok {
		payload.EnableHooks = resourceDataToBool(enableHooks)
	}
	if maxFailures, ok := d.GetOkExists("task.0.max_task_failures"); ok {
		payload.MaxTaskFailures = resourceDataToInt(maxFailures)
	}
	if d.Get("health_checker") != "" {
		health_checker := thermos_payload.DefaultHealthChecker()
		payload.HealthCheckConfig = health_checker
	}
	if maxConcurrency, ok := d.GetOkExists("task.0.max_concurrency"); ok {
		payload.Task.MaxConcurrency = resourceDataToInt(maxConcurrency)
	}

	//mesos_fetcher_uri
	if mesosFetcherUri, ok := d.GetOkExists("mesos_fetcher_uri"); ok {
		for _, value := range mesosFetcherUri.([]interface{}) {
			v := value.(map[string]interface{})
			extract := v["extract"].(bool)
			cache := v["cache"].(bool)
			uri := &aurora.MesosFetcherURI{
				Value:   v["uri"].(string),
				Extract: &extract,
				Cache:   &cache,
			}
			taskConfig.MesosFetcherUris[uri] = true
		}
	}
	// Process configuration
	for key, value := range d.Get("task.0.process").([]interface{}) {
		v := value.(map[string]interface{})
		//thermos_payload.Process{}
		proc := &thermos_payload.Process{
			Name:    v["name"].(string),
			Cmdline: v["cmdline"].(string),
			Daemon:  v["daemon"].(bool),
			Final:   v["final"].(bool),
		}

		if v["max_failures"] != "" {
			proc.MaxFailures = v["max_failures"].(int)
		}
		if v["ephemeral"].(bool) {
			proc.Ephemeral = v["ephemeral"].(bool)
		}
		if v["min_duration"] != "" {
			proc.MinDuration = v["min_duration"].(int)
		}
		// Process Logger
		for _, _l := range d.Get(fmt.Sprintf("task.0.process.%d.logger", key)).([]interface{}) {
			logger := _l.(map[string]interface{})
			if len(logger) > 0 {
				var rotate map[string]interface{}
				for _, v := range logger["rotate"].([]interface{}) {
					rotate = v.(map[string]interface{})
				}
				logSize := rotate["log_size"].(int)
				backups := rotate["backups"].(int)
				proc.Logger = thermos_payload.NewLogger(logger["destination"].(string), logger["mode"].(string),
					int64(logSize), int64(backups))
			}
		}
		payload.Task.AddProcess(proc)
	}

	// Resources
	cpu := d.Get("task.0.resources.0.cpu").(float64)
	disk_mb := d.Get("task.0.resources.0.disk_mb").(int)
	ram_mb := d.Get("task.0.resources.0.ram_mb").(int)
	payload.Task.Resources = &thermos_payload.Resources{
		CPU:  cpu,
		Disk: uint64(disk_mb * 1024 * 1024),
		RAM:  uint64(ram_mb * 1024 * 1024),
	}
	job.RAM(int64(ram_mb)).
		Disk(int64(disk_mb)).
		CPU(float64(cpu))

	// Executor Config
	job.ExecutorName(aurora.AURORA_EXECUTOR_NAME).
		ExecutorData(string(payload.String()))
	fmt.Println(payload.String())
	// TODO this is not working...
	d.Set("json", payload.String())

	// Services and Non-Service jobs must be created and started differently
	// Services are the only job type that can be "updated"
	switch {
	case taskConfig.IsService == true:
		err := runJobUpdate(job, meta)
		if err != nil {
			return err
		}
	case d.Get("cron_schedule").(string) != "":
		err := scheduleCronJob(job, meta)
		if err != nil {
			return err
		}
	case taskConfig.IsService != true && d.Get("cron_schedule").(string) == "":
		err := runJob(job, meta)
		if err != nil {
			return err
		}
	default:
		//TODO PANIC
	}
	d.SetId(fmt.Sprintf(
		"%s/%s/%s",
		d.Get("role").(string),
		d.Get("environment").(string),
		d.Get("name").(string)))

	return nil
}
func resourceAuroraJobExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}
func resourceAuroraJobRead(d *schema.ResourceData, meta interface{}) error {

	r := meta.(*Config).client
	_, result, err := r.GetJobs(d.Get("role").(string))
	if err != nil {
		return err
	}
	for jConfig := range result.GetConfigs() {
		if jConfig.Key.Name == d.Get("name").(string) && jConfig.Key.Environment == d.Get("environment").(string) {
			err = taskConfigToSchema(jConfig, d, meta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func resourceAuroraJobCreate(d *schema.ResourceData, meta interface{}) error {

	err := resourceAuroraJobUpdate(d, meta)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf(
		"%s/%s/%s",
		d.Get("role").(string),
		d.Get("environment").(string),
		d.Get("name").(string)))
	return nil
}

func resourceAuroraJobDelete(d *schema.ResourceData, meta interface{}) error {

	jKey := &aurora.JobKey{
		Role:        d.Get("role").(string),
		Environment: d.Get("environment").(string),
		Name:        d.Get("name").(string),
	}
	r := meta.(*Config).client
	fmt.Println("Killing job")

	monitor := realis.Monitor{r}
	resp, err := r.KillJob(jKey)
	if err != nil {
		log.Fatal(err)
	}

	if ok, err := monitor.Instances(jKey, 0, 5, 50); !ok || err != nil {
		log.Fatal("Unable to kill all instances of job")
	}
	fmt.Println(resp.String())
	d.SetId("")
	return nil
}

func runJobUpdate(job realis.Job, meta interface{}) error {

	config := meta.(*Config)
	settings := realis.NewUpdateSettings()

	resp, result, err := config.client.CreateService(job, settings)
	if err != nil {
		log.Println("error: ", err)
		log.Println("response: ", resp.String())
		log.Println("result: ", result.String())
		return err
	}
	monitor := realis.Monitor{config.client}
	fmt.Println(result.String())

	if ok, mErr := monitor.JobUpdate(*result.GetKey(), 5, 180); !ok || mErr != nil {
		_, err := config.client.AbortJobUpdate(*result.GetKey(), "Monitor timed out")
		_, err = config.client.KillJob(job.JobKey())
		if err != nil {
			log.Print(err)
			return err
		}
		log.Printf("ok: %v\n\n err: %v", ok, mErr)
	}
	return nil
}

func runJob(job realis.Job, meta interface{}) error {

	r := meta.(*Config).client
	monitor := realis.Monitor{Client: r}
	resp, err := r.CreateJob(job)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	fmt.Println(resp.String())

	if ok, mErr := monitor.Instances(job.JobKey(), job.GetInstanceCount(), 5, 50); !ok || mErr != nil {
		_, err := r.KillJob(job.JobKey())
		if err != nil {
			log.Fatalln(err)
		}
		log.Fatalf("ok: %v\n err: %v", ok, mErr)
		return mErr
	}
	return nil
}

func scheduleCronJob(job realis.Job, meta interface{}) error {

	r := meta.(*Config).client
	fmt.Println("Scheduling a Cron job")
	_, err := r.ScheduleCronJob(job)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func taskConfigToSchema(jConf *aurora.JobConfiguration, d *schema.ResourceData, meta interface{}) error {

	d.Set("service", jConf.TaskConfig.IsService)
	d.Set("priority", jConf.TaskConfig.Priority)
	d.Set("tier", jConf.TaskConfig.Tier)
	d.Set("contact", jConf.TaskConfig.ContactEmail)
	d.Set("max_task_failures", jConf.TaskConfig.MaxTaskFailures)
	d.Set("production", jConf.TaskConfig.Production)
	d.Set("instances", jConf.InstanceCount)
	d.Set("max_failure", jConf.TaskConfig.MaxTaskFailures)

	d.Set("cron_collision_policy", jConf.GetCronCollisionPolicy())
	if jConf.IsSetCronSchedule() {
		d.Set("cron_collision_policy", jConf.GetCronCollisionPolicy())
	}

	if jConf.TaskConfig.IsSetJob() {
		d.Set("name", jConf.TaskConfig.Job.Name)
		d.Set("environment", jConf.TaskConfig.Job.Environment)
		d.Set("role", jConf.TaskConfig.Job.Role)
	}
	if jConf.TaskConfig.IsSetExecutorConfig() {
		d.Set("executor_jConf.TaskConfig.name", jConf.TaskConfig.ExecutorConfig.Name)
		d.Set("executor_jConf.TaskConfig.data", jConf.TaskConfig.ExecutorConfig.Data)
	}

	if jConf.TaskConfig.IsSetPartitionPolicy() {
		d.Set("partition_policy.delay_secs", jConf.TaskConfig.PartitionPolicy.DelaySecs)
		d.Set("partition_policy.reschedule", jConf.TaskConfig.PartitionPolicy.Reschedule)
	}
	if jConf.TaskConfig.IsSetSlaPolicy() {
		switch {
		case jConf.TaskConfig.SlaPolicy.IsSetCoordinatorSlaPolicy():
			d.Set("sla_coordinator_policy.status_key", jConf.TaskConfig.SlaPolicy.CoordinatorSlaPolicy.StatusKey)
			d.Set("sla_coordinator_policy.coordinator_url", jConf.TaskConfig.SlaPolicy.CoordinatorSlaPolicy.CoordinatorUrl)
		case jConf.TaskConfig.SlaPolicy.IsSetCountSlaPolicy():
			d.Set("sla_count_policy.duration_secs", jConf.TaskConfig.SlaPolicy.CountSlaPolicy.DurationSecs)
			d.Set("sla_count_policy.count", jConf.TaskConfig.SlaPolicy.CountSlaPolicy.Count)
		case jConf.TaskConfig.SlaPolicy.IsSetPercentageSlaPolicy():
			d.Set("sla_percentage_policy.duration_secs", jConf.TaskConfig.SlaPolicy.PercentageSlaPolicy.DurationSecs)
			d.Set("sla_percentage_policy.percentage", jConf.TaskConfig.SlaPolicy.PercentageSlaPolicy.Percentage)
		}
	}
	if len(jConf.TaskConfig.Constraints) > 0 {
		i := 0
		for constraint := range jConf.TaskConfig.Constraints { // TODO finish this after rest
			d.Set(fmt.Sprintf("constraints.%d.name", i), constraint.Name)
			switch {
			case constraint.Constraint.IsSetValue():
				d.Set(fmt.Sprintf("constraints.%d", i), constraint.Constraint.Value)
			case constraint.Constraint.IsSetLimit():
				d.Set(fmt.Sprintf("constraints.%d", i), constraint.Constraint.Limit)
			}
		}
	}
	if jConf.TaskConfig.IsSetContainer() {
		switch {
		case jConf.TaskConfig.Container.IsSetMesos():
			if jConf.TaskConfig.Container.Mesos.IsSetImage() {
				switch {
				case jConf.TaskConfig.Container.Mesos.Image.IsSetAppc():
					d.Set("container.mesos.0.image.name", jConf.TaskConfig.Container.Mesos.Image.Appc.Name)
					d.Set("container.mesos.0.image.image_id", jConf.TaskConfig.Container.Mesos.Image.Appc.ImageId)
				case jConf.TaskConfig.Container.Mesos.Image.IsSetDocker():
					d.Set("container.mesos.0.image.name", jConf.TaskConfig.Container.Mesos.Image.Docker.Name)
					d.Set("container.mesos.0.image.tag", jConf.TaskConfig.Container.Mesos.Image.Docker.Tag)
				}
			}
			if jConf.TaskConfig.Container.Mesos.IsSetVolumes() {
				for key, vol := range jConf.TaskConfig.Container.Mesos.Volumes {
					d.Set(fmt.Sprintf("container.mesos.%d.volume.mode", key), vol.Mode)
					d.Set(fmt.Sprintf("container.mesos.%d.volume.container_path", key), vol.ContainerPath)
					d.Set(fmt.Sprintf("container.mesos.%d.volume.host_path", key), vol.HostPath)
				}
			}

		case jConf.TaskConfig.Container.IsSetDocker():
			d.Set("container.docker.image", jConf.TaskConfig.Container.Docker.Image)
			if jConf.TaskConfig.Container.Docker.IsSetParameters() {
				for key, param := range jConf.TaskConfig.Container.Docker.Parameters {
					d.Set(fmt.Sprintf("container.docker.parameters.%d.name", key), param.Name)
					d.Set(fmt.Sprintf("container.docker.parameters.%d.value", key), param.Value)
				}
			}
		}
	}

	tt := thermos_payload.ThermosTaskConfig{}
	tt.FromString(jConf.TaskConfig.ExecutorConfig.Data)
	d.Set("task.0.name", tt.Task.Name)
	d.Set("task.0.finalization_wait", tt.Task.FinalizationWait)
	d.Set("task.0.max_concurrency", tt.Task.MaxConcurrency)
	d.Set("task.0.max_failures", tt.Task.MaxFailures)

	d.Set("task.0.resources.0.cpu", tt.Task.Resources.CPU)
	d.Set("task.0.resources.0.disk_mb", tt.Task.Resources.Disk)
	d.Set("task.0.resources.0.ram_mb", tt.Task.Resources.RAM)

	for key, proc := range tt.Task.Processes {
		d.Set(fmt.Sprintf("task.0.process.%d.cmdline", key), proc.Cmdline)
		d.Set(fmt.Sprintf("task.0.process.%d.daemon", key), proc.Daemon)
		d.Set(fmt.Sprintf("task.0.process.%d.ephemeral", key), proc.Ephemeral)
		d.Set(fmt.Sprintf("task.0.process.%d.final", key), proc.Final)
		d.Set(fmt.Sprintf("task.0.process.%d.max_failures", key), proc.MaxFailures)
		d.Set(fmt.Sprintf("task.0.process.%d.min_duration", key), proc.MinDuration)
		d.Set(fmt.Sprintf("task.0.process.%d.name", key), proc.Name)
		if proc.Logger != nil {
			d.Set(fmt.Sprintf("task.0.process.%d.logger.mode", key), proc.Logger.Mode)
			d.Set(fmt.Sprintf("task.0.process.%d.logger.destination", key), proc.Logger.Destination)
			d.Set(fmt.Sprintf("task.0.process.%d.logger.rotate.backups", key), proc.Logger.Rotate.Backups)
			d.Set(fmt.Sprintf("task.0.process.%d.logger.rotate.log_size", key), proc.Logger.Rotate.LogSize)
		}
	}
	return nil
}
