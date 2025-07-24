# Keycloak Setup Guide

## Step-by-Step Configuration

### 1. Access Admin Console
- URL: `http://localhost:8081/admin`
- Username: `admin`
- Password: `admin`

### 2. Create DMS Realm
1. Click dropdown next to "Master" (top-left)
2. Click "Add realm"
3. Name: `dms`
4. Click "Create"

### 3. Create DMS Service Client
1. Go to "Clients" → "Create"
2. Client ID: `dms-service`
3. Client Protocol: `openid-connect`
4. Click "Save"

5. **In Settings tab (CRITICAL):**
   - Access Type: `confidential`
   - Direct Access Grants Enabled: `ON`
   - Service Accounts Enabled: `ON`
   - Valid Redirect URIs: `*`
   - Click "Save"

### 4. Create Test User
1. Go to "Users" → "Add user"
2. Username: `dms_admin`
3. Enabled: `ON`
4. Click "Save"

5. **Set Password (Credentials tab):**
   - New Password: `dms_admin_password`
   - Temporary: `OFF` (IMPORTANT!)
   - Click "Set Password" → Confirm

### 5. Verify Configuration
Test this URL in browser (should return JSON, not 404):
`http://localhost:8081/realms/dms/.well-known/openid_connect_configuration`

## Troubleshooting

If you still get 404 errors on OpenID Connect endpoints:

1. **Restart Keycloak completely:**
   ```bash
   docker-compose restart keycloak
   # Wait 2 minutes for full initialization
   ```

2. **Check logs:**
   ```bash
   docker-compose logs keycloak --tail=20
   ```

3. **Alternative: Fresh Keycloak install:**
   ```bash
   docker-compose down
   docker volume rm dms_keycloak_data
   docker-compose up -d
   # Wait 3 minutes, then redo configuration
   ```

## Testing Authentication

After setup is complete, test with:
```bash
curl -X POST http://localhost:8081/realms/dms/protocol/openid_connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=dms-service" \
  -d "username=dms_admin" \
  -d "password=dms_admin_password"
```

This should return a JWT token, not an error. 