package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// GetSystemDNSServers 获取系统DNS服务器
func GetSystemDNSServers() []string {
	// Windows系统
	if cmd := exec.Command("ipconfig", "/all"); cmd != nil {
		output, err := cmd.Output()
		if err == nil {
			// 匹配DNS服务器地址
			re := regexp.MustCompile(`DNS Servers[^:]*:\s*(\d+\.\d+\.\d+\.\d+)`)
			matches := re.FindAllStringSubmatch(string(output), -1)
			if len(matches) > 0 {
				var servers []string
				for _, match := range matches {
					if len(match) > 1 {
						servers = append(servers, match[1])
					}
				}
				if len(servers) > 0 {
					return servers
				}
			}
		}
	}

	// Linux/Mac系统
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err == nil && len(conf.Servers) > 0 {
		return conf.Servers
	}

	// 默认返回8.8.8.8
	return []string{"8.8.8.8"}
}

// DNS服务器列表
var DefaultDNSServers = GetSystemDNSServers()

// RecordTypeName 获取记录类型名称
func RecordTypeName(qtype uint16) string {
	switch qtype {
	case dns.TypeA:
		return "A"
	case dns.TypeAAAA:
		return "AAAA"
	case dns.TypeCNAME:
		return "CNAME"
	case dns.TypeMX:
		return "MX"
	case dns.TypeNS:
		return "NS"
	case dns.TypeTXT:
		return "TXT"
	case dns.TypeSOA:
		return "SOA"
	case dns.TypeSRV:
		return "SRV"
	case dns.TypeCAA:
		return "CAA"
	case dns.TypePTR:
		return "PTR"
	case dns.TypeDNSKEY:
		return "DNSKEY"
	case dns.TypeDS:
		return "DS"
	case dns.TypeNSEC:
		return "NSEC"
	case dns.TypeNSEC3:
		return "NSEC3"
	case dns.TypeRRSIG:
		return "RRSIG"
	case dns.TypeHINFO:
		return "HINFO"
	case dns.TypeMINFO:
		return "MINFO"
	case dns.TypeRP:
		return "RP"
	case dns.TypeAFSDB:
		return "AFSDB"
	case dns.TypeX25:
		return "X25"
	case dns.TypeISDN:
		return "ISDN"
	case dns.TypeRT:
		return "RT"
	case dns.TypeSIG:
		return "SIG"
	case dns.TypeKEY:
		return "KEY"
	case dns.TypePX:
		return "PX"
	case dns.TypeGPOS:
		return "GPOS"
	case dns.TypeLOC:
		return "LOC"
	case dns.TypeNXT:
		return "NXT"
	case dns.TypeNAPTR:
		return "NAPTR"
	case dns.TypeKX:
		return "KX"
	case dns.TypeCERT:
		return "CERT"
	case dns.TypeDNAME:
		return "DNAME"
	case dns.TypeOPT:
		return "OPT"
	case dns.TypeAPL:
		return "APL"
	case dns.TypeSSHFP:
		return "SSHFP"
	case dns.TypeIPSECKEY:
		return "IPSECKEY"
	case dns.TypeDHCID:
		return "DHCID"
	case dns.TypeNSEC3PARAM:
		return "NSEC3PARAM"
	case dns.TypeTLSA:
		return "TLSA"
	case dns.TypeSMIMEA:
		return "SMIMEA"
	case dns.TypeHIP:
		return "HIP"
	case dns.TypeCDS:
		return "CDS"
	case dns.TypeCDNSKEY:
		return "CDNSKEY"
	case dns.TypeOPENPGPKEY:
		return "OPENPGPKEY"
	case dns.TypeCSYNC:
		return "CSYNC"
	case dns.TypeZONEMD:
		return "ZONEMD"
	case dns.TypeSVCB:
		return "SVCB"
	case dns.TypeHTTPS:
		return "HTTPS"
	case dns.TypeSPF:
		return "SPF"
	case dns.TypeUINFO:
		return "UINFO"
	case dns.TypeUID:
		return "UID"
	case dns.TypeGID:
		return "GID"
	case dns.TypeNID:
		return "NID"
	case dns.TypeL32:
		return "L32"
	case dns.TypeL64:
		return "L64"
	case dns.TypeLP:
		return "LP"
	case dns.TypeEUI48:
		return "EUI48"
	case dns.TypeEUI64:
		return "EUI64"
	case dns.TypeURI:
		return "URI"
	case dns.TypeAVC:
		return "AVC"
	default:
		return fmt.Sprintf("TYPE%d", qtype)
	}
}

// DNSRecord DNS记录
type DNSRecord struct {
	Type  string
	Name  string
	TTL   uint32
	Value string
}

// ProbeResult 拨测结果
type ProbeResult struct {
	Domain    string
	DNSServer string
	Records   []DNSRecord
	Latency   time.Duration
	Error     error
}

