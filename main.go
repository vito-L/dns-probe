package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// 样式定义
var (
	warningStyle = &styleFunc{color: "11"}
	successStyle = &styleFunc{color: "10"}
	errorStyle   = &styleFunc{color: "9"}
	infoStyle    = &styleFunc{color: "12"}
)

type styleFunc struct {
	color string
}

func (s *styleFunc) Render(text string) string {
	return text
}

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
	Domain      string
	DNSServer   string
	Records     []DNSRecord
	Latency     time.Duration
	Error       error
	IsPolluted  bool
	PollutedIP  string
	RealIP      string
	DNSSECValid bool
}

// ProbeDNS 拨测单个DNS服务器（只查询A记录）
func ProbeDNS(domain, dnsServer string, checkDNSSEC ...bool) ProbeResult {
	// 检查是否是DoH服务器
	if strings.HasPrefix(dnsServer, "https://") {
		return probeDoH(domain, dnsServer)
	}

	// 检查是否是DoT服务器
	if strings.Contains(dnsServer, ":853") {
		return probeDoT(domain, dnsServer)
	}

	// 是否检查DNSSEC
	enableDNSSEC := false
	if len(checkDNSSEC) > 0 && checkDNSSEC[0] {
		enableDNSSEC = true
	}

	// 普通DNS查询
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	result := ProbeResult{
		Domain:    domain,
		DNSServer: dnsServer,
	}

	start := time.Now()

	// 查询A记录
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true
	m.CheckingDisabled = false

	// 只有在启用DNSSEC时才设置DO位
	if enableDNSSEC {
		m.SetEdns0(4096, true)
	}

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

	// 只有在启用DNSSEC时才检查DNSSEC状态
	if enableDNSSEC {
		if r.AuthenticatedData {
			result.DNSSECValid = true
		}

		// 检查是否有RRSIG记录
		for _, answer := range r.Answer {
			if _, ok := answer.(*dns.RRSIG); ok {
				result.DNSSECValid = true
				break
			}
		}
	}

	// 解析A记录结果
	for _, answer := range r.Answer {
		if _, ok := answer.(*dns.RRSIG); ok {
			continue // 跳过RRSIG记录
		}
		record := parseRecord(answer)
		result.Records = append(result.Records, record)
	}

	result.Latency = time.Since(start)
	return result
}

// probeDoH DNS over HTTPS查询
func probeDoH(domain, dohServer string) ProbeResult {
	result := ProbeResult{
		Domain:    domain,
		DNSServer: dohServer,
	}

	start := time.Now()

	// 构建DNS查询，设置DO位
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true
	m.SetEdns0(4096, true) // 设置DO位

	pack, err := m.Pack()
	if err != nil {
		result.Error = fmt.Errorf("failed to pack DNS message: %s", err.Error())
		result.Latency = time.Since(start)
		return result
	}

	// 获取代理设置
	proxyURL := getProxyURL()

	// 发送DoH请求
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	// 如果有代理，设置代理
	if proxyURL != nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	req, err := http.NewRequest("POST", dohServer, bytes.NewReader(pack))
	if err != nil {
		result.Error = fmt.Errorf("failed to create request: %s", err.Error())
		result.Latency = time.Since(start)
		return result
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := client.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("failed to send DoH request: %s", err.Error())
		result.Latency = time.Since(start)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Errorf("DoH request failed with status: %d", resp.StatusCode)
		result.Latency = time.Since(start)
		return result
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Errorf("failed to read response body: %s", err.Error())
		result.Latency = time.Since(start)
		return result
	}

	// 解析DNS响应
	r := new(dns.Msg)
	if err := r.Unpack(body); err != nil {
		result.Error = fmt.Errorf("failed to unpack DNS response: %s", err.Error())
		result.Latency = time.Since(start)
		return result
	}

	if r.Rcode != dns.RcodeSuccess {
		result.Error = fmt.Errorf("DNS query failed: %s", dns.RcodeToString[r.Rcode])
		result.Latency = time.Since(start)
		return result
	}

	// 检查DNSSEC验证状态
	if r.AuthenticatedData {
		result.DNSSECValid = true
	}

	// 检查是否有RRSIG记录
	for _, answer := range r.Answer {
		if _, ok := answer.(*dns.RRSIG); ok {
			result.DNSSECValid = true
			break
		}
	}

	// 解析A记录结果
	for _, answer := range r.Answer {
		if _, ok := answer.(*dns.RRSIG); ok {
			continue // 跳过RRSIG记录
		}
		record := parseRecord(answer)
		result.Records = append(result.Records, record)
	}

	result.Latency = time.Since(start)
	return result
}

