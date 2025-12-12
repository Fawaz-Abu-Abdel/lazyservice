# Minimal Dashboard - Enhanced Edition

## Overview
**Professional-grade dashboard** with **values-only updates** - just like `lazytime`, but with all the features of Pro Dashboard!

```
╔═══════════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║  ► ⚡ LazyService Pro - Minimal Dashboard ◄ │ 📊 12 Services ◄ │ 🕐 23:45:12 Thu Oct 31                     ║
╚═══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝

╔════════════════════ 🚀 Services ════════════════════╗ ╔═══════════ 📋 Service Details ════════════╗ ╔═══ 📈 Stats ════╗
║  ❯ ● 🐳 nginx-web-server                            ║ ║ Name: nginx-web-server                    ║ ║ Total:          ║
║    ● 🐳 redis-cache                                 ║ ║ Type: 🐳 docker                           ║ ║  12             ║
║    ● ⚡ node-api                                    ║ ║ Status: ● running                         ║ ║                 ║
║    ○ ⚙️ postgres-db                                ║ ║ ID: abc123def456...                       ║ ║ ● Running:      ║
╚════════════════════════════════════════════════════════╝ ║ Uptime: 2d 5h 23m                         ║ ║  8              ║
                                                         ║ Created: Oct 29, 18:20                    ║ ║                 ║
                                                         ║ Ports: 80, 443                            ║ ║ ● Stopped:      ║
                                                         ║ Image: nginx:latest                       ║ ║  4              ║
                                                         ╚═══════════════════════════════════════════════╝ ╚════════════════════╝
                                                         ╔═══════════ 📊 Live Metrics ═══════════════╗
                                                         ║ CPU  ████████████░░░░░░░░  45.2%         ║
                                                         ║ MEM  ██████████████░░░░░░  62.8%         ║
                                                         ║ NET↓ 1.2 MB/s                             ║
                                                         ║ NET↑ 450 KB/s                             ║
                                                         ║ DISK R:15 MB                              ║
                                                         ╚═══════════════════════════════════════════════╝
═══════════════════════════════════════════════════════════════════════════════════════════════════════════════
↑↓/jk Navigate │ R Refresh │ Q Quit
Last update: 23:45:12 │ ✨ Values-Only Updates (Borders Never Redraw)
```

## What Makes It Enhanced?

### 🎨 **Professional 3-Panel Layout**
- **Services List** (left) - 17 slots with icons and status
- **Details + Metrics** (center) - Full service info + live metrics
- **Statistics** (right) - Quick overview of service counts

### 📊 **All Pro Features**
✅ Service type icons (🐳 Docker, ☸️ K8s, ⚡ Process, ⚙️ Systemd)  
✅ Live metrics with progress bars  
✅ Network I/O and Disk stats  
✅ Uptime and creation time  
✅ Dynamic header with clock  
✅ Live statistics panel  
✅ Footer with last update time  

### ⚡ **Values-Only Updates** (Like lazytime)
✅ **Borders drawn ONCE** at startup  
✅ **Values updated** via cursor positioning  
✅ **Zero flicker** - no full screen redraws  
✅ **2-second auto-refresh** in background  
✅ **Manual refresh** with 'R' key  

## Technical Implementation

### ANSI Cursor Positioning
```go
// Draw borders ONCE
func drawInitialUI() {
    clearScreen()
    // Draw all panels with borders
    // Never called again!
}

// Update values ONLY
func updateMetrics() {
    moveCursor(18, 60)  // CPU line
    fmt.Print("CPU  " + cpuBar + " 45.2%")
    
    moveCursor(19, 60)  // Memory line  
    fmt.Print("MEM  " + memBar + " 62.8%")
    // Borders untouched!
}
```

### Update Flow
```
Background Loop (2s):
  ├─ updateHeader()      → Time + service count
  ├─ updateServiceList() → Service names/status
  ├─ updateDetails()     → Selected service info
  ├─ updateMetrics()     → CPU, MEM, NET, DISK
  ├─ updateStats()       → Total/Running/Stopped
  └─ updateFooter()      → Last update timestamp
```

### Keyboard Actions
```
k/K → Navigate up    (updates: list + details + metrics)
j/J → Navigate down  (updates: list + details + metrics)
R   → Manual refresh (updates: ALL values)
Q   → Quit
```

