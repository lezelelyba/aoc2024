MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
CLOUDPROVIDER ?= aws
ENVIRONMENT ?= prod
TF_PARAMS ?= 

.PHONY: all localci localcd localdestroy localclean init apply destroy output test 

all: init apply

# Run the CI pipeline (e.g. build docker images, run tests)
# TODO: testing is missing, requirements are missing
local: localci localcd
	@echo "All tasks completed"

localci:
	@echo "Building Docker images"
	./environments/local/build_docker_images.sh
	@echo "Done"

localrun: localci
	@echo "Running local Docker Image"
	docker run -p 8080:8080 advent2024.web &
	@echo "Done"

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

ifeq ($(CLOUDPROVIDER), aws)
TF_BOOTSTRAPDIR ?= environments/bootstrap/aws/terraform
TF_BACKEND_CONFIG ?= $(MAKEFILE_DIR)environments/aws/backend.json

ifeq ($(ENVIRONMENT), prod)
TF_DIR ?= $(MAKEFILE_DIR)environments/aws/prod/terraform
else ifeq ($(ENVIRONMENT), stage)
TF_DIR ?= $(MAKEFILE_DIR)environments/aws/stage/terraform
else ifeq ($(ENVIRONMENT), dev)
TF_DIR ?= $(MAKEFILE_DIR)environments/aws/dev/terraform
else ifeq ($(ENVIRONMENT), bootstrap)
TF_DIR ?= $(TF_BOOTSTRAPDIR)
else
$(error Unknown ENVIRONMENT: $(ENVIRONMENT))
endif

else ifeq ($(CLOUDPROVIDER), azure)
TF_BOOTSTRAPDIR ?= environments/bootstrap/azure/terraform
TF_BACKEND_CONFIG ?= $(MAKEFILE_DIR)environments/azure/backend.json

ifeq ($(ENVIRONMENT), prod)
TF_DIR ?= $(MAKEFILE_DIR)environments/azure/prod/terraform
else ifeq ($(ENVIRONMENT), stage)
TF_DIR ?= $(MAKEFILE_DIR)environments/azure/stage/terraform
else ifeq ($(ENVIRONMENT), dev)
TF_DIR ?= $(MAKEFILE_DIR)environments/azure/dev/terraform
else ifeq ($(ENVIRONMENT), bootstrap)
TF_DIR ?= $(TF_BOOTSTRAPDIR)
else
$(error Unknown ENVIRONMENT: $(ENVIRONMENT))
endif

else
$(error Unknown CLOUDPROVIDER: $(CLOUDPROVIDER))
endif
	
init:
ifeq ($(ENVIRONMENT), bootstrap)
	(cd ${TF_DIR}; terraform init ${TF_PARAMS})
else
	(cd ${TF_DIR}; terraform init -backend-config=${TF_BACKEND_CONFIG} ${TF_PARAMS})
endif
apply:
ifeq ($(ENVIRONMENT), bootstrap)
	(cd $(TF_DIR); terraform apply ${TF_PARAMS})
else
	(cd $(TF_DIR); terraform apply ${TF_PARAMS})
endif
destroy:
ifeq ($(ENVIRONMENT), bootstrap)
	(cd $(TF_DIR); terraform destroy ${TF_PARAMS})
else
	(cd $(TF_DIR); terraform destroy ${TF_PARAMS})
endif
output:
	(cd $(TF_DIR); terraform output)

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
	
	@echo "Running cli tests..."
	go test ./cmd/cli/...

	@echo "Running web tests..."
	go test ./cmd/web/...
