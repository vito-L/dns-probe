# DNS Probe Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/vito-L/dns-probe)](https://github.com/vito-L/dns-probe/releases)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-brightgreen)]()
[![Architecture](https://img.shields.io/badge/Architecture-amd64%20%7C%20arm64-orange)]()

> **[中文文档](README_CN.md)**

A DNS probing tool written in Go, supporting multiple platforms and architectures.

## ✨ Features

| Feature          | Description                                        |
| ---------------- | -------------------------------------------------- |
| 🔍 Basic Query     | Concurrent probing of multiple DNS servers                              |
| 📊 Record Types     | Supports A/AAAA/CNAME/MX/NS/TXT/SOA/SRV/CAA/PTR and more |
| 🎨 Formatted Output     | Terminal-friendly formatted output                                |
| 🌐 Auto Detection     | Automatically detects system DNS servers                              |
| 🔧 Customization      | Supports custom DNS servers                               |
| 💻 Cross-Platform      | Windows/Linux/macOS                       |
| 🏗️ Cross-Architecture     | amd64/arm64                               |
| 📝 JSON Output   | `--json` formatted output                            |
| 🛡️ DNS Pollution Detection | `--pollution` detects DNS pollution                     |
| 🔐 DNSSEC Validation | `--dnssec` verifies DNSSEC signatures                     |
| 🔒 DoH Support    | `--doh` DNS over HTTPS                    |
| 🔒 DoT Support    | `--dot` DNS over TLS                      |
| 📄 HTML Report   | `--html` generates visual report                          |
| 📁 Batch Query     | `--file` batch query domains                           |
| 📚 History     | `--history` view query history                        |

## 📦 Supported Systems

### Windows

| Architecture | Filename                           | Description            |
| ------------ | ---------------------------------- | ---------------------- |
| amd64        | `dns-probe-windows-amd64.exe`      | Windows 64-bit systems |

### Linux

| Architecture | Filename                     | Description                    |
| ------------ | ---------------------------- | ------------------------------ |
| amd64        | `dns-probe-linux-amd64`      | Most Linux distributions       |
| arm64        | `dns-probe-linux-arm64`      | ARM64 architecture (e.g., ARM servers) |

### Supported Linux Distributions

- Ubuntu / Debian
- CentOS / RHEL / Fedora
- Arch Linux
- openSUSE
- Alpine Linux
- And other mainstream distributions

## 🚀 Quick Start

### Download

```bash
# Linux amd64
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-amd64
chmod +x dns-probe-linux-amd64

# Linux arm64
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-arm64
chmod +x dns-probe-linux-arm64
```

### Basic Usage

```bash
# Query A record using system DNS
./dns-probe example.com

# Use specified DNS servers
./dns-probe example.com 8.8.8.8 114.114.114.114

# Query all record types
./dns-probe example.com --all
```

## 📖 Usage

### Command Line Arguments

| Argument               | Description         | Example                                   |
| ---------------------- | ------------------- | ----------------------------------------- |
| `<domain>`           | Domain to query     | `example.com`                        |
| `[DNS servers...]`    | Specify DNS servers   | `8.8.8.8 114.114.114.114`            |
| `--all`          | Query all record types   | `--all`                              |
| `--json`         | Output in JSON format   | `--json`                             |
| `--pollution`    | Detect DNS pollution    | `--pollution`                        |
| `--dnssec`       | Enable DNSSEC validation | `--dnssec`                           |
| `--doh <url>`    | Use DoH server   | `--doh https://dns.google/dns-query` |
| `--dot <server>` | Use DoT server   | `--dot dns.alidns.com:853`           |
| `--html <file>`    | Generate HTML report   | `--html report.html`                 |
| `--file <file>`    | Batch query domains from file | `--file domains.txt`                 |
| `--history`      | Show query history     | `--history`                          |

### Examples

#### Basic Query

```bash
# Using system DNS server
./dns-probe example.com

# Using specified DNS server
./dns-probe example.com 8.8.8.8

# Using multiple DNS servers
./dns-probe example.com 8.8.8.8 114.114.114.114 223.5.5.5
```

#### DNSSEC Validation

```bash
# DNSSEC validation using system DNS
./dns-probe example.com --dnssec

# DNSSEC validation using specified DNS
./dns-probe example.com --dnssec 8.8.8.8
```

#### DNS Pollution Detection

```bash
# Detect DNS pollution
./dns-probe example.com --pollution
```

Output example:

```
┌─ DNS Server: 114.114.114.114
│  Latency: 23 ms
│  ⚠️  DNS pollution detected
│  Polluted IP: 1.2.3.4
│  Real IP (foreign DNS): 142.250.73.78
│
│  Type         TTL    Value
│  ────────────────────────────────────────────────────────────
│  A          1m     1.2.3.4
└──────────────────────────────────────────────────────────────────
```

#### DoH/DoT Support

```bash
# Using DoH server
./dns-probe example.com --doh https://dns.google/dns-query

# Using DoT server
./dns-probe example.com --dot dns.alidns.com:853
```

#### JSON Output

```bash
# Output in JSON format
./dns-probe example.com --json

# Output in JSON format (all record types)
./dns-probe example.com --all --json
```

#### HTML Report

```bash
# Generate HTML report
./dns-probe example.com --html report.html
```

#### Batch Query

```bash
# Batch query domains from file
./dns-probe --file domains.txt

# Batch query with JSON output
./dns-probe --file domains.txt --json
```

#### Query History

```bash
# Show query history
./dns-probe --history
```

## 📊 Output Examples

### Default Output (A Record)

```
╔══════════════════════════════════════════════════════════════════╗
║                    DNS Probe Tool v1.0                          ║
╚══════════════════════════════════════════════════════════════════╝

  Domain: baidu.com
  Time:   2026-05-12 22:44:50

┌─ DNS Server: 8.8.8.8
│  Latency: 209 ms
│
│  Type         TTL    Value
│  ────────────────────────────────────────────────────────────
│  A          4m     111.63.65.247
│  A          4m     111.63.65.103
│  A          4m     110.242.74.102
│  A          4m     124.237.177.164
└──────────────────────────────────────────────────────────────────
```

### JSON Output

```json
{
  "domain": "baidu.com",
  "timestamp": "2026-05-12T22:55:09+08:00",
  "servers": [
    {
      "server": "8.8.8.8",
      "latency_ms": 267,
      "records": [
        {
          "type": "A",
          "name": "baidu.com.",
          "ttl": 580,
          "value": "124.237.177.164"
        }
      ]
    }
  ]
}
```

### Query History

```
Query History:
────────────────────────────────────────────────────────────
[1] 2026-05-12T23:14:32+08:00 - baidu.com
    DNS: 8.8.8.8, Latency: 205 ms
[2] 2026-05-12T23:15:01+08:00 - google.com
    DNS: 8.8.8.8, Latency: 192 ms
```

## 🔧 System DNS Servers

The tool automatically detects system-configured DNS servers:

- **Windows**: Retrieved from `ipconfig /all`
- **Linux**: Retrieved from `/etc/resolv.conf`

If detection fails, defaults to `8.8.8.8`

## 📋 Supported Record Types

When using the `--all` parameter, the following record types are supported:

| Record Type | Description     |
| ----------- | --------------- |
| A           | IPv4 address    |
| AAAA        | IPv6 address    |
| CNAME       | Canonical name  |
| MX          | Mail exchange    |
| NS          | Name server     |
| TXT         | Text record     |
| SOA         | Authority record |
| SRV         | Service record  |
| CAA         | Certificate authority |
| PTR         | Reverse lookup  |

## 🛠️ Building from Source

```bash
# Windows
go build -o dns-probe.exe .

# Linux amd64
GOOS=linux GOARCH=amd64 go build -o dns-probe-linux-amd64 .

# Linux arm64
GOOS=linux GOARCH=arm64 go build -o dns-probe-linux-arm64 .

# macOS amd64
GOOS=darwin GOARCH=amd64 go build -o dns-probe-darwin-amd64 .

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dns-probe-darwin-arm64 .
```

## 📦 Dependencies

- [Go](https://golang.org) 1.21+
- [miekg/dns](https://github.com/miekg/dns)

## 🧪 Tests

All test cases have passed:

```bash
$ go test -v ./...

=== RUN   TestRecordTypeName
--- PASS: TestRecordTypeName (0.00s)
=== RUN   TestFormatTTL
--- PASS: TestFormatTTL (0.00s)
=== RUN   TestGetSystemDNSServers
--- PASS: TestGetSystemDNSServers (0.03s)
=== RUN   TestProbeDNS
--- PASS: TestProbeDNS (0.21s)
=== RUN   TestProbeDNSInvalidDomain
--- PASS: TestProbeDNSInvalidDomain (0.21s)
=== RUN   TestProbeAll
--- PASS: TestProbeAll (0.22s)
=== RUN   TestProbeAllRecordTypes
--- PASS: TestProbeAllRecordTypes (2.15s)
=== RUN   TestFormatText
--- PASS: TestFormatText (0.01s)
=== RUN   TestFormatJSON
--- PASS: TestFormatJSON (0.00s)
=== RUN   TestReadDomainsFromFile
--- PASS: TestReadDomainsFromFile (0.01s)
=== RUN   TestReadDomainsFromFileNotFound
--- PASS: TestReadDomainsFromFileNotFound (0.00s)
=== RUN   TestSaveAndLoadHistory
--- PASS: TestSaveAndLoadHistory (0.02s)
=== RUN   TestFormatHistory
--- PASS: TestFormatHistory (0.00s)
=== RUN   TestProbeDoH
--- PASS: TestProbeDoH (1.30s)
=== RUN   TestProbeDoT
--- PASS: TestProbeDoT (0.04s)
=== RUN   TestDNSSECValidation
--- PASS: TestDNSSECValidation (0.20s)
=== RUN   TestPollutionDetection
--- PASS: TestPollutionDetection (0.48s)
=== RUN   TestMultipleDomains
--- PASS: TestMultipleDomains (0.46s)
=== RUN   TestFormatMultipleJSON
--- PASS: TestFormatMultipleJSON (0.00s)
PASS
ok  	github.com/vito-L/dns-probe	6.918s
```

## 📄 License

[MIT License](LICENSE)

## 🤝 Contributing

Issues and Pull Requests are welcome!

## ⭐ Star

If this project helps you, please give it a Star!
