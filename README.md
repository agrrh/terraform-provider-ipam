# Terraform IPAM Provider

Dead simple IPAM provider, which stores IP pools and allocations in local file.

Works for those who'd like to use Git repository as [Source of Truth](https://en.wikipedia.org/wiki/Single_source_of_truth) for IP management.

## Usage

```terraform
provider "ipam" {
  file = "mycorp.ipam.json"
}

resource "ipam_pool" "internal_network" {
  cidr = "10.0.0.0/8"
}

#
# Availability zone
#

resource "ipam_allocation" "internal_az1" {
  pool_id = ipam_pool.internal.id
  size    = 10
}

# Also make it a pool to allocate from this particular IP range
resource "ipam_pool" "internal_az1" {
  cidr = ipam_allocation.internal_az1.cidr
}

# Single hosts


resource "ipam_allocation" "host_foo" {
  pool_id = ipam_pool.internal_az1.id
  size    = 32
}

# ... and so on
```

Then you (your CI, preferrably) run `terraform apply` and commit resulting `mycorp.ipam.json` to the repo.

## Features 😅

- Approvals powered with GitHub PRs or any platform of your choice
- Audit as simple as `cat ipam.tf` or `terraform state`
- Logs with `git log`
