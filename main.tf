terraform {
  required_version = ">= 1.0"
}

locals {
  # Load metadata
  metadata = yamldecode(file("${path.module}/data/services/metadata.yaml"))

  # Load all YAML files from the services directory
  yaml_files = fileset("${path.module}/data/services", "*.yaml")

  # Filter out metadata.yaml
  service_files = [for f in local.yaml_files : f if f != "metadata.yaml"]

  # Load all services into a map
  services_raw = {
    for filename in local.service_files :
    trimsuffix(filename, ".yaml") => yamldecode(file("${path.module}/data/services/${filename}"))
  }

  # Create a map of service IDs to service data
  services_by_id = {
    for key, service in local.services_raw :
    service.id => service
  }

  # Extract all IP prefixes (all services, all IP versions)
  all_prefixes = flatten([
    for service in values(local.services_raw) :
    service.address_prefixes.all
  ])

  # Extract all IPv4 prefixes
  all_ipv4_prefixes = flatten([
    for service in values(local.services_raw) :
    service.address_prefixes.ipv4
  ])

  # Extract all IPv6 prefixes
  all_ipv6_prefixes = flatten([
    for service in values(local.services_raw) :
    service.address_prefixes.ipv6
  ])
}
