package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// -- table of bookmarks
	cursor      int
	table       table.Model
	tableActive bool
	//-- New Bookmark text input
	focusIndex int
	inputs     []textinput.Model
	addToggle  bool
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func initialModel() Model {
	f, err := os.Open("C:\\Users\\f8col\\OneDrive\\Desktop\\Projects\\EBM\\src\\bm.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var rows []table.Row
	for _, record := range records {
		rows = append(rows, table.Row{
			record[0],
			record[1],
		})
	}

	columns := []table.Column{
		{Title: "Name", Width: 6},
		{Title: "Path", Width: 90},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	var inputs []textinput.Model
	var ti = textinput.New()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.CharLimit = 20
	ti.Placeholder = "Name"
	ti.Focus()

	inputs = append(inputs, ti)
	ti = textinput.New()
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.CharLimit = 20
	ti.Placeholder = "Path"
	inputs = append(inputs, ti)

	return Model{
		cursor:      0,
		table:       t,
		tableActive: true,
		inputs:      inputs,
		focusIndex:  0,
		addToggle:   false,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmdAdd tea.Cmd
	var update bool = true
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "a":
			m.addToggle = true
			m.inputs[0].Focus()
			m.tableActive = false
			update = false
		case "enter":
			if m.addToggle {
				if m.focusIndex == 0 {
					m.focusIndex++
					m.inputs[m.focusIndex].Focus()
					if m.focusIndex > 0 {
						m.inputs[m.focusIndex-1].Blur()
					}
				} else {
					m.table.SetRows(append(m.table.Rows(), []string{m.inputs[0].Value(), m.inputs[1].Value()}))
					m.inputs[0].Reset()
					m.inputs[1].Reset()
					m.tableActive = true
					m.addToggle = false
				}

			}

		}
	}
	if m.tableActive {
		m.table.Focus()
	} else {
		m.table.Blur()
	}
	if update {
		cmdAdd = m.updateInputs(msg)
	}
	m.table, cmd = m.table.Update(msg)

	return m, tea.Batch(cmd, cmdAdd)
}

func (m Model) View() string {
	var sb strings.Builder
	sb.WriteString(baseStyle.Render(m.table.View()) + "\n")

	if m.addToggle {
		for i := range m.inputs {
			sb.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				sb.WriteRune('\n')
			}
		}
	}
	return sb.String()
}

func main() {

	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
