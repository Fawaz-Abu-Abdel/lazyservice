# AI Agent Prompt: Advanced Service Monitoring TUI Application

## Project Overview
Create a production-ready, cross-platform Terminal User Interface (TUI) application for unified service monitoring and management. The application should detect, monitor, and provide actionable insights for services across multiple platforms (Docker, Kubernetes, Systemd, bare processes) with real-time metrics, historical analytics, and interactive controls.

## Core Architecture

### 1. **Modular Plugin-Based Design**
```
project/
├── cmd/
│   └── app/
│       └── main.go                 # Entry point with graceful startup
├── internal/
│   ├── app/
│   │   ├── app.go                 # Core application orchestrator
│   │   ├── models.go              # Unified data models
│   │   └── registry.go            # Plugin registry pattern
│   ├── detectors/                 # Service discovery plugins
│   │   ├── interface.go           # Common detector interface
│   │   ├── docker.go              # Docker container detection
│   │   ├── kubernetes.go          # K8s pod detection
│   │   ├── systemd.go             # Systemd service detection
│   │   ├── process.go             # Bare process detection
│   │   ├── docker_compose.go      # Docker Compose projects
│   │   └── pm2.go                 # PM2 process manager
│   ├── metrics/                   # Metrics collection plugins
│   │   ├── interface.go           # Common collector interface
│   │   ├── docker_stats.go        # Docker metrics via API
│   │   ├── k8s_metrics.go         # K8s metrics server
│   │   ├── process_stats.go       # Process metrics via gopsutil
│   │   └── aggregator.go          # Metrics aggregation & analysis
│   ├── ui/
│   │   ├── dashboard.go           # Main TUI dashboard
│   │   ├── theme.go               # Customizable color themes
│   │   ├── components/            # Reusable UI components
│   │   │   ├── chart.go           # ASCII charts & sparklines
│   │   │   ├── table.go           # Enhanced data tables
│   │   │   ├── graph.go           # Dependency graphs
│   │   │   └── modal.go           # Modal dialogs
│   │   ├── views/                 # Different view modes
│   │   │   ├── list.go            # Service list view
│   │   │   ├── details.go         # Detailed service view
│   │   │   ├── logs.go            # Log viewer with search
│   │   │   ├── statistics.go      # Analytics dashboard
│   │   │   └── topology.go        # Service topology map
│   │   └── keybindings.go         # Customizable keybindings
│   ├── actions/
│   │   ├── restart.go             # Service restart logic
│   │   ├── logs.go                # Log streaming
│   │   ├── shell.go               # Interactive shell access
│   │   ├── scale.go               # Scale services (K8s/Docker)
│   │   └── export.go              # Export metrics/logs
│   ├── config/
│   │   ├── config.go              # Configuration management
│   │   └── defaults.go            # Sensible defaults
│   ├── storage/
│   │   ├── metrics_db.go          # Persistent metrics storage
│   │   └── cache.go               # In-memory caching
│   ├── alerts/
│   │   ├── rules.go               # Alert rules engine
│   │   ├── notifier.go            # Notification system
│   │   └── thresholds.go          # Threshold monitoring
│   └── logger/
│       └── logger.go              # Structured logging
├── pkg/                           # Public APIs (if needed)
├── configs/
│   ├── default.yaml               # Default configuration
│   └── themes/                    # Theme definitions
├── .github/
│   └── workflows/
│       ├── ci.yml                 # Automated testing
│       └── release.yml            # Automated releases
├── Dockerfile                     # Containerized deployment
├── docker-compose.yml             # Development environment
├── Makefile                       # Build automation
├── go.mod                         # Go dependencies
└── README.md                      # Comprehensive documentation
```

## Technical Requirements

### **Programming Language & Core Libraries**
- **Language**: Go 1.21+
- **TUI Framework**: `github.com/rivo/tview` (terminal UI) OR `github.com/charmbracelet/bubbletea` (modern alternative)
- **Docker SDK**: `github.com/docker/docker`
- **Kubernetes Client**: `k8s.io/client-go`
- **Process Monitoring**: `github.com/shirou/gopsutil/v3`
- **Charts/Visualization**: Custom ASCII art generators or `github.com/guptarohit/asciigraph`
- **Configuration**: `github.com/spf13/viper`
- **Logging**: `github.com/sirupsen/logrus` or `go.uber.org/zap`
- **Storage**: `github.com/dgraph-io/badger` (embedded DB for metrics history)

