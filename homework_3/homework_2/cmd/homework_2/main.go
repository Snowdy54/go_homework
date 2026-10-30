package main

import (
	"log"
	"net"
	"net/http"

	"github.com/Snowdy54/go_homework/tree/main/homework_2/internal/app"
	pb "github.com/Snowdy54/go_homework/tree/main/homework_2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	productService := app.NewProductService()
	pb.RegisterProductServer(server, productService)

	go func() {
		log.Println("listening on port :50051")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed connect to product-service")
	}
	defer conn.Close()

	productClient := pb.NewProductClient(conn)
	cartService := app.NewCartService(productClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", cartService.AddProduct)
	mux.HandleFunc("GET /user/{user_id}/cart", cartService.GetCart)

	log.Println("listening on port :8080")
	err2 := http.ListenAndServe(":8080", mux)
	if err2 != nil {
		log.Fatalf("failed to serve :8080 : %v", err2)
	}

}