## Usage

### Build
```powershell
go build -o lazyservice.exe ./cmd/lazyservice
```

### Run Minimal Dashboard
```powershell
.\lazyservice.exe -minimal
```

### Run Pro Dashboard (default)
```powershell
.\lazyservice.exe
```

## Architecture

```
┌──────────────────────────────────────────────────┐
│         Minimal Dashboard Enhanced               │
├──────────────────────────────────────────────────┤
│  1. enableVirtualTerminal() - ANSI on Windows    │
│  2. drawInitialUI() - borders ONCE               │
│  3. updateLoop() - values every 2s               │
│  4. Cursor positioning for updates               │
│  5. No border/layout redraws EVER                │
└──────────────────────────────────────────────────┘
```

### Panel Positions (Hardcoded)
- **Header**: Row 2
- **Services**: Rows 6-22, Cols 1-56
- **Details**: Rows 7-14, Cols 60-101
- **Metrics**: Rows 18-22, Cols 60-101
- **Stats**: Rows 7-22, Cols 105-119
- **Footer**: Rows 25-26

Each value update moves cursor to exact position and prints new value.

## Comparison: Enhanced vs Original

| Feature | Original Minimal | Enhanced Minimal |
|---------|-----------------|------------------|
| Panels | 2 (Services, Details) | 4 (Services, Details, Metrics, Stats) |
| Service Slots | 15 | 17 |
| Type Icons | ❌ | ✅ 🐳☸️⚡⚙️ |
| Statistics | ❌ | ✅ Total/Running/Stopped |
| Dynamic Header | ❌ | ✅ Time + Count |
| Footer Status | Static | ✅ Last update time |
| Disk I/O | ❌ | ✅ |
| Layout Width | 87 chars | 120 chars |
| Update Frequency | 2s | 2s |
| Values-Only | ✅ | ✅ |

## Benefits of Enhanced Design

✅ **All Pro features** - Nothing sacrificed  
✅ **Zero flicker** - Borders never redraw  
✅ **Fast updates** - Only changed values  
✅ **Professional look** - Icons, colors, layout  
✅ **Live clock** - Header updates every 2s  
✅ **Statistics** - Quick service overview  
✅ **Better spacing** - 120-char wide layout  

## Performance

### Memory Footprint
- **Initial draw**: ~120 lines × 120 chars = 14KB
- **Per update**: 20-50 values × average 30 chars = ~1.5KB

### CPU Usage
- **Background loop**: Minimal (cursor moves + prints)
- **Keyboard**: Instant response (3-5 value updates)
- **Network calls**: Same as Pro (app.GetServices)

### Flicker
- **0%** - Borders are drawn once and never touched

## When to Use Enhanced Minimal

### Perfect for:
✅ **Production monitoring** - Clean, flicker-free display  
✅ **Remote terminals** - Efficient bandwidth usage  
✅ **Legacy consoles** - Works on any ANSI terminal  
✅ **Always-on displays** - Zero burn-in risk  
✅ **Windows 10+** - Native ANSI support  

### Use Pro Dashboard instead if:
❌ Need mouse support  
❌ Want dynamic resizing  
❌ Prefer TUI framework abstractions  

## Code Example

```go
// Updating CPU metric (values only!)
func (d *MinimalDashboard) updateMetrics() {
    metrics := service.Metrics
    
    // Move to CPU line, update value
    d.moveCursor(18, 60)
    cpuBar := d.makeBar(metrics.CPUPercent, 20)
    fmt.Print("CPU  " + cpuBar + fmt.Sprintf(" %5.1f%%", metrics.CPUPercent))
    
    // Move to MEM line, update value
    d.moveCursor(19, 60)
    memBar := d.makeBar(metrics.MemoryPercent, 20)
    fmt.Print("MEM  " + memBar + fmt.Sprintf(" %5.1f%%", metrics.MemoryPercent))
    
    // Borders? Untouched! ✨
}
```

## Inspiration

This enhanced design takes the best of both worlds:
- **lazytime** - Values-only update technique
- **Pro Dashboard** - Professional features and layout

Result: **Production-grade monitoring** with **zero flicker**! 🚀

---

*"Draw once, update forever"* - The Minimal Dashboard philosophy
