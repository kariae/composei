GO=go
GOV=govendor
GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean
GOGET=$(GOV) fetch -v

fetch:
	@$(GOGET) github.com/urfave/cli
	@$(GOGET) gopkg.in/yaml.v2
	@$(GOGET) github.com/fatih/color
	@$(GOGET) github.com/coreos/go-semver/semver

build:
	@echo "start building..."
	$(GOINSTALL)
	@echo "Yay! build DONE!"

all: fetch build
