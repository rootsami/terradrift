variable "config_path" {
  description = "path to a kubernetes config file"
  default     = "Hii there"
}


resource "null_resource" "echo_apply" {
  #triggers = {
  #  config_contents = var.config_path)
  #  
  #}

  provisioner "local-exec" {
    command = "echo ${var.config_path}"

  }
}


terraform {
  required_version = ">=1.2.7"
}
