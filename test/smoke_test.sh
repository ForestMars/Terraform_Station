#!/bin/bash

set -e

# Define colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting smoke test for local file creation...${NC}"

# Change to the directory containing our Terraform config
cd "$(dirname "$0")/../tofu"

# Clean up any existing state
rm -f terraform.tfstate* || true
rm -f hello.txt || true

# Initialize Terraform
echo -e "\n${GREEN}Initializing Terraform...${NC}"
terraform init

# Apply the configuration
echo -e "\n${GREEN}Applying Terraform configuration...${NC}"
terraform apply -auto-approve

# Verify the file was created
echo -e "\n${GREEN}Verifying file creation...${NC}"
if [ ! -f "hello.txt" ]; then
    echo -e "${RED}Error: hello.txt was not created${NC}"
    exit 1
fi

# Verify the file content
EXPECTED_CONTENT="Hello, World from OpenTofu!"
ACTUAL_CONTENT=$(cat hello.txt)

if [ "$EXPECTED_CONTENT" != "$ACTUAL_CONTENT" ]; then
    echo -e "${RED}Error: File content does not match expected${NC}"
    echo "Expected: $EXPECTED_CONTENT"
    echo "Found:    $ACTUAL_CONTENT"
    exit 1
fi

echo -e "${GREEN}File created successfully with expected content!${NC}"

# Clean up
echo -e "\n${GREEN}Cleaning up...${NC}"
terraform destroy -auto-approve
rm -f terraform.tfstate* || true
rm -f hello.txt || true

echo -e "\n${GREEN}Smoke test passed successfully!${NC}"
