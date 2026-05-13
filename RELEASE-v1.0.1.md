# DNS Probe Tool v1.0.1

## 🐛 Bug Fixes

- 修复DNSSEC显示问题：只有在指定 `--dnssec` 参数时才显示DNSSEC验证结果
- 之前不指定 `--dnssec` 参数也会显示"DNSSEC验证通过"，现在已修复

## ✨ 改进

- 优化DNSSEC参数传递逻辑
- 更新README测试结果

## 📦 下载

| 平台 | 架构 | 文件 |
|------|------|------|
| Windows | amd64 | `dns-probe-windows-amd64.exe` |
| Linux | amd64 | `dns-probe-linux-amd64` |
| Linux | arm64 | `dns-probe-linux-arm64` |

## 🚀 使用方法

```bash
# 基本查询
dns-probe example.com

# DNSSEC验证（需要指定参数才会显示）
dns-probe example.com --dnssec

# 使用指定DNS服务器
dns-probe example.com --dnssec 8.8.8.8
```

## 🧪 测试

所有19个测试用例已通过。

## 相关链接

- [v1.0.0 Release](https://github.com/vito-L/dns-probe/releases/tag/v1.0.0)
