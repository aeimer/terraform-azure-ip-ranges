module "azure_ip_ranges" {
  source = "../.."
}

# Example: Create firewall rules for Azure Storage service
locals {
  storage_service = module.azure_ip_ranges.services["Storage"]

  # Get only IPv4 ranges for firewall rules
  storage_ipv4_ranges = local.storage_service.address_prefixes.ipv4
}

# Example output showing how to use this for firewall configuration
output "storage_ipv4_for_firewall" {
  description = "IPv4 ranges for Azure Storage that can be used in firewall rules"
  value       = local.storage_ipv4_ranges
}

output "storage_range_count" {
  description = "Number of IPv4 ranges for Azure Storage"
  value       = length(local.storage_ipv4_ranges)
}

# Example: Filter services by region
locals {
  # Get all service IDs that contain "WestEurope"
  west_europe_services = [
    for service_id in module.azure_ip_ranges.service_ids :
    service_id if can(regex("WestEurope", service_id))
  ]
}

output "west_europe_service_ids" {
  description = "Service IDs for West Europe region"
  value       = local.west_europe_services
}

# Example: Get combined IP ranges for multiple services
locals {
  # Combine IP ranges from multiple services
  management_services = ["AzureResourceManager", "AzurePortal", "AzureActiveDirectory"]

  management_ipv4_ranges = flatten([
    for service_id in local.management_services :
    module.azure_ip_ranges.services[service_id].address_prefixes.ipv4
    if contains(keys(module.azure_ip_ranges.services), service_id)
  ])
}

output "management_ipv4_combined" {
  description = "Combined IPv4 ranges for Azure management services"
  value       = local.management_ipv4_ranges
}
