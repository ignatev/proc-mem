provider "google" {
  region = "${var.region}"
  credentials = "${file(var.gcp_service_account)}"
  project = "${var.gcp_project_id}"
}

data "google_project" "current" {}

data "google_compute_default_service_account" "default" {}

resource "google_compute_network" "default" {
  name                    = "${var.network_name}"
  auto_create_subnetworks = "false"
}

resource "google_compute_subnetwork" "default" {
  name                     = "${var.network_name}"
  ip_cidr_range            = "${var.network_cidr}"
  network                  = "${google_compute_network.default.self_link}"
  region                   = "${var.region}"
  private_ip_google_access = true
}

resource "google_compute_firewall" "default" {
  name    = "demo-firewall"
  network = "${google_compute_network.default.name}"

  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "22", "5432"]
  }

}

resource "google_compute_instance" "default" {
  count                     = "${var.num_nodes}"
  name                      = "${var.name}-${count.index + 1}"
  zone                      = "${var.zone}"
  machine_type              = "${var.machine_type}"
  allow_stopping_for_update = true

  boot_disk {
    auto_delete = "${var.disk_auto_delete}"

    initialize_params {
      image = "${var.image_project}/${var.image_family}"
      size  = "${var.disk_size_gb}"
      type  = "${var.disk_type}"
    }
  }

  network_interface {
    subnetwork    = "${google_compute_subnetwork.default.name}"
    access_config = ["${var.access_config}"]
    address       = "${var.network_ip}"
  }

  metadata {
    sshKeys = "${var.gcp_ssh_user}:${file(var.gcp_ssh_pub_key_file)}"
  }

}