## Feature Requirements

### 🔍 **Enhanced Service Discovery**
1. **Multi-Platform Detection**:
   - Docker containers (local + remote daemons)
   - Kubernetes pods (all namespaces with filtering)
   - Docker Compose projects (grouped view)
   - Systemd services (Linux)
   - Bare processes (Node.js, Python, Java, Ruby, PHP, Go, Rust)
   - PM2 managed processes
   - Windows Services (on Windows)

2. **Smart Grouping & Filtering**:
   - Group by: platform, namespace, project, labels/tags
   - Filter by: status, resource usage thresholds, custom queries
   - Search: fuzzy search across service names and metadata
   - Tags: Custom user-defined tags for services

3. **Auto-Discovery Configuration**:
   - Configurable refresh intervals (1s - 60s)
   - Include/exclude patterns
   - Multi-host support (remote Docker/K8s clusters)

### 📊 **Comprehensive Metrics Collection**

1. **Real-Time Metrics**:
   - CPU: percentage, cores, throttling
   - Memory: usage, limit, swap, percentage, cache
   - Network: in/out bytes, packets, errors, connections
   - Disk I/O: read/write bytes, IOPS
   - Process: threads, file descriptors, open ports
   - Application: request rate, error rate, latency (p50, p95, p99)

2. **Historical Analytics**:
   - Store last 24 hours of metrics (configurable)
   - Statistical analysis: min, max, avg, percentiles
   - Trend detection (increasing/decreasing/stable)
   - Anomaly detection (spikes, drops)
   - Correlations between services

3. **Custom Metrics**:
   - HTTP endpoint polling (health checks, metrics endpoints)
   - Log parsing for custom metrics
   - Plugin system for custom collectors

### 🎨 **Beautiful & Intuitive TUI**

1. **Multiple View Modes**:
   - **List View**: Compact table of all services
   - **Detail View**: In-depth single service inspection
   - **Dashboard View**: Multi-panel overview with sparklines
   - **Logs View**: Real-time log streaming with search/filter
   - **Statistics View**: Charts and graphs of historical data
   - **Topology View**: Visual service dependency map
   - **Compare View**: Side-by-side service comparison

2. **Visual Elements**:
   - Color-coded status indicators (green/yellow/red with customizable thresholds)
   - ASCII sparklines for historical data (CPU, memory, network)
   - Progress bars for resource usage
   - Live-updating counters
   - Sortable and resizable columns
   - Smooth animations for transitions

3. **Themes**:
   - Built-in themes: Dark, Light, Dracula, Nord, Solarized
   - Custom theme support via YAML configuration
   - Adaptive colors based on terminal capabilities

4. **Layout**:
   - Responsive: Adapts to terminal size
   - Customizable panels (show/hide, resize, reorder)
   - Split views (horizontal/vertical)
   - Focus modes (fullscreen specific panels)

### 🚀 **Interactive Actions**

1. **Service Management**:
   - **Start/Stop/Restart**: Services with confirmation prompts
   - **Scale**: Adjust replicas (K8s/Docker Swarm)
   - **Shell Access**: Launch interactive shell in container/pod
   - **Logs**: Stream logs with tail, search, and export
   - **Inspect**: View raw configuration/metadata
   - **Delete**: Remove services with safety checks
   - **Edit**: Modify service configuration (opens $EDITOR)

2. **Batch Operations**:
   - Multi-select services (checkbox or tag-based)
   - Bulk restart/stop/start
   - Batch export logs or metrics

3. **Monitoring Actions**:
   - Pin services (always show at top)
   - Bookmark frequently accessed services
   - Create custom views (saved filter combinations)
   - Set alerts on specific services

### 📈 **Advanced Analytics & Reporting**

1. **Statistics Dashboard**:
   - Service uptime tracking
   - Resource utilization trends
   - Service health scores (composite metric)
   - Cost estimation (based on resource usage)
   - Comparison: current vs historical performance

2. **Reports**:
   - Generate reports (HTML, JSON, CSV)
   - Scheduled reports (daily/weekly summaries)
   - Export metrics to external systems (Prometheus format)

3. **Alerting System**:
   - Define alert rules (CPU > 80%, memory > 90%, service down)
   - Alert channels: desktop notifications, webhooks, email
   - Alert history and acknowledgment
   - Snooze/mute specific alerts

