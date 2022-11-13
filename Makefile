.PHONE: all
all: build/gopkgdep

.PHONY: build
build: build/gopkgdep

.PHONY: install
install:
	go install .

.PHONY: demo
demo: install doc/demo.gif doc/demo.svg doc/demo2.svg

.PHONY: lint
lint:
	golangci-lint run --max-same-issues 10

GO_SRCS = go.mod go.sum $(shell find . -name *.go)
build/gopkgdep: $(GO_SRCS)
	go build -o $@ .

doc/demo.gif: doc/demo.tape $(GO_SRCS)
	vhs < $<

doc/demo.svg: $(GO_SRCS)
	gopkgdep -dot | dot -Tsvg > $@.tmp && rm -f $@ && mv $@.tmp $@ || rm -f $@.tmp

doc/demo2.svg: $(GO_SRCS)
	gopkgdep -closed -dot ./internal/... | dot -Tsvg > $@.tmp && rm -f $@ && mv $@.tmp $@ || rm -f $@.tmp
