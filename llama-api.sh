#!/bin/bash

# 🚀 Llama API - Complete Setup & Service Management
# This script provides all setup, development, and maintenance functionality

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "${BLUE}$1${NC}"
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check OS
get_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "macos";;
        CYGWIN*)   echo "windows";;
        MINGW*)    echo "windows";;
        *)         echo "unknown";;
    esac
}

# Function to show main menu
show_main_menu() {
    clear
    echo "🚀 Llama API - Complete Management System"
    echo "========================================="
    echo
    echo "Choose an option:"
    echo
    echo "1.  🚀 Quick Start (30 seconds)"
    echo "2.  ⚙️  Complete Setup"
    echo "3.  🔄 Start Development Environment"
    echo "4.  🛑 Stop Development Environment"
    echo "5.  🧪 Test API"
    echo "6.  📊 Check Status"
    echo "7.  🧹 Cleaning & Maintenance"
    echo "8.  🐳 Docker Management"
    echo "9.  🔧 Development Tools"
    echo "10. 📋 Project Information"
    echo "0.  🚪 Exit"
    echo
    echo -n "Enter your choice: "
}

# Function to show setup submenu
show_setup_menu() {
    clear
    echo "⚙️  Setup & Installation"
    echo "========================"
    echo
    echo "Choose setup option:"
    echo
    echo "1.  🚀 Quick Setup (recommended)"
    echo "2.  🔧 Complete Environment Setup"
    echo "3.  📦 Install Go Dependencies"
    echo "4.  🤖 Install Ollama"
    echo "5.  🐳 Install Docker"
    echo "6.  📥 Pull Default Model"
    echo "7.  🔙 Back to Main Menu"
    echo
    echo -n "Enter your choice: "
}

# Function to show cleaning submenu
show_cleaning_menu() {
    clear
    echo "🧹 Cleaning & Maintenance"
    echo "========================="
    echo
    echo "Choose cleaning option:"
    echo
    echo "1.  🧹 Basic Clean (build artifacts)"
    echo "2.  🧽 Deep Clean (all generated files)"
    echo "3.  🐳 Docker Clean (Docker resources)"
    echo "4.  💥 Complete Clean (everything)"
    echo "5.  📁 Clean specific directories"
    echo "6.  🔍 Show what will be cleaned"
    echo "7.  📊 Show disk usage"
    echo "8.  🔙 Back to Main Menu"
    echo
    echo -n "Enter your choice: "
}

# Function to show Docker submenu
show_docker_menu() {
    clear
    echo "🐳 Docker Management"
    echo "==================="
    echo
    echo "Choose Docker option:"
    echo
    echo "1.  🚀 Start Full Stack (Ollama + API + DB)"
    echo "2.  🤖 Start Only Ollama (recommended)"
    echo "3.  🛑 Stop All Services"
    echo "4.  📝 View Logs"
    echo "5.  🔄 Rebuild Services"
    echo "6.  🔙 Back to Main Menu"
    echo
    echo -n "Enter your choice: "
}

# Function to show development tools submenu
show_dev_tools_menu() {
    clear
    echo "🔧 Development Tools"
    echo "==================="
    echo
    echo "Choose tool option:"
    echo
    echo "1.  📥 Install Development Tools"
    echo "2.  👀 Run with Hot Reloading"
    echo "3.  🔨 Build Application"
    echo "4.  ▶️  Run Application"
    echo "5.  🧪 Run Tests"
    echo "6.  🔍 Health Check"
    echo "7.  🔙 Back to Main Menu"
    echo
    echo -n "Enter your choice: "
}

# =============================================================================
# SETUP FUNCTIONS
# =============================================================================

# Function to install Go
install_go() {
    if command_exists go; then
        print_success "Go is already installed: $(go version)"
        return 0
    fi
    
    print_header "Installing Go..."
    OS=$(get_os)
    
    if [ "$OS" = "macos" ]; then
        if command_exists brew; then
            brew install go
        else
            print_error "Homebrew not found. Please install Homebrew first: https://brew.sh/"
            return 1
        fi
    elif [ "$OS" = "linux" ]; then
        GO_VERSION="1.21.0"
        wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        rm "go${GO_VERSION}.linux-amd64.tar.gz"
        
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
        export PATH=$PATH:/usr/local/go/bin
    else
        print_error "Unsupported OS: $OS"
        return 1
    fi
    
    print_success "Go installed successfully: $(go version)"
}

