GO?=go
APP=bin/scale-$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)$(shell $(GO) env GOEXE)
EVERYTHING=./...

.PHONY: run
run: $(APP)  ## run scale app
	cat sample_config.json | $<

.PHONY: build
build: $(APP)  ## build scale app
$(APP): $(shell find . -type f -iname '*.go')
	$(GO) build -o $@

.PHONY: build-all
build-all:  ## build scale app for all OS/architectures
	GOOS=linux   GOARCH=amd64 $(MAKE) build
	GOOS=linux   GOARCH=arm64 $(MAKE) build
	GOOS=windows GOARCH=amd64 $(MAKE) build
	GOOS=windows GOARCH=arm64 $(MAKE) build
	GOOS=darwin  GOARCH=amd64 $(MAKE) build
	GOOS=darwin  GOARCH=arm64 $(MAKE) build

.PHONY: test
test:  ## run automated tests
	$(GO) test $(EVERYTHING) -timeout 15s $(GOTEST_ARGS)

.PHONY: test-verbose
export GOTEST_ARGS
test-verbose: GOTEST_ARGS+=-v --count=1
test-verbose:  ## run tests with increased verbosity
	$(MAKE) test

.PHONY: bench
bench:  ## run benchmarks
	cd $(DIRECTORY) && $(GO) test -bench=. -count=3 -benchmem -benchtime=2s -run='^#'

.PHONY: fmt
fmt:  ## format Go code
	$(GO) fmt $(EVERYTHING)

STATICCHECK=$(shell $(GO) env GOPATH)/bin/staticcheck$(shell $(GO) env GOEXE)
$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@latest

ERRCHECK=$(shell $(GO) env GOPATH)/bin/errcheck$(shell $(GO) env GOEXE)
$(ERRCHECK):
	$(GO) install github.com/kisielk/errcheck@latest

.PHONY: lint
lint: $(STATICCHECK) $(ERRCHECK)
lint:  ## run linters
	$(GO) fmt $(EVERYTHING)
	$(GO) vet $(EVERYTHING)
	$(STATICCHECK) $(EVERYTHING)
	$(ERRCHECK) $(EVERYTHING)

.PHONY: ci
ci: lint
	git diff --exit-code --name-only

CLOC=$(shell $(GO) env GOPATH)/bin/gocloc$(shell $(GO) env GOEXE)
$(CLOC):
	$(GO) install github.com/hhatto/gocloc/cmd/gocloc@latest

.PHONY: loc
loc: $(CLOC)  ## count lines of code
	@echo "This directory: ./scale/"
	@$(CLOC) .
	@echo
	@echo
	@echo "Whole git repo: ../"
	@$(CLOC) --not-match-d='\.terraform' ..

.PHONY: help
help:  ## show help message
	@awk "$$HELP_AWK" $(lastword $(MAKEFILE_LIST)) | column -t -s '	'

define HELP_AWK
/  ##/ {
	if (stop) exit;

	match($$0, /^[^:]+:/);
	target = substr($$0, RSTART, RLENGTH-1);
	gsub(" ", ",", target);

	match($$0, /  ## /);
	help = substr($$0, RSTART + RLENGTH);

	printf("%s\t%s\n", target, help);
}
/HELP_AWK/ {
	stop = 1;
}
endef
export HELP_AWK
