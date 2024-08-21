# Define the proto source directory and output directory
PROTO_SRC_DIR := ../../saladin-eye-ai-protos
PROTO_OUT_DIR := .
PROTO_FILES = \
    media_service__file_info.proto \
    media_service__get_photo_upload_url_request.proto \
    media_service__get_photo_upload_url_response.proto \
    media_service__list_files_by_date_hour_request.proto \
    media_service__list_files_by_date_hour_response.proto \
    media_service.proto

# To generate Go and gRPC code from proto files
genproto:
	protoc --proto_path=$(PROTO_SRC_DIR) --go_out=$(PROTO_OUT_DIR) --go-grpc_out=$(PROTO_OUT_DIR) $(addprefix $(PROTO_SRC_DIR)/,$(PROTO_FILES))

cleanproto:
	rm -rf common/genproto
