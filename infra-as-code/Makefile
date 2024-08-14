all := services/pub services/sub

GO := go
GO_BUILD := $(GO) build
GO_RUN := $(GO) run
GO_TEST := $(GO) test
GO_FMT := $(GO) fmt
GO_LINT := golangci-lint run
GO_CLEAN := $(GO) clean
GO_MOD_TIDY := $(GO) mod tidy

# Default target
all: build

# Build the project
build:
	@echo "Building..."
	@for dir in $(all); do \
		$(MAKE) build -C $$dir || exit 1; \
	done

# Build the project containers
containers:
	@echo "Building docker containers..."
	@for dir in $(all); do \
		$(MAKE) container -C $$dir || exit 1; \
	done

# Test the project
test:
	@echo "Running tests..."
	@for dir in $(all); do \
		$(MAKE) test -C $$dir || exit 1; \
	done

# Format the code
fmt:
	@echo "Formatting code..."
	@for dir in $(all); do \
		$(MAKE) fmt -C $$dir || exit 1; \
	done

# Lint the code (requires golangci-lint to be installed)
lint:
	@echo "Linting code..."
	@for dir in $(all); do \
		$(MAKE) lint -C $$dir || exit 1; \
	done

# Clean the build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@for dir in $(all); do \
		$(MAKE) clean -C $$dir || exit 1; \
	done

# Tidy up the module
tidy:
	@echo "Tidying up go.mod files..."
	@for dir in $(all); do \
		$(MAKE) tidy -C $$dir || exit 1; \
	done

# Terraform commands
terraform-init:
	@cd terraform && terraform init

terraform-plan: containers
	@cd terraform && terraform plan

terraform-apply: containers
	@cd terraform && terraform apply

terraform-destroy:
	@cd terraform && terraform destroy


# Phony targets
.PHONY: all build containers run test fmt lint clean tidy terraform-init terraform-plan terraform-apply terraform-destroy
