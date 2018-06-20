include Makefile.defs

all: build

main:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo $(GOBUILD) ./cmd/main.go

vendor/vendor.json:
	govendor init
	govendor add -v +outside

./openapi/swagger.json: ./openapi/documentation/$(PROJECT_NAME)*.yaml
	swagger flatten ./openapi/documentation/$(PROJECT_NAME)-openapi.yaml $^ -o $@

.PHONY: full_clean
full_clean: clean
	rm -Rf ./vendor

.PHONY: freeze
freeze: vendor/vendor.json

.PHONY: prepare
prepare: vendor/vendor.json
	govendor sync -v -insecure


.PHONY: clean
clean:
	rm ./main
	rm ./openapi/swagger.json
	rm -Rf ./openapi/autogen

.PHONY: gen-openapi
gen-openapi:
	@make gen-openapi-server

.PHONY: gen-openapi-server
gen-openapi-server:
	mkdir -p $(OPENAPI_GEN_DIR)
	swagger generate server -t $(OPENAPI_GEN_DIR) -f $(OPENAPI_SPEC_FILE) -s server --exclude-main -A $(PROJECT_CAPITAL_NAME)

.PHONY: swagger
swagger: gen-openapi

.PHONY: build
build: main

.PHONY: serve
serve:
	swagger serve ./openapi/documentation/$(PROJECT_NAME)-openapi.yaml

.PHONY: run
run:
	./main

.PHONY: tests
tests:
	go test -cover `go list ./internal/...`
