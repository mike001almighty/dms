#!/bin/bash

echo "Testing Keycloak Authentication..."

# First, we need to get a token from Keycloak
echo "Getting JWT token from Keycloak..."

# For now, this will test basic token retrieval
# Note: We need to create a client first in Keycloak before this works
TOKEN_RESPONSE=$(curl -s -X POST \
  "http://localhost:8081/realms/dms/protocol/openid_connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=dms-service" \
  -d "username=dms_admin" \
  -d "password=dms_admin_password")

echo "Token response:"
echo $TOKEN_RESPONSE

# Extract token (if successful)
TOKEN=$(echo $TOKEN_RESPONSE | grep -o '"access_token":"[^"]*' | grep -o '[^"]*$')

if [ -n "$TOKEN" ]; then
    echo "✅ Successfully got JWT token!"
    echo "Token: $TOKEN"
    
    echo "Testing DMS API with token..."
    curl -X POST http://localhost:8080/documents \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "title": "Test Document",
        "extension": "txt",
        "description": "Test document created via authenticated API",
        "content": "This is a test document content"
      }'
else
    echo "❌ Failed to get token. Response:"
    echo $TOKEN_RESPONSE
fi 