# Function to install Docker
install_docker() {
    if command_exists docker; then
        print_success "Docker is already installed: $(docker --version)"
        return 0
    fi
    
    print_header "Installing Docker..."
    OS=$(get_os)
    
    if [ "$OS" = "macos" ]; then
        if command_exists brew; then
            brew install --cask docker
        else
            print_warning "Please install Docker Desktop for Mac from: https://www.docker.com/products/docker-desktop/"
            return 1
        fi
    elif [ "$OS" = "linux" ]; then
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        sudo usermod -aG docker $USER
        rm get-docker.sh
        print_warning "Please log out and back in for Docker group changes to take effect"
    else
        print_error "Unsupported OS: $OS"
        return 1
    fi
    
    print_success "Docker installed successfully"
}

# Function to install Ollama
install_ollama() {
    if command_exists ollama; then
        print_success "Ollama is already installed: $(ollama --version)"
        return 0
    fi
    
    print_header "Installing Ollama..."
    curl -fsSL https://ollama.ai/install.sh | sh
    
    print_success "Ollama installed successfully"
}

# Function to setup environment
setup_environment() {
    print_header "Setting up environment..."
    
    # Create .env file if it doesn't exist
    if [ ! -f .env ]; then
        print_header "Creating .env file from template..."
        cp env.example .env
        print_success ".env file created"
    else
        print_header ".env file already exists"
    fi
    
    # Create necessary directories
    mkdir -p bin
    mkdir -p logs
    
    print_success "Environment setup complete"
}

# Function to install Go dependencies
install_go_deps() {
    print_header "Installing Go dependencies..."
    
    if [ ! -f go.mod ]; then
        print_error "go.mod not found. Please run this script from the project root directory."
        return 1
    fi
    
    go mod download
    go mod tidy
    
    print_success "Go dependencies installed"
}

# Function to pull default model
pull_model() {
    print_header "Pulling default Llama model..."
    
    if ! command_exists ollama; then
        print_error "Ollama not installed. Please install Ollama first."
        return 1
    fi
    
    ollama pull llama2
    print_success "Default model pulled successfully"
}

# Function to build the application
build_app() {
    print_header "Building the application..."
    
    if [ ! -f main.go ]; then
        print_error "main.go not found. Please run this script from the project root directory."
        return 1
    fi
    
    go build -o bin/llama-api main.go
    
    if [ -f bin/llama-api ]; then
        print_success "Application built successfully: bin/llama-api"
    else
        print_error "Build failed"
        return 1
    fi
}

# =============================================================================
# DEVELOPMENT & SERVICE FUNCTIONS
# =============================================================================

# Function to start development environment
start_dev_environment() {
    print_header "Starting Llama API Development Environment"
    echo "============================================"
    echo
    
    # Check if Ollama is running
    echo -n "Checking Ollama status... "
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Running${NC}"
    else
        echo -e "${YELLOW}⚠ Not running${NC}"
        echo "Starting Ollama..."
        ollama serve > /dev/null 2>&1 &
        sleep 5
        
        # Check if Ollama started successfully
        if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
            echo -e "${GREEN}✓ Ollama started successfully${NC}"
        else
            echo -e "${RED}✗ Failed to start Ollama${NC}"
            echo "Please start Ollama manually: ollama serve"
            return 1
        fi
    fi
    
    # Check if default model is available
    echo -n "Checking default model... "
    if ollama list | grep -q "llama2"; then
        echo -e "${GREEN}✓ Available${NC}"
    else
        echo -e "${YELLOW}⚠ Not available${NC}"
        echo "Pulling default model (llama2)..."
        ollama pull llama2
        echo -e "${GREEN}✓ Model downloaded${NC}"
    fi
    
    # Check if .env exists
    if [ ! -f .env ]; then
        echo "Creating .env file from template..."
        cp env.example .env
        echo -e "${GREEN}✓ Environment file created${NC}"
    fi
    
    # Check if Go dependencies are installed
    if [ ! -f go.sum ]; then
        echo "Installing Go dependencies..."
        go mod download
        go mod tidy
        echo -e "${GREEN}✓ Dependencies installed${NC}"
    fi
    
    echo
    echo "🎯 Starting Llama API..."
    echo "Press Ctrl+C to stop the development environment"
    echo
    
    # Start the API
    go run main.go
}

