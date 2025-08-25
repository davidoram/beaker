package api

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/davidoram/beaker/internal/db"
	"github.com/davidoram/beaker/internal/utility"
	"github.com/davidoram/beaker/schemas"
	"github.com/nats-io/gnatsd/server"
	natsserver "github.com/nats-io/nats-server/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	server := RunServerOnPort(-1)
	defer server.Shutdown()

	nc, err := nats.Connect(server.Addr().String())
	require.NoError(t, err)
	defer nc.Close()

	pool := db.TestPostgresPool(t)
	defer pool.Close()

	compiler, err := utility.NewJSONSchemaCompiler(t.Context(), "../../schemas")
	require.NoError(t, err)

	app, err := StartNewApp(nc, pool, compiler)
	require.NoError(t, err)
	defer app.Stop()

	t.Run("test stockAddHandler", func(t *testing.T) {

		uniqueSku := fmt.Sprintf("sku-%d", time.Now().UnixNano())

		// Call the stockAddHandler with a valid request
		req := schemas.StockAddRequest{
			ProductSKU: uniqueSku,
			Quantity:   10,
		}
		reqBytes, err := json.Marshal(req)
		require.NoError(t, err)

		msg, err := nc.RequestWithContext(t.Context(), "stock.add", reqBytes)
		require.NoError(t, err)

		// Parse the response & check values
		resp := schemas.StockAddResponse{}
		err = json.Unmarshal(msg.Data, &resp)
		require.NoError(t, err)

		require.True(t, resp.OK)
		require.NotNil(t, resp.ProductSKU)
		require.NotNil(t, resp.Quantity)
		assert.Equal(t, uniqueSku, *resp.ProductSKU)
		assert.Equal(t, 10, *resp.Quantity)
	})

}

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}
