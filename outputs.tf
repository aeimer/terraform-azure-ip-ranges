# Metadata outputs
output "metadata" {
  description = "Metadata about the Azure Service Tags data including change number and cloud"
  value = {
    change_number = local.metadata.change_number
    cloud         = local.metadata.cloud
    service_count = local.metadata.service_count
    generated_at  = local.metadata.generated_at
  }
}

# All IP prefixes (all services combined)
output "all_prefixes" {
  description = "All IP address prefixes (both IPv4 and IPv6) from all Azure services"
  value       = local.all_prefixes
}

output "all_ipv4_prefixes" {
  description = "All IPv4 address prefixes from all Azure services"
  value       = local.all_ipv4_prefixes
}

output "all_ipv6_prefixes" {
  description = "All IPv6 address prefixes from all Azure services"
  value       = local.all_ipv6_prefixes
}

# Service information
output "services" {
  description = "Map of all Azure services by ID with their metadata and IP prefixes"
  value       = local.services_by_id
}

output "service_ids" {
  description = "List of all available Azure service IDs"
  value       = sort(keys(local.services_by_id))
}

# Statistics
output "prefix_counts" {
  description = "Count of IP prefixes by type"
  value = {
    total = length(local.all_prefixes)
    ipv4  = length(local.all_ipv4_prefixes)
    ipv6  = length(local.all_ipv6_prefixes)
  }
}
