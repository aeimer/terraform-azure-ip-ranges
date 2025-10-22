module "azure_ip_ranges" {
  source = "../.."
}

# Get IP ranges for a specific service (AzurePortal)
locals {
  azure_portal_service = module.azure_ip_ranges.services["AzurePortal"]
}

output "azure_portal_all_prefixes" {
  description = "All IP prefixes for AzurePortal"
  value       = local.azure_portal_service.address_prefixes.all
}

output "azure_portal_ipv4_prefixes" {
  description = "IPv4 prefixes for AzurePortal"
  value       = local.azure_portal_service.address_prefixes.ipv4
}

output "azure_portal_ipv6_prefixes" {
  description = "IPv6 prefixes for AzurePortal"
  value       = local.azure_portal_service.address_prefixes.ipv6
}

output "azure_portal_metadata" {
  description = "Metadata for AzurePortal service"
  value       = local.azure_portal_service.metadata
}
