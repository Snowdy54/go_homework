package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "github.com/Snowdy54/go_homework/tree/main/homework_2/pkg/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeProductClient struct {
	GetProducFunc func(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error)
}

func (f *fakeProductClient) GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error) {
	if f.GetProducFunc != nil {
		return f.GetProducFunc(ctx, in)
	}

	return &pb.GetProductResponse{
		Sku:   in.Sku,
		Name:  "Кроссовки",
		Price: 2400,
	}, nil
}

func TestHandler_CartServiceGetSuccess(t *testing.T) {
	t.Parallel()

	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
	}

	req := httptest.NewRequest("GET", "/user/1/cart", nil)
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{user_id}/cart", svc.GetCart)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, "{\"products\":[{\"sku\":1,\"count\":1,\"Product\":{\"name\":\"Кроссовки\",\"price\":2400}}],\"totalPrice\":2400}\n",
		rec.Body.String())
}

func TestHandler_CartServiceGetNotFound(t *testing.T) {
	t.Parallel()
	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
	}

	req := httptest.NewRequest("GET", "/user/10/cart", nil)
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{user_id}/cart", svc.GetCart)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_CartServiceGetBadRequest(t *testing.T) {
	t.Parallel()
	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
	}

	req := httptest.NewRequest("GET", "/user/-5/cart", nil)
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{user_id}/cart", svc.GetCart)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CartServicePostSuccess(t *testing.T) {
	t.Parallel()
	fakeClient := &fakeProductClient{}
	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
		ProductClient: fakeClient,
	}

	jsonBody := `{"count": 1}`

	req := httptest.NewRequest("POST", "/user/1/cart/1", strings.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", svc.AddProduct)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestHandler_CartServicePostInvaildUserId(t *testing.T) {
	t.Parallel()
	fakeClient := &fakeProductClient{}
	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
		ProductClient: fakeClient,
	}

	jsonBody := `{"count": 1}`

	req := httptest.NewRequest("POST", "/user/-5/cart/1", strings.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", svc.AddProduct)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CartServicePostInvaildSkuId(t *testing.T) {
	t.Parallel()
	fakeClient := &fakeProductClient{}
	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
		ProductClient: fakeClient,
	}

	jsonBody := `{"count": 1}`

	req := httptest.NewRequest("POST", "/user/1/cart/-5", strings.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", svc.AddProduct)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CartServicePostInvalidCount(t *testing.T) {
	t.Parallel()
	fakeClient := &fakeProductClient{}
	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Кроссовки", Price: 2400}}},
		},
		ProductClient: fakeClient,
	}

	jsonBody := `{"count": -5}`

	req := httptest.NewRequest("POST", "/user/1/cart/-5", strings.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", svc.AddProduct)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CartServicePostNotFound(t *testing.T) {
	t.Parallel()
	fakeClient := &fakeProductClient{GetProducFunc: func(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error) {
		return nil, status.Error(codes.NotFound, "product not found")
	}}

	svc := &CartService{
		Cart: map[int64]map[int64]*Cart{
			1: {1: &Cart{Sku: 1, Count: 1, Product: Product{Name: "Ботинки", Price: 2400}}},
		},
		ProductClient: fakeClient,
	}

	jsonBody := `{"count": 1}`

	req := httptest.NewRequest("POST", "/user/1/cart/1", strings.NewReader(jsonBody))
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", svc.AddProduct)
	mux.ServeHTTP(rec, req)

	require.Equal(t, http.StatusPreconditionFailed, rec.Code)
}
