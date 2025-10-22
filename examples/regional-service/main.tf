module "azure_ip_ranges" {
  source = "../.."
}

# Get IP ranges for a regional service (AzurePortal.SwedenCentral)
locals {
  # Check if the regional service exists
  has_sweden_central = contains(keys(module.azure_ip_ranges.services), "AzurePortal.SwedenCentral")

  sweden_portal_service = local.has_sweden_central ? module.azure_ip_ranges.services["AzurePortal.SwedenCentral"] : null
}

output "azure_portal_sweden_exists" {
  description = "Whether AzurePortal.SwedenCentral service exists"
  value       = local.has_sweden_central
}

output "azure_portal_sweden_ipv4" {
  description = "IPv4 prefixes for AzurePortal in Sweden Central"
  value       = local.has_sweden_central ? local.sweden_portal_service.address_prefixes.ipv4 : []
}

output "azure_portal_sweden_ipv6" {
  description = "IPv6 prefixes for AzurePortal in Sweden Central"
  value       = local.has_sweden_central ? local.sweden_portal_service.address_prefixes.ipv6 : []
}

output "azure_portal_sweden_all" {
  description = "All IP prefixes for AzurePortal in Sweden Central"
  value       = local.has_sweden_central ? local.sweden_portal_service.address_prefixes.all : []
}
