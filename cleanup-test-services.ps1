# LazyService Test Cleanup Script
# Stops and removes test services created by test-services.ps1

Write-Host "🧹 Cleaning up LazyService test services..." -ForegroundColor Yellow
Write-Host "====================================`n" -ForegroundColor Yellow

# Stop Python processes
Write-Host "Stopping Python HTTP servers..." -ForegroundColor Cyan
Get-Process python -ErrorAction SilentlyContinue | Where-Object {
    $_.CommandLine -like "*http.server*"
} | Stop-Process -Force
Write-Host "  ✓ Python servers stopped`n" -ForegroundColor Green

# Stop Node.js processes
Write-Host "Stopping Node.js HTTP servers..." -ForegroundColor Cyan
Get-Process node -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Write-Host "  ✓ Node.js servers stopped`n" -ForegroundColor Green

# Stop Docker containers
$dockerAvailable = Get-Command docker -ErrorAction SilentlyContinue
if ($dockerAvailable) {
    Write-Host "Stopping Docker test containers..." -ForegroundColor Cyan
    
    docker stop lazyservice-test-nginx 2>$null
    docker rm lazyservice-test-nginx 2>$null
    
    docker stop lazyservice-test-redis 2>$null
    docker rm lazyservice-test-redis 2>$null
    
    Write-Host "  ✓ Docker containers stopped and removed`n" -ForegroundColor Green
}

Write-Host "====================================`n" -ForegroundColor Yellow
Write-Host "✅ Cleanup complete!`n" -ForegroundColor Green
