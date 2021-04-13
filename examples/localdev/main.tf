provider "cockroach" {
  dsn = "postgresql://root@localhost:26257?sslmode=disable"
}

resource "cockroach_database" "helloworld" {
  name = "helloworld"
}
