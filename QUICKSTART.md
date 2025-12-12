# 🚀 Quick Start Guide

## Build and Run

```powershell
# Build the application
go build -o lazyservice.exe ./cmd/lazyservice

# Run it
./lazyservice.exe
```
## What You'll See

The app will show startup messages:
```
🚀 Starting LazyService...
✓ Docker detector registered
✓ Systemd detector registered
✓ Process detector registered

📊 3 detectors active. Starting dashboard...
```

Then the TUI dashboard will appear showing all detected services!

## Testing Without Docker

If you don't have Docker running, LazyService will still work and show:

### 1. Running Processes

The app automatically detects common service processes like:
- Python, Node.js, Java, Ruby, PHP
- Nginx, Apache, Redis, PostgreSQL, MySQL, MongoDB
- PM2 and other process managers

### 2. Start a Test Service

Open a new terminal and run:

```powershell
# Python HTTP server
python -m http.server 8000

# OR Node.js (if installed)
node -e "require('http').createServer((req,res)=>res.end('ok')).listen(3000)"
```

Switch back to LazyService - you should see it within 5 seconds!

## Keyboard Controls

- **↑/↓**: Navigate service list
- **Enter**: View detailed service info (currently displayed by default)
- **r**: Force refresh services
- **q**: Quit application

## What You'll See Per Service

### Service List (Left Panel)
- Service name with status emoji
- Service type (docker, kubernetes, systemd, process)
- Ports (if applicable)

### Service Details (Top Right)
- Name, Type, Status, ID
- Image and version (containers)
- Namespace (Kubernetes)
- Ports and health checks
- Environment variables
- Warnings (if any)

### Resource Usage (Bottom Right)
- **CPU**: Percentage and visual bar
- **Memory**: Usage, limit, and percentage
- **Network**: Upload/download (Docker only)
- **Disk I/O**: Read/write operations
- **Uptime**: How long the service has been running
- **History**: Sparklines showing last 5 minutes

## Example Output

```
╔════════════ Services ════════════╗
║(0) 🟢 python.exe                 ║
║    process                       ║
║(1) 🔵 nginx-web                  ║
║    docker • port 8080            ║
╚══════════════════════════════════╝
```

## With Docker

If Docker Desktop is running, you'll also see:

```powershell
# Run a test container
docker run -d -p 8080:80 --name test-nginx nginx

# Run another
docker run -d -p 6379:6379 --name test-redis redis
```

LazyService will show these containers with full metrics!

## Advanced

### Customize Process Detection

Edit `internal/detectors/process.go` to add more process names:

```go
namePatterns := []string{
    "node", "python", "java", "ruby", "php",
    "nginx", "apache", "redis", "postgres", "mysql",
    // Add your custom processes here:
    "myapp", "custom-service",
}
```

### Run in Background

```powershell
# Windows
Start-Process .\lazyservice.exe -WindowStyle Hidden
```

## Troubleshooting

See `TROUBLESHOOTING.md` for common issues and solutions.

## Need More Services?

The app auto-refreshes every 5 seconds. Just start new services and they'll appear automatically!