# Function to stop development environment
stop_dev_environment() {
    print_header "Stopping Llama API Development Environment"
    echo "============================================"
    echo
    
    # Stop Go processes
    echo -n "Stopping Go processes... "
    GO_PIDS=$(pgrep -f "go run main.go" || true)
    if [ -n "$GO_PIDS" ]; then
        echo "$GO_PIDS" | xargs kill -TERM 2>/dev/null || true
        sleep 2
        echo "$GO_PIDS" | xargs kill -KILL 2>/dev/null || true
        echo -e "${GREEN}✓ Stopped${NC}"
    else
        echo -e "${YELLOW}⚠ No Go processes found${NC}"
    fi
    
    # Stop Ollama if it was started by us
    echo -n "Checking Ollama status... "
    if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠ Still running${NC}"
        echo "Note: Ollama is still running. If you want to stop it, run: pkill -f 'ollama serve'"
        echo "Or keep it running for faster development restarts."
    else
        echo -e "${GREEN}✓ Not running${NC}"
    fi
    
    # Stop any other related processes
    echo -n "Cleaning up other processes... "
    pkill -f "air" 2>/dev/null || true
    pkill -f "llama-api" 2>/dev/null || true
    echo -e "${GREEN}✓ Done${NC}"
    
    echo
    echo -e "${GREEN}✅ Development environment stopped${NC}"
}

# Function to test API
test_api() {
    print_header "Testing Llama API endpoints..."
    echo "=================================="
    
    BASE_URL="http://localhost:8080/api/v1"
    
    # Wait for API to be ready
    echo "Waiting for API to be ready..."
    sleep 3
    
    # Test health endpoint
    echo -e "\n1. Testing health endpoint..."
    curl -s "${BASE_URL}/health" | jq '.' || echo "Health check failed"
    
    # Test models endpoint
    echo -e "\n2. Testing models endpoint..."
    curl -s "${BASE_URL}/llama/models" | jq '.' || echo "Models endpoint failed"
    
    echo -e "\n✅ API testing completed!"
}

# Function to check status
check_status() {
    print_header "🔍 Llama API Development Environment Status"
    echo "=========================================="
    echo
    
    # Check Go installation
    echo -n "Go: "
    if command_exists go > /dev/null; then
        echo -e "${GREEN}✓ Installed${NC} ($(go version | awk '{print $3}'))"
    else
        echo -e "${RED}✗ Not installed${NC}"
    fi
    
    # Check Go modules
    echo -n "Go Modules: "
    if [ -f go.mod ]; then
        echo -e "${GREEN}✓ Found${NC}"
        echo "  Module: $(grep '^module' go.mod | awk '{print $2}')"
        echo "  Go version: $(grep '^go' go.mod | awk '{print $2}')"
    else
        echo -e "${RED}✗ Not found${NC}"
    fi
    
    # Check dependencies
    echo -n "Dependencies: "
    if [ -f go.sum ]; then
        echo -e "${GREEN}✓ Installed${NC}"
    else
        echo -e "${YELLOW}⚠ Not installed${NC}"
    fi
    
    # Check binary
    echo -n "Binary: "
    if [ -f bin/llama-api ]; then
        echo -e "${GREEN}✓ Built${NC}"
        echo "  Size: $(ls -lh bin/llama-api | awk '{print $5}')"
        echo "  Modified: $(ls -lh bin/llama-api | awk '{print $6, $7, $8}')"
    else
        echo -e "${YELLOW}⚠ Not built${NC}"
    fi
    
    # Check environment file
    echo -n "Environment: "
    if [ -f .env ]; then
        echo -e "${GREEN}✓ Configured${NC}"
    else
        echo -e "${YELLOW}⚠ Not configured${NC}"
    fi
    
    echo
    
    # Check Ollama
    echo -n "Ollama: "
    if command_exists ollama > /dev/null; then
        echo -e "${GREEN}✓ Installed${NC}"
        
        # Check if Ollama is running
        if curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
            echo -e "  Status: ${GREEN}✓ Running${NC}"
            
            # Check available models
            echo -n "  Models: "
            MODEL_COUNT=$(ollama list | wc -l)
            if [ $MODEL_COUNT -gt 1 ]; then
                echo -e "${GREEN}✓ Available ($((MODEL_COUNT-1)) models)${NC}"
                echo "    Available models:"
                ollama list | tail -n +2 | awk '{print "      - " $1}'
            else
                echo -e "${YELLOW}⚠ No models available${NC}"
            fi
        else
            echo -e "  Status: ${RED}✗ Not running${NC}"
        fi
    else
        echo -e "${RED}✗ Not installed${NC}"
    fi
    
    echo
    
    # Check API
    echo -n "API Service: "
    if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Running${NC}"
    else
        echo -e "${RED}✗ Not running${NC}"
    fi
    
    echo
    
    # Summary and recommendations
    echo "📋 Summary & Recommendations:"
    echo "============================="
    
    if command_exists go > /dev/null && [ -f go.mod ] && [ -f go.sum ] && [ -f bin/llama-api ] && [ -f .env ]; then
        echo -e "${GREEN}✓ Basic setup complete${NC}"
    else
        echo -e "${YELLOW}⚠ Basic setup incomplete${NC}"
        echo "  Run: Complete Setup from main menu"
    fi
    
    if command_exists ollama > /dev/null && curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Ollama ready${NC}"
    else
        echo -e "${YELLOW}⚠ Ollama not ready${NC}"
        echo "  Run: Install Ollama from setup menu"
    fi
    
    if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ API ready${NC}"
    else
        echo -e "${YELLOW}⚠ API ready${NC}"
        echo "  Run: Start Development Environment from main menu"
    fi
}

