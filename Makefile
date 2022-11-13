.PHONE: all
all: build/gopkgdep demo

.PHONY: build
build: build/gopkgdep

.PHONY: demo
demo: doc/demo.gif

.PHONY: lint
lint:
	golangci-lint run --max-same-issues 10

build/gopkgdep: $(shell find . -name *.go)
	go build -o $@ .

doc/demo.gif: doc/demo.tape
	vhs < $<
