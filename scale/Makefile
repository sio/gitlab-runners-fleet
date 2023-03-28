GO?=go
APP=bin/scale$(shell $(GO) env GOEXE)
EVERYTHING=$(EVERYTHING)

.PHONY: run
run: $(APP)  ## run scaler app
	$<

.PHONY: build
build: $(APP)  ## build scaler app
$(APP): $(shell find . -type f -iname '*.go')
	$(GO) build -o $@

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
lint:  ## run linter for current day's solution
	$(GO) fmt $(EVERYTHING)
	$(GO) vet $(EVERYTHING)
	$(STATICCHECK) $(EVERYTHING)
	$(ERRCHECK) $(EVERYTHING)

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