# =============================================================================
# CLEANING FUNCTIONS
# =============================================================================

# Function to basic clean
basic_clean() {
    print_header "🧹 Performing basic clean..."
    
    local cleaned_items=()
    
    # Clean build artifacts
    if [ -d "bin" ]; then
        rm -rf bin/
        cleaned_items+=("bin/")
    fi
    
    if [ -f "coverage.out" ]; then
        rm -f coverage.out
        cleaned_items+=("coverage.out")
    fi
    
    if [ -d "tmp" ]; then
        rm -rf tmp/
        cleaned_items+=("tmp/")
    fi
    
    if [ ${#cleaned_items[@]} -eq 0 ]; then
        print_warning "No build artifacts found to clean"
    else
        print_success "Cleaned: ${cleaned_items[*]}"
    fi
}

# Function to deep clean
deep_clean() {
    print_header "🧽 Performing deep clean..."
    
    local cleaned_items=()
    
    # Basic clean first
    basic_clean
    
    # Clean logs
    if [ -d "logs" ]; then
        rm -rf logs/
        cleaned_items+=("logs/")
    fi
    
    # Clean log files
    for log_file in *.log; do
        if [ -f "$log_file" ]; then
            rm -f "$log_file"
            cleaned_items+=("$log_file")
        fi
    done
    
    # Clean Go cache
    print_header "Cleaning Go cache..."
    go clean -modcache 2>/dev/null || true
    go clean -cache 2>/dev/null || true
    go clean -testcache 2>/dev/null || true
    cleaned_items+=("Go cache")
    
    # Clean environment file (optional)
    if [ -f ".env" ]; then
        read -p "Remove .env file? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -f .env
            cleaned_items+=(".env")
        fi
    fi
    
    # Clean backup files
    for backup_file in *.bak *.backup *.tmp; do
        if [ -f "$backup_file" ]; then
            rm -f "$backup_file"
            cleaned_items+=("$backup_file")
        fi
    done
    
    # Clean macOS files
    if [ -f ".DS_Store" ]; then
        rm -f .DS_Store
        cleaned_items+=(".DS_Store")
    fi
    
    if [ ${#cleaned_items[@]} -eq 0 ]; then
        print_warning "No additional files found to clean"
    else
        print_success "Deep cleaned: ${cleaned_items[*]}"
    fi
}

# Function to clean Docker
clean_docker() {
    print_header "🐳 Cleaning Docker resources..."
    
    local cleaned_items=()
    
    # Stop and remove containers
    if command -v docker-compose > /dev/null && [ -f "docker-compose.yml" ]; then
        print_header "Stopping Docker services..."
        docker-compose down --volumes --remove-orphans 2>/dev/null || true
        cleaned_items+=("Docker services")
    fi
    
    # Clean Docker system
    if command -v docker > /dev/null; then
        print_header "Cleaning Docker system..."
        docker system prune -f 2>/dev/null || true
        docker image prune -f 2>/dev/null || true
        cleaned_items+=("Docker system")
    else
        print_warning "Docker not found"
    fi
    
    print_success "Docker cleaned: ${cleaned_items[*]}"
}

# Function to complete clean
complete_clean() {
    print_header "💥 Performing complete cleanup..."
    
    # Deep clean first
    deep_clean
    
    # Docker clean
    clean_docker
    
    # Additional cleanup
    print_header "Final cleanup..."
    
    # Remove any remaining temporary files
    find . -name "*.tmp" -type f -delete 2>/dev/null || true
    find . -name "*.swp" -type f -delete 2>/dev/null || true
    find . -name "*~" -type f -delete 2>/dev/null || true
    
    # Remove empty directories
    find . -type d -empty -delete 2>/dev/null || true
    
    print_success "✅ Complete cleanup finished!"
}

# Function to show clean preview
show_clean_preview() {
    print_header "🔍 Clean preview - Files that would be removed:"
    echo
    
    local items=()
    
    # Build artifacts
    if [ -d "bin" ]; then
        items+=("bin/ (build output)")
    fi
    
    if [ -f "coverage.out" ]; then
        items+=("coverage.out (test coverage)")
    fi
    
    if [ -d "tmp" ]; then
        items+=("tmp/ (temporary files)")
    fi
    
    if [ -d "logs" ]; then
        items+=("logs/ (log files)")
    fi
    
    # Log files
    for log_file in *.log; do
        if [ -f "$log_file" ]; then
            items+=("$log_file (log file)")
        fi
    done
    
    # Backup files
    for backup_file in *.bak *.backup *.tmp; do
        if [ -f "$backup_file" ]; then
            items+=("$backup_file (backup file)")
        fi
    done
    
    # macOS files
    if [ -f ".DS_Store" ]; then
        items+=(".DS_Store (macOS file)")
    fi
    
    if [ ${#items[@]} -eq 0 ]; then
        print_warning "No files found to clean"
    else
        echo "Files and directories to be cleaned:"
        for item in "${items[@]}"; do
            echo "  - $item"
        done
    fi
    
    echo
    print_warning "Note: This is a preview. No files have been deleted."
}

# Function to show disk usage
show_disk_usage() {
    print_header "📊 Disk usage information:"
    echo
    
    if command -v du > /dev/null; then
        echo "Current directory size:"
        du -sh . 2>/dev/null || echo "Unable to calculate directory size"
        
        echo
        echo "Largest files/directories:"
        du -sh * 2>/dev/null | sort -hr | head -10 || echo "Unable to list directory sizes"
    else
        print_warning "du command not available"
    fi
    
    echo
    if command -v df > /dev/null; then
        echo "Disk space available:"
        df -h . 2>/dev/null || echo "Unable to show disk space"
    else
        print_warning "df command not available"
    fi
}

# =============================================================================
# DOCKER FUNCTIONS
# =============================================================================

# Function to start Docker services
start_docker_services() {
    print_header "Starting Docker services..."
    
    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found"
        return 1
    fi
    
    docker-compose up -d
    print_success "Docker services started"
}

# Function to start only Ollama with Docker
start_ollama_docker() {
    print_header "Starting Ollama with Docker..."
    
    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found"
        return 1
    fi
    
    docker-compose up -d ollama
    print_success "Ollama started with Docker"
}

# Function to stop Docker services
stop_docker_services() {
    print_header "Stopping Docker services..."
    
    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found"
        return 1
    fi
    
    docker-compose down
    print_success "Docker services stopped"
}

# Function to view Docker logs
view_docker_logs() {
    print_header "Docker logs..."
    
    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found"
        return 1
    fi
    
    docker-compose logs -f
}

# Function to rebuild Docker services
rebuild_docker_services() {
    print_header "Rebuilding Docker services..."
    
    if [ ! -f "docker-compose.yml" ]; then
        print_error "docker-compose.yml not found"
        return 1
    fi
    
    docker-compose up -d --build
    print_success "Docker services rebuilt and started"
}

# =============================================================================
# DEVELOPMENT TOOLS FUNCTIONS
# =============================================================================

# Function to install development tools
install_dev_tools() {
    print_header "Installing development tools..."
    
    go install github.com/air-verse/air@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    print_success "Development tools installed"
}

# Function to run with hot reloading
run_hot_reload() {
    print_header "Starting with hot reloading..."
    
    if ! command -v air > /dev/null; then
        print_warning "Air not found. Installing..."
        install_dev_tools
    fi
    
    air
}

# Function to run tests
run_tests() {
    print_header "Running tests..."
    go test -v ./...
}

# Function to run health check
run_health_check() {
    print_header "Running health check..."
    
    echo "Checking API health..."
    curl -s http://localhost:8080/api/v1/health | jq . || echo "API not running"
    
    echo "Checking Ollama..."
    curl -s http://localhost:11434/api/tags > /dev/null && echo "Ollama: OK" || echo "Ollama: Not running"
}

# Function to show project info
show_project_info() {
    print_header "Llama API Project Information"
    echo "================================="
    echo
    
    echo "Go version: $(go version)"
    echo "Project path: $(pwd)"
    echo "Go modules: $(go list -m all | head -5)"
    
    echo
    echo "Available commands:"
    echo "  - ./llama-api.sh: This unified script (all functionality)"
    echo "  - make help: Show all Makefile commands"
    echo "  - make run: Run the application"
    echo "  - make test: Run tests"
    echo "  - make watch: Hot reloading development"
    
    echo
    echo "Makefile commands: make help"
}

# =============================================================================
# MENU HANDLERS
# =============================================================================

# Function to handle main menu selection
handle_main_menu() {
    case $1 in
        1)
            quick_start
            ;;
        2)
            show_setup_menu
            handle_setup_menu
            ;;
        3)
            start_dev_environment
            ;;
        4)
            stop_dev_environment
            ;;
        5)
            test_api
            ;;
        6)
            check_status
            ;;
        7)
            show_cleaning_menu
            handle_cleaning_menu
            ;;
        8)
            show_docker_menu
            handle_docker_menu
            ;;
        9)
            show_dev_tools_menu
            handle_dev_tools_menu
            ;;
        10)
            show_project_info
            ;;
        0)
            print_success "Goodbye! 👋"
            exit 0
            ;;
        *)
            print_error "Invalid option. Please try again."
            ;;
    esac
}

