resource "ipam_pool" "home" {
  cidr = "10.0.0.0/8"
}

resource "ipam_allocation" "foo" {
  pool_id = ipam_pool.home.id
  size    = 16
}
