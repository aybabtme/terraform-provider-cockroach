---
page_title: "cockroach_database Resource - terraform-provider-cockroach"
subcategory: ""
description: |-
  Database in a CockroachDB cluster.
---

# Resource `cockroach_database`

Database in a CockroachDB cluster.

## Example Usage

```terraform
resource "cockroach_database" "example" {
  name = "foo"
}
```

## Schema

### Required

- **name** (String) Name of the database.

### Optional

- **id** (String) The ID of this resource.
- **owner** (String) Owner of the database.


