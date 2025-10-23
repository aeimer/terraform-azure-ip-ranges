module "azure_ip_ranges" {
  # Use either main or a specific version tag which is updates by e.g. renovate bot
  source = "git::https://github.com/aeimer/terraform-azure-ip-ranges.git?ref=v20251022173335.0.0-ga9de95d-18724697633-1"
}

# Output all IP addresses (both IPv4 and IPv6)
output "all_prefixes" {
  description = "All Azure IP address prefixes"
  value       = module.azure_ip_ranges.all_prefixes
}

# Output only IPv4 addresses
output "ipv4_prefixes" {
  description = "All Azure IPv4 address prefixes"
  value       = module.azure_ip_ranges.all_ipv4_prefixes
}

# Output only IPv6 addresses
output "ipv6_prefixes" {
  description = "All Azure IPv6 address prefixes"
  value       = module.azure_ip_ranges.all_ipv6_prefixes
}

# Output metadata
output "metadata" {
  description = "Metadata about the Azure Service Tags"
  value       = module.azure_ip_ranges.metadata
}

# Output prefix counts
output "counts" {
  description = "Count of IP prefixes"
  value       = module.azure_ip_ranges.prefix_counts
}