### ⚙️ **Configuration & Extensibility**

1. **Configuration File** (YAML):
```yaml
refresh_interval: 5s
detectors:
  docker:
    enabled: true
    remote_hosts: ["tcp://remote:2376"]
  kubernetes:
    enabled: true
    contexts: ["prod", "staging"]
  process:
    enabled: true
    patterns: ["node", "python", "java"]

metrics:
  history_retention: 24h
  storage_path: "~/.lazyservice/metrics.db"

alerts:
  - name: "High CPU"
    condition: "cpu > 80"
    duration: "5m"
    notify: ["desktop", "webhook"]

ui:
  theme: "dracula"
  refresh_rate: "1s"
  panels: ["services", "details", "metrics"]

keybindings:
  quit: "q"
  refresh: "r"
  restart_service: "ctrl+r"
```

2. **Plugin System**:
   - Interface-based plugins for detectors and collectors
   - Dynamic loading of external plugins
   - Plugin marketplace/registry (future)

3. **API/CLI Mode**:
   - Headless mode for scripting
   - JSON output for automation
   - RESTful API for external integrations

### 🔐 **Security & Reliability**

1. **Permissions**:
   - Read-only mode (disable destructive actions)
   - Configurable action permissions
   - Audit log for all actions

2. **Error Handling**:
   - Graceful degradation (partial failures don't crash app)
   - Detailed error messages with recovery suggestions
   - Automatic retry for transient failures

3. **Performance**:
   - Efficient polling (adaptive refresh rates)
   - Connection pooling and caching
   - Minimal resource footprint (<50MB memory)
   - Background goroutines for non-blocking operations

### 📚 **Documentation & Quality**

1. **Code Quality**:
   - Unit tests (>80% coverage)
   - Integration tests for detectors
   - Benchmarks for performance-critical code
   - Linting (golangci-lint)

2. **Documentation**:
   - Comprehensive README with quickstart
   - Architecture decision records (ADR)
   - API documentation (GoDoc)
   - User guide with screenshots/GIFs
   - Troubleshooting guide
   - Contributing guidelines

3. **CI/CD**:
   - Automated testing on push
   - Multi-platform builds (Linux, macOS, Windows)
   - Automated releases with GoReleaser
   - Docker images published to registry

## Enhanced Features (Beyond Original)

### 🌟 **New Advanced Capabilities**

1. **Service Dependency Graph**:
   - Auto-detect dependencies (network connections, env vars)
   - Visual topology map (ASCII art or interactive)
   - Impact analysis (what breaks if service X stops)

2. **Time-Series Database Integration**:
   - Export to Prometheus, InfluxDB, Grafana
   - Historical query language
   - Metric retention policies

3. **AI-Powered Insights**:
   - Anomaly detection using statistical models
   - Predictive alerts (service likely to fail soon)
   - Resource optimization recommendations
   - Pattern recognition (similar incidents in past)

4. **Collaborative Features**:
   - Share service views (export/import configurations)
   - Team dashboards (remote sync)
   - Incident tracking integration (Jira, GitHub Issues)

5. **Advanced Log Management**:
   - Log aggregation across multiple services
   - Real-time log search (regex, full-text)
   - Log parsing and structured extraction
   - Log-based metrics and alerts

6. **Performance Profiling**:
   - Flame graphs for CPU profiling
   - Memory leak detection
   - Network latency tracing
   - Distributed tracing integration (Jaeger, Zipkin)

7. **Automation & Scripting**:
   - Record and replay action sequences
   - Scheduled tasks (e.g., restart every night)
   - Conditional actions (if CPU > 90%, restart)
   - Custom scripts triggered by events

8. **Multi-User Support**:
   - User profiles with preferences
   - Role-based access control
   - Shared configurations across team

9. **Cloud Provider Integration**:
   - AWS ECS/EKS service detection
   - Azure Container Instances
   - Google Cloud Run
   - Provider-specific metrics (CloudWatch, Stackdriver)

10. **Mobile Companion**:
    - Web-based dashboard (future)
    - Mobile app for alerts and basic actions
    - QR code for quick access

## Implementation Guidelines

### **Phase 1: Core Foundation (Week 1-2)**
- [ ] Project setup with modular architecture
- [ ] Core app orchestrator with plugin registry
- [ ] Basic Docker detector and metrics collector
- [ ] Simple TUI with service list view
- [ ] Configuration management

### **Phase 2: Multi-Platform Detection (Week 2-3)**
- [ ] Kubernetes detector
- [ ] Process detector (bare metal)
- [ ] Systemd detector
- [ ] Docker Compose project grouping
- [ ] Unified service model

### **Phase 3: Enhanced UI (Week 3-4)**
- [ ] Dashboard view with sparklines
- [ ] Detail view with comprehensive info
- [ ] Theme system
- [ ] Multiple layout modes
- [ ] Keyboard shortcuts

### **Phase 4: Actions & Interactivity (Week 4-5)**
- [ ] Service restart/stop/start
- [ ] Log viewer with streaming
- [ ] Shell access (docker exec, kubectl exec)
- [ ] Batch operations
- [ ] Safety confirmations

### **Phase 5: Analytics & Storage (Week 5-6)**
- [ ] Metrics history storage (BadgerDB)
- [ ] Statistics dashboard
- [ ] Charts and trend analysis
- [ ] Export functionality

### **Phase 6: Advanced Features (Week 6-8)**
- [ ] Alert system with rules engine
- [ ] Dependency graph visualization
- [ ] Advanced filtering and search
- [ ] Performance optimizations
- [ ] Plugin system

### **Phase 7: Polish & Release (Week 8-9)**
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] CI/CD pipeline
- [ ] Cross-platform builds
- [ ] Release v1.0.0

