package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/cockroachdb/cockroach-go/testserver"
	_ "github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgx/v4"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"cockroach": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	ctx := context.Background()
	ts, err := testserver.NewTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Stop()

	conn, err := pgx.Connect(ctx, ts.PGURL().String())
	require.NoError(t, err)
	require.NoError(t, conn.Ping(ctx))
	require.NoError(t, conn.Close(ctx))
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
