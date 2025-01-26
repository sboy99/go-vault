# Makefile for Cobra CLI Go project

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test

# Output binary configuration
BINARY_NAME = go-vault
BINARY_PATH = ./bin
MAIN_PATH = ./cmd/go-vault
COMPLETIONS_PATH = /tmp/completions

# Main targets
all: deps build

# Install and tidy dependencies using go.mod
deps:
	$(GOCMD) mod tidy

# Build the binary and place it in the specified directory
build:
	@mkdir -p $(BINARY_PATH)
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_NAME) -v $(MAIN_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_PATH)

# Run tests
test:
	$(GOTEST) -v ./...

# Development - watch and rebuild on file changes
dev:
	@mkdir -p $(BINARY_PATH)
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "Starting $(BINARY_NAME) in development mode..."
	@$(BINARY_PATH)/$(BINARY_NAME) serve

# Install the binary globally
install:
	@mkdir -p $(BINARY_PATH)
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_NAME) -v ./cmd/yourcobracli
	@echo "Binary installed to $(BINARY_PATH)/$(BINARY_NAME)"

# Generate shell completions for the CLI
completions:
	@echo "Generating shell completions..."
	@mkdir -p $(COMPLETIONS_PATH)
	@case "$$SHELL" in \
        */bash) \
            echo "Detected bash shell"; \
            $(BINARY_PATH)/$(BINARY_NAME) completion bash > $(COMPLETIONS_PATH)/$(BINARY_NAME).bash; \
            echo "Generated Bash completions"; \
       	    echo "source $(COMPLETIONS_PATH)/$(BINARY_NAME).bash" >> ~/.bashrc; \
            echo "Sourced completions, Reload the shell"; \
        ;; \
        */zsh) \
            echo "Detected zsh shell"; \
            $(BINARY_PATH)/$(BINARY_NAME) completion zsh > $(COMPLETIONS_PATH)/_$(BINARY_NAME); \
            echo "Generated Zsh completions"; \
       	    echo "source $(COMPLETIONS_PATH)/_$(BINARY_NAME)" >> ~/.zshrc; \
        	echo "Sourced completions, Reload the shell"; \
        ;; \
        */fish) \
            echo "Detected fish shell"; \
            $(BINARY_PATH)/$(BINARY_NAME) completion fish > $(COMPLETIONS_PATH)/$(BINARY_NAME).fish; \
        ;; \
        */pwsh) \
            echo "Detected powershell shell";\
            $(BINARY_PATH)/$(BINARY_NAME) completion powershell > $(COMPLETIONS_PATH)/$(BINARY_NAME).ps1; \
        ;; \
        *) \
            echo "Unsupported shell"; \
            exit 1; \
        ;; \
    esac

source:
	@echo "Sourcing shell completions..."
	@echo "source $(COMPLETIONS_PATH)/$(BINARY_NAME).bash" >> ~/.bashrc
	@echo "source $(COMPLETIONS_PATH)/_$(BINARY_NAME)" >> ~/.zshrc
	@echo "source $(COMPLETIONS_PATH)/$(BINARY_NAME).fish" >> ~/.config/fish/config.fish
	@echo "source $(COMPLETIONS_PATH)/$(BINARY_NAME).ps1" >> ~/.config/powershell/Microsoft.PowerShell_profile.ps1

# Ensure directories are created
.PHONY: all deps build clean test dev install completions
