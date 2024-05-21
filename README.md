# Terraform IPAM Provider

Dead simple IPAM provider, which stores pools and allocation in local file.

Works for those who'd like to use Git repository as Source of Truth for IP management.

## Usage

```terraform
provider "ipam" {
  file = "mycorp.ipam.json"
}

resource "ipam_pool" "mycorp_internal" {
  cidr = "10.0.0.0/8"
}

resource "ipam_allocation" "zoneA" {
  pool_id = ipam_pool.mycorp_internal.id
  size    = 16
}

resource "ipam_allocation" "zoneB" {
  pool_id = ipam_pool.mycorp_internal.id
  size    = 16
}

# ... and so on
```

Then you (your CI, preferrably) run `terraform apply` and commit resulting `mycorp.ipam.json` to the repo.

## Features 😅

- Approvals powered with GitHub PRs or any platform of your choice
- Audit as simple as `cat ipam.tf` or `terraform state`
- Logs with `git log`
