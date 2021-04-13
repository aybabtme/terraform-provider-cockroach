---
page_title: "cockroach Provider"
subcategory: ""
description: |-
  
---

# cockroach Provider



## Example Usage

```terraform
provider "cockroach" {
  dsn = "postgresql://root@localhost:26257?sslmode=disable"
}
```

## Schema

### Required

- **dsn** (String) DSN to connect to the Cockroach cluster.