// getProxyURL 获取代理URL
func getProxyURL() *url.URL {
	// 检查环境变量中的代理设置
	proxyEnv := os.Getenv("https_proxy")
	if proxyEnv == "" {
		proxyEnv = os.Getenv("http_proxy")
	}
	if proxyEnv == "" {
		proxyEnv = os.Getenv("HTTPS_PROXY")
	}
	if proxyEnv == "" {
		proxyEnv = os.Getenv("HTTP_PROXY")
	}

	if proxyEnv != "" {
		proxyURL, err := url.Parse(proxyEnv)
		if err == nil {
			return proxyURL
		}
	}

	return nil
}

// probeDoT DNS over TLS查询
func probeDoT(domain, dotServer string) ProbeResult {
	result := ProbeResult{
		Domain:    domain,
		DNSServer: dotServer,
	}

	start := time.Now()

	// 构建DNS查询，设置DO位
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true
	m.SetEdns0(4096, true) // 设置DO位

	// 获取代理设置
	proxyURL := getProxyURL()

	// 使用DNS客户端
	c := new(dns.Client)
	c.Net = "tcp-tls"
	c.Timeout = 30 * time.Second
	c.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         strings.Split(dotServer, ":")[0],
	}

	// 如果有代理，使用代理连接
	if proxyURL != nil {
		// 通过代理建立TCP连接
		conn, err := net.DialTimeout("tcp", proxyURL.Host, 30*time.Second)
		if err != nil {
			result.Error = fmt.Errorf("failed to connect to proxy: %s", err.Error())
			result.Latency = time.Since(start)
			return result
		}

		// 发送CONNECT请求
		connectReq := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", dotServer, dotServer)
		_, err = conn.Write([]byte(connectReq))
		if err != nil {
			result.Error = fmt.Errorf("failed to send CONNECT request: %s", err.Error())
			result.Latency = time.Since(start)
			return result
		}

		// 读取响应
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			result.Error = fmt.Errorf("failed to read CONNECT response: %s", err.Error())
			result.Latency = time.Since(start)
			return result
		}

		response := string(buf[:n])
		if !strings.Contains(response, "200") {
			result.Error = fmt.Errorf("CONNECT request failed: %s", response)
			result.Latency = time.Since(start)
			return result
		}

		// 建立TLS连接
		tlsConn := tls.Client(conn, c.TLSConfig)
		err = tlsConn.Handshake()
		if err != nil {
			result.Error = fmt.Errorf("TLS handshake failed: %s", err.Error())
			result.Latency = time.Since(start)
			return result
		}

		// 使用TLS连接发送DNS查询
		dnsConn := &dns.Conn{Conn: tlsConn}
		r, _, err := c.ExchangeWithConn(m, dnsConn)
		if err != nil {
			result.Error = fmt.Errorf("failed to query DoT server: %s", err.Error())
			result.Latency = time.Since(start)
			return result
		}

		if r.Rcode != dns.RcodeSuccess {
			result.Error = fmt.Errorf("DNS query failed: %s", dns.RcodeToString[r.Rcode])
			result.Latency = time.Since(start)
			return result
		}

		// 检查DNSSEC验证状态
		if r.AuthenticatedData {
			result.DNSSECValid = true
		}

		// 检查是否有RRSIG记录
		for _, answer := range r.Answer {
			if _, ok := answer.(*dns.RRSIG); ok {
				result.DNSSECValid = true
				break
			}
		}

		// 解析A记录结果
		for _, answer := range r.Answer {
			if _, ok := answer.(*dns.RRSIG); ok {
				continue // 跳过RRSIG记录
			}
			record := parseRecord(answer)
			result.Records = append(result.Records, record)
		}
	} else {
		// 直接连接
		r, _, err := c.Exchange(m, dotServer)
		if err != nil {
			result.Error = fmt.Errorf("failed to query DoT server: %s", err.Error())
			result.Latency = time.Since(start)
			return result
		}

		if r.Rcode != dns.RcodeSuccess {
			result.Error = fmt.Errorf("DNS query failed: %s", dns.RcodeToString[r.Rcode])
			result.Latency = time.Since(start)
			return result
		}

		// 检查DNSSEC验证状态
		if r.AuthenticatedData {
			result.DNSSECValid = true
		}

		// 检查是否有RRSIG记录
		for _, answer := range r.Answer {
			if _, ok := answer.(*dns.RRSIG); ok {
				result.DNSSECValid = true
				break
			}
		}

		// 解析A记录结果
		for _, answer := range r.Answer {
			if _, ok := answer.(*dns.RRSIG); ok {
				continue // 跳过RRSIG记录
			}
			record := parseRecord(answer)
			result.Records = append(result.Records, record)
		}
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
		record.Value = fmt.Sprintf("%s %v", v.NextDomain, v.TypeBitMap)
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
		record.Value = fmt.Sprintf("%s %v", v.NextDomain, v.TypeBitMap)
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
		record.Value = fmt.Sprintf("%d %d %v", v.Serial, v.Flags, v.TypeBitMap)
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
func ProbeAll(domain string, dnsServers []string, enableDNSSEC ...bool) []ProbeResult {
	var wg sync.WaitGroup
	results := make([]ProbeResult, len(dnsServers))

	for i, server := range dnsServers {
		wg.Add(1)
		go func(idx int, srv string) {
			defer wg.Done()
			results[idx] = ProbeDNS(domain, srv, enableDNSSEC...)
		}(i, server)
	}

	wg.Wait()
	return results
}

// ProbeAllWithPollutionCheck 并发拨测所有DNS服务器（带污染检测）
func ProbeAllWithPollutionCheck(domain string, dnsServers []string) []ProbeResult {
	// 先用国外DNS查询作为基准
国外DNS := []string{"8.8.8.8", "1.1.1.1"}
	benchmarkResults := ProbeAll(domain, 国外DNS)

	// 获取基准IP列表
	benchmarkIPs := make(map[string]bool)
	var realIP string
	for _, r := range benchmarkResults {
		if r.Error == nil {
			for _, rec := range r.Records {
				if rec.Type == "A" {
					benchmarkIPs[rec.Value] = true
					if realIP == "" {
						realIP = rec.Value
					}
				}
			}
		}
	}

	// 查询所有DNS服务器
	var wg sync.WaitGroup
	results := make([]ProbeResult, len(dnsServers))

	for i, server := range dnsServers {
		wg.Add(1)
		go func(idx int, srv string) {
			defer wg.Done()
			result := ProbeDNS(domain, srv)

			// 检测污染
			if result.Error == nil && len(benchmarkIPs) > 0 {
				for _, rec := range result.Records {
					if rec.Type == "A" && !benchmarkIPs[rec.Value] {
						result.IsPolluted = true
						result.PollutedIP = rec.Value
						result.RealIP = realIP
						break
					}
				}
			}

			results[idx] = result
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

		if r.IsPolluted {
			sb.WriteString(warningStyle.Render("│  ⚠️  检测到DNS污染"))
			sb.WriteString("\n")
			sb.WriteString(warningStyle.Render(fmt.Sprintf("│  被污染的IP: %s", r.PollutedIP)))
			sb.WriteString("\n")
			sb.WriteString(warningStyle.Render(fmt.Sprintf("│  真实IP（国外DNS）: %s", r.RealIP)))
			sb.WriteString("\n")
		}

		if r.DNSSECValid {
			sb.WriteString(successStyle.Render("│  ✅ DNSSEC验证通过"))
			sb.WriteString("\n")
		}

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

// JSONResult JSON输出格式
type JSONResult struct {
	Domain    string        `json:"domain"`
	Timestamp string        `json:"timestamp"`
	Servers   []JSONServer  `json:"servers"`
}

// JSONServer JSON服务器结果
type JSONServer struct {
	Server  string       `json:"server"`
	Latency int64        `json:"latency_ms"`
	Error   string       `json:"error,omitempty"`
	Records []JSONRecord `json:"records,omitempty"`
}

// JSONRecord JSON记录
type JSONRecord struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	TTL   uint32 `json:"ttl"`
	Value string `json:"value"`
}

// FormatJSON 格式化JSON输出
func FormatJSON(results []ProbeResult) string {
	jsonResult := JSONResult{
		Domain:    results[0].Domain,
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	for _, r := range results {
		server := JSONServer{
			Server:  r.DNSServer,
			Latency: r.Latency.Milliseconds(),
		}

		if r.Error != nil {
			server.Error = r.Error.Error()
		}

		for _, rec := range r.Records {
			server.Records = append(server.Records, JSONRecord{
				Type:  rec.Type,
				Name:  rec.Name,
				TTL:   rec.TTL,
				Value: rec.Value,
			})
		}

		jsonResult.Servers = append(jsonResult.Servers, server)
	}

	jsonBytes, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal JSON: %s"}`, err.Error())
	}

	return string(jsonBytes)
}

// ReadDomainsFromFile 从文件读取域名列表
func ReadDomainsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			domains = append(domains, line)
		}
	}

	return domains, scanner.Err()
}

