PROTOC_GENERATOR=gogoslick

proto_type_overrides=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types

PROTO_DIR ?= proto
PROTOS_DEST_DIR=protos

.phony: proto proto-build proto-clean

proto: proto-clean proto-build

proto-clean:
	@echo "[protoc-clean] removing $(PROTO_DIR)/$(PROTOS_DEST_DIR)"
	@rm -rf $(PROTO_DIR)/$(PROTOS_DEST_DIR)

proto-build:
	@echo "[protoc-build] generating protos in $(PROTO_DIR)/$(PROTOS_DEST_DIR)"
	@mkdir -p $(PROTO_DIR)/$(PROTOS_DEST_DIR)
	@protoc --$(PROTOC_GENERATOR)_out=plugins=grpc,$(proto_type_overrides),import_path=$(PROTOS_DEST_DIR):$(PROTO_DIR)/$(PROTOS_DEST_DIR) -Iproto proto/envelope.proto proto/events.proto
