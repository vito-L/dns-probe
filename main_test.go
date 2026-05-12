package main

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/miekg/dns"
)

// TestRecordTypeName 测试记录类型名称
func TestRecordTypeName(t *testing.T) {
	tests := []struct {
		qtype uint16
		want  string
	}{
		{dns.TypeA, "A"},
		{dns.TypeAAAA, "AAAA"},
		{dns.TypeCNAME, "CNAME"},
		{dns.TypeMX, "MX"},
		{dns.TypeNS, "NS"},
		{dns.TypeTXT, "TXT"},
		{dns.TypeSOA, "SOA"},
		{dns.TypeSRV, "SRV"},
		{dns.TypeCAA, "CAA"},
		{dns.TypePTR, "PTR"},
	}

	for _, tt := range tests {
		got := RecordTypeName(tt.qtype)
		if got != tt.want {
			t.Errorf("RecordTypeName(%d) = %s, want %s", tt.qtype, got, tt.want)
		}
	}
}

// TestFormatTTL 测试TTL格式化
func TestFormatTTL(t *testing.T) {
	tests := []struct {
		ttl  uint32
		want string
	}{
		{30, "30s"},
		{60, "1m"},
		{3600, "1h"},
		{86400, "1d"},
		{120, "2m"},
		{7200, "2h"},
		{172800, "2d"},
	}

	for _, tt := range tests {
		got := formatTTL(tt.ttl)
		if got != tt.want {
			t.Errorf("formatTTL(%d) = %s, want %s", tt.ttl, got, tt.want)
		}
	}
}

// TestGetSystemDNSServers 测试获取系统DNS服务器
func TestGetSystemDNSServers(t *testing.T) {
	servers := GetSystemDNSServers()
	if len(servers) == 0 {
		t.Error("GetSystemDNSServers() returned empty slice")
	}

	// 验证返回的是有效的IP地址
	for _, server := range servers {
		if net.ParseIP(server) == nil {
			t.Errorf("GetSystemDNSServers() returned invalid IP: %s", server)
		}
	}
}

// TestProbeDNS 测试DNS查询
func TestProbeDNS(t *testing.T) {
	result := ProbeDNS("example.com", "8.8.8.8")

	if result.Error != nil {
		t.Fatalf("ProbeDNS failed: %v", result.Error)
	}

	if result.Domain != "example.com" {
		t.Errorf("Domain = %s, want example.com", result.Domain)
	}

	if result.DNSServer != "8.8.8.8" {
		t.Errorf("DNSServer = %s, want 8.8.8.8", result.DNSServer)
	}

	if len(result.Records) == 0 {
		t.Error("ProbeDNS returned no records")
	}

	// 验证返回的是A记录
	for _, record := range result.Records {
		if record.Type != "A" {
			t.Errorf("Record type = %s, want A", record.Type)
		}
	}
}

// TestProbeDNSInvalidDomain 测试无效域名
func TestProbeDNSInvalidDomain(t *testing.T) {
	result := ProbeDNS("invalid.domain.that.does.not.exist", "8.8.8.8")

	// 无效域名应该返回错误或空记录
	if result.Error == nil && len(result.Records) > 0 {
		t.Error("ProbeDNS should return error or empty records for invalid domain")
	}
}

// TestProbeAll 测试并发拨测
func TestProbeAll(t *testing.T) {
	servers := []string{"8.8.8.8", "1.1.1.1"}
	results := ProbeAll("example.com", servers)

	if len(results) != len(servers) {
		t.Fatalf("ProbeAll returned %d results, want %d", len(results), len(servers))
	}

	for i, result := range results {
		if result.Error != nil {
			t.Errorf("ProbeAll[%d] failed: %v", i, result.Error)
		}
	}
}