// HistoryEntry 历史记录条目
type HistoryEntry struct {
	Timestamp string        `json:"timestamp"`
	Domain    string        `json:"domain"`
	Servers   []JSONServer  `json:"servers"`
}

// SaveHistory 保存历史记录
func SaveHistory(domain string, results []ProbeResult) error {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// 创建历史记录目录
	historyDir := filepath.Join(homeDir, ".dns-probe")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return err
	}

	// 读取现有历史记录
	historyFile := filepath.Join(historyDir, "history.json")
	var history []HistoryEntry

	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &history)
	}

	// 创建新的历史记录条目
	entry := HistoryEntry{
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Domain:    domain,
	}

	for _, r := range results {
		server := JSONServer{
			Server:  r.DNSServer,
			Latency: r.Latency.Milliseconds(),
		}

		if r.Error != nil {
			server.Error = r.Error.Error()
		}

		for _, rec := range r.Records {
			server.Records = append(server.Records, JSONRecord{
				Type:  rec.Type,
				Name:  rec.Name,
				TTL:   rec.TTL,
				Value: rec.Value,
			})
		}

		entry.Servers = append(entry.Servers, server)
	}

	// 添加到历史记录
	history = append(history, entry)

	// 保存历史记录
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// LoadHistory 加载历史记录
func LoadHistory() ([]HistoryEntry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFile := filepath.Join(homeDir, ".dns-probe", "history.json")
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return nil, err
	}

	var history []HistoryEntry
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, err
	}

	return history, nil
}

