UBINDIR ?= /usr/bin
DESTDIR ?=
EXTNAME := $(shell basename $(shell pwd))

# go tool nm ./luet | grep Commit
override LDFLAGS += -X "github.com/geaaru/rest-guard/pkg/specs.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
override LDFLAGS += -X "github.com/geaaru/rest-guard/pkg/specs.BuildCommit=$(shell git rev-parse HEAD)"

all: build install

build:
	CGO_ENABLED=0 go build -o luet-portage-converter -ldflags '$(LDFLAGS)'

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 $(EXTNAME) $(DESTDIR)/$(UBINDIR)/

.PHONY: deps
deps:
	go env
	# Installing dependencies...
	GO111MODULE=off go get golang.org/x/lint/golint
	GO111MODULE=off go get golang.org/x/tools/cmd/cover
	GO111MODULE=on go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
	GO111MODULE=off go get github.com/onsi/gomega/...
	ginkgo version

.PHONY: test
test:
	GO111MODULE=on go get github.com/onsi/ginkgo/v2/ginkgo
	GO111MODULE=off go get github.com/onsi/gomega/...
	ginkgo -r -race -flake-attempts 3 ./...

.PHONY: coverage
coverage:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: test-coverage
test-coverage:
	scripts/ginkgo.coverage.sh --codecov

.PHONY: clean
clean:
	-rm luet-portage-converter
	-rm -rf release/ dist/

.PHONY: goreleaser-snapshot
goreleaser-snapshot:
	rm -rf dist/ || true
	goreleaser release --debug --skip-publish  --skip-validate --snapshot

