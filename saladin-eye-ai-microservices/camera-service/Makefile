# Define the proto source directory and output directory
PROTO_SRC_DIR := ../../saladin-eye-ai-protos
PROTO_OUT_DIR := .

# To generate Go and gRPC code from proto files
genproto:
	protoc --proto_path=$(PROTO_SRC_DIR) --go_out=$(PROTO_OUT_DIR) --go-grpc_out=$(PROTO_OUT_DIR) $(PROTO_SRC_DIR)/*.proto
