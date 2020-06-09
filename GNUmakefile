TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=unleash

default: build

build: fmtcheck
	go install

install: build
	install -d -m 755 ~/.terraform.d/plugins
	install $(GOPATH)/bin/terraform-provider-unleash ~/.terraform.d/plugins

apply: install
	terraform init
	terraform $@

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint: tools.golangci-lint
	@echo "==> Checking source code against golangci-lint"
	@golangci-lint run ./$(PKG_NAME)

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build install apply test testacc vet fmt fmtcheck lint test-compile

#---------------
#-- tools
#---------------
.PHONY: tools.golangci-lint

tools.golangci-lint:
	GO111MODULE=off go install github.com/golangci/golangci-lint/cmd/golangci-lint
