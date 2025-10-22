module "azure_ip_ranges" {
  source = "../.."
}

# List all available service IDs
output "all_service_ids" {
  description = "List of all available Azure service IDs"
  value       = module.azure_ip_ranges.service_ids
}

output "service_count" {
  description = "Total number of Azure services"
  value       = length(module.azure_ip_ranges.service_ids)
}

# Filter services by pattern (e.g., all AzurePortal services)
output "azure_portal_services" {
  description = "All AzurePortal related service IDs"
  value       = [for id in module.azure_ip_ranges.service_ids : id if length(regexall("^AzurePortal", id)) > 0]
}

# Get services with IPv6 support
output "services_with_ipv6" {
  description = "Services that have IPv6 prefixes"
  value = [
    for id, service in module.azure_ip_ranges.services :
    id if service.address_prefixes.counts.ipv6 > 0
  ]
}
