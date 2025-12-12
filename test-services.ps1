# LazyService Test Script
# Starts some test services so you can see LazyService in action

Write-Host "🚀 LazyService Test Script" -ForegroundColor Green
Write-Host "====================================`n" -ForegroundColor Green

# Check if Python is available
$pythonAvailable = Get-Command python -ErrorAction SilentlyContinue
if ($pythonAvailable) {
    Write-Host "✓ Starting Python HTTP server on port 8000..." -ForegroundColor Yellow
    Start-Process python -ArgumentList "-m", "http.server", "8000" -WindowStyle Minimized
    Write-Host "  Python server started!" -ForegroundColor Green
} else {
    Write-Host "⚠  Python not found - skipping Python test server" -ForegroundColor DarkYellow
}

# Check if Node.js is available
$nodeAvailable = Get-Command node -ErrorAction SilentlyContinue
if ($nodeAvailable) {
    Write-Host "✓ Starting Node.js HTTP server on port 3000..." -ForegroundColor Yellow
    $nodeScript = "require('http').createServer((req,res)=>res.end('Hello from LazyService test!')).listen(3000,()=>console.log('Server running on port 3000'))"
    Start-Process node -ArgumentList "-e", $nodeScript -WindowStyle Minimized
    Write-Host "  Node.js server started!" -ForegroundColor Green
} else {
    Write-Host "⚠  Node.js not found - skipping Node.js test server" -ForegroundColor DarkYellow
}

# Check if Docker is available
$dockerAvailable = Get-Command docker -ErrorAction SilentlyContinue
if ($dockerAvailable) {
    Write-Host "✓ Checking Docker status..." -ForegroundColor Yellow
    $dockerRunning = docker ps 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  Docker is running!" -ForegroundColor Green
        
        # Start test containers
        Write-Host "`n✓ Starting test Docker containers..." -ForegroundColor Yellow
        
        docker run -d --name lazyservice-test-nginx -p 8080:80 nginx 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ nginx container started on port 8080" -ForegroundColor Green
        }
        
        docker run -d --name lazyservice-test-redis -p 6379:6379 redis 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ redis container started on port 6379" -ForegroundColor Green
        }
    } else {
        Write-Host "  ⚠  Docker is installed but not running" -ForegroundColor DarkYellow
        Write-Host "     Start Docker Desktop to see container metrics" -ForegroundColor DarkYellow
    }
} else {
    Write-Host "⚠  Docker not found - skipping Docker containers" -ForegroundColor DarkYellow
}

Write-Host "`n====================================`n" -ForegroundColor Green
Write-Host "✅ Test services started!" -ForegroundColor Green
Write-Host "`nNow run LazyService to see them:" -ForegroundColor Cyan
Write-Host "  .\lazyservice.exe`n" -ForegroundColor White

Write-Host "To stop test services later, run:" -ForegroundColor Yellow
Write-Host "  .\cleanup-test-services.ps1`n" -ForegroundColor White

Write-Host "Press any key to launch LazyService..." -ForegroundColor Cyan
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# Launch LazyService
if (Test-Path ".\lazyservice.exe") {
    .\lazyservice.exe
} else {
    Write-Host "`n⚠  lazyservice.exe not found!" -ForegroundColor Red
    Write-Host "Build it first with: go build -o lazyservice.exe ./cmd/lazyservice`n" -ForegroundColor Yellow
}
