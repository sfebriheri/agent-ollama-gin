#!/bin/bash

# Test script for Encyclopedia Agent API
# Make sure the server is running on localhost:8080

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080/api/v1"
SERVER_URL="http://localhost:8080"

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    
    log_info "Testing $name..."
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "${BASE_URL}${endpoint}")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    # Extract status code and body
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 201 ]; then
        log_success "$name passed (HTTP $http_code)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        log_error "$name failed (HTTP $http_code)"
        echo "$body"
    fi
    
    echo
}

# Check if server is running
check_server() {
    log_info "Checking if server is running..."
    if curl -s "$SERVER_URL/health" >/dev/null 2>&1; then
        log_success "Server is running on $SERVER_URL"
    else
        log_error "Server is not running on $SERVER_URL"
        log_info "Please start the server first: go run main.go"
        exit 1
    fi
}

# Check dependencies
check_dependencies() {
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        log_warning "jq is not installed. Installing jq..."
        if command -v brew &> /dev/null; then
            brew install jq
        elif command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y jq
        elif command -v yum &> /dev/null; then
            sudo yum install -y jq
        else
            log_error "Cannot install jq automatically. Please install it manually."
            exit 1
        fi
    fi
}

# Main test execution
main() {
    echo -e "${BLUE}🔍 Encyclopedia Agent API Test Suite${NC}"
    echo "=========================================="
    
    check_dependencies
    check_server
    
    # Test encyclopedia endpoints
    test_endpoint "Health Check" "GET" "/encyclopedia/health"
    test_endpoint "Sources" "GET" "/encyclopedia/sources"
    test_endpoint "Languages" "GET" "/encyclopedia/languages"
    
    test_endpoint "Search (Wikipedia)" "POST" "/encyclopedia/search" '{
        "query": "artificial intelligence",
        "source": "wikipedia",
        "max_results": 3,
        "language": "en"
    }'
    
    test_endpoint "Search (Britannica)" "POST" "/encyclopedia/search" '{
        "query": "machine learning",
        "source": "britannica",
        "max_results": 2,
        "language": "en"
    }'
    
    test_endpoint "Search (All Sources)" "POST" "/encyclopedia/search" '{
        "query": "quantum computing",
        "source": "all",
        "max_results": 5,
        "language": "en"
    }'
    
    test_endpoint "Article Retrieval" "POST" "/encyclopedia/article" '{
        "title": "Artificial Intelligence",
        "source": "wikipedia",
        "language": "en",
        "max_length": 1000
    }'
    
    test_endpoint "Prompt Generation" "POST" "/encyclopedia/prompt" '{
        "topic": "neural networks",
        "style": "educational",
        "length": "medium",
        "language": "en"
    }'
    
    # Test Llama endpoints
    log_info "Testing Llama LLM endpoints..."
    test_endpoint "Llama Models" "GET" "/llama/models"
    
    test_endpoint "Llama Chat" "POST" "/llama/chat" '{
        "messages": [
            {"role": "user", "content": "Hello, how are you?"}
        ],
        "temperature": 0.7,
        "max_tokens": 100
    }'
    
    log_success "All tests completed!"
    echo
    log_info "🌐 Web Interface: $SERVER_URL/examples/encyclopedia_interface.html"
    log_info "📚 API Documentation: $SERVER_URL/"
}

# Run main function
main "$@"