# Function to handle setup menu selection
handle_setup_menu() {
    while true; do
        read -r choice
        
        case $choice in
            1)
                quick_start
                break
                ;;
            2)
                complete_setup
                break
                ;;
            3)
                install_go_deps
                ;;
            4)
                install_ollama
                ;;
            5)
                install_docker
                ;;
            6)
                pull_model
                ;;
            7)
                break
                ;;
            *)
                print_error "Invalid choice. Please try again."
                ;;
        esac
        
        echo
        read -p "Press Enter to continue..."
        show_setup_menu
    done
}

# Function to handle cleaning menu selection
handle_cleaning_menu() {
    while true; do
        read -r choice
        
        case $choice in
            1)
                basic_clean
                ;;
            2)
                deep_clean
                ;;
            3)
                clean_docker
                ;;
            4)
                read -p "⚠️  This will remove ALL generated files and Docker resources. Continue? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    complete_clean
                else
                    print_warning "Complete cleanup cancelled"
                fi
                ;;
            5)
                clean_specific_directories
                ;;
            6)
                show_clean_preview
                ;;
            7)
                show_disk_usage
                ;;
            8)
                break
                ;;
            *)
                print_error "Invalid choice. Please try again."
                ;;
        esac
        
        echo
        read -p "Press Enter to continue..."
        show_cleaning_menu
    done
}

