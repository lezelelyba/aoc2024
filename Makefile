CLOUDPROVIDER ?= "aws"

.PHONY: all localci localcd localdestroy localclean bootstrap dev test

all: bootstrap dev

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


TF_INITPARAMS ?= 
TF_APPLYPARAMS ?= 

ifeq ($(CLOUDPROVIDER), "aws")
TF_BOOTSTRAPDIR ?= environments/aws_bootstrap/terraform
endif

bootstrapinit:
	(cd $(TF_BOOTSTRAPDIR); terraform init ${TF_INITPARAMS})

bootstrapapply:
	(cd $(TF_BOOTSTRAPDIR); terraform apply ${TF_APPLYPARAMS})

bootstrap: bootstrapinit bootstrapapply

bootstrapdestroy:
	(cd $(TF_BOOTSTRAPDIR); terraform destroy)
	

ifeq ($(CLOUDPROVIDER), "aws")
TF_DEVDIR ?= environments/aws/dev/terraform
endif

TF_DEVSSHPUBKEYPATH ?= ~/.ssh/aws.pem
TF_DEVSSHPRIVKEYPATH ?= ~/.ssh/aws.priv.pem

devinit:
	(cd ${TF_DEVDIR}; terraform init ${TF_INITPARAMS})

devapply:
	(cd $(TF_DEVDIR); terraform apply -var="sshpubkeypath=${TF_DEVSSHPUBKEYPATH}" -var="sshprivkeypath=${TF_DEVSSHPRIVKEYPATH}" ${TF_APPLYPARAMS})

dev: devinit devapply

devdestroy:
	(cd $(TF_DEVDIR); terraform destroy)

devoutput:
	(cd $(TF_DEVDIR); terraform output)

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
