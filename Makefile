# List of implements server
SERVERS = health/v1/health.proto
# Enable grpc gateway
GATEWAY = true

# List of gateway option, listing here:
# 	https://git.teko.vn/shared/kitchen/-/blob/master/grpc/gatewayopt/gatewayopt.go
# Example:
# 	GATEWAY_OPTIONS = DefaultMarshaler Redirect
GATEWAY_OPTIONS = ProtoJSONMarshaler

TARGET = bin
TARGET_BIN = rpc-runtime
GO_CMD_MAIN = cmd/main.go

RPC_FOLDER = shared/rpc

SERVER_PACKAGE_NAME = server
SERVER_OUT_FOLDER = rpcimpl

# To generate database model, accept both file_name.dbml & dbdiagram.io link
DBML = https://dbdiagram.io/d/5ec236e639d18f5553ff5aee
MODEL_OUT = model
MODEL_PACKAGE = model

####################  DOES NOT EDIT BELLOW  ############################
.PHONY = build generate all clean

GO_TOOLS = git.teko.vn/shared/rpc-framework/cmd/protoc-gen-rpc-server git.teko.vn/shared/rpc-framework/cmd/dbml-gen-go-model

$(GO_TOOLS):
	GOSUMDB=off go get -u $@

# support fresh install on osx, not sure if it can't run properly
install-osx: $(GO_TOOLS)
	brew install protobuf

# support fresh install on linux, not sure if it can't run properly
PROTOC_LINUX_VERSION = 3.11.4
PROTOC_LINUX_ZIP = protoc-$(PROTOC_LINUX_VERSION)-linux-x86_64.zip

install-linux: $(GO_TOOLS)
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_LINUX_VERSION)/$(PROTOC_LINUX_ZIP)
	sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local bin/protoc
	sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local 'include/*'
	rm -f $(PROTOC_LINUX_ZIP)

update-rpc:
	@echo '# update $(RPC_FOLDER)'
	@[ -d "$(RPC_FOLDER)" ] || git clone git@git.teko.vn:$(RPC_FOLDER).git $(RPC_FOLDER)
	@cd $(RPC_FOLDER) && git checkout master && git pull origin master

.ONESHELL:
common-update: update-rpc
	go mod edit -droprequire=rpc.tekoapis.com
	GOSUMDB=off go get -u go.tekoapis.com/kitchen/...
	GOSUMDB=off go get -u rpc.tekoapis.com/...

prepare:
	mkdir -p $(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME)

photon-server:
	@echo \# generating photon-server....
	protoc -I $(RPC_FOLDER)/proto \
		-I $(RPC_FOLDER)/.third_party/googleapis \
		-I $(RPC_FOLDER)/.third_party/envoyproxy \
		-I $(RPC_FOLDER)/.third_party/gogoprotobuf \
		-I $(RPC_FOLDER)/.third_party/gogoprotobuf \
		--rpc-server_out=gateway=$(GATEWAY),gateway_options="$(GATEWAY_OPTIONS)":$(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME) \
		$(SERVERS)

generate: prepare photon-server
	@echo \# source code is generated

build: generate
	go build -o $(TARGET)/$(TARGET_BIN) $(GO_CMD_MAIN)

run: generate
	go run $(GO_CMD_MAIN) server

migrate:
	echo \# make migrate name="$(name)"
	go run $(GO_CMD_MAIN) migrate create $(name)

migrate-up:
	go run $(GO_CMD_MAIN) migrate up

migrate-down-1:
	go run $(GO_CMD_MAIN) migrate down 1

model:
	rm -rf $(MODEL_OUT)/*.enum.go
	rm -rf $(MODEL_OUT)/*.table.go
	dbml-gen-go-model -f $(DBML) -o $(MODEL_OUT) -p $(MODEL_PACKAGE) 

all: common-update build
	@echo what is done is done!

clean:
	rm -rf $(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME)
	rm -rf $(TARGET)

test:
	go test ./...  -count=1 -v -cover -race

test-all-coverage:
	go test ./... -count=1 -race -coverprofile cover.out
	go tool cover -func cover.out
