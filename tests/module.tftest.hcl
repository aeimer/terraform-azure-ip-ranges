run "module_provides_data" {
  command = plan

  assert {
    condition     = length(output.all_prefixes) > 0
    error_message = "all_prefixes must not be empty"
  }

  assert {
    condition     = length(output.all_ipv4_prefixes) > 0
    error_message = "all_ipv4_prefixes must not be empty"
  }

  assert {
    condition     = length(output.all_ipv6_prefixes) > 0
    error_message = "all_ipv6_prefixes must not be empty"
  }

  assert {
    condition     = length(output.service_ids) > 0
    error_message = "service_ids must not be empty"
  }

  assert {
    condition     = output.prefix_counts.total == output.prefix_counts.ipv4 + output.prefix_counts.ipv6
    error_message = "prefix_counts.total must equal ipv4 + ipv6"
  }

  assert {
    condition     = output.prefix_counts.total == length(output.all_prefixes)
    error_message = "prefix_counts.total must match length of all_prefixes"
  }

  assert {
    condition     = output.metadata.cloud != ""
    error_message = "metadata.cloud must be set"
  }

  assert {
    condition     = output.metadata.service_count > 0
    error_message = "metadata.service_count must be greater than 0"
  }

  assert {
    condition     = length(output.service_ids) == output.metadata.service_count
    error_message = "number of service_ids must match metadata.service_count"
  }

  assert {
    condition     = length(output.empty_service_files) == 0
    error_message = "found empty YAML files in data/services (likely broken generator run): ${join(", ", output.empty_service_files)}"
  }
}
