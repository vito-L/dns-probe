# DNS Probe Tool v1.0.2

## 🌐 Internationalization

- Translated all user-facing strings to English (help text, error messages, output labels)
- Added English README (`README_EN.md`)
- Added cross-links between Chinese and English READMEs

## ✨ Improvements

- Help text and all CLI output now in English for broader accessibility
- Renamed internal variable `国外DNS` to `foreignDNS` for code consistency

## 📦 Download

| Platform | Architecture | File |
|----------|-------------|------|
| Windows  | amd64       | `dns-probe-windows-amd64.exe` |
| Linux    | amd64       | `dns-probe-linux-amd64` |
| Linux    | arm64       | `dns-probe-linux-arm64` |

## 🚀 Usage

```bash
# Basic query
dns-probe example.com

# DNSSEC validation
dns-probe example.com --dnssec

# DNS pollution detection
dns-probe example.com --pollution
```

## 🧪 Tests

All 19 test cases have passed.

## Related Links

- [v1.0.1 Release](https://github.com/vito-L/dns-probe/releases/tag/v1.0.1)
- [v1.0.0 Release](https://github.com/vito-L/dns-probe/releases/tag/v1.0.0)
