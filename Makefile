GO=go
GOV=govendor
GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean
GOGET=$(GOV) fetch -v

fetch:
  @$(GOGET) github.com/urfave/cli

build:
  @echo "start building..."
  $(GOINSTALL)
  @echo "Yay! build DONE!"

all: fetch build