# Function to handle Docker menu selection
handle_docker_menu() {
    while true; do
        read -r choice
        
        case $choice in
            1)
                start_docker_services
                ;;
            2)
                start_ollama_docker
                ;;
            3)
                stop_docker_services
                ;;
            4)
                view_docker_logs
                ;;
            5)
                rebuild_docker_services
                ;;
            6)
                break
                ;;
            *)
                print_error "Invalid choice. Please try again."
                ;;
        esac
        
        echo
        read -p "Press Enter to continue..."
        show_docker_menu
    done
}

# Function to handle development tools menu selection
handle_dev_tools_menu() {
    while true; do
        read -r choice
        
        case $choice in
            1)
                install_dev_tools
                ;;
            2)
                run_hot_reload
                ;;
            3)
                build_app
                ;;
            4)
                start_dev_environment
                ;;
            5)
                run_tests
                ;;
            6)
                run_health_check
                ;;
            7)
                break
                ;;
            *)
                print_error "Invalid choice. Please try again."
                ;;
        esac
        
        echo
        read -p "Press Enter to continue..."
        show_dev_tools_menu
    done
}

# =============================================================================
# MAIN FUNCTIONS
# =============================================================================

# Function to quick start
quick_start() {
    print_header "🚀 Llama API Quick Start"
    echo "========================"
    echo
    
    # Check if .env exists
    if [ ! -f .env ]; then
        echo "📝 Creating .env file..."
        cp env.example .env
        echo "✅ Environment file created"
    fi
    
    # Check if Go dependencies are installed
    if [ ! -d vendor ] && [ ! -f go.sum ]; then
        echo "📦 Installing Go dependencies..."
        go mod download
        go mod tidy
        echo "✅ Dependencies installed"
    fi
    
    # Check if binary exists
    if [ ! -f bin/llama-api ]; then
        echo "🔨 Building application..."
        mkdir -p bin
        go build -o bin/llama-api main.go
        echo "✅ Application built"
    fi
    
    # Check if Ollama is running
    if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
        echo "🤖 Starting Ollama..."
        echo "   Note: This will start Ollama in the background"
        echo "   To stop it later, run: pkill -f 'ollama serve'"
        ollama serve > /dev/null 2>&1 &
        sleep 5
        echo "✅ Ollama started"
    else
        echo "✅ Ollama is already running"
    fi
    
    # Check if default model exists
    if ! ollama list | grep -q "llama2"; then
        echo "📥 Pulling default model (llama2)..."
        ollama pull llama2
        echo "✅ Model downloaded"
    else
        echo "✅ Default model already available"
    fi
    
    echo
    echo "🎉 Ready to start development!"
    echo
    echo "Available commands:"
    echo "  ./llama-api.sh     - This unified script"
    echo "  make run           - Run the application"
    echo "  make test          - Run tests"
    echo
    echo "Quick test:"
    echo "  curl http://localhost:8080/api/v1/health"
    echo
    echo "To start the API, run: make run"
    echo "Or use this script: ./llama-api.sh (option 3)"
}

