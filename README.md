# 📟 lametric-cli - control LaMetric from your terminal

CLI tool for LaMetric TIME/SKY devices. Control your device from the command line.

## Agent Skills

Use with AI coding assistants:

```bash
npx skills add dedene/lametric-cli
```

Works with Claude Code, Cursor, Codex, and [35+ agents](https://github.com/vercel-labs/skills#supported-agents).

## Installation

### Homebrew (macOS/Linux)

```bash
brew install dedene/tap/lametric
```

### Go Install

```bash
go install github.com/dedene/lametric-cli/cmd/lametric@latest
```

### Binary Downloads

Download from [Releases](https://github.com/dedene/lametric-cli/releases).

## Setup

1. Get your API key from the LaMetric mobile app (Device Settings > API Key)
2. Run the setup wizard:

```bash
lametric setup
```

Or configure manually:

```bash
# Set API key
lametric auth set-key --device=living-room

# Or use environment variable
export LAMETRIC_API_KEY=your-api-key
export LAMETRIC_DEVICE=192.168.1.100
```

## Usage

### Notifications

```bash
# Simple notification
lametric notify "Hello World"

# With icon and sound
lametric notify "Build passed" --icon=checkmark --sound=positive1

# Critical alert (wakes device)
lametric notify "ALERT" --priority=critical --sound=alarm1

# Progress bar
lametric notify "Downloads" --goal=75/100 --icon=download

# Sparkline chart
lametric notify "CPU" --chart=10,25,50,30,45,80,60

# From stdin
echo "Pipeline complete" | lametric notify

# Wait for dismissal
lametric notify "Confirm?" --wait

# JSON output
lametric notify "Test" -j
```

### Device Control

```bash
# Device info
lametric device

# Display settings
lametric display get
lametric display brightness 50
lametric display mode auto

# Audio
lametric audio get
lametric audio volume 30

# Bluetooth
lametric bluetooth get
lametric bluetooth on
lametric bluetooth off
```

### App Control

```bash
# List apps
lametric app list

# Switch apps
lametric app next
lametric app prev

# Radio
lametric app radio play
lametric app radio stop
lametric app radio next

# Timer
lametric app timer set 5m
lametric app timer start
lametric app timer pause

# Stopwatch
lametric app stopwatch start
lametric app stopwatch reset
```

### Streaming

```bash
# Start streaming session
lametric stream start

# Send image
lametric stream image logo.png

# Send animated GIF
lametric stream gif animation.gif

# Pipe from ffmpeg
ffmpeg -i video.mp4 -vf "scale=37:8" -f rawvideo -pix_fmt rgb24 - | lametric stream pipe

# Stop streaming
lametric stream stop
```

### Discovery

```bash
# Find devices on network
lametric discover

# With timeout
lametric discover --timeout=5s
```

### Helpers

```bash
# List available sounds
lametric sounds

# List icon aliases
lametric icons
```

## Global Flags

| Flag            | Description                                |
| --------------- | ------------------------------------------ |
| `-d, --device`  | Device name or IP (env: `LAMETRIC_DEVICE`) |
| `-j, --json`    | Output JSON                                |
| `--plain`       | Output plain TSV (for scripting)           |
| `--no-color`    | Disable colors (env: `NO_COLOR`)           |
| `-v, --verbose` | Verbose logging                            |

## Configuration

Config file: `~/.config/lametric-cli/config.yaml`

```yaml
default_device: living-room

devices:
  living-room:
    ip: 192.168.1.100
  bedroom:
    ip: 192.168.1.101
```

API keys are stored securely in your system keyring.

## Environment Variables

| Variable           | Description                 |
| ------------------ | --------------------------- |
| `LAMETRIC_API_KEY` | API key (overrides keyring) |
| `LAMETRIC_DEVICE`  | Default device IP or alias  |
| `NO_COLOR`         | Disable colored output      |

## License

MIT
