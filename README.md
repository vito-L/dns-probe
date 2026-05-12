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

## 许可证

MIT License
