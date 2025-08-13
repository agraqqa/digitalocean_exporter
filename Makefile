EXECUTABLE ?= digitalocean_exporter
IMAGE ?= metalmatze/$(EXECUTABLE)
GO := CGO_ENABLED=0 go
DATE := $(shell date -u '+%FT%T%z')

LDFLAGS += -X main.Version=$(DRONE_TAG)
LDFLAGS += -X main.Revision=$(DRONE_COMMIT)
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./...)

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf dist/

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) install golang.org/x/lint/golint@latest; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	@for PKG in $(PACKAGES); do $(GO) test -cover $$PKG || exit 1; done;

.PHONY: test-short
test-short:
	$(GO) test -v -short ./...

.PHONY: test-coverage
test-coverage:
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration:
	$(GO) test -v -run TestIntegration ./...

.PHONY: docker-build
docker-build: build
	docker build -t $(IMAGE):latest .

.PHONY: podman-build
podman-build: build
	podman build -t $(IMAGE):latest .

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run all tests"
	@echo "  test-short     - Run tests excluding integration tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-integration - Run only integration tests"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  podman-build   - Build Podman image"

$(EXECUTABLE): $(wildcard *.go)
	$(GO) build -v -ldflags '-w $(LDFLAGS)'

.PHONY: build
build: $(EXECUTABLE)

.PHONY: install
install:
	$(GO) install -v -ldflags '-w $(LDFLAGS)'

.PHONY: release
release:
	@which gox > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	CGO_ENABLED=0 gox -verbose -ldflags '-w $(LDFLAGS)' -osarch '!darwin/386' -output="dist/$(EXECUTABLE)-${DRONE_TAG}-{{.OS}}-{{.Arch}}"
