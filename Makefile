# Makefile for go-vault (CLI + server)

GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test

BINARY_CLI = go-vault
BINARY_SERVER = go-vault-server
BINARY_PATH = ./bin
CLI_PATH = ./cmd/cli
SERVER_PATH = ./cmd/server
COMPLETIONS_PATH = /tmp/completions

all: deps build

deps:
	$(GOCMD) mod tidy

build: build-cli build-server

build-cli:
	@mkdir -p $(BINARY_PATH)
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_CLI) -v $(CLI_PATH)

build-server:
	@mkdir -p $(BINARY_PATH)
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_SERVER) -v $(SERVER_PATH)

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_PATH)

test:
	$(GOTEST) -v ./...

dev: build-server
	@echo "Built $(BINARY_PATH)/$(BINARY_SERVER). Run: $(BINARY_PATH)/$(BINARY_SERVER)"

install: build
	@echo "Binaries installed to $(BINARY_PATH)/"

completions: build-cli
	@echo "Generating shell completions..."
	@mkdir -p $(COMPLETIONS_PATH)
	@case "$$SHELL" in \
        */bash) \
            $(BINARY_PATH)/$(BINARY_CLI) completion bash > $(COMPLETIONS_PATH)/$(BINARY_CLI).bash; \
            echo "source $(COMPLETIONS_PATH)/$(BINARY_CLI).bash" >> ~/.bashrc; \
        ;; \
        */zsh) \
            $(BINARY_PATH)/$(BINARY_CLI) completion zsh > $(COMPLETIONS_PATH)/_$(BINARY_CLI); \
            echo "source $(COMPLETIONS_PATH)/_$(BINARY_CLI)" >> ~/.zshrc; \
        ;; \
        */fish) \
            $(BINARY_PATH)/$(BINARY_CLI) completion fish > $(COMPLETIONS_PATH)/$(BINARY_CLI).fish; \
            echo "source $(COMPLETIONS_PATH)/$(BINARY_CLI).fish" >> ~/.config/fish/config.fish; \
        ;; \
        */pwsh) \
            $(BINARY_PATH)/$(BINARY_CLI) completion powershell > $(COMPLETIONS_PATH)/$(BINARY_CLI).ps1; \
            echo "source $(COMPLETIONS_PATH)/$(BINARY_CLI).ps1" >> ~/.config/powershell/Microsoft.PowerShell_profile.ps1; \
        ;; \
        *) \
            echo "Unsupported shell"; \
            exit 1; \
        ;; \
    esac

.PHONY: all deps build build-cli build-server clean test dev install completions
