# Minimal Dashboard - Values-Only Updates

## Overview

The **Minimal Dashboard** uses ANSI escape sequences to update **only values**, never redrawing borders or layout - just like the `lazytime` project.

## How It Works

### 1. **Border Drawn Once**
```go
drawInitialUI()  // Called ONCE at startup
```
- All borders, titles, and layout are drawn once
- Never touched again during runtime

### 2. **Values Updated with Cursor Positioning**
```go
moveCursor(row, col)  // Move to specific position
fmt.Print(value)      // Update ONLY the value
```

Instead of redrawing the entire screen, we:
- Move cursor to exact position
- Print only the new value
- Leave borders/layout untouched

### 3. **ANSI Escape Codes**
- `\033[2J\033[H` - Clear screen
- `\033[5;10H` - Move cursor to row 5, column 10
- `\033[?25l` - Hide cursor
- `\033[92m` - Green color

## Usage

### Run with Minimal Dashboard
```powershell
# Build
go build -o lazyservice.exe ./cmd/lazyservice

# Run minimal dashboard (values-only updates)
.\lazyservice.exe -minimal

# Run pro dashboard (default, with tview)
.\lazyservice.exe
```

### Keyboard Controls
- **k/K** - Navigate up
- **j/J** - Navigate down
- **r/R** - Refresh (values only!)
- **q/Q** - Quit

## Architecture

```
┌─────────────────────────────────────────┐
│          ANSI Escape Codes              │
├─────────────────────────────────────────┤
│  1. Draw borders once (initial)         │
│  2. Track cursor positions              │
│  3. Update values at positions          │
│  4. Background loop (2s interval)       │
└─────────────────────────────────────────┘
```

### Code Structure

```go
type MinimalDashboard struct {
    // Cursor positions for each panel
    serviceListStartRow int
    detailsStartRow     int
    metricsStartRow     int
}

// Draw ONCE
func drawInitialUI() {
    clearScreen()
    // Draw all borders
    // Never called again
}

// Update VALUES only
func updateServiceList() {
    moveCursor(row, col)
    fmt.Print(newValue)
    // No border redraw!
}
```

## Comparison: Pro vs Minimal

| Feature | Pro Dashboard | Minimal Dashboard |
|---------|--------------|-------------------|
| Library | tview (TUI framework) | Direct ANSI codes |
| Borders | Redrawn with tview | Drawn once |
| Values | SetText() triggers redraw | Cursor position + print |
| Mouse | Supported | No |
| Updates | Framework-managed | Manual cursor control |
| Flicker | Possible on some terminals | None |
| Performance | Good (with caching) | Excellent |

## Technical Details

### ANSI Cursor Positioning
```go
// Move to row 10, column 5
fmt.Printf("\033[%d;%dH", 10, 5)

// Update value at that position
fmt.Print("CPU: 45.2%")
```

### Virtual Terminal on Windows
```go
stdout := windows.Handle(os.Stdout.Fd())
var mode uint32
windows.GetConsoleMode(stdout, &mode)
windows.SetConsoleMode(stdout, mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
```

This enables ANSI codes on Windows 10+.

## Benefits

✅ **Zero flicker** - borders never redraw  
✅ **Minimal updates** - only changed values  
✅ **Lightweight** - no TUI framework overhead  
✅ **Precise control** - exact cursor positioning  
✅ **Fast** - direct terminal writes  

## Example: lazytime Comparison

```go
// lazytime approach (TIME ONLY)
func updateTime(t time.Time) {
    fmt.Print("\033[5;1H")  // Move to time line
    fmt.Print("Time: " + t.Format("15:04:05"))
}

// lazyservice minimal (ALL VALUES)
func updateMetrics() {
    fmt.Print("\033[17;42H")  // Move to CPU line
    fmt.Print("CPU: " + cpuValue)
    
    fmt.Print("\033[18;42H")  // Move to MEM line
    fmt.Print("MEM: " + memValue)
}
```

Same concept - cursor positioning for value updates!

## When to Use

### Use Minimal Dashboard when:
- ❌ Running on terminals with flicker issues
- ✅ Need absolute minimal updates
- ✅ Want maximum performance
- ✅ Don't need mouse support

### Use Pro Dashboard when:
- ✅ Want rich TUI features
- ✅ Need mouse support
- ✅ Prefer framework-managed rendering
- ✅ Want hover effects

## Implementation Reference

Based on the excellent `lazytime` project pattern:
1. Draw static UI once
2. Track positions
3. Update values with cursor moves
4. Background refresh loop

This is the **classic terminal approach** - fast, efficient, and flicker-free!
