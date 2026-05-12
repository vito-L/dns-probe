# DNS Probe Tool

一个用Go语言编写的DNS拨测工具，支持多平台、多架构。

## 功能特性

- 并发拨测多个DNS服务器
- 支持查询A记录（默认）和所有记录类型（--all）
- 美化输出格式
- 自动获取系统DNS服务器
- 支持自定义DNS服务器
- 跨平台支持（Windows/Linux/macOS）
- 跨架构支持（amd64/arm64）
- JSON格式输出（--json）
- TUI交互界面（--tui）
- DNS污染检测（--pollution）
- HTML报告生成（--html）
- 批量域名查询（--file）
- 查询历史记录（--history）
- DNSSEC验证
- DoH/DoT支持

## 支持的系统

### Windows

| 架构             | 文件名                           | 说明            |
| -------------- | ----------------------------- | ------------- |
| amd64 (x86_64) | `dns-probe-windows-amd64.exe` | Windows 64位系统 |

### Linux

| 架构              | 文件名                     | 说明             |
| --------------- | ----------------------- | -------------- |
| amd64 (x86_64)  | `dns-probe-linux-amd64` | 大多数Linux发行版    |
| arm64 (aarch64) | `dns-probe-linux-arm64` | ARM64架构（如国产系统） |

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

## 使用方法

### 基本用法

```bash
# 查询域名（使用系统DNS服务器，只查询A记录）
./dns-probe example.com

# 查询域名（指定DNS服务器）
./dns-probe example.com 8.8.8.8 114.114.114.114

# 查询所有记录类型
./dns-probe example.com --all

# 使用指定DNS服务器查询所有记录类型
./dns-probe example.com 8.8.8.8 --all
```

### JSON格式输出

```bash
# 输出JSON格式
./dns-probe example.com --json

# 输出JSON格式（所有记录类型）
./dns-probe example.com --all --json
```

### TUI交互界面

```bash
# 启动TUI交互界面
./dns-probe example.com --tui
```

### DNS污染检测

```bash
# 检测DNS污染
./dns-probe example.com --pollution
```

### HTML报告生成

```bash
# 生成HTML报告
./dns-probe example.com --html report.html
```

### 批量域名查询

```bash
# 批量查询文件中的域名
./dns-probe --file domains.txt

# 批量查询并输出JSON格式
./dns-probe --file domains.txt --json
```

### 查询历史记录

```bash
# 显示查询历史
./dns-probe --history
```

### Linux系统使用

```bash
# 下载对应架构的版本
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-amd64

# 添加执行权限
chmod +x dns-probe-linux-amd64

# 运行
./dns-probe-linux-amd64 example.com
```

### 国产系统使用

#### 银河麒麟 (Kylin) / 统信UOS / 深度Deepin

```bash
# 这些系统通常是amd64架构
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-amd64
chmod +x dns-probe-linux-amd64
./dns-probe-linux-amd64 example.com
```

#### ARM64架构的国产系统

```bash
# 如果系统是ARM64架构（如飞腾、鲲鹏处理器）
wget https://github.com/vito-L/dns-probe/releases/latest/download/dns-probe-linux-arm64
chmod +x dns-probe-linux-arm64
./dns-probe-linux-arm64 example.com
```

## 输出示例

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

### DNS污染检测

```
┌─ DNS服务器: 8.8.8.8
│  查询耗时: 192 ms
│  ⚠️  检测到DNS污染
│
│  类型         TTL    值
│  ────────────────────────────────────────────────────────────
│  A          1m     142.250.73.142
└──────────────────────────────────────────────────────────────────
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

## 命令行参数

| 参数 | 说明 |
|------|------|
| `<域名>` | 要查询的域名 |
| `[DNS服务器...]` | 指定DNS服务器（可选） |
| `--all` | 查询所有记录类型 |
| `--json` | 输出JSON格式 |
| `--tui` | 启动TUI交互界面 |
| `--pollution` | 检测DNS污染 |
| `--html <文件>` | 生成HTML报告 |
| `--file <文件>` | 批量查询文件中的域名 |
| `--history` | 显示查询历史 |

## 系统DNS服务器

工具会自动获取系统配置的DNS服务器：

- Windows: 从`ipconfig /all`获取
- Linux: 从`/etc/resolv.conf`获取

如果获取失败，默认使用`8.8.8.8`

## 支持的记录类型

使用`--all`参数时，支持查询以下记录类型：

- A (IPv4地址)
- AAAA (IPv6地址)
- CNAME (别名)
- MX (邮件交换)
- NS (域名服务器)
- TXT (文本记录)
- SOA (权威记录)
- SRV (服务记录)
- CAA (证书授权)
- PTR (反向解析)

## 编译

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

## 依赖

- Go 1.21+
- github.com/miekg/dns
- github.com/charmbracelet/bubbletea
- github.com/charmbracelet/lipgloss

## 许可证

MIT License
