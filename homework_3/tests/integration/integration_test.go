//go:build api

package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	pb "github.com/Snowdy54/go_homework/tree/main/homework_2/pkg/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestEnvironment struct {
	HTTPClient *http.Client
	RestUrl    string
	gRPCClient pb.ProductClient
	conn       *grpc.ClientConn
}

func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	grpcClient := pb.NewProductClient(conn)

	httpClient := &http.Client{Timeout: 5 * time.Second}

	return &TestEnvironment{
		HTTPClient: httpClient,
		RestUrl:    "http://127.0.0.1:8080",
		gRPCClient: grpcClient,
		conn:       conn,
	}
}

func (t *TestEnvironment) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}

func (t *TestEnvironment) GenerateRandomUserId() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(r.Intn(1000000) + 2000000)
}

func (e *TestEnvironment) GetCart(t *testing.T, userId int64) *http.Response {
	t.Helper()
	url := fmt.Sprintf("%s/user/%d/cart", e.RestUrl, userId)
	resp, err := e.HTTPClient.Get(url)
	require.NoError(t, err)

	return resp
}

func (e *TestEnvironment) PostAddProduct(t *testing.T, userId int64, sku int64, count int) *http.Response {
	t.Helper()

	reqBody, err := json.Marshal(map[string]interface{}{"count": count})
	require.NoError(t, err)

	url := fmt.Sprintf("%s/user/%d/cart/%d", e.RestUrl, userId, sku)
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	require.NoError(t, err)

	resp, err := e.HTTPClient.Do(req)
	require.NoError(t, err)

	return resp
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	env := SetupTestEnvironment(t)
	defer env.Close()

	ctx := context.Background()

	t.Run("Happy path - успешная работа сервиса", func(t *testing.T) {
		userId := env.GenerateRandomUserId()

		existingSku := 1076963

		grpcResp, err := env.gRPCClient.GetProduct(ctx, &pb.GetProductRequest{Sku: int64(existingSku)})
		require.NoError(t, err)
		require.NotNil(t, grpcResp)

		postResp := env.PostAddProduct(t, userId, int64(existingSku), 2)
		defer postResp.Body.Close()

		require.Equal(t, http.StatusOK, postResp.StatusCode)

		getResp := env.GetCart(t, userId)
		defer getResp.Body.Close()
		require.Equal(t, http.StatusOK, getResp.StatusCode)

		var cartResponseBody map[string]interface{}
		err2 := json.NewDecoder(getResp.Body).Decode(&cartResponseBody)

		require.NoError(t, err2)

		expectedTotalPrice := float64(grpcResp.Price * 2)
		require.Equal(t, expectedTotalPrice, cartResponseBody["totalPrice"])
	})

	t.Run("Error test (400) - некоректный аргумент sku", func(t *testing.T) {
		userId := env.GenerateRandomUserId()
		invalidSku := -5

		resp := env.PostAddProduct(t, userId, int64(invalidSku), 1)
		defer resp.Body.Close()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Error test (400) - некорретный аргумент userId", func(t *testing.T) {
		userId := 99999

		resp := env.GetCart(t, int64(userId))
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Error test (404) - несуществующий товар", func(t *testing.T) {
		userId := env.GenerateRandomUserId()
		sku := 1

		postResp := env.PostAddProduct(t, userId, int64(sku), 1)
		defer postResp.Body.Close()

		require.Equal(t, http.StatusPreconditionFailed, postResp.StatusCode)
	})
}
