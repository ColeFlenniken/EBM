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
	table table.Model
	//-- New Bookmark text input
	focusIndex int
	inputs     []textinput.Model
	inputMode  bool
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
		table:      t,
		inputs:     inputs,
		focusIndex: 0,
		inputMode:  false,
	}
}
func saveFile(bookmarks []table.Row) error {
	f, err := os.Create("C:\\Users\\f8col\\OneDrive\\Desktop\\Projects\\EBM\\src\\bm.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var sb strings.Builder
	for i := range len(bookmarks) {
		sb.WriteString(bookmarks[i][0] + "," + bookmarks[i][1] + "\n")
	}
	_, err = f.Write([]byte(sb.String()))
	if err != nil {
		log.Fatal(err)

	}
	return err
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

func (m *Model) inputModeUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return tea.Quit
		case "enter":
			log.Default().Print(m.inputs[m.focusIndex].Value())
			m.AdvanceInput()

		}
	}
	//placeholder to remove compile error
	return nil
}
func (m *Model) AdvanceInput() bool {
	if m.focusIndex < len(m.inputs)-1 && strings.Trim(m.inputs[m.focusIndex].Value(), " ") != "" {
		m.inputs[m.focusIndex].Blur()
		m.focusIndex += 1
		m.inputs[m.focusIndex].Focus()
		return false
	}

	if m.inputs[m.focusIndex].Value() == "" {
		return false
	}
	m.table.SetRows(append(m.table.Rows(), []string{m.inputs[0].Value(), m.inputs[1].Value()}))
	saveFile(m.table.Rows())
	m.ResetInput()
	m.table.Focus()
	m.focusIndex = 0
	m.inputMode = false
	return true

}
func (m *Model) ResetInput() {
	for i := range len(m.inputs) {
		m.inputs[i].Reset()
		m.inputs[i].Blur()
	}
	m.inputMode = false
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	var update bool = true

	if m.inputMode {
		var cmdAdd tea.Cmd = m.inputModeUpdate(msg)
		if cmdAdd != nil {
			return m, tea.Quit
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "a":
				m.inputMode = true
				m.inputs[0].Focus()
				m.table.Blur()
				update = false
			case "r":
				if len(m.table.Rows()) > 0 {
					var currNdx int = m.table.Cursor()
					m.table.SetRows(append(m.table.Rows()[:m.table.Cursor()], m.table.Rows()[m.table.Cursor()+1:]...))
					if currNdx == len(m.table.Rows()) {
						m.table.SetCursor(currNdx - 1)
					}
					saveFile(m.table.Rows())

				}
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	if update {
		cmd = tea.Batch(cmd, m.updateInputs(msg))
	}

	return m, cmd
}

func (m Model) View() string {

	var sb strings.Builder
	sb.WriteString(baseStyle.Render(m.table.View()) + "\n")
	sb.WriteString("a: add bookmark r: delete bookmark enter:  submit addition\n")
	if m.inputMode {
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
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
