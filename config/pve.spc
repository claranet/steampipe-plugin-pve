connection "pve" {
  plugin = "pve"

  # Proxmox VE API URL (required)
  # api_url = "https://pve.example.com:8006/api2/json"

  # Authentication: API Token (recommended)
  # Create a token in Datacenter -> Permissions -> API Tokens
  # Format: user@realm!tokenname
  # api_token_id     = "root@pam!steampipe"
  # api_token_secret  = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

  # Authentication: Username/Password (alternative)
  # user     = "root@pam"
  # password = "your-password"

  # Skip TLS certificate verification (for self-signed certs)
  # tls_insecure = true
}

# Multi-cluster example: query multiple PVE clusters with aggregator
#
# connection "pve_dc1" {
#   plugin   = "pve"
#   api_url  = "https://pve-dc1.example.com:8006/api2/json"
#   api_token_id     = "monitor@pve!steampipe"
#   api_token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
#   tls_insecure     = true
# }
#
# connection "pve_dc2" {
#   plugin   = "pve"
#   api_url  = "https://pve-dc2.example.com:8006/api2/json"
#   api_token_id     = "monitor@pve!steampipe"
#   api_token_secret = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
#   tls_insecure     = true
# }
#
# connection "pve_all" {
#   plugin      = "pve"
#   type        = "aggregator"
#   connections = ["pve_dc1", "pve_dc2"]
# }
