variable "contact" {
  default = "test@test.com"
}

variable "environment" {
  default = "prod"
}

resource "aurora_role" "test" {
  name = "terraform-example-role"
  cpu = 6.0
  ram_mb = 256
  disk_mb = 1024
}

data "consul_keys" "my_hello_world" {
  datacenter = "nyc1"
  token = "abcd"

  key {
    name = "instance_count"
    path = "service/catalog/${aurora_role.test.name}/${var.environment}/my-hello-world"
    default = "5"
  }
}

resource "aurora_task" "hello" {
  name = "my-hello-world"
  role = "${aurora_role.test.name}"
  environment = "${var.environment}"
  contact = "${var.contact}"

  instances = "${data.consul_keys.my_hello_world.var.instance_count}"
  service = true

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
    }
  }
}