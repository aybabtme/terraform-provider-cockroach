package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackc/pgx/v4"
)

func dataSourceDatabase() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample data source in the Terraform provider scaffolding.",

		ReadContext: dataSourceDatabaseRead,

		Schema: map[string]*schema.Schema{
			dbNameAttr: {
				// This description is used by the documentation generator and the language server.
				Description: "Name of the database.",
				Type:        schema.TypeString,
				Required:    true,
			},
			dbOwnerAttr: {
				Description: "Owner of the database.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func dataSourceDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*pgx.Conn)
	if err := conn.Ping(ctx); err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	var (
		id    int
		owner string
	)
	err := conn.QueryRow(ctx, `SELECT id, owner FROM crdb_internal.databases WHERE name = $1`, name).Scan(
		&id,
		&owner,
	)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))
	if err := d.Set(dbNameAttr, name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(dbOwnerAttr, owner); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