## Code Quality Standards

1. **Clean Code Principles**:
   - Single Responsibility Principle for all modules
   - Interface-driven design (dependency injection)
   - Minimal cyclomatic complexity (<10 per function)
   - Descriptive naming (no abbreviations)

2. **Error Handling**:
   - Always return errors, never panic (except fatal startup errors)
   - Wrap errors with context using `fmt.Errorf("%w")`
   - Log errors at appropriate levels
   - Provide user-friendly error messages in UI

3. **Concurrency**:
   - Use context for cancellation and timeouts
   - Proper mutex usage for shared state
   - Channel-based communication preferred
   - Worker pools for parallel operations

4. **Testing**:
   - Table-driven tests
   - Mocks for external dependencies
   - Integration tests for critical paths
   - Benchmarks for performance-sensitive code

## Success Criteria

The application should:
1. ✅ Run on Linux, macOS, and Windows without modification
2. ✅ Detect and monitor 100+ services without performance degradation
3. ✅ Update UI smoothly at 1-second intervals
4. ✅ Use <50MB memory in typical usage
5. ✅ Handle network failures gracefully (retry + fallback)
6. ✅ Have <500ms startup time
7. ✅ Be usable with keyboard only (full accessibility)
8. ✅ Have comprehensive documentation
9. ✅ Be production-ready with proper error handling
10. ✅ Be extensible for future enhancements

## Bonus Enhancements

- 🎁 Web UI companion (optional HTTP server mode)
- 🎁 Terminal recordings (record TUI sessions)
- 🎁 Export to Grafana dashboards
- 🎁 Integration with popular DevOps tools (Slack, PagerDuty)
- 🎁 Container security scanning (CVE detection)
- 🎁 Cost optimization recommendations
- 🎁 Service SLA tracking and reporting
- 🎁 Chaos engineering mode (simulate failures)

---

## Final Notes for AI Agent

**Approach**:
- Start with a working MVP (Phases 1-3) before adding advanced features
- Prioritize reliability and performance over feature count
- Follow Go best practices and idioms consistently
- Write self-documenting code with clear comments for complex logic
- Create a delightful user experience (smooth animations, helpful messages)
- Make it easy to extend (future plugins, themes, integrations)

**Design Philosophy**:
- "It just works" - sensible defaults, minimal configuration
- "Beautiful and fast" - smooth TUI, instant feedback
- "Production-ready" - proper error handling, logging, testing
- "Developer-friendly" - clean code, good documentation

**Inspiration**:
- Take inspiration from: `lazydocker`, `k9s`, `htop`, `bottom`, `glances`
- But create something unique and more comprehensive
- Modern UI/UX principles applied to terminal applications

**Remember**: This tool will be used daily by developers and ops engineers. Make it reliable, fast, and delightful to use. Every detail matters!
