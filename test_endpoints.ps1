# DMS Endpoint Testing Script

Write-Host "=== DMS Endpoint Testing ===" -ForegroundColor Green

# Test 1: Health endpoints (no auth required)
Write-Host "`n1. Testing basic health endpoint..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
    Write-Host "✓ Basic health check successful:" -ForegroundColor Green
    $healthResponse | ConvertTo-Json
} catch {
    Write-Host "✗ Basic health check failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n2. Testing detailed health endpoint..." -ForegroundColor Yellow
try {
    $detailedHealthResponse = Invoke-RestMethod -Uri "http://localhost:8080/health/detailed" -Method GET
    Write-Host "✓ Detailed health check successful:" -ForegroundColor Green
    $detailedHealthResponse | ConvertTo-Json
} catch {
    Write-Host "✗ Detailed health check failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Protected endpoints (require auth)
Write-Host "`n3. Testing document creation with Basic Auth..." -ForegroundColor Yellow

# Create Basic Auth header
$username = "dms_admin"
$password = "dms_admin_password"
$base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes("${username}:${password}"))

$headers = @{
    "Authorization" = "Basic $base64AuthInfo"
    "Content-Type" = "application/json"
}

$documentData = @{
    title = "Test Document"
    extension = "txt"
    description = "Test document created via PowerShell"
    content = "This is test content from PowerShell script"
} | ConvertTo-Json

Write-Host "Attempting to create document with credentials: $username" -ForegroundColor Cyan

try {
    $createResponse = Invoke-RestMethod -Uri "http://localhost:8080/documents" -Method POST -Headers $headers -Body $documentData
    Write-Host "✓ Document creation successful:" -ForegroundColor Green
    $createResponse | ConvertTo-Json
    
    # Test 4: Get the created document
    if ($createResponse.id) {
        Write-Host "`n4. Testing document retrieval..." -ForegroundColor Yellow
        $getHeaders = @{
            "Authorization" = "Basic $base64AuthInfo"
        }
        try {
            $getResponse = Invoke-RestMethod -Uri "http://localhost:8080/documents/$($createResponse.id)" -Method GET -Headers $getHeaders
            Write-Host "✓ Document retrieval successful:" -ForegroundColor Green
            $getResponse | ConvertTo-Json
        } catch {
            Write-Host "✗ Document retrieval failed: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
} catch {
    Write-Host "✗ Document creation failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        try {
            $responseStream = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($responseStream)
            $responseBody = $reader.ReadToEnd()
            Write-Host "Response body: $responseBody" -ForegroundColor Red
        } catch {
            Write-Host "Could not read error response body" -ForegroundColor Red
        }
    }
}

# Test 5: Test with wrong credentials
Write-Host "`n5. Testing with wrong credentials (should fail)..." -ForegroundColor Yellow
$wrongBase64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes("wrong:credentials"))
$wrongHeaders = @{
    "Authorization" = "Basic $wrongBase64AuthInfo"
    "Content-Type" = "application/json"
}

try {
    $wrongResponse = Invoke-RestMethod -Uri "http://localhost:8080/documents" -Method POST -Headers $wrongHeaders -Body $documentData
    Write-Host "✗ This should have failed but didn't!" -ForegroundColor Red
} catch {
    Write-Host "✓ Correctly rejected wrong credentials: $($_.Exception.Message)" -ForegroundColor Green
}

Write-Host "`n=== Testing Complete ===" -ForegroundColor Green