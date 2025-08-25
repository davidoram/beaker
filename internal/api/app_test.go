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

	t.Run("add stock", func(t *testing.T) {

		uniqueSku := fmt.Sprintf("sku-%d", time.Now().UnixNano())

		resp := addStock(t, nc, uniqueSku, 10)

		require.True(t, resp.OK)
		assert.Equal(t, uniqueSku, *resp.ProductSKU)
		assert.Equal(t, 10, *resp.Quantity)
	})

	t.Run("add lots of stock", func(t *testing.T) {

		uniqueSku := fmt.Sprintf("sku-%d", time.Now().UnixNano())

		for i := 0; i < 100; i++ {
			resp := addStock(t, nc, uniqueSku, 25)
			require.True(t, resp.OK)
			assert.Equal(t, uniqueSku, *resp.ProductSKU)
			assert.Equal(t, 25*(i+1), *resp.Quantity)
		}
	})

	t.Run("remove stock", func(t *testing.T) {

		uniqueSku := fmt.Sprintf("sku-%d", time.Now().UnixNano())

		addStock(t, nc, uniqueSku, 10)
		resp := removeStock(t, nc, uniqueSku, 7)

		require.True(t, resp.OK)
		assert.Equal(t, uniqueSku, *resp.ProductSKU)
		assert.Equal(t, 3, *resp.Quantity)
	})

	t.Run("remove stock below zero", func(t *testing.T) {

		uniqueSku := fmt.Sprintf("sku-%d", time.Now().UnixNano())

		addStock(t, nc, uniqueSku, 10)
		resp := removeStock(t, nc, uniqueSku, 11)

		require.False(t, resp.OK)
		assert.Equal(t, fmt.Sprintf("stock level cannot go below zero for %s", uniqueSku), *resp.Error)
	})

	t.Run("get stock balance", func(t *testing.T) {

		uniqueSku := fmt.Sprintf("sku-%d", time.Now().UnixNano())

		resp := getStock(t, nc, uniqueSku)
		require.True(t, resp.OK)
		assert.Equal(t, uniqueSku, *resp.ProductSKU)
		assert.Equal(t, 0, *resp.Quantity)

		addStock(t, nc, uniqueSku, 133)

		resp = getStock(t, nc, uniqueSku)
		require.True(t, resp.OK)
		assert.Equal(t, uniqueSku, *resp.ProductSKU)
		assert.Equal(t, 133, *resp.Quantity)

		removeStock(t, nc, uniqueSku, 131)

		resp = getStock(t, nc, uniqueSku)
		require.True(t, resp.OK)
		assert.Equal(t, uniqueSku, *resp.ProductSKU)
		assert.Equal(t, 2, *resp.Quantity)
	})
}

func addStock(t *testing.T, nc *nats.Conn, uniqueSku string, quantity int) schemas.StockAddResponse {
	// Call the stockAddHandler with a valid request
	req := schemas.StockAddRequest{
		ProductSKU: uniqueSku,
		Quantity:   quantity,
	}
	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)

	msg, err := nc.RequestWithContext(t.Context(), "stock.add", reqBytes)
	require.NoError(t, err)

	// Parse the response & check values
	resp := schemas.StockAddResponse{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	return resp
}

func removeStock(t *testing.T, nc *nats.Conn, uniqueSku string, quantity int) schemas.StockRemoveResponse {
	// Call the stockAddHandler with a valid request
	req := schemas.StockRemoveRequest{
		ProductSKU: uniqueSku,
		Quantity:   quantity,
	}
	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)

	msg, err := nc.RequestWithContext(t.Context(), "stock.remove", reqBytes)
	require.NoError(t, err)

	// Parse the response & check values
	resp := schemas.StockRemoveResponse{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	return resp
}

func getStock(t *testing.T, nc *nats.Conn, uniqueSku string) schemas.StockGetResponse {
	// Call the stockGetHandler with a valid request
	req := schemas.StockGetRequest{
		ProductSKU: uniqueSku,
	}
	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)

	msg, err := nc.RequestWithContext(t.Context(), "stock.get", reqBytes)
	require.NoError(t, err)

	// Parse the response & check values
	resp := schemas.StockGetResponse{}
	err = json.Unmarshal(msg.Data, &resp)
	require.NoError(t, err)

	return resp
}

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}
