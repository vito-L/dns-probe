# DNS Probe Tool v1.0.0

## 🎉 首次发布

DNS Probe Tool 是一个用 Go 语言编写的 DNS 拨测工具，支持多平台、多架构。

## ✨ 功能特性

- 🔍 **基本查询** - 并发拨测多个 DNS 服务器
- 📊 **记录类型** - 支持 A/AAAA/CNAME/MX/NS/TXT/SOA/SRV/CAA/PTR 等
- 🎨 **美化输出** - 终端友好的格式化输出
- 🌐 **自动检测** - 自动获取系统 DNS 服务器
- 🔧 **自定义** - 支持自定义 DNS 服务器
- 💻 **跨平台** - Windows/Linux/macOS
- 🏗️ **跨架构** - amd64/arm64
- 📝 **JSON 输出** - `--json` 格式化输出
- 🛡️ **DNS 污染检测** - `--pollution` 检测 DNS 污染
- 🔐 **DNSSEC 验证** - `--dnssec` 验证 DNSSEC 签名
- 🔒 **DoH 支持** - `--doh` DNS over HTTPS
- 🔒 **DoT 支持** - `--dot` DNS over TLS
- 📄 **HTML 报告** - `--html` 生成可视化报告
- 📁 **批量查询** - `--file` 批量查询域名
- 📚 **历史记录** - `--history` 查看查询历史

## 📦 下载

| 平台 | 架构 | 文件 |
|------|------|------|
| Windows | amd64 | `dns-probe-windows-amd64.exe` |
| Linux | amd64 | `dns-probe-linux-amd64` |
| Linux | arm64 | `dns-probe-linux-arm64` |

## 🚀 快速开始

### Windows

```powershell
# 下载并运行
.\dns-probe-windows-amd64.exe example.com
```

### Linux

```bash
# 下载
wget https://github.com/vito-L/dns-probe/releases/download/v1.0.0/dns-probe-linux-amd64

# 添加执行权限
chmod +x dns-probe-linux-amd64

# 运行
./dns-probe-linux-amd64 example.com
```

## 📖 使用示例

```bash
# 基本查询
dns-probe example.com

# 使用指定 DNS 服务器
dns-probe example.com 8.8.8.8 114.114.114.114

# 查询所有记录类型
dns-probe example.com --all

# DNSSEC 验证
dns-probe example.com --dnssec 8.8.8.8

# DNS 污染检测
dns-probe example.com --pollution

# DoH 查询
dns-probe example.com --doh https://dns.google/dns-query

# DoT 查询
dns-probe example.com --dot dns.alidns.com:853

# JSON 输出
dns-probe example.com --json

# 生成 HTML 报告
dns-probe example.com --html report.html

# 批量查询
dns-probe --file domains.txt

# 查看历史
dns-probe --history
```

## 🔧 命令行参数

| 参数 | 说明 |
|------|------|
| `<域名>` | 要查询的域名 |
| `[DNS服务器...]` | 指定 DNS 服务器 |
| `--all` | 查询所有记录类型 |
| `--json` | 输出 JSON 格式 |
| `--pollution` | 检测 DNS 污染 |
| `--dnssec` | 启用 DNSSEC 验证 |
| `--doh <url>` | 使用 DoH 服务器 |
| `--dot <server>` | 使用 DoT 服务器 |
| `--html <文件>` | 生成 HTML 报告 |
| `--file <文件>` | 批量查询文件中的域名 |
| `--history` | 显示查询历史 |

## 📋 测试结果

所有测试用例已通过：

```
=== RUN   TestRecordTypeName
--- PASS: TestRecordTypeName (0.00s)
=== RUN   TestFormatTTL
--- PASS: TestFormatTTL (0.00s)
=== RUN   TestGetSystemDNSServers
--- PASS: TestGetSystemDNSServers (0.02s)
=== RUN   TestProbeDNS
--- PASS: TestProbeDNS (0.19s)
...
PASS
ok  	github.com/vito-L/dns-probe	5.264s
```

## 📄 许可证

[MIT License](LICENSE)
