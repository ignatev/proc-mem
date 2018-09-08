variable "region" {
  default = "us-central1"
}

variable "zone" {
  default = "us-central1-b"
}

variable "network_cidr" {
  default = "10.127.0.0/20"
}

variable "network_name" {
  default = "chef-managed-net"
}


variable "name" {
  description = "Name prefix for the nodes"
  default     = "chef-managed"
}

variable "num_nodes" {
  description = "Number of nodes to create"
  default     = 2
}

variable "image_family" {
  default = "ubuntu-1604-lts"
}

variable "image_project" {
  default = "ubuntu-os-cloud"
}

variable "disk_auto_delete" {
  description = "Whether or not the disk should be auto-deleted."
  default     = true
}

variable "disk_type" {
  description = "The GCE disk type. Can be either pd-ssd, local-ssd, or pd-standard."
  default     = "pd-ssd"
}

variable "disk_size_gb" {
  description = "The size of the image in gigabytes."
  default     = 30
}

variable "network_ip" {
  description = "Set the network IP of the instance. Useful only when num_nodes is equal to 1."
  default     = ""
}

variable "machine_type" {
    default = "n1-standard-1"
}

variable "access_config" {
    default = {}
}

variable "gcp_service_account" {
    default = ".gcp-creds"
}

variable "gcp_ssh_user" {
    default = "chef-demo"
}

variable "gcp_ssh_pub_key_file" {
    default = "gcp-key.pub"
}

variable "gcp_project_id" {
    default = "eternal-ship-213716"
}