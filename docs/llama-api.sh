#!/bin/bash

# setup.sh - Development environment setup for agent-ollama-gin

set -e

echo "🚀 Setting up development environment for agent-ollama-gin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        exit 1
    fi
    print_success "Go is installed: $(go version)"
}

# Check if Node.js is installed (for Genkit)
check_node() {
    if ! command -v node &> /dev/null; then
        print_warning "Node.js not found. Genkit may require Node.js for some features."
    else
        print_success "Node.js is installed: $(node --version)"
    fi
}

# Install Go dependencies
install_go_deps() {
    print_status "Installing Go dependencies..."
    go mod download
    go mod tidy
    print_success "Go dependencies installed"
}

# Install Firebase Genkit CLI
install_genkit() {
    print_status "Checking for Firebase Genkit CLI..."
    
    if command -v genkit &> /dev/null; then
        print_success "Genkit CLI already installed: $(genkit --version)"
    else
        print_status "Installing Firebase Genkit CLI..."
        curl -sL cli.genkit.dev | bash
        
        # Add to PATH if needed
        if ! command -v genkit &> /dev/null; then
            print_warning "Genkit installed but not in PATH. You may need to restart your terminal or run:"
            echo "export PATH=\$PATH:\$HOME/.local/bin"
        else
            print_success "Genkit CLI installed successfully"
        fi
    fi
}

# Install Air for hot reload
install_air() {
    print_status "Installing Air for hot reload..."
    go install github.com/air-verse/air@latest
    print_success "Air installed"
}

# Install other development tools
install_dev_tools() {
    print_status "Installing development tools..."
    
    # golangci-lint for linting
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    # gofumpt for formatting
    go install mvdan.cc/gofumpt@latest
    
    print_success "Development tools installed"
}

# Create necessary directories
create_directories() {
    print_status "Creating necessary directories..."
    
    mkdir -p bin
    mkdir -p logs
    mkdir -p config
    mkdir -p templates
    mkdir -p static/{css,js,images}
    
    print_success "Directories created"
}

# Create .air.toml configuration
create_air_config() {
    if [ ! -f .air.toml ]; then
        print_status "Creating Air configuration..."
        cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
        print_success "Air configuration created"
    fi
}

# Create development environment file
create_env_file() {
    if [ ! -f .env.example ]; then
        print_status "Creating example environment file..."
        cat > .env.example << 'EOF'
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Ollama Configuration
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama2

# Database Configuration (if using)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=agent_ollama
DB_USER=postgres
DB_PASSWORD=password

# Logging
LOG_LEVEL=info

# Genkit Configuration
GENKIT_ENV=dev
EOF
        print_success "Example environment file created"
    fi
}

# Initialize Genkit if needed
init_genkit() {
    print_status "Checking Genkit initialization..."
    
    if [ ! -f "genkit.config.js" ] && [ ! -f "genkit.config.ts" ]; then
        print_status "Initializing Genkit..."
        genkit init --non-interactive --template=blank
        print_success "Genkit initialized"
    else
        print_success "Genkit already initialized"
    fi
}

# Main setup function
main() {
    print_status "Starting setup process..."
    
    check_go
    check_node
    install_go_deps
    install_genkit
    install_air
    install_dev_tools
    create_directories
    create_air_config
    create_env_file
    
    # Initialize Genkit (optional - comment out if not needed)
    # init_genkit
    
    print_success "✅ Development environment setup complete!"
    echo
    echo "Next steps:"
    echo "1. Copy .env.example to .env and configure your settings"
    echo "2. Make sure Ollama is running: ollama serve"
    echo "3. Start development server: make watch"
    echo "4. Or use Genkit: genkit start"
    echo
    echo "Available commands:"
    echo "  make build       - Build the application"
    echo "  make run         - Run the application"
    echo "  make watch       - Start with hot reload"
    echo "  genkit start     - Start Genkit development server"
    echo "  genkit --help    - Show Genkit help"
}

# Run main function
main "$@"