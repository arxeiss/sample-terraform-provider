.PHONY: lint fix_lint install_tools terraform_build start_server

lint:
	golangci-lint run ./...

fix_lint:
	golangci-lint run --fix

install_tools:
	@echo Installing tools from tools.go
	go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install

terraform_build:
	@echo Build and replace terraform-provider-sdk
	@cd provider && go build -o terraform-provider-sdk
	@mkdir -p ./provider/config/terraform.d/plugins/terraform.kutac.cz/superdupercloud/sdk/0.1.0/$$(go env GOHOSTOS)_$$(go env GOHOSTARCH)/
	@mv ./provider/terraform-provider-sdk ./provider/config/terraform.d/plugins/terraform.kutac.cz/superdupercloud/sdk/0.1.0/$$(go env GOHOSTOS)_$$(go env GOHOSTARCH)/
	@rm -f ./provider/config/.terraform.lock.hcl
	cd ./provider/config/ && terraform init -backend=false

start_server:
	cd server && CGO_ENABLED=1 go run .