# Function to complete setup
complete_setup() {
    print_header "Running complete setup..."
    
    # Install dependencies
    install_go
    install_docker
    install_ollama
    
    # Setup environment
    setup_environment
    install_go_deps
    
    # Build and setup
    build_app
    pull_model
    
    print_success "✅ Complete setup finished!"
}

# Function to clean specific directories
clean_specific_directories() {
    print_header "📁 Cleaning specific directories..."
    
    echo "Available directories to clean:"
    echo "1. bin/ (build output)"
    echo "2. tmp/ (temporary files)"
    echo "3. logs/ (log files)"
    echo "4. Custom directory"
    
    read -p "Enter choice (1-4): " choice
    
    case $choice in
        1)
            if [ -d "bin" ]; then
                rm -rf bin/
                print_success "Cleaned bin/ directory"
            else
                print_warning "bin/ directory not found"
            fi
            ;;
        2)
            if [ -d "tmp" ]; then
                rm -rf tmp/
                print_success "Cleaned tmp/ directory"
            else
                print_warning "tmp/ directory not found"
            fi
            ;;
        3)
            if [ -d "logs" ]; then
                rm -rf logs/
                print_success "Cleaned logs/ directory"
            else
                print_warning "logs/ directory not found"
            fi
            ;;
        4)
            read -p "Enter directory path to clean: " custom_dir
            if [ -d "$custom_dir" ]; then
                read -p "Remove directory '$custom_dir'? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    rm -rf "$custom_dir"
                    print_success "Cleaned directory: $custom_dir"
                else
                    print_warning "Directory cleanup cancelled"
                fi
            else
                print_error "Directory '$custom_dir' not found"
            fi
            ;;
        *)
            print_error "Invalid choice"
            ;;
    esac
}

# =============================================================================
# MAIN EXECUTION
# =============================================================================

# Main function
main() {
    # Check if running from project root
    if [ ! -f main.go ] || [ ! -f go.mod ]; then
        print_error "Please run this script from the project root directory"
        exit 1
    fi
    
    # Main loop
    while true; do
        show_main_menu
        read -r choice
        
        if [ -z "$choice" ]; then
            print_warning "Please enter a valid option"
            sleep 1
            continue
        fi
        
        handle_main_menu "$choice"
        
        echo
        read -p "Press Enter to continue..."
    done
}

# Run main function
main "$@"
