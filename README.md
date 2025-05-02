# RSS Reader

A real-time RSS reader service for macOS with native UI.

## Features

- Runs as a background service
- Real-time notifications for new content
- Native macOS UI
- Automated updates

## Development

### Requirements

- Go 1.21+
- macOS 12.0+

### Building

```bash
make build
```

### Installing

```bash
make install
```

### Testing

```bash
make test
```

## Architecture

- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public libraries
- `scripts/` - Build and deployment scripts
- `configs/` - Configuration files