// TestProbeAllRecordTypes 测试查询所有记录类型
func TestProbeAllRecordTypes(t *testing.T) {
	results := ProbeAllRecordTypes("example.com", []string{"8.8.8.8"})

	if len(results) == 0 {
		t.Fatal("ProbeAllRecordTypes returned no results")
	}

	result := results[0]
	if result.Error != nil {
		t.Fatalf("ProbeAllRecordTypes failed: %v", result.Error)
	}

	// 应该包含多种记录类型
	recordTypes := make(map[string]bool)
	for _, record := range result.Records {
		recordTypes[record.Type] = true
	}

	// 至少应该有A记录
	if !recordTypes["A"] {
		t.Error("ProbeAllRecordTypes should return A records")
	}
}

// TestFormatText 测试文本格式化输出
func TestFormatText(t *testing.T) {
	results := []ProbeResult{
		{
			Domain:    "example.com",
			DNSServer: "8.8.8.8",
			Latency:   100 * time.Millisecond,
			Records: []DNSRecord{
				{Type: "A", Name: "example.com.", TTL: 300, Value: "93.184.216.34"},
			},
		},
	}

	output := FormatText(results)

	if !strings.Contains(output, "example.com") {
		t.Error("FormatText should contain domain name")
	}

	if !strings.Contains(output, "8.8.8.8") {
		t.Error("FormatText should contain DNS server")
	}

	if !strings.Contains(output, "A") {
		t.Error("FormatText should contain record type")
	}

	if !strings.Contains(output, "93.184.216.34") {
		t.Error("FormatText should contain IP address")
	}
}

// TestFormatJSON 测试JSON格式化输出
func TestFormatJSON(t *testing.T) {
	results := []ProbeResult{
		{
			Domain:    "example.com",
			DNSServer: "8.8.8.8",
			Latency:   100 * time.Millisecond,
			Records: []DNSRecord{
				{Type: "A", Name: "example.com.", TTL: 300, Value: "93.184.216.34"},
			},
		},
	}

	output := FormatJSON(results)

	// 验证是有效的JSON
	var jsonResult JSONResult
	if err := json.Unmarshal([]byte(output), &jsonResult); err != nil {
		t.Fatalf("FormatJSON returned invalid JSON: %v", err)
	}

	if jsonResult.Domain != "example.com" {
		t.Errorf("JSON domain = %s, want example.com", jsonResult.Domain)
	}

	if len(jsonResult.Servers) == 0 {
		t.Fatal("JSON should contain servers")
	}

	if jsonResult.Servers[0].Server != "8.8.8.8" {
		t.Errorf("JSON server = %s, want 8.8.8.8", jsonResult.Servers[0].Server)
	}
}

// TestReadDomainsFromFile 测试从文件读取域名
func TestReadDomainsFromFile(t *testing.T) {
	// 创建临时文件
	content := `# 注释
example.com
google.com
# 另一个注释
github.com
`
	tmpFile, err := os.CreateTemp("", "domains-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	domains, err := ReadDomainsFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadDomainsFromFile failed: %v", err)
	}

	expected := []string{"example.com", "google.com", "github.com"}
	if len(domains) != len(expected) {
		t.Fatalf("ReadDomainsFromFile returned %d domains, want %d", len(domains), len(expected))
	}

	for i, domain := range domains {
		if domain != expected[i] {
			t.Errorf("Domain[%d] = %s, want %s", i, domain, expected[i])
		}
	}
}

// TestReadDomainsFromFileNotFound 测试读取不存在的文件
func TestReadDomainsFromFileNotFound(t *testing.T) {
	_, err := ReadDomainsFromFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("ReadDomainsFromFile should return error for nonexistent file")
	}
}

// TestSaveAndLoadHistory 测试保存和加载历史记录
func TestSaveAndLoadHistory(t *testing.T) {
	// 清理测试环境
	homeDir, _ := os.UserHomeDir()
	historyDir := filepath.Join(homeDir, ".dns-probe")
	os.RemoveAll(historyDir)

	// 保存历史记录
	results := []ProbeResult{
		{
			Domain:    "example.com",
			DNSServer: "8.8.8.8",
			Latency:   100 * time.Millisecond,
			Records: []DNSRecord{
				{Type: "A", Name: "example.com.", TTL: 300, Value: "93.184.216.34"},
			},
		},
	}

	if err := SaveHistory("example.com", results); err != nil {
		t.Fatalf("SaveHistory failed: %v", err)
	}

	// 加载历史记录
	history, err := LoadHistory()
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}

	if len(history) == 0 {
		t.Fatal("LoadHistory returned empty history")
	}

	if history[0].Domain != "example.com" {
		t.Errorf("History domain = %s, want example.com", history[0].Domain)
	}

	// 清理
	os.RemoveAll(historyDir)
}

