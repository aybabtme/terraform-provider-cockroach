package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
)

const (
	dbNameAttr = "name"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Database in a CockroachDB cluster.",

		CreateContext: resourceDatabaseCreate,
		ReadContext:   resourceDatabaseRead,
		UpdateContext: resourceDatabaseUpdate,
		DeleteContext: resourceDatabaseDelete,

		Schema: map[string]*schema.Schema{
			dbNameAttr: {
				Description: "Name of the database.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*pgx.Conn)
	if err := conn.Ping(ctx); err != nil {
		return diag.FromErr(err)
	}
	name := d.Get(dbNameAttr).(string)

	if name == "" {
		return diag.Errorf("database name can't be an empty string")
	}

	_, err := conn.Exec(ctx, `CREATE DATABASE `+pq.QuoteIdentifier(name))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceDatabaseRead(ctx, d, meta)
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*pgx.Conn)
	if err := conn.Ping(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*pgx.Conn)
	if err := conn.Ping(ctx); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(dbNameAttr) {
		oraw, nraw := d.GetChange(dbNameAttr)
		o := oraw.(string)
		n := nraw.(string)
		if n == "" {
			return diag.Errorf("database name can't be an empty string")
		}
		_, err := conn.Exec(ctx,
			`ALTER DATABASE `+
				pq.QuoteIdentifier(o)+
				` RENAME TO `+
				pq.QuoteIdentifier(n),
		)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(n)
	}

	return resourceDatabaseRead(ctx, d, meta)
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*pgx.Conn)
	if err := conn.Ping(ctx); err != nil {
		return diag.FromErr(err)
	}

	return diag.Errorf("not implemented")
}
