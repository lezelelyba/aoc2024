MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
CLOUDPROVIDER ?= aws
ENVIRONMENT ?= prod
TF_INITPARAMS ?= 
TF_APPLYPARAMS ?= 
TF_DESTROYPARAMS ?= 
TF_SSHPUBKEYPATH ?= ~/.ssh/aws.pem
TF_SSHPRIVKEYPATH ?= ~/.ssh/aws.priv.pem

.PHONY: all localci localcd localdestroy localclean bootstrap init apply destroy output test 

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
TF_BACKEND_CONFIG ?= $(MAKEFILE_DIR)environments/aws/backend.json
else ifeq ($(ENVIRONMENT), dev)
TF_DIR ?= $(MAKEFILE_DIR)environments/aws/dev/terraform
TF_BACKEND_CONFIG ?= $(MAKEFILE_DIR)environments/aws/backend.json
else
$(error Unknown ENVIRONMENT: $(ENVIRONMENT))
endif

else
$(error Unknown CLOUDPROVIDER: $(CLOUDPROVIDER))
endif

bootstrapinit:
	(cd $(TF_BOOTSTRAPDIR); terraform init ${TF_INITPARAMS})

bootstrapapply:
	(cd $(TF_BOOTSTRAPDIR); terraform apply ${TF_APPLYPARAMS})

bootstrap: bootstrapinit bootstrapapply

bootstrapdestroy:
	(cd $(TF_BOOTSTRAPDIR); terraform destroy ${TF_DESTROY_PARAMS})
	
init:
	(cd ${TF_DIR}; terraform init -backend-config=${TF_BACKEND_CONFIG} ${TF_INITPARAMS})
apply:
	(cd $(TF_DIR); terraform apply -var="sshpubkeypath=${TF_SSHPUBKEYPATH}" -var="sshprivkeypath=${TF_SSHPRIVKEYPATH}" ${TF_APPLYPARAMS})
destroy:
	(cd $(TF_DIR); terraform destroy -var="sshpubkeypath=${TF_SSHPUBKEYPATH}" -var="sshprivkeypath=${TF_SSHPRIVKEYPATH}" ${TF_DESTROYPARAMS})
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
