# Terraform Azure IP Ranges Module

A Terraform/OpenTofu module that provides all IP address ranges used by Azure services. The data is automatically updated nightly from Microsoft's official ServiceTags JSON.

## Features

- **Comprehensive IP ranges**: Access all Azure service IP addresses (IPv4 and IPv6)
- **Service-specific filtering**: Get IP ranges for specific Azure services by ID
- **IP version filtering**: Filter by IPv4 only, IPv6 only, or both
- **Automatic updates**: Data is updated nightly via GitHub Actions
- **No external dependencies**: All data is pre-generated and stored as YAML files

## Usage

### Get All IP Addresses

```hcl
module "azure_ip_ranges" {
  source = "github.com/aeimer/terraform-azure-ip-ranges"
}

output "all_azure_ips" {
  value = module.azure_ip_ranges.all_prefixes
}

output "ipv4_only" {
  value = module.azure_ip_ranges.all_ipv4_prefixes
}

output "ipv6_only" {
  value = module.azure_ip_ranges.all_ipv6_prefixes
}
```

### Get IP Ranges for a Specific Service

```hcl
module "azure_ip_ranges" {
  source = "github.com/aeimer/terraform-azure-ip-ranges"
}

locals {
  azure_portal = module.azure_ip_ranges.services["AzurePortal"]
}

output "azure_portal_ipv4" {
  value = local.azure_portal.address_prefixes.ipv4
}

output "azure_portal_ipv6" {
  value = local.azure_portal.address_prefixes.ipv6
}
```

### Get Regional Service IP Ranges

```hcl
module "azure_ip_ranges" {
  source = "github.com/aeimer/terraform-azure-ip-ranges"
}

# Regional services use the format: ServiceName.RegionName
locals {
  sweden_portal = module.azure_ip_ranges.services["AzurePortal.SwedenCentral"]
}

output "sweden_portal_ips" {
  value = local.sweden_portal.address_prefixes.all
}
```

### List All Available Services

```hcl
module "azure_ip_ranges" {
  source = "github.com/aeimer/terraform-azure-ip-ranges"
}

output "all_service_ids" {
  value = module.azure_ip_ranges.service_ids
}

output "azure_portal_services" {
  value = [
    for id in module.azure_ip_ranges.service_ids :
    id if length(regexall("^AzurePortal", id)) > 0
  ]
}
```

## Module Outputs

| Name | Description |
|------|-------------|
| `metadata` | Metadata about the ServiceTags data (change number, cloud, service count, generated date) |
| `all_prefixes` | All IP address prefixes from all Azure services (both IPv4 and IPv6) |
| `all_ipv4_prefixes` | All IPv4 address prefixes from all Azure services |
| `all_ipv6_prefixes` | All IPv6 address prefixes from all Azure services |
| `services` | Map of all services by ID with their metadata and address prefixes |
| `service_ids` | Sorted list of all available Azure service IDs |
| `prefix_counts` | Count of IP prefixes by type (total, ipv4, ipv6) |

## Service Structure

Each service in the `services` output has the following structure:

```hcl
{
  id = "AzurePortal"
  name = "AzurePortal"
  metadata = {
    change_number        = 53
    region              = ""
    platform            = "Azure"
    system_service      = "AzurePortal"
    network_features    = ["API", "NSG", "UDR", "FW"]
    global_change_number = 373
    cloud               = "Public"
  }
  address_prefixes = {
    all   = ["4.145.74.52/30", "2603:1000:4::10c/126", ...]
    ipv4  = ["4.145.74.52/30", ...]
    ipv6  = ["2603:1000:4::10c/126", ...]
    counts = {
      total = 324
      ipv4  = 200
      ipv6  = 124
    }
  }
}
```

## Examples

See the [examples](./examples/) directory for complete examples:

- [all-ip-addresses](./examples/all-ip-addresses/) - Get all Azure IP addresses
- [specific-service](./examples/specific-service/) - Get IP ranges for a specific service
- [regional-service](./examples/regional-service/) - Get IP ranges for a regional service
- [list-services](./examples/list-services/) - List and filter available services
- [firewall-rules](./examples/firewall-rules/) - Practical example for firewall configuration

## How It Works

1. **Data Source**: Microsoft publishes ServiceTags JSON at https://www.microsoft.com/en-us/download/details.aspx?id=56519
2. **Nightly Updates**: A GitHub Action runs every night at 2:00 AM UTC
3. **YAML Generation**: A Go script converts the JSON into individual YAML files per service
4. **Terraform Module**: The module reads these YAML files and provides structured outputs

## Data Updates

The module data is automatically updated through GitHub Actions:

- **Schedule**: Runs nightly at 2:00 AM UTC
- **Process**:
  1. Scrapes the Microsoft download page for the latest JSON URL
  2. Downloads and validates the new ServiceTags JSON
  3. Compares change numbers with the current version
  4. If changes detected, regenerates all YAML files
  5. Commits changes directly to the main branch
- **Manual Trigger**: You can manually trigger the workflow from the Actions tab

## Build the Generator

```bash
# Build the Go generator
cd generate
go build -o generator .
```

## Contributing

Contributions are welcome!
Please open an issue or submit a pull request.

## Maintainer

Alexander Eimer ([@aeimer](https://github.com/aeimer))
