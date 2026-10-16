.PHONY: version
version:
	go version
	protoc --version

.PHONY: generate
generate:
	protoc -I ./ -I ./api/third_party \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./pkg \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	api/product_service.proto

.PHONY: run-all
run-all:
	go run ./cmd/homework_2