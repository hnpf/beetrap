# Beetrap
Beetrap is a lightweight honeypot that mimics common network services (SSH, FTP, HTTP) and logs connection attempts in a real-time TUI.

## Features
- mimics SSH (port 2222), FTP (port 2121), and HTTP (port 8080) with realistic service banners
- captures first-line payloads from connecting clients
- real-time TUI with per-service connection counts

## Installation
requires Go 1.22+.

1. **Clone the repo:**
```bash
    git clone https://github.com/hnpf/beetrap.git
    cd beetrap
```
2. **Run:**
```bash
    go run ./cmd/beetrap
```

> to bind standard ports (22, 21, 80) without root:
> ```bash
> go build -o beetrap ./cmd/beetrap
> sudo setcap cap_net_bind_service=ep ./beetrap
> ```

## Usage
test locally with `nc` or `curl` against ports 2222, 2121, or 8080. connection attempts show up live in the TUI with timestamp, service, remote address, and any payload sent.

press `q` or `Ctrl+C` to quit.

## License
MIT. see [LICENSE](LICENSE).
