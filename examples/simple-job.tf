

resource aurora_role "www-data" {
  name = "www-data"
  cpu = 32.5
  ram_mb = 128
  disk_mb = 256000
  gpu = 0
}

variable "cluster_name" {
  default = "devcluster"
}

resource aurora_job "hello" {
  name = "hello"
  role = "${aurora_role.www-data.name}"
  cluster = "${var.cluster_name}"
  environment = "prod"
  task {
    resources {
      cpu = 1.0
      ram_mb = 128
      disk_mb = 128
    }
    process {
      name = "hello"
      cmdline =<<EOF
      while true; do
      echo hello world
      sleep 10
      done
      EOF
    }
  }
}