// ProbeDNS 拨测单个DNS服务器（只查询A记录）
func ProbeDNS(domain, dnsServer string) ProbeResult {
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	result := ProbeResult{
		Domain:    domain,
		DNSServer: dnsServer,
	}

	start := time.Now()

	// 只查询A记录
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, net.JoinHostPort(dnsServer, "53"))
	if err != nil {
		result.Error = err
		result.Latency = time.Since(start)
		return result
	}

	if r.Rcode != dns.RcodeSuccess {
		result.Error = fmt.Errorf("DNS query failed: %s", dns.RcodeToString[r.Rcode])
		result.Latency = time.Since(start)
		return result
	}

	// 解析A记录结果
	for _, answer := range r.Answer {
		record := parseRecord(answer)
		result.Records = append(result.Records, record)
	}

	result.Latency = time.Since(start)
	return result
}

// ProbeDNSAll 拨测单个DNS服务器（查询所有记录类型）
func ProbeDNSAll(domain, dnsServer string) ProbeResult {
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	result := ProbeResult{
		Domain:    domain,
		DNSServer: dnsServer,
	}

	start := time.Now()

	// 先查询A记录测试连通性
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, net.JoinHostPort(dnsServer, "53"))
	if err != nil {
		result.Error = err
		result.Latency = time.Since(start)
		return result
	}

	if r.Rcode != dns.RcodeSuccess {
		result.Error = fmt.Errorf("DNS query failed: %s", dns.RcodeToString[r.Rcode])
		result.Latency = time.Since(start)
		return result
	}

	// 解析A记录结果
	for _, answer := range r.Answer {
		record := parseRecord(answer)
		result.Records = append(result.Records, record)
	}

	// 逐个查询其他记录类型
	otherTypes := []uint16{
		dns.TypeAAAA,
		dns.TypeCNAME,
		dns.TypeMX,
		dns.TypeNS,
		dns.TypeTXT,
		dns.TypeSOA,
		dns.TypeSRV,
		dns.TypeCAA,
		dns.TypePTR,
	}

	for _, qtype := range otherTypes {
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(domain), qtype)
		m.RecursionDesired = true

		r, _, err := c.Exchange(m, net.JoinHostPort(dnsServer, "53"))
		if err != nil {
			continue
		}

		if r.Rcode != dns.RcodeSuccess {
			continue
		}

		for _, answer := range r.Answer {
			record := parseRecord(answer)
			result.Records = append(result.Records, record)
		}
	}

	result.Latency = time.Since(start)
	return result
}

