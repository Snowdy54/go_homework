package app

import (
	"context"
	"testing"

	pb "github.com/Snowdy54/go_homework/tree/main/homework_2/pkg/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandler_GetProductSuccess(t *testing.T) {
	t.Parallel()
	svc := &ProductService{
		Products: map[int64]*Product{
			1: {Name: "Кружка", Price: 750},
		},
		Sku: 1,
	}

	ctx := context.Background()
	req := &pb.GetProductRequest{
		Sku: 1,
	}

	resp, err := svc.GetProduct(ctx, req)

	require.NoError(t, err)

	require.Equal(t, svc.Sku, resp.Sku)

	require.Equal(t, svc.Products[1].Name, resp.Name)

	require.Equal(t, svc.Products[1].Price, resp.Price)
}

func TestHandler_GetProductNotFound(t *testing.T) {
	t.Parallel()
	svc := &ProductService{
		Products: map[int64]*Product{
			1: {Name: "Котелок", Price: 3500},
		},
		Sku: 1,
	}

	ctx := context.Background()
	req := &pb.GetProductRequest{
		Sku: 3,
	}

	resp, err := svc.GetProduct(ctx, req)

	require.Error(t, err)

	require.Nil(t, resp)

	status, ok := status.FromError(err)

	require.True(t, ok)

	require.Equal(t, codes.NotFound, status.Code())
}

func TestHandler_GetProductInvalidSku(t *testing.T) {
	t.Parallel()
	svc := &ProductService{
		Products: map[int64]*Product{
			1: {Name: "Кроссовки", Price: 3500},
		},
		Sku: 1,
	}

	ctx := context.Background()
	req := &pb.GetProductRequest{
		Sku: -5,
	}

	resp, err := svc.GetProduct(ctx, req)

	require.Nil(t, resp)

	status, ok := status.FromError(err)

	require.True(t, ok)

	require.Equal(t, codes.InvalidArgument, status.Code())
}
