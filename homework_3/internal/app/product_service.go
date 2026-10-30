package app

import (
	"context"
	"encoding/json"
	"log"
	"os"

	pb "github.com/Snowdy54/go_homework/tree/main/homework_2/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb.ProductServer

	Products map[int64]*Product
	Sku      int64
}

type Product struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

func NewProductService() *ProductService {
	data, err := os.Open("./data.json")
	if err != nil {
		log.Fatalf("file error: %v", err)
	}
	defer data.Close()

	var products map[int64]*Product

	decoder := json.NewDecoder(data)
	err2 := decoder.Decode(&products)
	if err2 != nil {
		log.Fatalf("file error: %v", err)
	}

	return &ProductService{
		Products: products,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if in.Sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid sku")
	}

	product, ok := s.Products[in.Sku]
	if !ok {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &pb.GetProductResponse{Sku: in.Sku, Name: product.Name, Price: product.Price}, nil
}
