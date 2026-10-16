package app

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	pb "github.com/Snowdy54/go_homework/tree/main/homework_2/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartService struct {
	ProductClient pb.ProductClient
	Cart          map[int64]map[int64]*Cart `json:"cart"`
}

type Cart struct {
	Sku     int64 `json:"sku"`
	Count   int32 `json:"count"`
	Product Product
}

func NewCartService(client pb.ProductClient) *CartService {
	return &CartService{
		ProductClient: client,
		Cart:          make(map[int64]map[int64]*Cart),
	}
}

type AddProductRequest struct {
	Count int32
}

func (s *CartService) AddProduct(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userId <= 0 {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	skuId, err2 := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err2 != nil || skuId <= 0 {
		http.Error(w, "invalid sku_id", http.StatusBadRequest)
		return
	}

	var req AddProductRequest
	decoder := json.NewDecoder(r.Body)
	err3 := decoder.Decode(&req)
	if err3 != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.Count <= 0 {
		http.Error(w, "quantity must be grater than zero", http.StatusBadRequest)
		return
	}

	productResp, err := s.ProductClient.GetProduct(r.Context(), &pb.GetProductRequest{
		Sku: skuId,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "product for this sku not found", http.StatusPreconditionFailed)
			return
		}
		http.Error(w, "system error", http.StatusInternalServerError)
		return
	}

	_, ok := s.Cart[userId]
	if !ok {
		s.Cart[userId] = make(map[int64]*Cart)
	}

	product, ok := s.Cart[userId][skuId]
	if ok {
		product.Count += int32(req.Count)
	} else {
		s.Cart[userId][skuId] = &Cart{
			Sku:   skuId,
			Count: int32(req.Count),
			Product: Product{
				Name:  productResp.Name,
				Price: productResp.Price,
			},
		}
	}

	w.WriteHeader(http.StatusOK)
}

type GetCartResponse struct {
	Products   []Cart `json:"products"`
	TotalPrice int32  `json:"totalPrice"`
}

func (s *CartService) GetCart(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userId <= 0 {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	userCart, ok := s.Cart[userId]
	if !ok || len(userCart) == 0 {
		http.Error(w, "cart not found or empty", http.StatusNotFound)
		return
	}

	items := make([]Cart, 0, len(userCart))
	var totalPrice int32

	for _, product := range userCart {
		quantity := int32(product.Count)
		price := int32(product.Product.Price)

		items = append(items, Cart{
			Sku:   product.Sku,
			Count: quantity,
			Product: Product{
				Name:  product.Product.Name,
				Price: price,
			},
		})
		totalPrice += quantity * price
	}

	sort.Slice(items, func(i int, j int) bool {
		return items[i].Sku < items[j].Sku
	})

	resp := GetCartResponse{
		Products:   items,
		TotalPrice: totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
