package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Tabs          []string
	activeTab     int
	urlInput      string
	queues        []string
	selectedQueue int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			// Save the last download
			return m, tea.Quit
		case "right", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		case "down":
			if m.activeTab == 0 {
				m.selectedQueue = min(m.selectedQueue+1, len(m.queues)-1)
			}
			return m, nil
		case "up":
			if m.activeTab == 0 {
				m.selectedQueue = max(m.selectedQueue-1, 0)
			}
			return m, nil
		case "enter":
			if m.activeTab == 0 {
				// Handle the URL input and queue selection here
				fmt.Printf("URL: %s, Selected Queue: %s\n", m.urlInput, m.queues[m.selectedQueue])
				m.urlInput = ""
			}

			return m, nil
		case "backspace":
			if m.activeTab == 0 && len(m.urlInput) > 0 {
				m.urlInput = m.urlInput[:len(m.urlInput)-1]
			}
		default:
			if m.activeTab == 0 {
				m.urlInput += keypress
			}
			return m, nil
		}
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

func (m model) View() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	// The content inside the window
	content := strings.Builder{}

	if m.activeTab == 0 {
		// URL
		urlInput := fmt.Sprintf("URL: %s", m.urlInput)
		content.WriteString(urlInput)
		content.WriteString("\n\n")

		// Queues list
		for i, queue := range m.queues {
			if i == m.selectedQueue {
				content.WriteString(fmt.Sprintf("> %s\n", queue))
			} else {
				content.WriteString(fmt.Sprintf(" %s\n", queue))
			}
		}
	} else {
		// content.WriteString(m.TabContent[m.activeTab])
	}
	// fmt.Println("------------------------------\n", content.String())
	doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(content.String()))
	return docStyle.Render(doc.String())
}

func Uimain() {
	tabs := []string{"New Download", "Queues", "Downloads"}
	queues := []string{"Queue 1", "Queue 2", "Queue 3"}
	m := model{Tabs: tabs, queues: queues}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
