.PHONY: tfbuild

tfbuild:
	@echo Build and replace terraform-provider-sdk
	@cd provider && go build -o terraform-provider-sdk
	@mkdir -p ./provider/config/terraform.d/plugins/terraform.kutac.cz/superdupercloud/sdk/0.1.0/$$(go env GOHOSTOS)_$$(go env GOHOSTARCH)/
	@mv ./provider/terraform-provider-sdk ./provider/config/terraform.d/plugins/terraform.kutac.cz/superdupercloud/sdk/0.1.0/$$(go env GOHOSTOS)_$$(go env GOHOSTARCH)/
	@rm -f ./provider/config/.terraform.lock.hcl
	cd ./provider/config/ && terraform init -backend=false