// FormatHistory 格式化历史记录输出
func FormatHistory(history []HistoryEntry) string {
	var sb strings.Builder

	sb.WriteString("查询历史:\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	for i, entry := range history {
		sb.WriteString(fmt.Sprintf("[%d] %s - %s\n", i+1, entry.Timestamp, entry.Domain))
		for _, server := range entry.Servers {
			sb.WriteString(fmt.Sprintf("    DNS: %s, 耗时: %d ms\n", server.Server, server.Latency))
		}
	}

	return sb.String()
}

// ProbeMultipleDomains 批量拨测多个域名
func ProbeMultipleDomains(domains []string, dnsServers []string, queryAll bool) []JSONResult {
	var results []JSONResult

	for _, domain := range domains {
		var probeResults []ProbeResult
		if queryAll {
			probeResults = ProbeAllRecordTypes(domain, dnsServers)
		} else {
			probeResults = ProbeAll(domain, dnsServers)
		}

		jsonResult := JSONResult{
			Domain:    domain,
			Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		}

		for _, r := range probeResults {
			server := JSONServer{
				Server:  r.DNSServer,
				Latency: r.Latency.Milliseconds(),
			}

			if r.Error != nil {
				server.Error = r.Error.Error()
			}

			for _, rec := range r.Records {
				server.Records = append(server.Records, JSONRecord{
					Type:  rec.Type,
					Name:  rec.Name,
					TTL:   rec.TTL,
					Value: rec.Value,
				})
			}

			jsonResult.Servers = append(jsonResult.Servers, server)
		}

		results = append(results, jsonResult)
	}

	return results
}

// FormatMultipleJSON 格式化多个域名的JSON输出
func FormatMultipleJSON(results []JSONResult) string {
	jsonBytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to marshal JSON: %s"}`, err.Error())
	}

	return string(jsonBytes)
}

