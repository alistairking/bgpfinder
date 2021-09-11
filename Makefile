GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINDIR=./bin
CLI=bgpf

all: cli

dev: mod-tidy codegen pkg cli test

cli:
	mkdir -p $(BINDIR)
	$(GOBUILD) -o $(BINDIR)/$(CLI) -v ./cmd/$(CLI)

pkg:
	$(GOBUILD) ./

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINDIR)

mod-tidy:
	$(GOMOD) tidy

codegen:
	go install github.com/alvaroloes/enumer@latest
	go generate ./...

run: cli
	$(BINDIR)/$(CLI)
