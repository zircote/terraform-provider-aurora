package aurora

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAuroraResourceJob_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuroraResourceJob_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("aurora_job.hello",
						"name", "my-hello-world"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"role", "www-data"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"environment", "prod"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"contact", "test@test.com"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"instances", "1"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"tier", "preemptible"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"task.0.resources.#", "1"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"task.0.resources.0.cpu", "0.1"),
					resource.TestCheckResourceAttr("aurora_job.hello",
						"task.0.process.#", "1"),
					//resource.TestCheckResourceAttr("aurora_job.hello",
					//	"task.0.process.0.logger.#", "1"),
					//resource.TestCheckResourceAttr("aurora_job.hello",
					//	"task.0.process.0.logger.0.destination", "both"),
					//resource.TestCheckResourceAttr("aurora_job.hello",
					//	"task.0.process.0.logger.0.rotate.#", "1"),
					//resource.TestCheckResourceAttr("aurora_job.hello",
					//	"task.0.process.0.logger.0.rotate.0.log_size", "100"),
					//resource.TestCheckResourceAttr("aurora_job.hello",
					//	"task.0.json", testAuroraDataSourceThermosPayloadExpected_basic),
				),
			},
		},
	})
}

var testAuroraResourceJob_basic = `
variable "contact" {
  default = "test@test.com"
}


resource "aurora_job" "hello" {
  name = "my-hello-world"
  role = "www-data"
  environment = "prod"
  contact = "${var.contact}"

  instances = 1
  service = true
  tier = "preemptible"

  task {
    name = "hello_world"
    resources {
      cpu = 0.1
      ram_mb = "128"
      disk_mb = "40"
    }
    process {
      name = "hello"
      cmdline = <<EOF
while true; do
 echo hello
 sleep 10
done
EOF
      max_failures = 0
      daemon = false
      ephemeral = false
      min_duration = 5
      final = false
    }
  }
}

`

var testAuroraDataSourceThermosPayloadExpected_basic = `{"task":{"name":"my-hello-world","processes":[],"finalization_wait":30,"max_failures":1,"resources":{},"constraints":{"order":null}},"cluster":"develcluster","role":"www-data","environment":"prod","name":"my-hello-world","health_check_config":{"initial_interval_secs":5,"health_checker":{"http":{"endpoint":"/health","expected_response":"ok"}},"interval_secs":30,"timeout_secs":1,"max_consecutive_failures":1},"service":true,"max_task_failures":1,"cron_collision_policy":"KILL_EXISTING","lifecycle":{"http":{"graceful_shutdown_endpoint":"/quitquitquit","port":"health","shutdown_endpoint":"/abortabortabort"}},"tier":"preemtible","finalization_wait":30}`
