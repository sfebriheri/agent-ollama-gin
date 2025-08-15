#!/bin/bash

# Encyclopedia Agent CLI Demo Script
# This script demonstrates the capabilities of the encyclopedia CLI

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Encyclopedia Agent CLI Demo${NC}"
echo "=================================="
echo

# Check if CLI exists
if [ ! -f "./bin/encyclopedia" ]; then
    echo -e "${YELLOW}⚠️  CLI not found. Building it first...${NC}"
    make build-cli
    echo
fi

# Check if server is running
echo -e "${BLUE}🔍 Checking if server is running...${NC}"
if ! curl -s http://localhost:8080/api/v1/encyclopedia/health >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Server is not running. Please start it first:${NC}"
    echo "   go run main.go"
    echo
    echo -e "${BLUE}💡 You can also run this demo after starting the server${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Server is running!${NC}"
echo

# Demo 1: Show available sources
echo -e "${BLUE}📚 Demo 1: Available Sources${NC}"
echo "--------------------------------"
./bin/encyclopedia sources
echo

# Demo 2: Show supported languages
echo -e "${BLUE}🌍 Demo 2: Supported Languages${NC}"
echo "------------------------------------"
./bin/encyclopedia languages
echo

# Demo 3: Search for articles
echo -e "${BLUE}🔍 Demo 3: Search for Articles${NC}"
echo "----------------------------------"
echo "Searching for 'artificial intelligence' across all sources..."
./bin/encyclopedia search "artificial intelligence" all en 3
echo

# Demo 4: Get specific article
echo -e "${BLUE}📖 Demo 4: Get Specific Article${NC}"
echo "-----------------------------------"
echo "Retrieving article about 'Machine Learning'..."
./bin/encyclopedia article "Machine Learning" wikipedia en 1000
echo

# Demo 5: Generate AI prompt
echo -e "${BLUE}✍️  Demo 5: Generate AI Prompt${NC}"
echo "--------------------------------"
echo "Generating educational prompt about 'neural networks'..."
./bin/encyclopedia prompt "neural networks" educational medium en
echo

# Demo 6: Health check
echo -e "${BLUE}🏥 Demo 6: Service Health${NC}"
echo "----------------------------"
./bin/encyclopedia health
echo

# Demo 7: Interactive mode preview
echo -e "${BLUE}🎯 Demo 7: Interactive Mode Preview${NC}"
echo "----------------------------------------"
echo "The CLI also supports interactive mode:"
echo "  ./bin/encyclopedia"
echo
echo "Available interactive commands:"
echo "  • search <query> [source] [language] [max_results]"
echo "  • article <title> [source] [language] [max_length]"
echo "  • prompt <topic> [style] [length] [language]"
echo "  • sources, languages, health, web, clear, help, quit"
echo

# Demo 8: Command line usage examples
echo -e "${BLUE}💻 Demo 8: Command Line Usage Examples${NC}"
echo "-------------------------------------------"
echo "Single commands for automation:"
echo "  ./bin/encyclopedia search 'quantum computing' wikipedia en 5"
echo "  ./bin/encyclopedia article 'Quantum Computing' wikipedia en 2000"
echo "  ./bin/encyclopedia prompt 'blockchain' academic long en"
echo

# Demo 9: Integration examples
echo -e "${BLUE}🔗 Demo 9: Integration Examples${NC}"
echo "----------------------------------"
echo "Save search results to file:"
echo "  ./bin/encyclopedia search 'AI' > ai_results.txt"
echo
echo "Generate prompt and extract keywords:"
echo "  ./bin/encyclopedia prompt 'machine learning' | grep -i 'learning'"
echo
echo "Batch processing:"
echo "  for topic in 'AI' 'ML' 'DL'; do"
echo "    ./bin/encyclopedia search \"\$topic\" wikipedia en 3"
echo "  done"
echo

echo -e "${GREEN}🎉 CLI Demo Completed!${NC}"
echo
echo -e "${BLUE}📚 Next Steps:${NC}"
echo "1. Try interactive mode: ./bin/encyclopedia"
echo "2. Explore different topics and sources"
echo "3. Generate prompts in various styles"
echo "4. Integrate with your scripts and workflows"
echo
echo -e "${BLUE}🌐 Web Interface:${NC}"
echo "Open: http://localhost:8080/examples/encyclopedia_interface.html"
echo
echo -e "${BLUE}📖 Documentation:${NC}"
echo "See CLI_USAGE.md for detailed usage information"
echo
echo -e "${GREEN}Happy exploring with your Encyclopedia Agent CLI! 🚀${NC}"
