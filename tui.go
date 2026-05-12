package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TUI样式
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12"))
)

// TUIModel TUI模型
type TUIModel struct {
	domain     string
	dnsServers []string
	queryAll   bool
	results    []ProbeResult
	quitting   bool
	loading    bool
}

// NewTUIModel 创建TUI模型
func NewTUIModel(domain string, dnsServers []string, queryAll bool) TUIModel {
	return TUIModel{
		domain:     domain,
		dnsServers: dnsServers,
		queryAll:   queryAll,
		loading:    true,
	}
}

// Init 初始化
func (m TUIModel) Init() tea.Cmd {
	return tea.Batch(
		m.probe(),
	)
}

// probe 执行拨测
func (m TUIModel) probe() tea.Cmd {
	return func() tea.Msg {
		var results []ProbeResult
		if m.queryAll {
			results = ProbeAllRecordTypes(m.domain, m.dnsServers)
		} else {
			results = ProbeAll(m.domain, m.dnsServers)
		}
		return probeDoneMsg{results: results}
	}
}

// probeDoneMsg 拨测完成消息
type probeDoneMsg struct {
	results []ProbeResult
}

// Update 更新
func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case probeDoneMsg:
		m.results = msg.results
		m.loading = false
		return m, nil
	}
	return m, nil
}

// View 视图
func (m TUIModel) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	// 标题
	sb.WriteString(titleStyle.Render("╔══════════════════════════════════════════════════════════════════╗"))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("║                    DNS Probe Tool v1.0                          ║"))
	sb.WriteString("\n")
	sb.WriteString(titleStyle.Render("╚══════════════════════════════════════════════════════════════════╝"))
	sb.WriteString("\n\n")

	// 域名信息
	sb.WriteString(infoStyle.Render(fmt.Sprintf("  域名: %s\n", m.domain)))
	sb.WriteString("\n")

	// 加载中
	if m.loading {
		sb.WriteString(warningStyle.Render("  ⏳ 查询中...\n"))
		return sb.String()
	}

	// 结果
	for _, r := range m.results {
		sb.WriteString(fmt.Sprintf("┌─ DNS服务器: %s\n", r.DNSServer))
		sb.WriteString(fmt.Sprintf("│  查询耗时: %d ms\n", r.Latency.Milliseconds()))

		if r.IsPolluted {
			sb.WriteString(warningStyle.Render("│  ⚠️  检测到DNS污染\n"))
		}

		if r.Error != nil {
			sb.WriteString(errorStyle.Render(fmt.Sprintf("│  ❌ 错误: %s\n", r.Error.Error())))
			sb.WriteString("└──────────────────────────────────────────────────────────────────\n\n")
			continue
		}

		if len(r.Records) == 0 {
			sb.WriteString("│  (无记录)\n")
			sb.WriteString("└──────────────────────────────────────────────────────────────────\n\n")
			continue
		}

		sb.WriteString("│\n")
		sb.WriteString(fmt.Sprintf("│  %-10s %-6s %s\n", "类型", "TTL", "值"))
		sb.WriteString("│  " + strings.Repeat("─", 60) + "\n")

		for _, rec := range r.Records {
			ttlStr := formatTTL(rec.TTL)
			sb.WriteString(fmt.Sprintf("│  %-10s %-6s %s\n", rec.Type, ttlStr, rec.Value))
		}

		sb.WriteString("└──────────────────────────────────────────────────────────────────\n\n")
	}

	sb.WriteString(infoStyle.Render("  按 q 退出\n"))

	return sb.String()
}

// RunTUI 运行TUI界面
func RunTUI(domain string, dnsServers []string, queryAll bool) {
	m := NewTUIModel(domain, dnsServers, queryAll)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %s\n", err.Error())
		os.Exit(1)
	}
}
