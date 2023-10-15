resource "ipam_allocation" "foo" {
  provider = ipam.lan

  pool_id = ipam_pool.lan.id
  size    = 16
}

resource "ipam_allocation" "bar" {
  provider = ipam.lan

  pool_id = ipam_pool.lan.id
  size    = 24
}

resource "ipam_allocation" "baz" {
  provider = ipam.lan

  pool_id = ipam_pool.lan.id
  size    = 32
}
