Terraform Provider
==================


Apache Aurora

# Install

```bash
git clone https://github.com/zircote/terraform-provider-aurora

cd terraform-provider-aurora

go build -o ${TF_PROJECT_DIR}/terraform.d/plugins/${ARCH}/terraform-provider-aurora

cd ${TF_PROJECT_DIR}

terraform init

```


# Aurora Role/Quota Example

Using the Apache Aurora projects Vagrant box, the following example below demonstrates how to set and manage the quota for the `www-data` role.

```hcl-terraform
provider "aurora" {
  server_url = "http://192.168.33.7:8081/scheduler"

}

resource "aurora_role" "www-data" {
  count = 1
  role = "www-data"
  cpu = 12.0
  ram_mb = 256
  disk_mb = 1024
}

```

# Apache Aurora Job Example

```hcl-terraform
variable "contact" {
  default = "test@test.com"
}

resource "aurora_role" "test" {
  name = "terraform-test-role"
  cpu = 6.0
  ram_mb = 256
  disk_mb = 1024
}


resource "aurora_task" "hello" {
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
```

# External Storage of Instance Counts

```hcl-terraform
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
```