// parseRecord 解析DNS记录
func parseRecord(answer dns.RR) DNSRecord {
	record := DNSRecord{
		Type: RecordTypeName(answer.Header().Rrtype),
		Name: answer.Header().Name,
		TTL:  answer.Header().Ttl,
	}

	switch v := answer.(type) {
	case *dns.A:
		record.Value = v.A.String()
	case *dns.AAAA:
		record.Value = v.AAAA.String()
	case *dns.CNAME:
		record.Value = v.Target
	case *dns.MX:
		record.Value = fmt.Sprintf("%d %s", v.Preference, v.Mx)
	case *dns.NS:
		record.Value = v.Ns
	case *dns.TXT:
		record.Value = strings.Join(v.Txt, " ")
	case *dns.SOA:
		record.Value = fmt.Sprintf("%s %s %d %d %d %d %d",
			v.Ns, v.Mbox, v.Serial, v.Refresh, v.Retry, v.Expire, v.Minttl)
	case *dns.SRV:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Priority, v.Weight, v.Port, v.Target)
	case *dns.CAA:
		record.Value = fmt.Sprintf("%d %s \"%s\"", v.Flag, v.Tag, v.Value)
	case *dns.PTR:
		record.Value = v.Ptr
	case *dns.DNSKEY:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Flags, v.Protocol, v.Algorithm, v.PublicKey)
	case *dns.DS:
		record.Value = fmt.Sprintf("%d %d %d %s", v.KeyTag, v.Algorithm, v.DigestType, v.Digest)
	case *dns.NSEC:
		record.Value = fmt.Sprintf("%s %s", v.NextDomain, v.TypeBitMap)
	case *dns.NSEC3:
		record.Value = fmt.Sprintf("%d %d %d %s %s", v.Hash, v.Flags, v.Iterations, v.Salt, v.NextDomain)
	case *dns.RRSIG:
		record.Value = fmt.Sprintf("%d %d %d %d %d %d %s",
			v.TypeCovered, v.Algorithm, v.Labels, v.OrigTtl,
			v.Expiration, v.Inception, v.SignerName)
	case *dns.HINFO:
		record.Value = fmt.Sprintf("%s %s", v.Cpu, v.Os)
	case *dns.RP:
		record.Value = fmt.Sprintf("%s %s", v.Mbox, v.Txt)
	case *dns.AFSDB:
		record.Value = fmt.Sprintf("%d %s", v.Subtype, v.Hostname)
	case *dns.X25:
		record.Value = v.PSDNAddress
	case *dns.ISDN:
		record.Value = v.Address
	case *dns.RT:
		record.Value = fmt.Sprintf("%d %s", v.Preference, v.Host)
	case *dns.SIG:
		record.Value = fmt.Sprintf("%d %d %d %d %d %d %s",
			v.TypeCovered, v.Algorithm, v.Labels, v.OrigTtl,
			v.Expiration, v.Inception, v.SignerName)
	case *dns.KEY:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Flags, v.Protocol, v.Algorithm, v.PublicKey)
	case *dns.PX:
		record.Value = fmt.Sprintf("%d %s %s", v.Preference, v.Map822, v.Mapx400)
	case *dns.GPOS:
		record.Value = fmt.Sprintf("%s %s %s", v.Longitude, v.Latitude, v.Altitude)
	case *dns.LOC:
		record.Value = fmt.Sprintf("%d %d %d %d %d %d %d",
			v.Version, v.Size, v.HorizPre, v.VertPre,
			v.Latitude, v.Longitude, v.Altitude)
	case *dns.NXT:
		record.Value = fmt.Sprintf("%s %s", v.NextDomain, v.TypeBitMap)
	case *dns.NAPTR:
		record.Value = fmt.Sprintf("%d %d \"%s\" \"%s\" \"%s\" %s",
			v.Order, v.Preference, v.Flags, v.Service, v.Regexp, v.Replacement)
	case *dns.KX:
		record.Value = fmt.Sprintf("%d %s", v.Preference, v.Exchanger)
	case *dns.CERT:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Type, v.KeyTag, v.Algorithm, v.Certificate)
	case *dns.DNAME:
		record.Value = v.Target
	case *dns.OPT:
		record.Value = fmt.Sprintf("UDP size: %d", v.UDPSize())
	case *dns.APL:
		record.Value = fmt.Sprintf("%v", v.Prefixes)
	case *dns.SSHFP:
		record.Value = fmt.Sprintf("%d %d %s", v.Algorithm, v.Type, v.FingerPrint)
	case *dns.IPSECKEY:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Precedence, v.GatewayType, v.Algorithm, v.PublicKey)
	case *dns.DHCID:
		record.Value = v.Digest
	case *dns.NSEC3PARAM:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Hash, v.Flags, v.Iterations, v.Salt)
	case *dns.TLSA:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Usage, v.Selector, v.MatchingType, v.Certificate)
	case *dns.SMIMEA:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Usage, v.Selector, v.MatchingType, v.Certificate)
	case *dns.HIP:
		record.Value = fmt.Sprintf("%d %s %s", v.PublicKeyAlgorithm, v.Hit, v.PublicKey)
	case *dns.CDS:
		record.Value = fmt.Sprintf("%d %d %d %s", v.KeyTag, v.Algorithm, v.DigestType, v.Digest)
	case *dns.CDNSKEY:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Flags, v.Protocol, v.Algorithm, v.PublicKey)
	case *dns.OPENPGPKEY:
		record.Value = v.PublicKey
	case *dns.CSYNC:
		record.Value = fmt.Sprintf("%d %d %s", v.Serial, v.Flags, v.TypeBitMap)
	case *dns.ZONEMD:
		record.Value = fmt.Sprintf("%d %d %d %s", v.Scheme, v.Hash, v.Serial, v.Digest)
	case *dns.SVCB:
		record.Value = fmt.Sprintf("%d %s %s", v.Priority, v.Target, v.Value)
	case *dns.HTTPS:
		record.Value = fmt.Sprintf("%d %s %s", v.Priority, v.Target, v.Value)
	case *dns.SPF:
		record.Value = strings.Join(v.Txt, " ")
	case *dns.UINFO:
		record.Value = v.Uinfo
	case *dns.UID:
		record.Value = fmt.Sprintf("%d", v.Uid)
	case *dns.GID:
		record.Value = fmt.Sprintf("%d", v.Gid)
	case *dns.NID:
		record.Value = fmt.Sprintf("%d %d", v.Preference, v.NodeID)
	case *dns.L32:
		record.Value = fmt.Sprintf("%d %s", v.Preference, v.Locator32)
	case *dns.L64:
		record.Value = fmt.Sprintf("%d %d", v.Preference, v.Locator64)
	case *dns.LP:
		record.Value = fmt.Sprintf("%d %s", v.Preference, v.Fqdn)
	case *dns.EUI48:
		record.Value = v.String()
	case *dns.EUI64:
		record.Value = v.String()
	case *dns.URI:
		record.Value = fmt.Sprintf("%d %d \"%s\"", v.Priority, v.Weight, v.Target)
	case *dns.AVC:
		record.Value = strings.Join(v.Txt, " ")
	default:
		record.Value = answer.String()
	}

	return record
}

