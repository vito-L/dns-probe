# DNS Probe Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/vito-L/dns-probe)](https://github.com/vito-L/dns-probe/releases)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-brightgreen)]()
[![Architecture](https://img.shields.io/badge/Architecture-amd64%20%7C%20arm64-orange)]()

一个用Go语言编写的DNS拨测工具，支持多平台、多架构。

## ✨ 功能特性

| 功能          | 说明                                        |
| ----------- | ----------------------------------------- |
| 🔍 基本查询     | 并发拨测多个DNS服务器                              |
| 📊 记录类型     | 支持A/AAAA/CNAME/MX/NS/TXT/SOA/SRV/CAA/PTR等 |
| 🎨 美化输出     | 终端友好的格式化输出                                |
| 🌐 自动检测     | 自动获取系统DNS服务器                              |
| 🔧 自定义      | 支持自定义DNS服务器                               |
| 💻 跨平台      | Windows/Linux/macOS                       |
| 🏗️ 跨架构     | amd64/arm64                               |
| 📝 JSON输出   | `--json` 格式化输出                            |
| 🛡️ DNS污染检测 | `--pollution` 检测DNS污染                     |
| 🔐 DNSSEC验证 | `--dnssec` 验证DNSSEC签名                     |
| 🔒 DoH支持    | `--doh` DNS over HTTPS                    |
| 🔒 DoT支持    | `--dot` DNS over TLS                      |
| 📄 HTML报告   | `--html` 生成可视化报告                          |
| 📁 批量查询     | `--file` 批量查询域名                           |
| 📚 历史记录     | `--history` 查看查询历史                        |

## 📦 支持的系统

### Windows

| 架构    | 文件名                           | 说明            |
| ----- | ----------------------------- | ------------- |
| amd64 | `dns-probe-windows-amd64.exe` | Windows 64位系统 |

### Linux

| 架构    | 文件名                     | 说明             |
| ----- | ----------------------- | -------------- |
| amd64 | `dns-probe-linux-amd64` | 大多数Linux发行版    |
| arm64 | `dns-probe-linux-arm64` | ARM64架构（如国产系统） |

### 支持的Linux发行版

- Ubuntu / Debian
- CentOS / RHEL / Fedora
- 银河麒麟 (Kylin)
- 统信UOS
- 深度Deepin
- 中标麒麟
- 红旗Linux
- Arch Linux
- openSUSE
- Alpine Linux
- 以及其他主流发行版

## 🚀 快速开始

### 下载

```bash
# Linux amd64
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-amd64
chmod +x dns-probe-linux-amd64

# Linux arm64
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-arm64
chmod +x dns-probe-linux-arm64
```

### 基本用法

```bash
# 使用系统DNS服务器查询A记录
./dns-probe example.com

# 使用指定DNS服务器
./dns-probe example.com 8.8.8.8 114.114.114.114

# 查询所有记录类型
./dns-probe example.com --all
```

## 📖 使用方法

### 命令行参数

| 参数               | 说明         | 示例                                   |
| ---------------- | ---------- | ------------------------------------ |
| `<域名>`           | 要查询的域名     | `example.com`                        |
| `[DNS服务器...]`    | 指定DNS服务器   | `8.8.8.8 114.114.114.114`            |
| `--all`          | 查询所有记录类型   | `--all`                              |
| `--json`         | 输出JSON格式   | `--json`                             |
| `--pollution`    | 检测DNS污染    | `--pollution`                        |
| `--dnssec`       | 启用DNSSEC验证 | `--dnssec`                           |
| `--doh <url>`    | 使用DoH服务器   | `--doh https://dns.google/dns-query` |
| `--dot <server>` | 使用DoT服务器   | `--dot dns.alidns.com:853`           |
| `--html <文件>`    | 生成HTML报告   | `--html report.html`                 |
| `--file <文件>`    | 批量查询文件中的域名 | `--file domains.txt`                 |
| `--history`      | 显示查询历史     | `--history`                          |

### 使用示例

#### 基本查询

```bash
# 使用系统DNS服务器
./dns-probe example.com

# 使用指定DNS服务器
./dns-probe example.com 8.8.8.8

# 使用多个DNS服务器
./dns-probe example.com 8.8.8.8 114.114.114.114 223.5.5.5
```

#### DNSSEC验证

```bash
# 使用系统DNS服务器进行DNSSEC验证
./dns-probe example.com --dnssec

# 使用指定DNS服务器进行DNSSEC验证
./dns-probe example.com --dnssec 8.8.8.8
```

#### DNS污染检测

```bash
# 检测DNS污染
./dns-probe example.com --pollution
```

输出示例：

```
┌─ DNS服务器: 114.114.114.114
│  查询耗时: 23 ms
│  ⚠️  检测到DNS污染
│  被污染的IP: 1.2.3.4
│  真实IP（国外DNS）: 142.250.73.78
│
│  类型         TTL    值
│  ────────────────────────────────────────────────────────────
│  A          1m     1.2.3.4
└──────────────────────────────────────────────────────────────────
```

#### DoH/DoT支持

```bash
# 使用DoH服务器
./dns-probe example.com --doh https://dns.google/dns-query

# 使用DoT服务器
./dns-probe example.com --dot dns.alidns.com:853
```

#### JSON输出

```bash
# 输出JSON格式
./dns-probe example.com --json

# 输出JSON格式（所有记录类型）
./dns-probe example.com --all --json
```

#### HTML报告

```bash
# 生成HTML报告
./dns-probe example.com --html report.html
```

#### 批量查询

```bash
# 批量查询文件中的域名
./dns-probe --file domains.txt

# 批量查询并输出JSON格式
./dns-probe --file domains.txt --json
```

#### 查询历史

```bash
# 显示查询历史
./dns-probe --history
```

## 📊 输出示例

### 默认输出（A记录）

```
╔══════════════════════════════════════════════════════════════════╗
║                    DNS Probe Tool v1.0                          ║
╚══════════════════════════════════════════════════════════════════╝

  域名: baidu.com
  时间: 2026-05-12 22:44:50

┌─ DNS服务器: 8.8.8.8
│  查询耗时: 209 ms
│
│  类型         TTL    值
│  ────────────────────────────────────────────────────────────
│  A          4m     111.63.65.247
│  A          4m     111.63.65.103
│  A          4m     110.242.74.102
│  A          4m     124.237.177.164
└──────────────────────────────────────────────────────────────────
```

### JSON输出

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

### 查询历史

```
查询历史:
────────────────────────────────────────────────────────────
[1] 2026-05-12T23:14:32+08:00 - baidu.com
    DNS: 8.8.8.8, 耗时: 205 ms
[2] 2026-05-12T23:15:01+08:00 - google.com
    DNS: 8.8.8.8, 耗时: 192 ms
```

## 🔧 系统DNS服务器

工具会自动获取系统配置的DNS服务器：

- **Windows**: 从`ipconfig /all`获取
- **Linux**: 从`/etc/resolv.conf`获取

如果获取失败，默认使用`8.8.8.8`

## 📋 支持的记录类型

使用`--all`参数时，支持查询以下记录类型：

| 记录类型  | 说明     |
| ----- | ------ |
| A     | IPv4地址 |
| AAAA  | IPv6地址 |
| CNAME | 别名     |
| MX    | 邮件交换   |
| NS    | 域名服务器  |
| TXT   | 文本记录   |
| SOA   | 权威记录   |
| SRV   | 服务记录   |
| CAA   | 证书授权   |
| PTR   | 反向解析   |

## 🛠️ 编译

如果需要从源码编译：

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

## 📦 依赖

- [Go](https://golang.org) 1.21+
- [miekg/dns](https://github.com/miekg/dns)

## 🧪 测试

所有测试用例已通过：

```bash
$ go test -v ./...

=== RUN   TestRecordTypeName
--- PASS: TestRecordTypeName (0.00s)
=== RUN   TestFormatTTL
--- PASS: TestFormatTTL (0.00s)
=== RUN   TestGetSystemDNSServers
--- PASS: TestGetSystemDNSServers (0.02s)
=== RUN   TestProbeDNS
--- PASS: TestProbeDNS (0.19s)
=== RUN   TestProbeDNSInvalidDomain
--- PASS: TestProbeDNSInvalidDomain (0.21s)
=== RUN   TestProbeAll
--- PASS: TestProbeAll (0.25s)
=== RUN   TestProbeAllRecordTypes
--- PASS: TestProbeAllRecordTypes (2.00s)
=== RUN   TestFormatText
--- PASS: TestFormatText (0.01s)
=== RUN   TestFormatJSON
--- PASS: TestFormatJSON (0.00s)
=== RUN   TestReadDomainsFromFile
--- PASS: TestReadDomainsFromFile (0.01s)
=== RUN   TestReadDomainsFromFileNotFound
--- PASS: TestReadDomainsFromFileNotFound (0.00s)
=== RUN   TestSaveAndLoadHistory
--- PASS: TestSaveAndLoadHistory (0.01s)
=== RUN   TestFormatHistory
--- PASS: TestFormatHistory (0.00s)
=== RUN   TestProbeDoH
--- PASS: TestProbeDoH (0.19s)
=== RUN   TestProbeDoT
--- PASS: TestProbeDoT (0.05s)
=== RUN   TestDNSSECValidation
--- PASS: TestDNSSECValidation (0.20s)
=== RUN   TestPollutionDetection
--- PASS: TestPollutionDetection (0.42s)
=== RUN   TestMultipleDomains
--- PASS: TestMultipleDomains (0.41s)
=== RUN   TestFormatMultipleJSON
--- PASS: TestFormatMultipleJSON (0.00s)
PASS
ok      github.com/vito-L/dns-probe    5.264s
```

## 📄 许可证

[MIT License](LICENSE)

## 🤝 贡献

欢迎提交Issue和Pull Request！

## ⭐ Star

如果这个项目对你有帮助，请给个Star支持一下！