// FormatHTML 格式化HTML输出
func FormatHTML(results []ProbeResult) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DNS Probe Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 30px;
        }
        h1 {
            color: #333;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        .info {
            color: #666;
            margin-bottom: 20px;
        }
        .server-card {
            background: #f8f9fa;
            border-radius: 6px;
            padding: 15px;
            margin-bottom: 15px;
            border-left: 4px solid #667eea;
        }
        .server-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        .server-name {
            font-weight: bold;
            color: #333;
        }
        .latency {
            color: #666;
            font-size: 14px;
        }
        .pollution {
            background: #fff3cd;
            color: #856404;
            padding: 5px 10px;
            border-radius: 4px;
            font-size: 14px;
            margin-bottom: 10px;
        }
        .error {
            background: #f8d7da;
            color: #721c24;
            padding: 5px 10px;
            border-radius: 4px;
            font-size: 14px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 10px;
        }
        th, td {
            padding: 8px 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background: #667eea;
            color: white;
        }
        tr:hover {
            background: #f5f5f5;
        }
        .type {
            font-weight: bold;
            color: #667eea;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>DNS Probe Report</h1>
        <div class="info">
            <p>域名: ` + results[0].Domain + `</p>
            <p>时间: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
        </div>
`)

	for _, r := range results {
		sb.WriteString(`        <div class="server-card">
            <div class="server-header">
                <span class="server-name">` + r.DNSServer + `</span>
                <span class="latency">` + fmt.Sprintf("%d ms", r.Latency.Milliseconds()) + `</span>
            </div>
`)

		if r.IsPolluted {
			sb.WriteString(`            <div class="pollution">⚠️ 检测到DNS污染</div>
`)
		}

		if r.Error != nil {
			sb.WriteString(`            <div class="error">❌ 错误: ` + r.Error.Error() + `</div>
`)
		} else if len(r.Records) > 0 {
			sb.WriteString(`            <table>
                <tr>
                    <th>类型</th>
                    <th>TTL</th>
                    <th>值</th>
                </tr>
`)
			for _, rec := range r.Records {
				ttlStr := formatTTL(rec.TTL)
				sb.WriteString(`                <tr>
                    <td class="type">` + rec.Type + `</td>
                    <td>` + ttlStr + `</td>
                    <td>` + rec.Value + `</td>
                </tr>
`)
			}
			sb.WriteString(`            </table>
`)
		} else {
			sb.WriteString(`            <p>(无记录)</p>
`)
		}

		sb.WriteString(`        </div>
`)
	}

	sb.WriteString(`    </div>
</body>
</html>`)

	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "用法: dns-probe <域名> [DNS服务器...] [--all] [--json] [--pollution] [--dnssec] [--doh <url>] [--dot <server>] [--html <文件>] [--file <文件>] [--history]\n")
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com                            # 使用系统DNS服务器查询A记录\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com 8.8.8.8                    # 使用指定DNS服务器查询A记录\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com 8.8.8.8 114.114.114.114    # 使用多个DNS服务器查询A记录\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --all                      # 查询所有记录类型\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --json                     # 输出JSON格式\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --pollution                # 检测DNS污染\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --dnssec                   # 使用系统DNS服务器进行DNSSEC验证\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --dnssec 8.8.8.8           # 使用指定DNS服务器进行DNSSEC验证\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --doh https://dns.google/dns-query  # 使用DoH服务器\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --dot dns.alidns.com:853   # 使用DoT服务器\n")
		fmt.Fprintf(os.Stderr, "  dns-probe example.com --html report.html         # 生成HTML报告\n")
		fmt.Fprintf(os.Stderr, "  dns-probe --file domains.txt                     # 批量查询文件中的域名\n")
		fmt.Fprintf(os.Stderr, "  dns-probe --file domains.txt --json              # 批量查询并输出JSON格式\n")
		fmt.Fprintf(os.Stderr, "  dns-probe --history                              # 显示查询历史\n")
		fmt.Fprintf(os.Stderr, "  dns-probe --history                      # 显示查询历史\n")
		fmt.Fprintf(os.Stderr, "\n系统DNS服务器:\n")
		for _, s := range DefaultDNSServers {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
		os.Exit(1)
	}

	domain := ""
	dnsServers := []string{} // 初始化为空，后面根据情况填充
	useDefaultDNS := true   // 是否使用默认DNS服务器
	queryAll := false
	outputJSON := false
	filePath := ""
	checkPollution := false
	enableDNSSEC := false
	dohServer := ""
	dotServer := ""
	htmlFile := ""
	showHistory := false

	// 解析参数
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--all" {
			queryAll = true
		} else if os.Args[i] == "--json" {
			outputJSON = true
		} else if os.Args[i] == "--pollution" {
			checkPollution = true
		} else if os.Args[i] == "--dnssec" {
			enableDNSSEC = true
		} else if os.Args[i] == "--doh" {
			if i+1 < len(os.Args) {
				dohServer = os.Args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "错误: --doh 参数需要指定URL\n")
				os.Exit(1)
			}
		} else if os.Args[i] == "--dot" {
			if i+1 < len(os.Args) {
				dotServer = os.Args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "错误: --dot 参数需要指定服务器\n")
				os.Exit(1)
			}
		} else if os.Args[i] == "--history" {
			showHistory = true
		} else if os.Args[i] == "--html" {
			if i+1 < len(os.Args) {
				htmlFile = os.Args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "错误: --html 参数需要指定文件名\n")
				os.Exit(1)
			}
		} else if os.Args[i] == "--file" {
			if i+1 < len(os.Args) {
				filePath = os.Args[i+1]
				i++
			} else {
				fmt.Fprintf(os.Stderr, "错误: --file 参数需要指定文件名\n")
				os.Exit(1)
			}
		} else if i == 1 && !strings.HasPrefix(os.Args[i], "--") {
			domain = os.Args[i]
		} else {
			// 用户指定了DNS服务器，不再使用默认DNS
			if useDefaultDNS {
				dnsServers = []string{}
				useDefaultDNS = false
			}
			dnsServers = append(dnsServers, os.Args[i])
		}
	}

	// 如果用户没有指定DNS服务器，使用默认DNS
	if useDefaultDNS {
		dnsServers = DefaultDNSServers
	}

	// 显示历史记录
	if showHistory {
		history, err := LoadHistory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 读取历史记录失败: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Print(FormatHistory(history))
		return
	}

	// 批量查询模式
	if filePath != "" {
		domains, err := ReadDomainsFromFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 读取文件失败: %s\n", err.Error())
			os.Exit(1)
		}

		if len(domains) == 0 {
			fmt.Fprintf(os.Stderr, "错误: 文件中没有域名\n")
			os.Exit(1)
		}

		results := ProbeMultipleDomains(domains, dnsServers, queryAll)
		if outputJSON {
			fmt.Println(FormatMultipleJSON(results))
		} else {
			for _, result := range results {
				fmt.Printf("域名: %s\n", result.Domain)
				for _, server := range result.Servers {
					fmt.Printf("  DNS服务器: %s\n", server.Server)
					fmt.Printf("  查询耗时: %d ms\n", server.Latency)
					if server.Error != "" {
						fmt.Printf("  错误: %s\n", server.Error)
					} else {
						for _, rec := range server.Records {
							fmt.Printf("    %-10s %-6d %s\n", rec.Type, rec.TTL, rec.Value)
						}
					}
				}
				fmt.Println()
			}
		}
		return
	}

	// 单域名查询模式
	if domain == "" {
		fmt.Fprintf(os.Stderr, "错误: 请指定域名或使用 --file 参数\n")
		os.Exit(1)
	}

	// 如果指定了DoH或DoT服务器，添加到DNS服务器列表
	if dohServer != "" {
		dnsServers = []string{dohServer}
	}
	if dotServer != "" {
		dnsServers = []string{dotServer}
	}

	// 如果启用DNSSEC，设置DO位
	if enableDNSSEC {
		// DNSSEC验证会在ProbeDNS函数中自动处理
		// 通过SetEdns0设置DO位
	}

	var results []ProbeResult
	if checkPollution {
		results = ProbeAllWithPollutionCheck(domain, dnsServers)
	} else if queryAll {
		results = ProbeAllRecordTypes(domain, dnsServers)
	} else {
		results = ProbeAll(domain, dnsServers, enableDNSSEC)
	}

	// 保存历史记录
	if err := SaveHistory(domain, results); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 保存历史记录失败: %s\n", err.Error())
	}

	// HTML报告
	if htmlFile != "" {
		html := FormatHTML(results)
		err := os.WriteFile(htmlFile, []byte(html), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 写入HTML文件失败: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("HTML报告已生成: %s\n", htmlFile)
		return
	}

	if outputJSON {
		fmt.Println(FormatJSON(results))
	} else {
		fmt.Print(FormatText(results))
	}
}
