---
page_title: "cockroach_database Data Source - terraform-provider-cockroach"
subcategory: ""
description: |-
  Sample data source in the Terraform provider scaffolding.
---

# Data Source `cockroach_database`

Sample data source in the Terraform provider scaffolding.

## Example Usage

```terraform
data "cockroach_database" "example" {
  name = "foo"
}
```

## Schema

### Required

- **name** (String) Name of the database.

### Optional

- **id** (String) The ID of this resource.


