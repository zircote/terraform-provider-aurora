variable "contact" {
  default = "test@test.com"
}
provider "aurora" {
  server_url = "http://127.0.0.1:8081"
}

resource "aurora_role" "test" {
  name = "terraform-test-role"
  cpu = 6.0
  ram_mb = 256
  disk_mb = 1024
}

//
resource "aurora_job" "hello" {
  name = "my-hello-world"
  role = "${aurora_role.test.name}"
  environment = "prod"
  cluster = "devcluster"
  contact = "${var.contact}"

  instances = 3
  service = true
  tier = "preemtible"

  task {
    name = "hello_world"
    resources {
      cpu = 1.0
      ram_mb = 8
      disk_mb = 4
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
      ephemeral = true
      min_duration = 5
      final = false
      logger {
        destination = "stderr"
        mode = "standard"
        rotate {
          log_size = 100
          backups = 5
        }
      }
    }
    process {
      name = "finalizer"
      cmdline = "echo DONE; sleep 10"
      final = true
    }
  }
}