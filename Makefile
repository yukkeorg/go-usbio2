ifeq ($(OS), Windows_NT)
	VERSION := $(shell git describe --exact-match --tags 2>nil)
else
	VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
endif

COMMIT := $(shell git rev-parse --short HEAD)
LDFLAGS := $(LDFLAGS) -X main.commit=$(COMMIT)

ifdef VERSION
	LDFLAGS += -X main.version=$(VERSION)
endif


export GO_BUILD=env go build -ldflags "$(LDFLAGS)"

SOURCES := $(shell find . -name '+.go' -not -name '*_test.go') go.mod go.sum


bin/usbio2-config: $(SOURCES)
	$(GO_BUILD) -o $@ ./cmd/$(shell basename "$@")
