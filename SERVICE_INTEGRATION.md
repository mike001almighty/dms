# Service Integration Guide

## How Other Services Can Authenticate with DMS

Your DMS is now secured with Keycloak JWT authentication. Here's how other services can integrate:

## 🔧 **For Backend Services**

### Option 1: Service Account (Recommended)
```bash
# Get service account token
curl -X POST http://localhost:8081/realms/dms/protocol/openid_connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=dms-service" \
  -d "client_secret=YOUR_CLIENT_SECRET"
```

### Option 2: User Password Grant (Development Only)
```bash
# Get user token
curl -X POST http://localhost:8081/realms/dms/protocol/openid_connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=dms-service" \
  -d "username=dms_admin" \
  -d "password=dms_admin_password"
```

## 📄 **Using DMS API with JWT**

### Create Document
```bash
curl -X POST http://localhost:8080/documents \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Document",
    "extension": "pdf",
    "description": "Document from external service",
    "content": "Document content here"
  }'
```

### Get Document
```bash
curl -X GET http://localhost:8080/documents/DOCUMENT_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Delete Document
```bash
curl -X DELETE http://localhost:8080/documents/DOCUMENT_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 💻 **Programming Examples**

### JavaScript/Node.js
```javascript
const axios = require('axios');

// Get token
async function getToken() {
  const response = await axios.post('http://localhost:8081/realms/dms/protocol/openid_connect/token', 
    new URLSearchParams({
      grant_type: 'client_credentials',
      client_id: 'dms-service',
      client_secret: 'YOUR_CLIENT_SECRET'
    })
  );
  return response.data.access_token;
}

// Use DMS API
async function createDocument(title, content) {
  const token = await getToken();
  const response = await axios.post('http://localhost:8080/documents', {
    title,
    extension: 'txt',
    description: 'Created from Node.js service',
    content
  }, {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  return response.data;
}
```

### Python
```python
import requests

def get_token():
    response = requests.post('http://localhost:8081/realms/dms/protocol/openid_connect/token', 
      data={
        'grant_type': 'client_credentials',
        'client_id': 'dms-service',
        'client_secret': 'YOUR_CLIENT_SECRET'
      }
    )
    return response.json()['access_token']

def create_document(title, content):
    token = get_token()
    response = requests.post('http://localhost:8080/documents',
      json={
        'title': title,
        'extension': 'txt', 
        'description': 'Created from Python service',
        'content': content
      },
      headers={'Authorization': f'Bearer {token}'}
    )
    return response.json()
```

### Go
```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type TokenResponse struct {
    AccessToken string `json:"access_token"`
}

func getToken() (string, error) {
    data := url.Values{}
    data.Set("grant_type", "client_credentials")
    data.Set("client_id", "dms-service")
    data.Set("client_secret", "YOUR_CLIENT_SECRET")
    
    resp, err := http.PostForm("http://localhost:8081/realms/dms/protocol/openid_connect/token", data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var tokenResp TokenResponse
    json.NewDecoder(resp.Body).Decode(&tokenResp)
    return tokenResp.AccessToken, nil
}

func createDocument(title, content string) error {
    token, err := getToken()
    if err != nil {
        return err
    }
    
    doc := map[string]string{
        "title": title,
        "extension": "txt",
        "description": "Created from Go service",
        "content": content,
    }
    
    jsonData, _ := json.Marshal(doc)
    req, _ := http.NewRequest("POST", "http://localhost:8080/documents", bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer " + token)
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    fmt.Println("Document created successfully")
    return nil
}
```

## 🛡️ **Security Notes**

1. **Production**: Use client credentials, not password grants
2. **HTTPS**: Always use HTTPS in production
3. **Token Storage**: Securely store and refresh tokens
4. **Client Secrets**: Keep client secrets secure
5. **Token Expiry**: Handle token expiration gracefully

## 📋 **Required Keycloak Setup for New Services**

1. Create new client in Keycloak admin console
2. Set "Service Accounts Enabled" to ON
3. Get client secret from "Credentials" tab
4. Assign appropriate roles for tenant access

## 🔗 **Service Endpoints**

- **Keycloak**: `http://localhost:8081`
- **DMS API**: `http://localhost:8080`
- **Health Check**: `http://localhost:8080/health`

## 📞 **Support**

For issues with authentication integration, check:
1. JWT token validity
2. Keycloak client configuration  
3. DMS service logs: `docker-compose logs app` 