IMAGE_NAME ?= astronomerio/cs-event-api

GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(shell git describe --tags --always)
GIT_IMPORT=github.com/astronomerio/cs-event-api/pkg/version
GOLDFLAGS=-X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)

VERSION ?= SNAPSHOT-$(GIT_COMMIT)

build:
	go build -ldflags '$(GOLDFLAGS)' -tags static -o event-api main.go

build-image:
	docker build -t $(IMAGE_NAME):$(VERSION) .

tag-latest:
	docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest

push-image:
	docker push $(IMAGE_NAME):$(VERSION)

install: build
	mkdir -p $(DESTDIR)
	cp event-api $(DESTDIR)

uninstall:
	rm -rf $(DESTDIR)
