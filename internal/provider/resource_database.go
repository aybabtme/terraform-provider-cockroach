package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
)

const (
	dbNameAttr  = "name"
	dbOwnerAttr = "owner"
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
			dbOwnerAttr: {
				Description: "Owner of the database.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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

	var id int
	err = conn.QueryRow(ctx, `SELECT id FROM crdb_internal.databases WHERE name = $1`, name).Scan(
		&id,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceDatabaseRead(ctx, d, meta)
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*pgx.Conn)
	if err := conn.Ping(ctx); err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	var (
		name  string
		owner string
	)
	err := conn.QueryRow(ctx, `SELECT name, owner FROM crdb_internal.databases WHERE id = $1`, id).Scan(
		&name,
		&owner,
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(dbNameAttr, name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(dbOwnerAttr, owner); err != nil {
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

	name := d.Get(dbNameAttr).(string)

	if name == "" {
		return diag.Errorf("database name can't be an empty string")
	}

	_, err := conn.Exec(ctx, `DROP DATABASE `+pq.QuoteIdentifier(name))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
