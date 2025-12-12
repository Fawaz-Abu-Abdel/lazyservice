# 🎯 LazyService

A powerful TUI (Terminal User Interface) dashboard that provides unified monitoring of all your services across different platforms - Docker, Kubernetes, Systemd, and bare processes - with real-time metrics, statistics, and request analytics.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## ✨ Features

### 🔍 Unified Service Discovery
- **Multi-platform Support**: Automatically detects services from:
  - 🐳 Docker containers
  - ☸️ Kubernetes pods
  - ⚙️ Systemd services
  - 💻 Bare processes (Node.js, Python, Java, etc.)
- **Auto-refresh**: Service list updates every 5 seconds
- **Smart Grouping**: Services organized by type and namespace

### 📊 Real-time Metrics
- **CPU Usage**: Percentage and core utilization
- **Memory**: Usage, limits, and percentage
- **Network**: Inbound/outbound traffic
- **Disk I/O**: Read/write operations
- **Uptime**: Service runtime tracking
- **Historical Data**: Visual sparklines showing last 5 minutes of metrics

### 🎨 Beautiful TUI
- Clean, intuitive interface built with `tview`
- Color-coded service status indicators
- Real-time updating dashboard
- Detailed service inspection view
- ASCII charts and sparklines

### 🚀 Quick Actions
- **[L]ogs**: View service logs (coming soon)
- **[R]estart**: Restart services (coming soon)
- **[S]hell**: Access service shell (coming soon)
- **[E]dit**: Edit service configuration (coming soon)
- **[D]elete**: Remove services (coming soon)
- **[Q]uit**: Exit application

## 📋 Prerequisites

- **Go 1.21+** for building
- **Docker** (optional): For Docker container monitoring
- **Kubernetes** (optional): For K8s pod monitoring with kubectl configured
- **Systemd** (optional): For systemd service monitoring (Linux only)

## 🛠️ Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/lazyservice.git
cd lazyservice

# Install dependencies
go mod download

# Build the application
go build -o lazyservice ./cmd/lazyservice

# Run
./lazyservice
```

### Quick Install

```bash
# Install directly
go install github.com/yourusername/lazyservice/cmd/lazyservice@latest
```

## 🚀 Usage

### Basic Usage

Simply run the binary:

```bash
lazyservice
```

The dashboard will automatically detect and display all available services.

### Keyboard Shortcuts

- **↑/↓**: Navigate service list
- **Enter**: View detailed service information
- **r**: Manually refresh service list
- **q**: Quit application

## 📁 Project Structure

```
lazyservice/
├── cmd/
│   └── lazyservice/
│       └── main.go              # Entry point
├── internal/
│   ├── app/
│   │   ├── app.go              # Core application logic
│   │   └── models.go           # Data models
│   ├── ui/
│   │   ├── dashboard.go        # Main dashboard UI
│   │   └── components/         # Reusable UI components
│   │       └── chart.go        # Chart utilities
│   ├── detectors/
│   │   ├── docker.go           # Docker service detection
│   │   ├── kubernetes.go       # K8s service detection
│   │   ├── systemd.go          # Systemd service detection
│   │   └── process.go          # Process detection
│   └── metrics/
│       ├── docker_stats.go     # Docker metrics collection
│       ├── k8s_metrics.go      # Kubernetes metrics
│       └── process_stats.go    # Process metrics
├── go.mod
└── README.md
```

## 🔧 Configuration

LazyService works out of the box with sensible defaults. However, you can customize process detection by modifying the patterns in `detectors/process.go`.

### Docker

Ensure Docker daemon is running and your user has permission to access the Docker socket.

### Kubernetes

LazyService will automatically use your kubectl configuration:
- In-cluster config (if running inside a pod)
- `~/.kube/config` (for local development)

For metrics, ensure the Kubernetes Metrics Server is installed:
```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

### Systemd

Requires `systemctl` to be available (Linux only). Works automatically with proper permissions.

## 📊 Metrics Collection

LazyService collects the following metrics per platform:

| Metric | Docker | Kubernetes | Systemd | Process |
|--------|--------|------------|---------|---------|
| CPU % | ✅ | ✅ | ⚠️ | ✅ |
| Memory | ✅ | ✅ | ⚠️ | ✅ |
| Network | ✅ | ⚠️ | ❌ | ❌ |
| Disk I/O | ✅ | ⚠️ | ❌ | ✅ |
| Uptime | ✅ | ✅ | ⚠️ | ✅ |

- ✅ Full support
- ⚠️ Partial support (depends on platform availability)
- ❌ Not available

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Roadmap

- [ ] Log viewer integration
- [ ] Service restart/stop/start actions
- [ ] Shell access to containers/pods
- [ ] Service configuration editing
- [ ] Export metrics to files
- [ ] Alerting system
- [ ] Custom filters and search
- [ ] Docker Compose project grouping
- [ ] PM2 process manager integration
- [ ] HTTP endpoint monitoring
- [ ] Service dependency graph

## 🐛 Known Issues

- Systemd metrics collection is limited on some platforms
- Kubernetes network metrics require metrics-server
- Windows systemd support is not available (systemd doesn't run on Windows)

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [tview](https://github.com/rivo/tview) - Terminal UI library
- [gopsutil](https://github.com/shirou/gopsutil) - Process utilities
- [Docker SDK](https://github.com/docker/docker) - Docker client
- [Kubernetes client-go](https://github.com/kubernetes/client-go) - Kubernetes client

## 📧 Contact

For questions, suggestions, or issues, please open an issue on GitHub.

---