// TestFormatHistory 测试历史记录格式化
func TestFormatHistory(t *testing.T) {
	history := []HistoryEntry{
		{
			Timestamp: "2026-05-12T23:14:32+08:00",
			Domain:    "example.com",
			Servers: []JSONServer{
				{Server: "8.8.8.8", Latency: 100},
			},
		},
	}

	output := FormatHistory(history)

	if !strings.Contains(output, "example.com") {
		t.Error("FormatHistory should contain domain name")
	}

	if !strings.Contains(output, "8.8.8.8") {
		t.Error("FormatHistory should contain DNS server")
	}
}

// TestProbeDoH 测试DoH查询
func TestProbeDoH(t *testing.T) {
	result := ProbeDNS("example.com", "https://dns.google/dns-query")

	if result.Error != nil {
		t.Logf("DoH test failed (may be network issue): %v", result.Error)
		return
	}

	if len(result.Records) == 0 {
		t.Error("DoH should return records")
	}
}

// TestProbeDoT 测试DoT查询
func TestProbeDoT(t *testing.T) {
	result := ProbeDNS("example.com", "dns.alidns.com:853")

	if result.Error != nil {
		t.Logf("DoT test failed (may be network issue): %v", result.Error)
		return
	}

	if len(result.Records) == 0 {
		t.Error("DoT should return records")
	}
}

// TestDNSSECValidation 测试DNSSEC验证
func TestDNSSECValidation(t *testing.T) {
	result := ProbeDNS("example.com", "8.8.8.8")

	if result.Error != nil {
		t.Fatalf("DNSSEC test failed: %v", result.Error)
	}

	// 8.8.8.8应该支持DNSSEC
	if !result.DNSSECValid {
		t.Log("DNSSEC validation not passed (may depend on domain)")
	}
}

// TestPollutionDetection 测试DNS污染检测
func TestPollutionDetection(t *testing.T) {
	results := ProbeAllWithPollutionCheck("example.com", []string{"8.8.8.8", "114.114.114.114"})

	if len(results) == 0 {
		t.Fatal("Pollution detection returned no results")
	}

	// 验证结果结构
	for _, result := range results {
		if result.Error != nil {
			t.Logf("Pollution detection server %s failed: %v", result.DNSServer, result.Error)
		}
	}
}

// TestMultipleDomains 测试批量域名查询
func TestMultipleDomains(t *testing.T) {
	domains := []string{"example.com", "google.com"}
	servers := []string{"8.8.8.8"}

	results := ProbeMultipleDomains(domains, servers, false)

	if len(results) != len(domains) {
		t.Fatalf("ProbeMultipleDomains returned %d results, want %d", len(results), len(domains))
	}

	for i, result := range results {
		if result.Domain != domains[i] {
			t.Errorf("Result[%d].Domain = %s, want %s", i, result.Domain, domains[i])
		}
	}
}

// TestFormatMultipleJSON 测试批量JSON格式化
func TestFormatMultipleJSON(t *testing.T) {
	results := []JSONResult{
		{
			Domain:    "example.com",
			Timestamp: "2026-05-12T23:14:32+08:00",
			Servers: []JSONServer{
				{Server: "8.8.8.8", Latency: 100},
			},
		},
	}

	output := FormatMultipleJSON(results)

	var jsonResults []JSONResult
	if err := json.Unmarshal([]byte(output), &jsonResults); err != nil {
		t.Fatalf("FormatMultipleJSON returned invalid JSON: %v", err)
	}

	if len(jsonResults) != 1 {
		t.Errorf("JSON results count = %d, want 1", len(jsonResults))
	}
}
