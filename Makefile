IMAGE_NAME ?= astronomerinc/ap-event-api

GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(shell git describe --tags --always)
GIT_IMPORT=github.com/astronomerio/event-api/pkg/version
GOLDFLAGS=-X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)

VERSION ?= SNAPSHOT-$(GIT_COMMIT)

# Set default for make.
.DEFAULT_GOAL := build-image

.PHONY: build
build:
	go build -ldflags '$(GOLDFLAGS)' -tags static -o event-api main.go

.PHONY: install
install: build
	mkdir -p $(DESTDIR)
	cp event-api $(DESTDIR)

.PHONY: uninstall
uninstall:
	rm -rf $(DESTDIR)

.PHONY: build-image
build-image:
	docker build -t $(IMAGE_NAME):$(VERSION) .

.PHONY: test-image
test-image: build-image
	docker run ${IMAGE_NAME}:${VERSION} go test ./...
