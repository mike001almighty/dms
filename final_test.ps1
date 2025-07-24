Write-Host "Final Authentication Test" -ForegroundColor Green
Write-Host "=========================" -ForegroundColor Green

$tokenUrl = "http://localhost:8081/realms/dms/protocol/openid_connect/token"

Write-Host "`nStep 1: Getting JWT token from Keycloak..." -ForegroundColor Yellow

$body = @{
    grant_type = "password"
    client_id = "dms-service"
    username = "dms_admin"
    password = "dms_admin_password"
}

try {
    $tokenResponse = Invoke-RestMethod -Uri $tokenUrl -Method POST -Body $body -ContentType "application/x-www-form-urlencoded"
    Write-Host "✅ Successfully obtained JWT token!" -ForegroundColor Green
    
    Write-Host "`nStep 2: Testing DMS API with real JWT..." -ForegroundColor Yellow
    
    $headers = @{
        "Authorization" = "Bearer $($tokenResponse.access_token)"
        "Content-Type" = "application/json"
    }
    
    $documentData = @{
        title = "Production Test Document"
        extension = "pdf"
        description = "Document created with real Keycloak authentication"
        content = "This document was created using proper JWT authentication"
    } | ConvertTo-Json
    
    $dmsResponse = Invoke-RestMethod -Uri "http://localhost:8080/documents" -Method POST -Body $documentData -Headers $headers
    
    Write-Host "🎉 SUCCESS! Full authentication working!" -ForegroundColor Green
    Write-Host "Document ID: $($dmsResponse.id)" -ForegroundColor Cyan
    Write-Host "Title: $($dmsResponse.title)" -ForegroundColor Cyan
    Write-Host "Tenant: $($dmsResponse.tenant_id)" -ForegroundColor Cyan
    
    Write-Host "`n✅ AUTHENTICATION INTEGRATION COMPLETE!" -ForegroundColor Green
    Write-Host "Your DMS is now secured with Keycloak JWT authentication." -ForegroundColor White
    
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $errorBody = $reader.ReadToEnd()
        Write-Host "Error details: $errorBody" -ForegroundColor Yellow
        
        if ($errorBody -match "404") {
            Write-Host "`n💡 The Keycloak realm still needs to be configured." -ForegroundColor Cyan
            Write-Host "Follow the setup guide in keycloak_setup_guide.md" -ForegroundColor White
        }
    }
} 