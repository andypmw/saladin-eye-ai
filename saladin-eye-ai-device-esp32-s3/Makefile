# Define the proto files will be used to generate C/C++
PROTO_FILES = \
    media_service__get_photo_upload_url_request.proto \
    media_service__get_photo_upload_url_response.proto

# Define the source directory containing the proto files
PROTO_SRC_DIR := ../saladin-eye-ai-protos

# Define the output directory for the compiled files
PROTO_OUT_DIR := src/genproto

# Define the Python virtual environment directory
VENV_DIR = .venv

# Define the nanopb directory
NANOPB_DIR = .nanopb

# Define the nanopb source
NANOPB_DOWNLOAD_URL = https://jpa.kapsi.fi/nanopb/download/nanopb-0.4.8.tar.gz

nanopb-clean:
	rm -rf $(NANOPB_DIR)
	rm -rf $(VENV_DIR)

nanopb-setup:
	mkdir -p $(NANOPB_DIR)
	@curl -L $(NANOPB_DOWNLOAD_URL) -o $(NANOPB_DIR)/nanopb.tar.gz
	@tar -xzf $(NANOPB_DIR)/nanopb.tar.gz -C $(NANOPB_DIR)
	rm $(NANOPB_DIR)/nanopb.tar.gz
	@python3 -m venv $(VENV_DIR)
	@$(VENV_DIR)/bin/pip3 install --upgrade pip
	@$(VENV_DIR)/bin/pip3 install setuptools protobuf
	chmod +x $(VENV_DIR)/bin/activate

genproto: $(VENV_DIR)/bin/activate
	@mkdir -p $(PROTO_OUT_DIR)
	. $(VENV_DIR)/bin/activate && \
	python3 .nanopb/nanopb/generator/nanopb_generator.py \
		--proto-path=$(PROTO_SRC_DIR) \
		--output-dir=$(PROTO_OUT_DIR) \
		$(addprefix $(PROTO_SRC_DIR)/,$(PROTO_FILES))

cleanproto:
	rm -rf src/genproto
