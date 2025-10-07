.PHONY: localci localcd localdestroy localclean devbootstrap devbootstrapdestroy test

# Run the CI pipeline (e.g. build docker images, run tests)
# TODO: testing is missing, requirements are missing
local: localci localcd
	@echo "All tasks completed"

localci:
	@echo "Building Docker images"
	./environments/local/build_docker_images.sh
	@echo "Completed"

# Run the CD pipeline (e.g. deploy with ansible, terraform, etc.)
localcd:
	@echo "Spining local deployment"
	./environments/local/build_local.sh
	@echo "Done"
	@echo "Configuring local deployment"
	./environments/local/configure_run_local.sh
	@echo "Done"

localdestroy:
	@echo "Destroying local environment"
	./environments/local/build_local.sh destroy
	@echo "Done"

# Clean build artifacts (if any)
localclean:
	go clean
	docker rmi advent2024.web || true
	docker rmi advent2024.cli || true

TF_INITVARS ?= ""
TF_BOOTSTRAPDIR ?= environments/aws_bootstrap/terraform

awsbootstrap:
	(cd $(TF_BOOTSTRAPDIR); terraform init ${TF_INITVARS}; terraform apply)

awsbootstrapdestroy:
	(cd $(TF_BOOTSTRAPDIR); terraform destroy)

TF_DEVDIR ?= environments/aws/dev/terraform

awsdev:
	(cd $(TF_DEVDIR); terraform init ${TF_INITVARS}; terraform apply)

awsdevdestroy:
	(cd $(TF_DEVDIR); terraform destroy)

# Run tests (unit, integration)
test:
	@echo "Running unit tests..."
	@go work sync
	@for mod in ./pkg/*/; do \
		if [ -f "$$mod/go.mod" ]; then \
			echo "Testing $$mod..."; \
			go test "$$mod/..."; \
		fi \
	done

