#!/bin/bash

echo "Testing Basic Authentication..."

# Test with username and password directly to DMS API
echo "Testing DMS API with Basic Auth..."

curl -X POST http://localhost:8080/documents \
  -u "dms_admin:dms_admin_password" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Document", 
    "extension": "txt",
    "description": "Test document created via Basic Auth",
    "content": "This is a test document content"
  }'

echo -e "\n\nTesting document retrieval with Basic Auth..."
curl -X GET http://localhost:8080/documents/1 \
  -u "dms_admin:dms_admin_password"

echo -e "\n\nTesting health endpoint (no auth required)..."
curl -X GET http://localhost:8080/health 