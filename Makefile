GO ?= godep go
ifdef CIRCLE_ARTIFACTS
  COVERAGEDIR = $(CIRCLE_ARTIFACTS)
endif

all: clean build
godep:
	go get github.com/tools/godep
godep-save:
	godep save ./...
build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -o bin/canon-intervalometer
fmt:
	$(GO) fmt ./...
test:
	$(GO) test ./...
bench:
	$(GO) test -bench ./...
run:
	bin/canon-intervalometer
clean:
	$(GO) clean
	rm -f bin/canon-intervalometer