// ProbeAll 并发拨测所有DNS服务器（只查询A记录）
func ProbeAll(domain string, dnsServers []string) []ProbeResult {
	var wg sync.WaitGroup
	results := make([]ProbeResult, len(dnsServers))

	for i, server := range dnsServers {
		wg.Add(1)
		go func(idx int, srv string) {
			defer wg.Done()
			results[idx] = ProbeDNS(domain, srv)
		}(i, server)
	}

	wg.Wait()
	return results
}

// ProbeAllRecordTypes 并发拨测所有DNS服务器（查询所有记录类型）
func ProbeAllRecordTypes(domain string, dnsServers []string) []ProbeResult {
	var wg sync.WaitGroup
	results := make([]ProbeResult, len(dnsServers))

	for i, server := range dnsServers {
		wg.Add(1)
		go func(idx int, srv string) {
			defer wg.Done()
			results[idx] = ProbeDNSAll(domain, srv)
		}(i, server)
	}

	wg.Wait()
	return results
}

// FormatText 格式化输出（美化版）
func FormatText(results []ProbeResult) string {
	var sb strings.Builder

	// 标题
	sb.WriteString("╔══════════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║                    DNS Probe Tool v1.0                          ║\n")
	sb.WriteString("╚══════════════════════════════════════════════════════════════════╝\n\n")

	sb.WriteString(fmt.Sprintf("  域名: %s\n", results[0].Domain))
	sb.WriteString(fmt.Sprintf("  时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// 按服务器分组显示
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("┌─ DNS服务器: %s\n", r.DNSServer))
		sb.WriteString(fmt.Sprintf("│  查询耗时: %d ms\n", r.Latency.Milliseconds()))

		if r.Error != nil {
			sb.WriteString(fmt.Sprintf("│  ❌ 错误: %s\n", r.Error.Error()))
			sb.WriteString("└──────────────────────────────────────────────────────────────────\n\n")
			continue
		}

		if len(r.Records) == 0 {
			sb.WriteString("│  (无记录)\n")
			sb.WriteString("└──────────────────────────────────────────────────────────────────\n\n")
			continue
		}

		// 按记录类型排序
		sort.Slice(r.Records, func(i, j int) bool {
			return r.Records[i].Type < r.Records[j].Type
		})

		sb.WriteString("│\n")
		sb.WriteString(fmt.Sprintf("│  %-10s %-6s %s\n", "类型", "TTL", "值"))
		sb.WriteString("│  " + strings.Repeat("─", 60) + "\n")

		for _, rec := range r.Records {
			ttlStr := formatTTL(rec.TTL)
			sb.WriteString(fmt.Sprintf("│  %-10s %-6s %s\n", rec.Type, ttlStr, rec.Value))
		}

		sb.WriteString("└──────────────────────────────────────────────────────────────────\n\n")
	}

	return sb.String()
}

// formatTTL 格式化TTL
func formatTTL(ttl uint32) string {
	if ttl < 60 {
		return fmt.Sprintf("%ds", ttl)
	} else if ttl < 3600 {
		return fmt.Sprintf("%dm", ttl/60)
	} else if ttl < 86400 {
		return fmt.Sprintf("%dh", ttl/3600)
	} else {
		return fmt.Sprintf("%dd", ttl/86400)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "用法: dns-probe <域名> [DNS服务器...] [--all]\n")
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com                    # 使用系统DNS服务器查询A记录\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com 8.8.8.8            # 使用指定DNS服务器查询A记录\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com 8.8.8.8 114.114.114.114  # 使用多个DNS服务器查询A记录\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --all              # 使用系统DNS服务器查询所有记录类型\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com 8.8.8.8 --all      # 使用指定DNS服务器查询所有记录类型\n")
		fmt.Fprintf(os.Stderr, "\n系统DNS服务器:\n")
		for _, s := range DefaultDNSServers {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
		os.Exit(1)
	}

	domain := os.Args[1]
	dnsServers := DefaultDNSServers
	queryAll := false

	// 解析参数
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--all" {
			queryAll = true
		} else {
			// 如果不是--all参数，则认为是DNS服务器
			if i == 2 {
				dnsServers = []string{}
			}
			dnsServers = append(dnsServers, os.Args[i])
		}
	}

	var results []ProbeResult
	if queryAll {
		results = ProbeAllRecordTypes(domain, dnsServers)
	} else {
		results = ProbeAll(domain, dnsServers)
	}

	fmt.Print(FormatText(results))
}
