# Golang demo app for exposing memory consumption from linux host

Application connsists of two parts: collect and expose.

## Expose

This script reads name and vmrss values from /proc/<pid>/status, stores them in a map.
Gathered information exposed using web-server on 8080 port in Prometheus format (tagged by process name and hostname)
There are free memory metric, free memory percentage metric and per-proc consumption metrics.

## Collect

This script reads hostname lists from plain-text file and starts scraping metrics for each host.
Collected metrics are send to the statsd server.

## Demo

There are 7 hosts deployed in GCP using Terraform and Ansible.

* Graphite web interface: http://35.232.157.60/dashboard/per-proc
* Raw metrics: http://35.232.157.60:8080/metrics

### Notes

Some values hardcoded due to demo purposes.