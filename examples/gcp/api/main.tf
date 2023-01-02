variable "config_path" {
  description = "path to a kubernetes config file"
  default     = "Hii there"
}

variable "env_path" {
  description = "a source to check hcl and env files"

}

resource "null_resource" "echo_apply" {


  provisioner "local-exec" {
    command = "echo ${var.config_path}"

  }
}

resource "null_resource" "echo_env" {


  provisioner "local-exec" {
    command = "echo ${var.env_path} here"

  }
}

resource "null_resource" "echo_diff" {


  provisioner "local-exec" {
    command = "echo ${var.env_path} here"

  }
}
